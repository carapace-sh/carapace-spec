package spec

import (
	_ "embed"
	"testing"

	"gopkg.in/yaml.v3"
)

//go:embed example/run.yaml
var runSpec string

func TestRun(t *testing.T) {
	var command Command
	if err := yaml.Unmarshal([]byte(runSpec), &command); err != nil {
		t.Error(err)
	}

	cmd := command.ToCobra()
	cmd.SetArgs([]string{"script", "one", "two", "three"})
	if err := cmd.Execute(); err != nil {
		t.Error(err)
	}
}
