package command

import "github.com/spf13/cobra"

type run struct {
	r interface {
		parse() func(cmd *cobra.Command, args []string) error
	}
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
	return func(cmd *cobra.Command, args []string) error { return nil }
}

func Alias(command string, args ...string) run {
	return run{r: alias(append([]string{command}, args...))}
}
func Macro(s string) run {
	return run{r: macro(s)}
}
func Shebang(s string) run {
	return run{r: shebang(s)}
}
