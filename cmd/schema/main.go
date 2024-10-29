package main

import (
	"fmt"
	"os"
	"path/filepath"

	spec "github.com/carapace-sh/carapace-spec"
	"github.com/invopop/jsonschema"
)

func main() {
	schema := jsonschema.Reflect(&spec.Command{})
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
