package spec

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	shlex "github.com/carapace-sh/carapace-shlex"
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
	case "alias": // legacy
		return r.parseAlias()
	default:
		return nil // TODO handle the error somehow (log or give feedback)
	}
}

func (r run) parseAlias() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
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

		mCmd := ""
		mArgs := make([]string, 0)
		if err := yaml.Unmarshal([]byte(r), &mArgs); err != nil {
			return err
		}
		if len(mArgs) < 1 {
			return fmt.Errorf("malformed alias: %#v", r)
		}

		mCmd = mArgs[0]
		mArgs = mArgs[1:]

		var err error
		for index, arg := range mArgs {
			if mArgs[index], err = context.Envsubst(arg); err != nil {
				return err
			}
		}

		execCmd := execlog.Command(mCmd, append(mArgs, args...)...)
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

		mCmd := ""
		mArgs := make([]string, 0)

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

		switch mCmd {
		case "nu", "pwsh":
			// nu and pwsh handle arguments after `-c` differently (https://github.com/PowerShell/PowerShell/issues/13832)

			suffix := ".nu"
			if mCmd == "pwsh" {
				suffix = ".ps1"
			}

			path, err := os.CreateTemp(os.TempDir(), "carapace-spec_run*"+suffix)
			if err != nil {
				return err
			}
			defer os.Remove(path.Name())

			if err = os.WriteFile(path.Name(), []byte(script), 0700); err != nil {
				return err
			}
			if err := path.Close(); err != nil {
				return err
			}
			mArgs = append(mArgs, path.Name())

		default:
			mArgs = append(mArgs, "-c", script, "--")
		}

		execCmd := exec.Command(mCmd, append(mArgs, args...)...)
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

type shebang struct {
	Command string   // interpreter
	Args    []string // optional arguments (deriving from the standard and allowing more than one)
	Script  string   // script without shebang header for compability
}

func (r run) parseShebang() (*shebang, error) {
	firstLine, script, ok := strings.Cut(string(r), "\n")
	if !ok {
		return nil, errors.New("missing shebang header")
	}

	re := regexp.MustCompile(`^#!(?P<command>[^ ]+)( (?P<arg>.*))?$`)
	matches := re.FindStringSubmatch(strings.TrimSpace(firstLine))
	if matches == nil {
		return nil, errors.New("invalid shebang header")
	}

	shebang := &shebang{
		Command: matches[1],
		Args:    []string{},
		Script:  script,
	}
	if matches[3] != "" {
		tokens, err := shlex.Split(matches[3])
		if err != nil {
			return nil, err
		}
		shebang.Args = tokens.Words().Strings() // optional args
	}

	return shebang, nil
}

func (r run) parseScript() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// TODO currently duplicated in each run type
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

		shebang, err := r.parseShebang()
		if err != nil {
			return err
		}

		file, err := os.CreateTemp(os.TempDir(), "carapace-spec_run_")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		os.WriteFile(file.Name(), []byte(shebang.Script), 0600) // TODO make only readable by current user

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
