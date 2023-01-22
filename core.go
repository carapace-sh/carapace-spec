package spec

import (
	"runtime"
	"strings"

	"github.com/rsteube/carapace"
)

func init() {
	// modifiers added as dummy for completeness
	addCoreMacro("chdir", MacroI(func(s string) carapace.Action { return carapace.ActionValues() }))
	addCoreMacro("list", MacroI(func(s string) carapace.Action { return carapace.ActionValues() }))
	addCoreMacro("multiparts", MacroI(func(s string) carapace.Action { return carapace.ActionValues() }))
	addCoreMacro("nospace", MacroI(func(s string) carapace.Action { return carapace.ActionValues() }))
	addCoreMacro("uniquelist", MacroI(func(s string) carapace.Action { return carapace.ActionValues() }))

	addCoreMacro("directories", MacroN(carapace.ActionDirectories))
	addCoreMacro("files", MacroV(carapace.ActionFiles))
	addCoreMacro("message", MacroI(func(s string) carapace.Action { return carapace.ActionMessage(s) }))
	addCoreMacro("noflag", MacroN(func() carapace.Action { return carapace.ActionValues() }).NoFlag())
	addCoreMacro("spec", MacroI(ActionSpec).NoFlag())

	addCoreMacro("", MacroI(func(s string) carapace.Action {
		if runtime.GOOS == "windows" {
			return shell("pwsh", s)
		}
		return shell("sh", s)
	}))
	addCoreMacro("bash", MacroI(func(s string) carapace.Action { return shell("bash", s) }))
	addCoreMacro("elvish", MacroI(func(s string) carapace.Action { return shell("elvish", s) }))
	addCoreMacro("fish", MacroI(func(s string) carapace.Action { return shell("fish", s) }))
	addCoreMacro("ion", MacroI(func(s string) carapace.Action { return shell("ion", s) }))
	addCoreMacro("nu", MacroI(func(s string) carapace.Action { return shell("nu", s) }))
	addCoreMacro("osh", MacroI(func(s string) carapace.Action { return shell("osh", s) }))
	addCoreMacro("pwsh", MacroI(func(s string) carapace.Action { return shell("pwsh", s) }))
	addCoreMacro("sh", MacroI(func(s string) carapace.Action { return shell("sh", s) }))
	addCoreMacro("xonsh", MacroI(func(s string) carapace.Action { return shell("xonsh", s) }))
	addCoreMacro("zsh", MacroI(func(s string) carapace.Action { return shell("zsh", s) }))
}

func shell(shell, command string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		if runtime.GOOS == "windows" &&
			shell != "elvish" &&
			shell != "nu" &&
			shell != "pwsh" &&
			shell != "xonsh" {
			return carapace.ActionMessage("unsupported shell [%v]: %v", runtime.GOOS, shell)
		}

		return carapace.ActionExecCommand(shell, "-c", command)(func(output []byte) carapace.Action {
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

}
