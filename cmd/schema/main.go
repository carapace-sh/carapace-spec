package main

import (
	"fmt"
	"os"
	"path/filepath"

	spec "github.com/carapace-sh/carapace-spec"
	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/invopop/jsonschema"
)

func main() {
	schema := jsonschema.Reflect(&spec.Command{})

	// patch for dual type in FlagSet
	s, _ := schema.Definitions["Command"].Properties.Get("flags")
	*s = jsonschema.Schema{
		Ref:         "#/$defs/FlagSet",
		Description: s.Description,
	}
	s, _ = schema.Definitions["Command"].Properties.Get("persistentflags")
	*s = jsonschema.Schema{
		Ref:         "#/$defs/FlagSet",
		Description: s.Description,
	}

	delete(schema.Definitions, "Flag")
	schema.Definitions["FlagSet"] = &jsonschema.Schema{
		Type: "object",
		AdditionalProperties: &jsonschema.Schema{
			OneOf: []*jsonschema.Schema{
				jsonschema.Reflect(&command.Extended{}).Definitions["Extended"],
				{Type: "string"},
			},
		},
	}

	m, err := schema.MarshalJSON()
	if err != nil {
		panic(err.Error())
	}

	switch len(os.Args) {
	case 2:
		if path := os.Args[1]; filepath.Base(path) == "schema.json" {
			os.WriteFile(path, m, os.ModePerm)
		}
	default:
		fmt.Println(string(m))
	}
}
