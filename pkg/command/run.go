package command

import (
	"errors"
	"regexp"
	"strings"

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

func (r *run) UnmarshalYAML(b []byte) error {
	var m any
	if err := yaml.Unmarshal(b, &m); err != nil {
		return err
	}

	switch m := m.(type) {
	case []string:
		r.runnable = alias(m)
	case string:
		switch {
		case strings.HasPrefix(m, "$"):
			r.runnable = macro(m)
		case strings.HasPrefix(m, "#!"):
			r.runnable = shebang(m)
		case strings.HasPrefix(m, "["):
			// TODO legacy alias
		}
	default:
		// invalid
	}
	return nil // TODO invalid?
}

type alias []string

func (a alias) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error { return nil }
}

type macro string

func (m macro) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error { return nil }
}

type shebang string

func (s shebang) parse() func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		sb, _, _ := strings.Cut(string(s), "\n")
		r := regexp.MustCompile(`^#!(?P<command>[^ ]+)( (?P<arg>.*))?$`)

		matches := r.FindStringSubmatch(sb)
		if matches == nil {
			return errors.New("TODO") // TODO
		}

		return nil
	}
}

func Alias(command string, args ...string) run {
	return run{alias(append([]string{command}, args...))}
}
func Macro(s string) run {
	return run{macro(s)}
}
func Shebang(s string) run {
	return run{shebang(s)}
}
