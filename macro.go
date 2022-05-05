package spec

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
	"gopkg.in/yaml.v3"
)

type Macro struct {
	Func      func(string) carapace.Action
	Signature func() string
}

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
		return m.Func(matches["arg"])
	}
}

func MacroN(f func() carapace.Action) Macro {
	return Macro{
		Func: func(s string) carapace.Action {
			return f()
		},
		Signature: func() string { return "" },
	}
}

func MacroI[T any](f func(t T) carapace.Action) Macro {
	return Macro{
		Func: func(s string) carapace.Action {
			var t T
			if err := yaml.Unmarshal([]byte(s), &t); err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return f(t)
		},
		Signature: func() string {
			var t interface{} = new(T)
			//if elem := reflect.TypeOf(t).Elem(); elem.Kind() == reflect.Slice {
			// TODO slice member
			//println(reflect.TypeOf(elem.Elem()).Kind().String())
			//t = reflect.New(reflect.TypeOf(elem.Elem()))
			//}

			out, err := yaml.Marshal(t)
			if err != nil {
				return err.Error()
			}
			lines := strings.Split(string(out), "\n")

			if reflect.ValueOf(t).Elem().Kind() == reflect.Struct {
				return fmt.Sprintf("{%v}", strings.Join(lines[:len(lines)-1], ", "))
			} else {
				return fmt.Sprintf("%v", strings.Join(lines[:len(lines)-1], ", "))
			}
		},
	}
}

func MacroVarI[T any](f func(s ...T) carapace.Action) Macro {
	return Macro{
		Func: func(s string) carapace.Action {
			if s == "" {
				return f()
			}

			var t []T
			if err := yaml.Unmarshal([]byte(s), &t); err != nil {
				return carapace.ActionMessage(fmt.Sprintf("malformed macro arg: '%v', expected '%v'", s, reflect.TypeOf(t)))
			}
			return f(t...)
		},
		Signature: func() string { return "TODO" },
	}
}
