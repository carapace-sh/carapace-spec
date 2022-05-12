package spec

import (
	"strings"

	"github.com/rsteube/carapace"
)

func init() {
	addCoreMacro("directories", MacroN(carapace.ActionDirectories))
	addCoreMacro("noflag", MacroN(func() carapace.Action { return carapace.ActionValues() }).NoFlag())
	addCoreMacro("files", MacroV(carapace.ActionFiles))
	addCoreMacro("message", MacroI(carapace.ActionMessage))
	addCoreMacro("spec", MacroI(ActionSpec).NoFlag())
	addCoreMacro("", MacroI(func(s string) carapace.Action {
		return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return carapace.ActionExecCommand("sh", "-c", s)(func(output []byte) carapace.Action {
				lines := strings.Split(string(output), "\n")
				vals := make([]string, 0)
				for _, line := range lines {
					if line != "" {
						vals = append(vals, parseValue(line)...)
					}
				}
				return carapace.ActionStyledValuesDescribed(vals...)
			}).Invoke(c).ToA()
		})
	}))
}
