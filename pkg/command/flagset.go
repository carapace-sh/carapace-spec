package command

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type FlagSet map[string]Flag

type Extended struct {
	Description string `yaml:"description,omitempty" json:"description,omitempty" jsonschema_description:"Description of the flag"`
	Default     string `yaml:"default,omitempty" json:"default,omitempty" jsonschema_description:"Default value of the flag"`
	OptDefault  string `yaml:"optdefault,omitempty" json:"optdefault,omitempty" jsonschema_description:"Default value when the optional flag is present without an argument"`
	Deprecated  string `yaml:"deprecated,omitempty" json:"deprecated,omitempty" jsonschema_description:"Deprecation message of the flag"`
	Nargs       int    `yaml:"nargs,omitempty" json:"nargs,omitempty" jsonschema_description:"Amount of arguments consumed"`
}

func (fs FlagSet) MarshalYAML() (any, error) {
	m := make(map[string]any)

	for _, f := range fs {
		switch {
		case f.Default != "" || f.OptDefault != "" || f.Deprecated != "" || f.Nargs != 0:
			m[f.format()] = Extended{
				Description: f.Description,
				Default:     f.Default,
				OptDefault:  f.OptDefault,
				Deprecated:  f.Deprecated,
				Nargs:       f.Nargs,
			}
		default:
			m[f.format()] = f.Description
		}
	}
	return m, nil
}

func (fs *FlagSet) UnmarshalYAML(value *yaml.Node) error {
	m := make(map[string]any)
	if err := value.Decode(&m); err != nil {
		return err
	}

	flagSet := make(FlagSet)
	for k, v := range m {
		switch v := v.(type) {
		case string:
			f, err := parseFlag(k, v)
			if err != nil {
				return err
			}
			flagSet[f.Name()] = *f // TODO ref?

		case map[string]any:
			f, err := parseFlag(k, "")
			if err != nil {
				return err
			}
			f.Description, _ = v["description"].(string)
			f.Default, _ = v["default"].(string)
			f.OptDefault, _ = v["optdefault"].(string)
			f.Deprecated, _ = v["deprecated"].(string)
			f.Nargs, _ = v["nargs"].(int)

			flagSet[f.Name()] = *f // TODO ref?

		default:
			return errors.New("invalid type for FlagSet")
		}
	}
	*fs = flagSet
	return nil
}
