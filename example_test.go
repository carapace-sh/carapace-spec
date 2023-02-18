package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
)

//go:embed example/runnable.yaml
var runnable string

func TestRunnable(t *testing.T) {
	sandbox.Package(t, "github.com/rsteube/carapace-spec/cmd/carapace-spec")(func(s *sandbox.Sandbox) {
		s.Files("runnable.yaml", runnable)

		s.Run("runnable.yaml", "ex").
			Expect(carapace.ActionValues("export"))

		s.Run("runnable.yaml", "export", "runnable", "").
			Expect(carapace.ActionValuesDescribed(
				"sub1", "alias",
				"sub2", "shell",
				"sub3", "shell with flags").
				Tag("commands"))
	})
}
