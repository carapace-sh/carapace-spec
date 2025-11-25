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

	// TODO
	// flagSet := make(FlagSet)
	// for k, v := range m {

	// }
	return nil
}
