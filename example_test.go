package spec

import (
	"bytes"
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace/pkg/assert"
	"github.com/carapace-sh/carapace/pkg/sandbox"
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

func runnableSpec(t *testing.T, spec string) func(func(r runnable)) {
	return func(f func(r runnable)) {
		var command Command
		if err := yaml.Unmarshal([]byte(spec), &command); err != nil {
			t.Error(err)
		}
		f(runnable{t, command})
	}
}

type runnable struct {
	t       *testing.T
	command Command
}

type runnableResult struct {
	t      *testing.T
	actual string
}

func (r runnable) Run(args ...string) runnableResult {
	cmd := r.command.ToCobra()
	var stdout, stderr bytes.Buffer
	cmd.SetArgs(args)
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	if err := cmd.Execute(); err != nil {
		r.t.Error(err)
	}
	return runnableResult{r.t, stdout.String()}
}

func (r runnableResult) Expect(expected string) {
	assert.Equal(r.t, expected, r.actual)
}
