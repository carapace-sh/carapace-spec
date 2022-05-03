package spec

import "github.com/rsteube/carapace"

var macros = make(map[string]func(string) carapace.Action)

func AddMacro(s string, f func(string) carapace.Action) {
	macros[s] = f
}
