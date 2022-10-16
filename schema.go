package spec

import (
	"github.com/invopop/jsonschema"
)

// Schema returns a json schema with currently registered macros
func Schema() (string, error) {
	r := new(jsonschema.Reflector)
	if err := r.AddGoComments("github.com/rsteube/carapace-spec", "./"); err != nil {
		return "", err
	}
	schema := r.Reflect(&Command{})
	out, err := schema.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(out), nil
}
