package spec

import (
	"testing"

	"github.com/rsteube/carapace"
)

type Arg struct {
	Name   string
	Option bool
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

	signature = MacroVarI(func(a ...Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `[{name: "", option: false}]`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroVarI(func(b ...bool) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `[false]`; signature != expected {
		t.Errorf("should be: %v", expected)
	}

	signature = MacroI(func(s string) carapace.Action { return carapace.ActionValues() }).Signature()
	if expected := `""`; signature != expected {
		t.Errorf("should be: %v", expected)
	}
}
