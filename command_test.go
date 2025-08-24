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
			).Usage("string flag"))

		s.Run("parsing", "interspersed", "--bool", "p1", "--string", "s1", "--", "").
			Expect(carapace.ActionValues(
				"d1",
				"dash1",
			))

		s.Run("parsing", "non-interspersed", "--bool", "p1", "--").
			Expect(carapace.ActionValues())

		s.Run("parsing", "disabled", "--").
			Expect(carapace.ActionValues())

		s.Run("flags", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--bool", "bool flag", style.Default,
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("persistentflags", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--bool", "bool flag", style.Default,
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("persistentflags", "subcommand", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--bool", "bool flag", style.Default,
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("persistentflags", "--bool", "subcommand", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("exclusiveflags", "--").
			Expect(carapace.ActionStyledValuesDescribed(
				"--bool", "bool flag", style.Default,
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("longhand flags"))

		s.Run("exclusiveflags", "--bool", "--").
			Expect(carapace.ActionValues().
				NoSpace('.').
				Tag("longhand flags"))

		s.Run("run", "shell", "--color=").
			Expect(carapace.ActionValues(
				"always",
				"auto",
				"never",
			).Prefix("--color=").
				Usage("colored diff"))
	})
}
