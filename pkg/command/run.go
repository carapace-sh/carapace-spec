package command

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/execlog"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type Run string

func Alias(s string, args ...string) (Run, error) {
	m, err := yaml.Marshal(append([]string{s}, args...)) // TODO ensure this is single line style
	if err != nil {
		return "", err
	}
	return Run(m), nil
}

func (r Run) Type() string { // TODO return custom type?
	switch {
	case strings.HasPrefix(string(r), "$"):
		return "macro"
	case strings.HasPrefix(string(r), "#!"):
		return "shebang" // shebang or script?
	case strings.HasPrefix(string(r), "["): // legacy
		return "alias"
	default:
		return ""
	}
}

func (r *Run) UnmarshalYAML(value *yaml.Node) error {
	if err := value.Decode(r); err == nil {
		return nil
	}

	var alias []string
	if err := value.Decode(&alias); err != nil {
		return err
	}

	m, err := yaml.Marshal(alias)
	if err != nil {
		return err
	}

	run := Run(m)
	r = &run
	return nil
}

func (r Run) Parse() func(cmd *cobra.Command, args []string) error {
	switch r.Type() {
	case "macro":
		return r.parseMacro()
	case "shebang":
		return r.parseShebang()
	case "alias": // legacy
		return r.parseAlias()
	default:
		return nil // TODO handle the error somehow (log or give feedback)
	}
}

func (r Run) parseAlias() func(cmd *cobra.Command, args []string) error {
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

func (r Run) parseMacro() func(cmd *cobra.Command, args []string) error {
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

func (r Run) parseShebang() func(cmd *cobra.Command, args []string) error {
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

		sb, _, _ := strings.Cut(string(r), "\n")
		re := regexp.MustCompile(`^#!(?P<command>[^ ]+)( (?P<arg>.*))?$`)

		matches := re.FindStringSubmatch(sb)
		if matches == nil {
			return errors.New("invalid shebang header") // TODO
		}

		file, err := os.CreateTemp(os.TempDir(), "carapace-spec_run")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		os.WriteFile(file.Name(), []byte(r), os.ModePerm) // TODO make only readable by current user

		scriptArgs := make([]string, 0)
		if matches[3] != "" {
			scriptArgs = append(scriptArgs, matches[3])
		}
		scriptArgs = append(scriptArgs, file.Name())
		scriptArgs = append(scriptArgs, args...)

		scriptCmd := execlog.Command(matches[1], scriptArgs...)
		scriptCmd.Stdout = cmd.OutOrStdout()
		scriptCmd.Stderr = cmd.ErrOrStderr()
		scriptCmd.Stdin = cmd.InOrStdin()
		scriptCmd.Env = context.Env
		// TODO support dir potentially modified by `$chdir()` modifier
		return scriptCmd.Run()
	}
}
