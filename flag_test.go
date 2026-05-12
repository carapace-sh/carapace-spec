package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/internal/pflagfork"
	commandspec "github.com/carapace-sh/carapace-spec/pkg/command"
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
	cmd, err := Command(commandspec.Command{
		Name: "flag-attributes",
		Flags: commandspec.FlagSet{
			"value": {
				Longhand:            "value",
				Shorthand:           "v",
				Description:         "flag with extended attributes",
				Value:               true,
				Optarg:              true,
				Default:             "configured",
				OptDefault:          "implicit",
				Deprecated:          "use --other",
				ShorthandDeprecated: "use --value",
				Delimiter:           ":",
			},
		},
	}).ToCobraE()
	if err != nil {
		t.Fatal(err)
	}

	flag := cmd.Flags().Lookup("value")
	if flag == nil {
		t.Fatal("expected value flag")
	}
	if flag.DefValue != "configured" {
		t.Fatalf("expected default %q, got %q", "configured", flag.DefValue)
	}
	if flag.NoOptDefVal != "implicit" {
		t.Fatalf("expected optional default %q, got %q", "implicit", flag.NoOptDefVal)
	}
	if flag.Deprecated != "use --other" {
		t.Fatalf("expected deprecation %q, got %q", "use --other", flag.Deprecated)
	}
	if flag.ShorthandDeprecated != "use --value" {
		t.Fatalf("expected shorthand deprecation %q, got %q", "use --value", flag.ShorthandDeprecated)
	}
	if actual := (pflagfork.Flag{Flag: flag}).OptargDelimiter(); actual != ':' {
		t.Fatalf("expected delimiter %q, got %q", ':', actual)
	}
}
