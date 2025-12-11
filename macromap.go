package spec

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"unicode"

	"github.com/carapace-sh/carapace-spec/pkg/macro"
)

type MacroMap macro.MacroMap[Macro]

// TODO experimental - internal use
func (m MacroMap) Format(pkg string) (string, error) {
	imports := make([]string, 0)
	macros := make([]string, 0)
	for _, name := range slices.Sorted(maps.Keys(m)) {
		macro := m[name]
		macroPkg, macroFunction, ok := strings.Cut(macro.Function, "#")
		if !ok {
			return "", fmt.Errorf("missing function: %#v", macro)
		}

		var macroType string
		if arg := macro.Args; strings.Contains(arg, ",") {
			macros = append(macros, "// TODO unsupported signature: "+macro.Args)
			continue
		} else if arg == "" {
			macroType = "MacroN"
		} else if strings.Contains(arg, "...") {
			macroType = "MacroV"
		} else {
			macroType = "MacroI"
		}

		imports = append(imports, fmt.Sprintf("%s %q", varName(name), macroPkg))
		macros = append(macros, fmt.Sprintf(`%q: {
	Name: %q,
	Description: %q,
	Example: %q,
	Function: %q,
	Macro: %s,
},`,
			macro.Name,
			macro.Name,
			macro.Description,
			macro.Example,
			macro.Function,
			fmt.Sprintf("spec.%s(%s.%s).Macro", macroType, varName(name), macroFunction),
		))
	}

	return fmt.Sprintf(`package %s

import(%s
	spec "github.com/carapace-sh/carapace-spec"
)

var Macros = spec.MacroMap{%s}
`, pkg, strings.Join(imports, "\n"), strings.Join(macros, "\n")), nil
}

func varName(name string) string {
	if name == "go" {
		return "_go"
	}
	if unicode.IsDigit([]rune(name)[0]) {
		name = "_" + name
	}
	return strings.NewReplacer(
		"-", "_",
		".", "_",
	).Replace(name)
}
