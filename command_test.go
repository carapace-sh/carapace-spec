package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

//go:embed example/command.yaml
var commandSpec string

func TestCommand(t *testing.T) {
	sandboxSpec(t, commandSpec)(func(s *sandbox.Sandbox) {
		s.Run("name").
			Expect(carapace.ActionValues(
				"name",
			).Tag("other commands"))

		s.Run("usage", "").
			Expect(carapace.ActionValues().
				Usage("usage [-F file | -D dir]... [-f format] profile").
				Tag("other commands"))

		s.Run("description").
			Expect(carapace.ActionValuesDescribed(
				"description", "with description",
			).Tag("other commands"))

		s.Run("group").
			Expect(carapace.ActionStyledValues(
				"group", style.Blue,
			).Tag("grouped commands"))

		s.Run("a").
			Expect(carapace.ActionValues(
				"a",
				"al",
				"aliases",
			).Tag("other commands"))

		s.Run("hidden").
			Expect(carapace.ActionValues())

		s.Run("hidden", "").
			Expect(carapace.ActionValues(
				"p1",
				"positional1",
			))

		s.Run("parsing", "interspersed", "--bool", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			))

		s.Run("parsing", "interspersed", "--bool", "one", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("parsing", "interspersed", "--bool", "one", "--string", "").
			Expect(carapace.ActionValues().
				Usage("string flag"))

		s.Run("parsing", "interspersed", "--bool", "one", "--string", "", "").
			Expect(carapace.ActionValues(
				"one",
				"two",
				"three",
			))

		s.Run("parsing", "interspersed", "--bool", "p1", "--string", "s1", "--", "").
			Expect(carapace.ActionValues())

		s.Run("parsing", "non-interspersed", "--bool", "p1", "--").
			Expect(carapace.ActionValues())

		s.Run("parsing", "disabled", "--").
			Expect(carapace.ActionValues())

		s.Run("flags", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-b", "bool flag", style.Default,
					"-o", "shorthand and longhand with optional argument", style.Yellow,
					"-v", "shorthand with value", style.Blue,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--repeatable", "longhand repeatable", style.Default,
					"--optarg", "shorthand and longhand with optional argument", style.Yellow,
					"--required", "longhand required", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("persistentflags", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-p", "persistent flag", style.Default,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--persistent", "persistent flag", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("persistentflags", "subcommand", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-l", "local flag", style.Default,
					"-p", "persistent flag", style.Default,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--local", "local flag", style.Default,
					"--persistent", "persistent flag", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("persistentflags", "-p", "subcommand", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed(
					"-l", "local flag", style.Default,
				).Tag("shorthand flags"),
				carapace.ActionStyledValuesDescribed(
					"--local", "local flag", style.Default,
				).Tag("longhand flags"),
			).ToA().NoSpace('.'))

		s.Run("exclusiveflags", "--").
			Expect(carapace.ActionValuesDescribed(
				"--add", "add package",
				"--delete", "delete package",
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("exclusiveflags", "--add", "").
			Expect(carapace.ActionValues())

		s.Run("exclusiveflags", "--delete", "").
			Expect(carapace.ActionValues())
	})
}
