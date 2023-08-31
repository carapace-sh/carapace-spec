package command

import (
	"errors"

	"gopkg.in/yaml.v3"
)

type Parsing int

const (
	DEFAULT          Parsing = iota // INTERSPERSED but allows implicit changes
	INTERSPERSED                    // mixed flags and positional arguments
	NON_INTERSPERSED                // flag parsing stopped after first positional argument
	DISABLED                        // flag parsing disabled
)

func (p Parsing) MarshalYAML() (interface{}, error) {
	switch p {
	case DEFAULT:
		return "", nil
	case INTERSPERSED:
		return "interspersed", nil
	case NON_INTERSPERSED:
		return "non-interspersed", nil
	case DISABLED:
		return "disabled", nil
	default:
		return "", errors.New("unknown parsing mode")
	}
}

func (p *Parsing) UnmarshalYAML(value *yaml.Node) error {
	switch value.Value {
	case "":
		*p = DEFAULT
	case "interspersed":
		*p = INTERSPERSED
	case "non-interspersed":
		*p = NON_INTERSPERSED
	case "disabled":
		*p = DISABLED
	default:
		return errors.New("unknown parsing mode")
	}
	return nil
}
