package spec

import (
	_ "embed"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/internal/pflagfork"
	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/spf13/pflag"
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

func TestAddFlagToAppliesExtendedAttributes(t *testing.T) {
	fset := pflag.NewFlagSet("test", pflag.ContinueOnError)
	err := addFlagTo(command.Flag{
		Longhand:            "flag",
		Shorthand:           "f",
		Description:         "extended flag",
		Value:               true,
		Optarg:              true,
		Default:             "default value",
		OptDefault:          "optional default",
		Deprecated:          "use --other instead",
		ShorthandDeprecated: "use --flag instead",
	}, fset)
	if err != nil {
		t.Fatal(err)
	}

	flag := fset.Lookup("flag")
	if flag.DefValue != "default value" {
		t.Fatalf("expected DefValue to be applied, got %q", flag.DefValue)
	}
	if flag.Value.String() != "default value" {
		t.Fatalf("expected default value to be applied, got %q", flag.Value.String())
	}
	if flag.NoOptDefVal != "optional default" {
		t.Fatalf("expected NoOptDefVal to be applied, got %q", flag.NoOptDefVal)
	}
	if flag.Deprecated != "use --other instead" {
		t.Fatalf("expected Deprecated to be applied, got %q", flag.Deprecated)
	}
	if flag.ShorthandDeprecated != "use --flag instead" {
		t.Fatalf("expected ShorthandDeprecated to be applied, got %q", flag.ShorthandDeprecated)
	}
}

func TestAddFlagToAppliesDelimiter(t *testing.T) {
	fset := pflag.NewFlagSet("test", pflag.ContinueOnError)
	err := addFlagTo(command.Flag{
		Longhand:  "flag",
		Value:     true,
		Optarg:    true,
		Delimiter: ":",
	}, fset)
	if err != nil {
		t.Fatal(err)
	}
	if delimiter := (pflagfork.Flag{Flag: fset.Lookup("flag")}).OptargDelimiter(); delimiter != ':' {
		t.Fatalf("expected delimiter to be applied, got %q", delimiter)
	}
}
