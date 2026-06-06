package spec

import (
	"reflect"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/pkg/macro"
)

// TODO lots of overlapping `macro` terms. Macro[I|V|N] don't match the embedded generic `macro.Macro` well.
type Macro struct {
	Name        string                       `json:"name"`
	Description string                       `json:"description,omitempty"`
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

// AddMacro adds a custom macro with explicit name
func AddMacro(s string, m Macro, opts ...string) {
	if len(opts) > 0 {
		m.Description = opts[0]
	}
	if len(opts) > 1 {
		m.Example = strings.Join(opts[1:], "\n")
	}
	macros["_."+s] = m
}

// AddMacroI adds a custom macro inferring the name from the function.
// Strips the ".../actions/" path prefix and "Action" function prefix.
// First string arg is description, any further strings are joined with "\n" as example.
func AddMacroI[T any](f func(t T) carapace.Action, opts ...string) {
	AddMacro(macroName(f), MacroI(f), opts...)
}

// AddMacroV adds a custom macro inferring the name from the function.
// Strips the ".../actions/" path prefix and "Action" function prefix.
// First string arg is description, any further strings are joined with "\n" as example.
func AddMacroV[T any](f func(t ...T) carapace.Action, opts ...string) {
	AddMacro(macroName(f), MacroV(f), opts...)
}

// macroName extracts the function name and strips prefixes.
// For carapace-bin: "github.com/carapace-sh/carapace-bin/pkg/actions/tools/git.ActionRefs" -> "tools.git.Refs"
// For carapace-spec: "github.com/carapace-sh/carapace-spec.ActionSpec" -> "Spec"
func macroName(f any) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()

	// strip everything up to and including the last "/actions/" package
	// e.g. "github.com/carapace-sh/carapace-bin/pkg/actions/tools/git" -> "tools/git"
	if idx := strings.LastIndex(name, "/actions/"); idx >= 0 {
		name = name[idx+9:] // +9 to skip "/actions/"
	} else if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}

	// strip "Action" prefix from function name
	name = strings.TrimPrefix(name, "Action")

	// replace remaining slashes with dots for nested package paths
	name = strings.ReplaceAll(name, "/", ".")

	return name
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
