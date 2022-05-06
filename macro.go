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
			if reflect.TypeOf(t).Kind() == reflect.String {
				reflect.ValueOf(&t).Elem().SetString(s)
			} else {
				if err := yaml.Unmarshal([]byte(s), &t); err != nil {
					return carapace.ActionMessage(err.Error())
				}
			}
			return f(t)
		},
		Signature: func() string { return signature(new(T)) },
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
		Signature: func() string { return fmt.Sprintf("[%v]", signature(new(T))) },
	}
}

func signature(i interface{}) string {
	elem := reflect.ValueOf(i).Elem()
	switch elem.Kind() {
	case reflect.Struct:
		out, err := yaml.Marshal(i)
		if err != nil {
			return err.Error()
		}
		lines := strings.Split(string(out), "\n")
		return fmt.Sprintf("{%v}", strings.Join(lines[:len(lines)-1], ", "))

	case reflect.Slice:
		ptr := reflect.New(elem.Type().Elem()).Interface()
		return fmt.Sprintf("[%v]", signature(ptr))

	case reflect.String:
		return `""`

	default:
		return fmt.Sprintf("%v", reflect.ValueOf(i).Elem())
	}
}
