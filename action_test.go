package spec

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

func TestActionSpec(t *testing.T) {
	cmd := &cobra.Command{}
	carapace.Gen(cmd).PositionalCompletion(
		ActionSpec("./example/example.yaml"),
	)

	var stdout bytes.Buffer
	cmd.SetOut(&stdout)
	cmd.SetArgs([]string{"_carapace", "export", ""})
	cmd.Execute()

	if !strings.Contains(stdout.String(), "sub1") {
		t.Error("should contain sub1 subcommand")
	}
}
