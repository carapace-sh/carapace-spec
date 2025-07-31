package macro

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type MacroMap[T any] map[string]T

type UnknownError struct{ s string }

func (e *UnknownError) Error() string { return e.s }

func (m MacroMap[T]) Lookup(s string) (*T, error) {
	r := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`)

	matches := r.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("malformed macro: %#v", s)
	}

	if m, ok := m[matches[1]]; ok {
		return &m, nil
	}
	return nil, &UnknownError{fmt.Sprintf("unknown macro: %#v", s)}
}

type Macro[T any] struct {
	f func(string) (*T, error)
	s func() string
}

type Default[T any] interface {
	Default() T
}

func (m Macro[T]) Parse(s string) (*T, error) {
	r := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`)
	matches := r.FindStringSubmatch(s)
	if matches == nil {
		return nil, fmt.Errorf("malformed macro: '%v'", s)
	}
	return m.f(matches[3])
}

func (m Macro[T]) Signature() string { return m.s() }

// MacroN creates a macro without an argument
func MacroN[T any](f func() (*T, error)) Macro[T] {
	return Macro[T]{
		f: func(s string) (*T, error) {
			return f()
		},
		s: func() string { return "" },
	}
}

// MacroI creates a macro with an argument
func MacroI[A, T any](f func(arg A) (*T, error)) Macro[T] {
	return Macro[T]{
		f: func(s string) (*T, error) {
			var arg A
			switch reflect.TypeOf(arg).Kind() {
			case reflect.String:
				reflect.ValueOf(&arg).Elem().SetString(s)

			default:
				if err := yaml.Unmarshal([]byte(s), &arg); err != nil {
					return nil, err
				}
				if s == "" {
					if v, ok := any(arg).(Default[A]); ok {
						arg = v.Default()
					}
				}
			}
			return f(arg)
		},
		s: func() string { return signature(new(A)) },
	}
}

// MacroV creates a macro with a variable argument
func MacroV[A, T any](f func(args ...A) (*T, error)) Macro[T] {
	return Macro[T]{
		f: func(s string) (*T, error) {
			if s == "" {
				return f()
			}

			var args []A
			if err := yaml.Unmarshal([]byte(s), &args); err != nil {
				return nil, fmt.Errorf("malformed macro arg: '%v', expected '%v'", s, reflect.TypeOf(args))
			}
			return f(args...)
		},
		s: func() string { return fmt.Sprintf("[%v]", signature(new(A))) },
	}
}

func signature(i any) string {
	elem := reflect.ValueOf(i).Elem()
	switch elem.Kind() {
	case reflect.Struct:
		out, err := yaml.Marshal(i)
		if err != nil {
			return err.Error()
		}
		lines := strings.Split(string(out), "\n") // TODO enforcing `flow` style by patching struct tags possible?
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
