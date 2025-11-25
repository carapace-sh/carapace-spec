package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

//go:embed example/flag.yaml
var flagSpec string

func TestFlag(t *testing.T) {
	sandboxSpec(t, flagSpec)(func(s *sandbox.Sandbox) {
		s.Run("--hidden").
			Expect(carapace.ActionValues().
				NoSpace('.'))

		s.Run("--hidden-arg", "").
			Expect(carapace.ActionValues(
				"h1",
				"h2",
			).Usage("hidden with argument"))

		s.Run("--hidden-opt=").
			Expect(carapace.ActionValues(
				"ho1",
				"ho2",
			).Prefix("--hidden-opt=").
				Usage("hidden with optional argument"))

		s.Run("--repeatable", "--repeat").
			Expect(carapace.ActionValuesDescribed(
				"--repeatable", "repeatable",
			).NoSpace('.').
				Tag("longhand flags"))
	})
}
