package command

import (
	"gopkg.in/yaml.v3"
)

type FlagSet map[string]Flag

type extendedFlag struct {
	Description string
	Nargs       int
}

func (fs FlagSet) MarshalYAML() (any, error) {
	m := make(map[string]any)

	for _, f := range fs {
		switch {
		case f.Nargs != 0: // TODO other values causing extended version
			m[f.format()] = extendedFlag{
				Description: f.Description,
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
			flagSet[f.format()] = *f // TODO ref?

		case map[string]any:
			f, err := parseFlag(k, "")
			if err != nil {
				return err
			}
			f.Description, _ = v["description"].(string)
			f.Nargs, _ = v["nargs"].(int)

			flagSet[f.format()] = *f // TODO ref?

		default:
			panic("TODO fail") // TODO:w
		}
	}
	*fs = flagSet
	return nil
}
