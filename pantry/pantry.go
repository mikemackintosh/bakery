package pantry

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/mikemackintosh/bakery/cli"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

type PantryInterface interface {
	Parse(*hcl.EvalContext) error
	Bake()
	Baked()
	Ready() bool
	GetDependencies() []string
	ValidateNotIf() bool
	ValidateOnlyIf() bool
}

var dependsOn = &hcldec.AttrSpec{
	Name:     "depends_on",
	Required: false,
	Type:     cty.String,
}

var defaultSpec = &hcldec.ObjectSpec{
	"depends_on": &hcldec.AttrSpec{
		Name:     "depends_on",
		Required: false,
		Type:     cty.String,
	},
	"not_if": &hcldec.AttrSpec{
		Name:     "not_if",
		Required: false,
		Type:     cty.String,
	},
	"only_if": &hcldec.AttrSpec{
		Name:     "only_if",
		Required: false,
		Type:     cty.String,
	},
	"user": &hcldec.AttrSpec{
		Name:     "user",
		Required: false,
		Type:     cty.String,
	},
}

// NewPantrySpec appends the default spec fields to a pantry item spec
func NewPantrySpec(spec *hcldec.ObjectSpec) *hcldec.ObjectSpec {
	// Loop through the default spec and append it to
	// the provided spec.
	for k, v := range *defaultSpec {
		(*spec)[k] = v
	}

	return spec
}

type PantryItem struct {
	DependsOn string  `json:"depends_on"`
	NotIf     *string `json:"not_if"`
	OnlyIf    *string `json:"only_if"`
	User      *string `json:"user"`
	IsPrepped bool
	IsBaked   bool
}

func (p *PantryItem) Baked() {
	p.IsBaked = true
}

func (p *PantryItem) Ready() bool {
	return p.IsBaked
}

func (p *PantryItem) GetDependencies() []string {
	var out []string
	if len(p.DependsOn) > 0 {
		deps := strings.Split(p.DependsOn, ",")
		for _, d := range deps {
			out = append(out, d)
		}
	}
	return out
}

func (p *PantryItem) Populate(cfg cty.Value, obj interface{}) error {
	cli.Debug(cli.DEBUG3, "\t#=> Populating Config", cfg)
	cli.Debug(cli.DEBUG2, "\t#=> Populating Receiving Object", obj)
	out, err := json.Marshal(ctyjson.SimpleJSONValue{Value: cfg})
	if err != nil {
		return err
	}

	cli.Debug(cli.DEBUG2, "\t#=> Compiled Object", string(out))
	err = json.Unmarshal(out, obj)
	if err != nil {
		return fmt.Errorf("Error compiling configuartion: %s", err)
	}

	cli.Debug(cli.DEBUG2, "\t#=> Resulting Object", obj)
	return nil
}

// ValidateNotIf returns true if met, or false if criteria failes
func (p *PantryItem) ValidateNotIf() bool {
	// TODO: clean this up and refactor it to make it re-usable
	if p.NotIf != nil {
		o, err := RunCommand([]string{"sh", "-c", *p.NotIf})
		if err != nil {
			cli.Debug(cli.ERROR, fmt.Sprintf("\t-> Error running not_if %s, response: %s", *p.NotIf, o.FormattedString()), err)
			return true
		}

		if o.ExitCode == 0 {
			cli.Debug(cli.INFO, "\t-> Skipping due to matched not_if", nil)
			return true
		}
	}

	return false
}

// ValidateOnlyIf returns true if met, or false if criteria failes
func (p *PantryItem) ValidateOnlyIf() bool {
	if p.OnlyIf != nil {
		o, err := RunCommand([]string{"sh", "-c", *p.OnlyIf})
		if err != nil {
			cli.Debug(cli.ERROR, fmt.Sprintf("\t-> Error running only_if %s, response: %s", *p.NotIf, o.FormattedString()), err)
			return true
		}

		if o.ExitCode != 0 {
			cli.Debug(cli.INFO, "\t-> Skipping due to matched only_if", nil)
			return true
		}
	}

	return false
}
