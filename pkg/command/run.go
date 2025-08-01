package command

import (
	"errors"
	"os"
	"regexp"
	"strings"

	"github.com/carapace-sh/carapace/pkg/execlog"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type runnable interface {
	parse() func(cmd *cobra.Command, args []string) error
}

type run struct{ runnable }

func (r run) Parse() func(cmd *cobra.Command, args []string) error {
	if r.runnable == nil {
		return nil
	}
	return r.runnable.parse()
}

func (r run) MarshalYAML() ([]byte, error) {
	return yaml.Marshal(r.runnable)
}

func (r *run) UnmarshalYAML(value *yaml.Node) error {
	var a []string
	if err := value.Decode(&a); err == nil {
		if len(a) > 0 {
			r.runnable = alias(a)
		}
		return nil
	}

	var s string
	if err := value.Decode(&s); err == nil {
		switch {
		case strings.HasPrefix(s, "$"):
			r.runnable = macro(s)
			return nil
		case strings.HasPrefix(s, "#!"):
			r.runnable = script(s)
			return nil
		case strings.HasPrefix(s, "["):
			// TODO legacy alias
			if err := yaml.Unmarshal([]byte(s), &a); err == nil {
				r.runnable = alias(a)
				return nil
			}
		}
	}
	return errors.New("invalid type")
}

type alias []string

func (a alias) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error { return nil }
}

type macro string

func (m macro) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error { return nil }
}

type script string

func (s script) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		sb, _, _ := strings.Cut(string(s), "\n")
		r := regexp.MustCompile(`^#!(?P<command>[^ ]+)( (?P<arg>.*))?$`)

		matches := r.FindStringSubmatch(sb)
		if matches == nil {
			return errors.New("invalid shebang header") // TODO
		}

		file, err := os.CreateTemp(os.TempDir(), "carapace-spec_run")
		if err != nil {
			return err
		}
		defer os.Remove(file.Name())

		os.WriteFile(file.Name(), []byte(s), os.ModePerm) // TODO make only readable by current user

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
		return scriptCmd.Run()
	}
}

func Alias(command string, args ...string) run {
	return run{alias(append([]string{command}, args...))}
}
func Macro(s string) run {
	return run{macro(s)}
}
func Script(s string) run {
	return run{script(s)}
}
