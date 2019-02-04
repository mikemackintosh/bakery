package pantry

import (
	"encoding/json"

	"github.com/hashicorp/hcl2/hcldec"
	"github.com/mikemackintosh/bakery/cli"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

/*
type PantryInterface interface {
	Parse()
	Prepare()
	Bake()
}
*/

// dependsOn is the global depends on definition
var dependsOn = &hcldec.AttrSpec{
	Name:     "depends_on",
	Required: false,
	Type:     cty.String,
}

type PantryItem struct{}

func (p PantryItem) Populate(cfg cty.Value, obj interface{}) error {
	cli.Debug(cli.DEBUG, "\t->Populating Config", cfg)
	cli.Debug(cli.DEBUG, "\t->Populating Receiving Object", obj)
	out, err := json.Marshal(ctyjson.SimpleJSONValue{cfg})
	if err != nil {
		return err
	}

	cli.Debug(cli.DEBUG, "\t->Compiled Object", string(out))
	err = json.Unmarshal(out, obj)
	if err != nil {
		return err
	}
	return nil
}
