package spec

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
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

			mArgs = append(mArgs, "-c", matches[3], "--")

		default:
			return fmt.Errorf("malformed macro: %#v", r)
		}

		execCmd := exec.Command(mCmd, append(mArgs, args...)...)
		execCmd.Stdin = cmd.InOrStdin()
		execCmd.Stdout = cmd.OutOrStdout()
		execCmd.Stderr = cmd.ErrOrStderr()
		if err := execCmd.Run(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				os.Exit(exitErr.ProcessState.ExitCode())
			}
			return err
		}
		return nil
	}
}
