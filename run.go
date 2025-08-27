package spec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/internal/shebang"
	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/carapace-sh/carapace-spec/pkg/macro"
	"github.com/carapace-sh/carapace/pkg/execlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type run string

func (r run) Parse() func(cmd *cobra.Command, args []string) error {
	switch command.Run(r).Type() {
	case "macro":
		return r.parseMacro()
	case "script":
		return r.parseScript()
	case "alias":
		return r.parseAlias()
	default:
		return nil // TODO handle the error somehow (log or give feedback)
	}
}

func (r run) parseAlias() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		alias := make([]string, 0)
		if err := yaml.Unmarshal([]byte(r), &alias); err != nil {
			return err
		}
		if len(alias) < 1 {
			return fmt.Errorf("malformed alias: %#v", r)
		}

		context := r.context(cmd, args)
		var err error
		for index, arg := range alias[1:] {
			if alias[index+1], err = context.Envsubst(arg); err != nil {
				return err
			}
		}

		execCmd := execlog.Command(alias[0], append(alias[1:], args...)...)
		execCmd.Stdin = cmd.InOrStdin()
		execCmd.Stdout = cmd.OutOrStdout()
		execCmd.Stderr = cmd.ErrOrStderr()
		execCmd.Env = context.Env
		if err := execCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ProcessState.ExitCode())
			}
			return err
		}
		return nil
	}
}

func (r run) parseMacro() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		context := r.context(cmd, nil)
		context.Args = args // force context.Args contain all args (ignore Value)

		splitted := strings.Split(string(r), " ||| ")

		m, err := macro.MacroMap[Macro]{
			"": MacroI(func(s string) carapace.Action {
				if runtime.GOOS == "windows" {
					return runAction(cmd, "pwsh", s)
				}
				return runAction(cmd, "sh", s)
			}),
			"bash":   MacroI(func(s string) carapace.Action { return runAction(cmd, "bash", s) }),
			"cmd":    MacroI(func(s string) carapace.Action { return runAction(cmd, "cmd", s) }),
			"elvish": MacroI(func(s string) carapace.Action { return runAction(cmd, "elvish", s) }),
			"fish":   MacroI(func(s string) carapace.Action { return runAction(cmd, "fish", s) }),
			"ion":    MacroI(func(s string) carapace.Action { return runAction(cmd, "ion", s) }),
			"nu":     MacroI(func(s string) carapace.Action { return runAction(cmd, "nu", s) }),
			"osh":    MacroI(func(s string) carapace.Action { return runAction(cmd, "osh", s) }),
			"pwsh":   MacroI(func(s string) carapace.Action { return runAction(cmd, "pwsh", s) }),
			"sh":     MacroI(func(s string) carapace.Action { return runAction(cmd, "sh", s) }),
			"xonsh":  MacroI(func(s string) carapace.Action { return runAction(cmd, "xonsh", s) }),
			"zsh":    MacroI(func(s string) carapace.Action { return runAction(cmd, "zsh", s) }),
		}.Lookup(splitted[0])
		if err != nil {
			return err
		}

		action := m.Parse(splitted[0])
		for _, s := range splitted[1:] {
			if !strings.HasPrefix(s, "$chdir(") { // TODO only chdir modifier accepted at the moment
				return errors.New("invalid modifier")
			}
			action = modifier{action}.Parse(s)
		}

		action.Invoke(context) // run the command
		return nil             // TODO return err from ActionMessage?
	}
}

func (r run) context(cmd *cobra.Command, args []string) carapace.Context {
	context := carapace.NewContext(args...)
	cmd.Flags().VisitAll(func(f *pflag.Flag) { // VisitAll as Visit() skips changed persistent flags of parent commands
		if f.Changed {
			if slice, ok := f.Value.(pflag.SliceValue); ok {
				context.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), strings.Join(slice.GetSlice(), ","))
			} else {
				context.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
			}
		}
	})
	return context
}

func (r run) parseScript() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		shebang, err := shebang.Parse(string(r))
		if err != nil {
			return err
		}

		pattern := "carapace-spec_run_*"
		switch strings.TrimSuffix(filepath.Base(shebang.Command), ".exe") {
		case "cmd":
			pattern += ".cmd"
			shebang.Args = append(shebang.Args, "/c")
		case "pwsh":
			pattern += ".ps1"
			shebang.Args = append(shebang.Args, "-f")
		}

		file, err := os.CreateTemp(os.TempDir(), pattern)
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		if err := os.WriteFile(file.Name(), []byte(shebang.Script), 0600); err != nil {
			return err
		}

		context := r.context(cmd, args)
		scriptArgs := append(shebang.Args, file.Name())
		scriptArgs = append(scriptArgs, args...)

		scriptCmd := execlog.Command(shebang.Command, scriptArgs...)
		scriptCmd.Stdout = cmd.OutOrStdout()
		scriptCmd.Stderr = cmd.ErrOrStderr()
		scriptCmd.Stdin = cmd.InOrStdin()
		scriptCmd.Env = context.Env
		// TODO support dir potentially modified by `$chdir()` modifier
		return scriptCmd.Run()
	}
}

func runAction(cmd *cobra.Command, shell, command string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		args, err := shellArgs(shell, command, c.Args...) // c.Args contains all args here (as they can be empty)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		execCmd := execlog.Command(shell, args...)
		execCmd.Stdin = cmd.InOrStdin()    // TODO yuck
		execCmd.Stdout = cmd.OutOrStdout() // TODO yuck
		execCmd.Stderr = cmd.ErrOrStderr() // TODO yuck
		execCmd.Dir = c.Dir
		execCmd.Env = c.Env
		if err := execCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ProcessState.ExitCode())
			}
			println(err.Error())
			os.Exit(1) // TODO return ActionMessage?
			// return carapace.ActionMessage(err.Error())
		}
		return carapace.ActionValues()
	})
}
