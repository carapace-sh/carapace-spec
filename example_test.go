package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
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

func sandboxSpec(t *testing.T, spec string) (f func(func(s *sandbox.Sandbox))) {
	var command Command
	if err := yaml.Unmarshal([]byte(spec), &command); err != nil {
		panic(err.Error())
	}
	return sandbox.Command(t, func() *cobra.Command {
		return command.ToCobra()
	})
}
