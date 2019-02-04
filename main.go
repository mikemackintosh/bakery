package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/mikemackintosh/bakery/cli"
	"github.com/mikemackintosh/bakery/config"
	"github.com/mikemackintosh/bakery/pantry"
	"github.com/zclconf/go-cty/cty"
)

type Variable struct {
	Name    string         `hcl:"name,label"`
	Default hcl.Attributes `hcl:"default,remain"`
}

type Bakery struct {
	Variables []*Variable   `hcl:"variable,block"`
	Dmgs      []*pantry.Dmg `hcl:"dmg,block"`
}

func main() {
	flag.Parse()

	// Load the configuration file
	err := config.NewFromFile(cli.FlagConfig)
	if err != nil {
		cli.ErrorAndExit(err)
	}

	// Make the temp file directory
	// TODO: refactor this out
	err = os.MkdirAll(config.Registry.TempDir, 0755)
	if err != nil {
		cli.ErrorAndExit(err)
	}

	// Start parsing
	p := hclparse.NewParser()
	file, diags := p.ParseHCLFile(cli.FlagRecipe)
	if len(diags) != 0 {
		for _, diag := range diags {
			fmt.Printf("- %s\n", diag)
		}
		return
	}

	body := file.Body

	var bakery Bakery
	diags = gohcl.DecodeBody(body, nil, &bakery)
	if len(diags) != 0 {
		for _, diag := range diags {
			fmt.Printf("decoding - %s\n", diag)
		}
		return
	}

	variables := map[string]cty.Value{}
	for _, v := range bakery.Variables {
		if len(v.Default) == 0 {
			continue
		}

		val, diags := v.Default["default"].Expr.Value(nil)
		if len(diags) != 0 {
			for _, diag := range diags {
				fmt.Printf("decoding - %s\n", diag)
			}
			return
		}

		variables[v.Name] = val
	}

	evalContext := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var": cty.ObjectVal(nil),
		},
	}

	//var runList map[string]interface{}
	for _, entry := range bakery.Dmgs {
		entry.Parse(evalContext)
		entry.Bake()
	}
}
