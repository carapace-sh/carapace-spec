package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

//go:embed example/modifier.yaml
var modifierSpec string

func TestModifier(t *testing.T) {
	sandboxSpec(t, modifierSpec)(func(s *sandbox.Sandbox) {
		for _, command := range []string{"generic", "specific"} {

			s.Run(command, "--filter", "").
				Expect(carapace.ActionValues(
					"one",
					"three",
				).Usage("$filter"))

			s.Run(command, "--list", "").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Usage("$list"))

			s.Run(command, "--list", "one,").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Prefix("one,").
					Usage("$list"))

			s.Run(command, "--multiparts", "").
				Expect(carapace.ActionValues(
					"one/",
				).NoSpace('/').
					Usage("$multiparts"))

			s.Run(command, "--multiparts", "one/").
				Expect(carapace.ActionValues(
					"two/",
				).NoSpace('/').
					Prefix("one/").
					Usage("$multiparts"))

			s.Run(command, "--nospace", "").
				Expect(carapace.ActionValues(
					"one",
					"two/",
					"three,",
				).NoSpace('/', ',').
					Usage("$nospace"))

			s.Run(command, "--retain", "").
				Expect(carapace.ActionValues(
					"two",
				).Usage("$retain"))

			s.Run(command, "--split", "").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Suffix(" ").
					Usage("$split"))

			s.Run(command, "--split", "one ").
				Expect(carapace.ActionValues(
					"two",
					"three",
				).NoSpace().
					Prefix("one ").
					Suffix(" ").
					Usage("$split"))

			s.Run(command, "--splitp", "").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Suffix(" ").
					Usage("$splitp"))

			s.Run(command, "--splitp", "one ").
				Expect(carapace.ActionValues(
					"two",
					"three",
				).NoSpace().
					Prefix("one ").
					Suffix(" ").
					Usage("$splitp"))

			s.Run(command, "--splitp", "one two | ").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Prefix("one two | ").
					Suffix(" ").
					Usage("$splitp"))

			s.Run(command, "--splitp", "one two | one ").
				Expect(carapace.ActionValues(
					"two",
					"three",
				).NoSpace().
					Prefix("one two | one ").
					Suffix(" ").
					Usage("$splitp"))

			s.Run(command, "--style", "").
				Expect(carapace.ActionStyledValues(
					"one", style.Underlined,
					"two", style.Underlined,
					"three", style.Underlined,
				).Usage("$style"))

			s.Run(command, "--suffix", "").
				Expect(carapace.ActionValues(
					"apple",
					"melon",
					"orange",
				).Suffix("juice").
					Usage("$suffix"))

			s.Run(command, "--suppress", "").
				Expect(carapace.ActionValues().Usage("$suppress"))

			s.Run(command, "--tag", "").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).Tag("numbers").
					Usage("$tag"))

			s.Run(command, "--uniquelist", "").
				Expect(carapace.ActionValues(
					"one",
					"two",
					"three",
				).NoSpace().
					Usage("$uniquelist"))

			s.Run(command, "--uniquelist", "one,").
				Expect(carapace.ActionValues(
					"two",
					"three",
				).NoSpace().
					Prefix("one,").
					Usage("$uniquelist"))

			s.Run(command, "--usage", "").
				Expect(carapace.ActionValues().Usage("custom"))
		}
	})
}
