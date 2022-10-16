package spec

import (
	"github.com/invopop/jsonschema"
)

// Schema returns a json schema with currently registered macros
func Schema() (string, error) {
	schema := jsonschema.Reflect(&Command{})
	out, err := schema.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
