package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"gopkg.in/yaml.v3"
)

//go:embed example/run.yaml
var runSpec string

func TestRunAlias(t *testing.T) {
	var command Command
	if err := yaml.Unmarshal([]byte(runSpec), &command); err != nil {
		t.Error(err)
	}

	cmd := command.ToCobra()
	cmd.SetArgs([]string{"alias", "one"})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("alias", "").
			Expect(carapace.ActionValues("one", "two").
				Usage("alias ARG"))
	})
}

func TestRunLegacy(t *testing.T) {
	var command Command
	if err := yaml.Unmarshal([]byte(runSpec), &command); err != nil {
		t.Error(err)
	}

	cmd := command.ToCobra()
	cmd.SetArgs([]string{"legacy", "one"})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("legacy", "").
			Expect(carapace.ActionValues("one", "two").
				Usage("legacy ARG"))
	})
}

func TestRunMacro(t *testing.T) {
	var command Command
	if err := yaml.Unmarshal([]byte(runSpec), &command); err != nil {
		t.Error(err)
	}

	cmd := command.ToCobra()
	cmd.SetArgs([]string{"macro", "one"})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("macro", "").
			Expect(carapace.ActionValues("one", "two").
				Usage("macro ARG"))
	})
}

func TestRunScript(t *testing.T) {
	var command Command
	if err := yaml.Unmarshal([]byte(runSpec), &command); err != nil {
		t.Error(err)
	}

	cmd := command.ToCobra()
	cmd.SetArgs([]string{"script", "one"})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("script", "").
			Expect(carapace.ActionValues("one", "two").
				Usage("script ARG"))
	})
}
