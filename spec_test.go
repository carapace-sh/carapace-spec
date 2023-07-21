package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/pflag"
)

//go:embed example/example.yaml
var example string

//go:embed example/nonposix.yaml
var nonposix string

func TestPosix(t *testing.T) {
	sandboxSpec(t, example)(func(s *sandbox.Sandbox) {
		s.Run("sub1", "--styled", "c").
			Expect(carapace.ActionStyledValuesDescribed(
				"cyan", "cyan", style.Cyan,
			).Usage("styled values"))

		s.Run("sub1", "--optarg=").
			Expect(carapace.ActionValues(
				"first",
				"second",
				"third",
			).Prefix("--optarg=").
				Usage("optarg flag"))

		s.Run("sub1", "--list", "a,b,").
			Expect(carapace.ActionValues(
				"a",
				"b",
				"c",
				"d",
			).Prefix("a,b,").
				NoSpace().
				Usage("list flag"))

		s.Run("sub1", "--repeatable", "--repeatable", "").
			Expect(carapace.Batch(
				carapace.ActionValues(
					"pos1A",
					"pos1B",
				),
				carapace.ActionStyledValuesDescribed(
					"subsub1", "sub sub command", style.Blue,
				).Tag("group3 commands"),
			).ToA())

		s.Run("sub1", "--", "").
			Expect(carapace.ActionValues(
				"dash1",
				"dash2",
			))

		s.Run("sub1", "--persistent", "p").
			Expect(carapace.ActionValues(
				"pos1A",
				"pos1B",
			))

		s.Run("sub1", "--env", "C_").
			Expect(carapace.ActionValues(
				"C_VALUE=C_",
			).Usage("env"))

		s.Run("sub1", "--sty").
			Expect(carapace.ActionValuesDescribed(
				"--styled", "styled values",
			).NoSpace('.').
				Style(style.Carapace.FlagArg).
				Tag("flags"))

		s.Run("hidden", "").
			Expect(carapace.ActionValues(
				"hPos1",
				"hPos2",
			))

		s.Run("hidden", "--hidden", "").
			Expect(carapace.ActionValues(
				"first",
				"second",
				"third",
			).Usage("hidden flag"))
	})
}

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
				Tag("flags"))

		s.Run("--m").
			Expect(carapace.ActionValuesDescribed(
				"--mixed", "mixed repeatable",
			).NoSpace('.').
				Style(style.Carapace.FlagNoArg).
				Tag("flags"))

		s.Run("-opt=").
			Expect(carapace.ActionValues(
				"1",
				"2",
				"3",
			).Prefix("-opt=").
				Usage("both nonposix"))
	})
}
