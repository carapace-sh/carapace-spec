package spec

import (
	_ "embed"
	"os/exec"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

//go:embed example/core.yaml
var coreSpec string

func TestCore(t *testing.T) {
	sandboxSpec(t, coreSpec)(func(s *sandbox.Sandbox) {
		for _, shell := range []string{
			"bash",
			"cmd",
			"elvish",
			"fish",
			"nu",
			"osh",
			"pwsh",
			"sh",
			"xonsh",
			"zsh",
		} {
			t.Run(shell, func(t *testing.T) {
				if _, err := exec.LookPath(shell); err != nil {
					t.Skip(err.Error())
				}
				s.Run(shell, "").
					Expect(carapace.ActionValues(
						"one",
						"two",
					))
			})
		}
	})
}
