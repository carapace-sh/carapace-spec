package command

import (
	"testing"

	"github.com/carapace-sh/carapace/pkg/assert"
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
			Longhand:            "complex",
			Shorthand:           "c",
			Description:         "some complex flag",
			Value:               true,
			Optarg:              true,
			Nargs:               2,
			Default:             "default-value",
			OptDefault:          "optional-default",
			Deprecated:          "use --string",
			ShorthandDeprecated: "use --complex",
			Delimiter:           ":",
		},
	}

	expected := `--string*!: some string flag
-c, --complex?:
    description: some complex flag
    nargs: 2
    default: default-value
    optdefault: optional-default
    deprecated: use --string
    shorthanddeprecated: use --complex
    delimiter: ':'
`
	m, err := yaml.Marshal(fs)
	if err != nil {
		t.Error(err)
	}

	assert.Equal(t, string(m), expected)

	var actual FlagSet
	if err := yaml.Unmarshal(m, &actual); err != nil {
		t.Error(err)
	}
	assert.Equal(t, actual, fs)
}
