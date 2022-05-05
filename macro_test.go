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
	// TODO verify
	signature := MacroI(func(a Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	println(signature)

	signature = MacroI(func(a []Arg) carapace.Action { return carapace.ActionValues() }).Signature()
	println(signature)

	signature = MacroI(func(b bool) carapace.Action { return carapace.ActionValues() }).Signature()
	println(signature)
}
