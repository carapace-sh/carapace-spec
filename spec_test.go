package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
	"github.com/spf13/pflag"
)

//go:embed example/nonposix.yaml
var nonposix string

func skipNonFork(t *testing.T) {
	if fs := (flagSet{pflag.NewFlagSet("test", pflag.PanicOnError)}); !fs.IsFork() {
		t.Skip("skip nonposix tests with spf13/pflag")
	}
}

func TestNonposix(t *testing.T) {
	skipNonFork(t)

	sandboxSpec(t, nonposix)(func(s *sandbox.Sandbox) {
		s.Run("a").
			Expect(carapace.ActionValues("a"))

		s.Run("-s").
			Expect(carapace.ActionValuesDescribed(
				"-styled", "nonposix shorthand").
				NoSpace('.').
				Style(style.Carapace.FlagArg).
				Tag("shorthand flags"))

		s.Run("--m").
			Expect(carapace.ActionValuesDescribed(
				"--mixed", "mixed repeatable",
			).NoSpace('.').
				Style(style.Carapace.FlagNoArg).
				Tag("longhand flags"))

		s.Run("-opt=").
			Expect(carapace.ActionValues(
				"1",
				"2",
				"3",
			).Prefix("-opt=").
				Usage("both nonposix"))
	})
}
