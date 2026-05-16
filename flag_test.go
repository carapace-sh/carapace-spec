package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/carapace-sh/carapace/pkg/sandbox"
)

//go:embed example/flag.yaml
var flagSpec string

func TestFlag(t *testing.T) {
	sandboxSpec(t, flagSpec)(func(s *sandbox.Sandbox) {
		s.Run("--hidden").
			Expect(carapace.ActionValues().
				NoSpace('.'))

		s.Run("--hidden-arg", "").
			Expect(carapace.ActionValues(
				"h1",
				"h2",
			).Usage("hidden with argument"))

		s.Run("--hidden-opt=").
			Expect(carapace.ActionValues(
				"ho1",
				"ho2",
			).Prefix("--hidden-opt=").
				Usage("hidden with optional argument"))

		s.Run("--repeatable", "--repeat").
			Expect(carapace.ActionValuesDescribed(
				"--repeatable", "repeatable",
			).NoSpace('.').
				Tag("longhand flags"))
	})
}

func TestExtendedFlagAttributes(t *testing.T) {
	specCmd := Command(command.Command{
		Name: "test",
		Flags: command.FlagSet{
			"mode": command.Flag{
				Longhand:            "mode",
				Shorthand:           "m",
				Description:         "select mode",
				Default:             "fast",
				OptDefault:          "auto",
				Deprecated:          "use --speed",
				ShorthandDeprecated: "use --mode",
				Value:               true,
				Optarg:              true,
			},
		},
	})

	cmd, err := specCmd.ToCobraE()
	if err != nil {
		t.Fatal(err)
	}

	flag := cmd.Flags().Lookup("mode")
	if flag.DefValue != "fast" {
		t.Errorf("expected default %q, got %q", "fast", flag.DefValue)
	}
	if flag.Value.String() != "fast" {
		t.Errorf("expected flag value %q, got %q", "fast", flag.Value.String())
	}
	if flag.NoOptDefVal != "auto" {
		t.Errorf("expected optional default %q, got %q", "auto", flag.NoOptDefVal)
	}
	if flag.Deprecated != "use --speed" {
		t.Errorf("expected deprecation %q, got %q", "use --speed", flag.Deprecated)
	}
	if flag.ShorthandDeprecated != "use --mode" {
		t.Errorf("expected shorthand deprecation %q, got %q", "use --mode", flag.ShorthandDeprecated)
	}
}
