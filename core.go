package spec

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
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
	addCoreMacro("executables", MacroV(carapace.ActionExecutables))
	addCoreMacro("message", MacroI(func(s string) carapace.Action { return carapace.ActionMessage(s) }))
	// TODO is there still use for this? addCoreMacro("noflag", MacroN(func() carapace.Action { return carapace.ActionValues() }).NoFlag())
	addCoreMacro("spec", MacroI(ActionSpec))

	addCoreMacro("", MacroI(func(s string) carapace.Action {
		if runtime.GOOS == "windows" {
			return shell("cmd", s)
		}
		return shell("sh", s)
	}))
	addCoreMacro("bash", MacroI(func(s string) carapace.Action { return shell("bash", s) }))
	addCoreMacro("cmd", MacroI(func(s string) carapace.Action { return shell("cmd", s) }))
	addCoreMacro("elvish", MacroI(func(s string) carapace.Action { return shell("elvish", s) }))
	addCoreMacro("fish", MacroI(func(s string) carapace.Action { return shell("fish", s) }))
	// addCoreMacro("ion", MacroI(func(s string) carapace.Action { return shell("ion", s) }))
	addCoreMacro("nu", MacroI(func(s string) carapace.Action { return shell("nu", s) }))
	addCoreMacro("osh", MacroI(func(s string) carapace.Action { return shell("osh", s) }))
	addCoreMacro("pwsh", MacroI(func(s string) carapace.Action { return shell("pwsh", s) }))
	addCoreMacro("sh", MacroI(func(s string) carapace.Action { return shell("sh", s) }))
	addCoreMacro("xonsh", MacroI(func(s string) carapace.Action { return shell("xonsh", s) }))
	addCoreMacro("zsh", MacroI(func(s string) carapace.Action { return shell("zsh", s) }))
}

func shell(shell, command string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		args, err := shellArgs(shell, command, c.Args...)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		return carapace.ActionExecCommand(shell, args...)(func(output []byte) carapace.Action {
			lines := strings.Split(string(output), "\n")
			batch := carapace.Batch()
			for _, line := range lines {
				switch shell {
				case "cmd", "pwsh":
					line = strings.TrimSuffix(line, "\r")
				}
				if line != "" {
					batch = append(batch, parseValue(line))
				}
			}
			return batch.ToA()
		}).Invoke(c).ToA()
	})
}

func shellArgs(shell, command string, arguments ...string) ([]string, error) {
	if runtime.GOOS == "windows" &&
		shell != "cmd" &&
		shell != "elvish" &&
		shell != "nu" &&
		shell != "pwsh" &&
		shell != "xonsh" {
		return nil, fmt.Errorf("unsupported shell [%v]: %v", runtime.GOOS, shell)
	}

	args := []string{"-c"}
	switch shell {
	case "cmd":
		args[0] = "/c"
		args = append(args, command)
		args = append(args, arguments...)
	case "nu":
		args = append(args, fmt.Sprintf("def --wrapped main [...args] { %v }; main %v", command, shlex.Join(arguments)))
	case "pwsh":
		args = append(args, command)
		args = append(args, arguments...)
	default:
		args = append(args, command)
		args = append(args, "--")
		args = append(args, arguments...)
	}
	return args, nil
}
