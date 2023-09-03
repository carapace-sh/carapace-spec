package spec

import (
	_ "embed"
	"testing"

	"github.com/rsteube/carapace/pkg/sandbox"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func sandboxSpec(t *testing.T, spec string) (f func(func(s *sandbox.Sandbox))) {
	var command Command
	if err := yaml.Unmarshal([]byte(spec), &command); err != nil {
		panic(err.Error())
	}
	return sandbox.Command(t, func() *cobra.Command {
		return command.ToCobra()
	})
}
