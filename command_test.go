package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/rsteube/carapace/pkg/style"
)

//go:embed example/interspersed.yaml
var interspersedSpec string

func TestInterspersed(t *testing.T) {
	sandboxSpec(t, interspersedSpec)(func(s *sandbox.Sandbox) {
		s.Run("--bool", "").
			Expect(carapace.ActionValues(
				"four",
				"five",
				"six",
			))

		s.Run("--bool", "-").
			Expect(carapace.ActionStyledValuesDescribed(
				"--string", "string flag", style.Blue,
			).NoSpace('.').
				Tag("flags"))

		s.Run("--bool", "four", "-").
			Expect(carapace.ActionValues())

	})
}
