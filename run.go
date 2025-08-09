package spec

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/internal/shebang"
	"github.com/carapace-sh/carapace-spec/pkg/command"
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
			if alias[index], err = context.Envsubst(arg); err != nil {
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
		context := r.context(cmd, args)

		// TODO parse and execute using the core exec macros

		mCmd := ""

		matches := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`).FindStringSubmatch(string(r))
		if matches == nil {
			return fmt.Errorf("malformed macro: %#v", r)
		}

		script, err := context.Envsubst(matches[3])
		if err != nil {
			return err
		}

		mCmd = "sh"
		switch matches[1] {
		case "":
			if runtime.GOOS == "windows" {
				mCmd = "pwsh"
			}

		case "bash", "elvish", "fish", "ion", "nu", "osh", "pwsh", "sh", "xonsh", "zsh":
			mCmd = matches[1]

		default:
			return fmt.Errorf("unknown macro: %#v", matches[1])
		}

		args, err = shellArgs(mCmd, script, args...)
		if err != nil {
			return err
		}
		execCmd := exec.Command(mCmd, args...)
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

		file, err := os.CreateTemp(os.TempDir(), "carapace-spec_run_")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		os.WriteFile(file.Name(), []byte(shebang.Script), 0600)

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

// TODO  func runAction(shell, command string) carapace.Action {
// 	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
// 		args, err := shellArgs(shell, command, c.Args...)
// 		if err != nil {
// 			return carapace.ActionMessage(err.Error())
// 		}

// 		execCmd := execlog.Command(shell, args...)
// 		execCmd.Stdin = os.Stdin
// 		execCmd.Stdout = os.Stdout
// 		execCmd.Stderr = os.Stderr
// 		execCmd.Env = c.Env
// 		if err := execCmd.Run(); err != nil {
// 			if exitErr, ok := err.(*exec.ExitError); ok {
// 				os.Exit(exitErr.ProcessState.ExitCode())
// 			}
// 			return carapace.ActionMessage(err.Error())
// 		}
// 		return carapace.ActionValues()
// 	})
// }
