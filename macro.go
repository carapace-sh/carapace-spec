package spec

import (
	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/pkg/macro"
)

type Macro struct {
	macro macro.Macro[carapace.Action]
}

func (m Macro) Parse(s string) carapace.Action {
	a, err := m.macro.Parse(s)
	if err != nil {
		return carapace.ActionMessage(err.Error())
	}
	return *a
}

func (m Macro) Signature() string {
	return m.macro.Signature()
}

var macros = make(macro.MacroMap[Macro])

func addCoreMacro(s string, m Macro) {
	macros[s] = m
}

// AddMacro adds a custom macro
func AddMacro(s string, m Macro) {
	macros["_."+s] = m
}

func MacroN(f func() carapace.Action) Macro {
	return Macro{
		macro: macro.MacroN(func() (*carapace.Action, error) {
			a := f()
			return &a, nil
		}),
	}
}

func MacroI[T any](f func(t T) carapace.Action) Macro {
	return Macro{
		macro: macro.MacroI(func(t T) (*carapace.Action, error) {
			a := f(t)
			return &a, nil
		}),
	}
}

func MacroV[T any](f func(t ...T) carapace.Action) Macro {
	return Macro{
		macro: macro.MacroV(func(t ...T) (*carapace.Action, error) {
			a := f(t...)
			return &a, nil
		}),
	}
}
