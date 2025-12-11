package spec

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/pkg/macro"
)

// TODO lots of overlapping `macro` terms. Macro[I|V|N] don't match the embedded generic `macro.Macro` well.
type Macro struct {
	Name        string                       `json:"name"`
	Description string                       `json:"descriptions,omitempty"`
	Example     string                       `json:"example,omitempty"`
	Function    string                       `json:"function,omitempty"`
	Args        string                       `json:"args,omitempty"` // TODO shouldn't be necessary
	Macro       macro.Macro[carapace.Action] `json:"-"`              // TODO public?
}

func (m Macro) Parse(s string) carapace.Action {
	a, err := m.Macro.Parse(s)
	if err != nil {
		return carapace.ActionMessage(err.Error())
	}
	return *a
}

func (m Macro) Signature() string {
	return m.Macro.Signature()
}

var macros = make(macro.MacroMap[Macro])

func addCoreMacro(s string, m Macro) {
	macros[s] = m
}

// AddMacro adds a custom macro
func AddMacro(s string, m Macro) {
	// TODO is the underscore prefix still a good idea?
	macros["_."+s] = m
}

func MacroN(f func() carapace.Action) Macro {
	return Macro{
		Macro: macro.MacroN(func() (*carapace.Action, error) {
			a := f()
			return &a, nil
		}),
	}
}

func MacroI[T any](f func(t T) carapace.Action) Macro {
	return Macro{
		Macro: macro.MacroI(func(t T) (*carapace.Action, error) {
			a := f(t)
			return &a, nil
		}),
	}
}

func MacroV[T any](f func(t ...T) carapace.Action) Macro {
	return Macro{
		Macro: macro.MacroV(func(t ...T) (*carapace.Action, error) {
			a := f(t...)
			return &a, nil
		}),
	}
}
