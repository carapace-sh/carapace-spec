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
		s.Run("--filter", "").
			Expect(carapace.ActionValues(
				"one",
				"three",
			).Usage("$filter"))

		s.Run("--list", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Usage("$list"))

		s.Run("--list", "one,").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Prefix("one,").
				Usage("$list"))

		s.Run("--multiparts", "").
			Expect(carapace.ActionValues(
				"one/",
			).NoSpace('/').
				Usage("$multiparts"))

		s.Run("--multiparts", "one/").
			Expect(carapace.ActionValues(
				"two/",
			).NoSpace('/').
				Prefix("one/").
				Usage("$multiparts"))

		s.Run("--nospace", "").
			Expect(carapace.ActionValues(
				"one",
				"two/",
				"three,",
			).NoSpace('/', ',').
				Usage("$nospace"))

		s.Run("--retain", "").
			Expect(carapace.ActionValues(
				"two",
			).Usage("$retain"))

		s.Run("--split", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Suffix(" ").
				Usage("$split"))

		s.Run("--split", "one ").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).NoSpace().
				Prefix("one ").
				Suffix(" ").
				Usage("$split"))

		s.Run("--splitp", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Suffix(" ").
				Usage("$splitp"))

		s.Run("--splitp", "one ").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).NoSpace().
				Prefix("one ").
				Suffix(" ").
				Usage("$splitp"))

		s.Run("--splitp", "one two | ").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Prefix("one two | ").
				Suffix(" ").
				Usage("$splitp"))

		s.Run("--splitp", "one two | one ").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).NoSpace().
				Prefix("one two | one ").
				Suffix(" ").
				Usage("$splitp"))

		s.Run("--style", "").
			Expect(carapace.ActionStyledValues(
				"one", style.Underlined,
				"two", style.Underlined,
				"three", style.Underlined,
			).Usage("$style"))

		s.Run("--suffix", "").
			Expect(carapace.ActionValues(
				"apple",
				"melon",
				"orange",
			).Suffix("juice").
				Usage("$suffix"))

		s.Run("--suppress", "").
			Expect(carapace.ActionValues().Usage("$suppress"))

		s.Run("--tag", "").
			Expect(carapace.Batch(
				carapace.ActionValues("two").Tag("even numbers"),
				carapace.ActionValues("one", "three").Tag("odd numbers"),
			).ToA().
				Usage("$tag"))

		s.Run("--uniquelist", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			).NoSpace().
				Usage("$uniquelist"))

		s.Run("--uniquelist", "one,").
			Expect(carapace.ActionValues(
				"two",
				"three",
			).NoSpace().
				Prefix("one,").
				Usage("$uniquelist"))

		s.Run("--usage", "").
			Expect(carapace.ActionValues().Usage("custom"))
	})
}
