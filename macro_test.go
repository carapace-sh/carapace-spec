package spec

import (
	"testing"

	"github.com/rsteube/carapace"
)

type Arg struct {
	Name   string
	Option bool
}

func (a Arg) Default() Arg {
	a.Option = true
	return a
}

func TestSignature(t *testing.T) {
	signature := MacroI(func(a Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `{name: "", option: false}`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroI(func(a []Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `[{name: "", option: false}]`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroI(func(b bool) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `false`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroV(func(a ...Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `[{name: "", option: false}]`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroV(func(b ...bool) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `[false]`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroI(func(s string) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `""`; signature != expected {
		t.Errorf("should be: %v", expected)
	}
}

func TestDefault(t *testing.T) {
	var actual Arg
	m := MacroI(func(a Arg) carapace.Action { actual = a; return carapace.ActionValues() })

	if m.f(""); !actual.Option {
		t.Error("should be true (default)")
	}

	if m.f("{option: false}"); actual.Option {
		t.Error("should be false")
	}
}
