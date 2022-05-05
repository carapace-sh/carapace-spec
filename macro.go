package spec

import (
	"fmt"
	"reflect"
	"regexp"

	"github.com/rsteube/carapace"
	"gopkg.in/yaml.v3"
)

type Macro func(string) carapace.Action

var macros = make(map[string]Macro)

func addCoreMacro(s string, m Macro) {
	macros[s] = m
}

func AddMacro(s string, m Macro) {
	macros["_"+s] = m
}

func parseMacro(s string) carapace.Action {
	r := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`)
	if !r.MatchString(s) {
		return carapace.ActionMessage(fmt.Sprintf("malformed macro: '%v'", s))
	}

	matches := findNamedMatches(r, s)
	if m, ok := macros[matches["macro"]]; !ok {
		return carapace.ActionMessage(fmt.Sprintf("unknown macro: '%v'", s))
	} else {
		return m(matches["arg"])
	}
}

func MacroN(f func() carapace.Action) Macro {
	return func(s string) carapace.Action {
		return f()
	}
}

func MacroI[T any](f func(t T) carapace.Action) Macro {
	return func(s string) carapace.Action {
		var t T
		if err := yaml.Unmarshal([]byte(s), &t); err != nil {
			return carapace.ActionMessage(err.Error())
		}
		return f(t)
	}
}

func MacroVarI[T any](f func(s ...T) carapace.Action) Macro {
	return func(s string) carapace.Action {
		if s == "" {
			return f()
		}

		var t []T
		if err := yaml.Unmarshal([]byte(s), &t); err != nil {
			return carapace.ActionMessage(fmt.Sprintf("malformed macro arg: '%v', expected '%v'", s, reflect.TypeOf(t)))
		}
		return f(t...)
	}
}
