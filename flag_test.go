package spec

import (
	_ "embed"
	"reflect"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/spf13/pflag"
)

func TestParseFlag(t *testing.T) {
	test := func(
		usage string,
		id string,
		expected *flag,
	) {
		expected.usage = usage // skip usage test
		f, err := parseFlag(id, usage)
		if err != nil {
			t.Error(err.Error())
		}
		if !reflect.DeepEqual(f, expected) {
			t.Error(usage)
			t.Logf("expected: %#v", expected)
			t.Logf("actual:   %#v", f)
		}

		flagSet := pflag.NewFlagSet("test", pflag.PanicOnError)
		f.addTo(flagSet)
	}

	test("shorthand-only", "-s", &flag{
		shorthand: "s",
	})

	test("shorthand-only slice", "-s*", &flag{
		shorthand: "s",
		slice:     true,
	})

	test("shorthand-only slice value", "-s*=", &flag{
		shorthand: "s",
		slice:     true,
		value:     true,
	})

	test("shorthand-only value", "-s=", &flag{
		shorthand: "s",
		value:     true,
	})

	test("shorthand-only optarg", "-s?", &flag{
		shorthand: "s",
		value:     true,
		optarg:    true,
	})

	test("longhand-only", "--long", &flag{
		longhand: "long",
	})

	test("longhand-only value", "--long=", &flag{
		longhand: "long",
		value:    true,
	})

	test("longhand-only slice optarg", "--long*?", &flag{
		longhand: "long",
		optarg:   true,
		slice:    true,
		value:    true,
	})

	test("both", "-s, --long", &flag{
		shorthand: "s",
		longhand:  "long",
	})

	test("both value", "-s, --long=", &flag{
		shorthand: "s",
		longhand:  "long",
		value:     true,
	})

	test("both optarg", "-s, --long?", &flag{
		shorthand: "s",
		longhand:  "long",
		value:     true,
		optarg:    true,
	})

	test("both slice optarg", "-s, --long*?", &flag{
		shorthand: "s",
		longhand:  "long",
		value:     true,
		slice:     true,
		optarg:    true,
	})

	test("nonposix shorthand", "-short", &flag{
		shorthand: "short",
	})

	test("nonposix shorthand optarg", "-short?", &flag{
		shorthand: "short",
		value:     true,
		optarg:    true,
	})

	test("nonposix both", "-short, -long*", &flag{
		shorthand:       "short",
		longhand:        "long",
		slice:           true,
		nameAsShorthand: true,
	})

	test("nonposix mixed", "-short, --long", &flag{
		shorthand:       "short",
		longhand:        "long",
		nameAsShorthand: false,
	})
}

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
				Tag("flags"))
	})
}
