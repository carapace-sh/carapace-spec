package command

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestFlagSet(t *testing.T) {
	fs := FlagSet{
		"string": Flag{
			Longhand:    "string",
			Description: "some string flag",
			Required:    true,
			Repeatable:  true,
		},
		"complex": Flag{
			Longhand:    "complex",
			Shorthand:   "c",
			Description: "some complex flag",
			Value:       true,
			Nargs:       2,
		},
	}

	m, err := yaml.Marshal(fs)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(m))

	var actual FlagSet
	if err := yaml.Unmarshal(m, &actual); err != nil {
		t.Error(err)
	}

	fmt.Println("====")
	fmt.Printf("%#v", actual)
}
