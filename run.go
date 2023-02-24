package spec

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type run string

func (r run) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		mCmd := ""
		mArgs := make([]string, 0)

		switch {
		case strings.HasPrefix(string(r), "["):
			if err := yaml.Unmarshal([]byte(r), &mArgs); err != nil {
				return err
			}
			if len(mArgs) < 1 {
				return fmt.Errorf("malformed alias: %#v", r)
			}

			mCmd = mArgs[0]
			mArgs = mArgs[1:]

		case strings.HasPrefix(string(r), "$"):
			matches := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`).FindStringSubmatch(string(r))
			if matches == nil {
				return fmt.Errorf("malformed macro: %#v", r)
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

			if mCmd != "pwsh" {
				mArgs = append(mArgs, "-c", matches[3], "--")
			} else {
				// pwsh handles arguments after `-c` differently (https://github.com/PowerShell/PowerShell/issues/13832)
				path, err := os.CreateTemp(os.TempDir(), "carapace-spec_run*.ps1")
				if err != nil {
					return err
				}
				defer os.Remove(path.Name())

				if err = os.WriteFile(path.Name(), []byte(matches[3]), 0700); err != nil {
					return err
				}
				if err := path.Close(); err != nil {
					return err
				}
				mArgs = append(mArgs, path.Name())

			}

		default:
			return fmt.Errorf("malformed macro: %#v", r)
		}

		context := carapace.NewContext(args...)
		cmd.Flags().Visit(func(f *pflag.Flag) {
			if slice, ok := f.Value.(pflag.SliceValue); ok {
				context.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), strings.Join(slice.GetSlice(), ","))
			} else {
				context.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
			}
		})
		var err error
		for index, mArg := range mArgs {
			mArgs[index], err = context.Envsubst(mArg)
			if err != nil {
				return err
			}
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
