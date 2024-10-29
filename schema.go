package spec

import (
	_ "embed"
)

//go:embed schema.json
var schema string

//go:generate go run -C cmd/schema . ../../schema.json
func Schema() string {
	return schema
}
