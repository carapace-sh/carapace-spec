package spec

import (
	"bytes"
	_ "embed"
	"strings"
	"testing"

	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

//go:embed example/example.yaml
var example string

//go:embed example/nonposix.yaml
var nonposix string

func TestSpec(t *testing.T) {
	if out := execute(t, example, "example", "sub1", "--styled", ""); !strings.Contains(out, "cyan") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--optarg="); !strings.Contains(out, "second") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--list", "a,b,"); !strings.Contains(out, "a,b,c") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--repeatable", "--repeatable", ""); !strings.Contains(out, "pos1A") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--", ""); !strings.Contains(out, "dash1") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--persistent", ""); !strings.Contains(out, "pos1") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "--env", "C_"); !strings.Contains(out, "C_VALUE=C_") {
		t.Error(out)
	}

	if out := execute(t, example, "example", "sub1", "", ""); !strings.Contains(out, "action.go") {
		t.Error(out)
	}

	if out := execute(t, example, "example", ""); !strings.Contains(out, "sub1") {
		t.Error(out)
	}
}

func TestSpecNonposix(t *testing.T) {
	if fs := (flagSet{pflag.NewFlagSet("test", pflag.PanicOnError)}); !fs.IsFork() {
		t.Skip("skip nonposix tests with spf13/pflag")
	}

	if out := execute(t, nonposix, "nonposix", ""); !strings.Contains(out, "a") {
		t.Error(out)
	}

	if out := execute(t, nonposix, "nonposix", "-"); !strings.Contains(out, "-styled") {
		t.Error(out)
	}

	if out := execute(t, nonposix, "nonposix", "-mx", "--"); !strings.Contains(out, "--mixed") {
		t.Error(out)
	}

	if out := execute(t, nonposix, "nonposix", "-opt="); !strings.Contains(out, "1") {
		t.Error(out)
	}
}

func execute(t *testing.T, spec string, args ...string) string {
	var stdout, stderr bytes.Buffer
	var c Command
	if err := yaml.Unmarshal([]byte(spec), &c); err != nil {
		t.Error(err.Error())
	}
	cmd, err := c.ToCobraE()
	if err != nil {
		t.Error(err.Error())
	}
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)
	cmd.SetArgs(append([]string{"_carapace", "export"}, args...))
	if err := cmd.Execute(); err != nil {
		t.Error(err.Error())
	}
	return stdout.String()
}
