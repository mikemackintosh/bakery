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
	DependsOn string `json:"depends_on"`
	NotIf     string `json:"not_if"`
	OnlyIf    string `json:"only_if"`
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
