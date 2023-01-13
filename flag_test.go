package spec

import (
	"reflect"
	"testing"

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
		shorthand: "short",
		longhand:  "long",
		slice:     true,
		nonposix:  true,
	})

	test("nonposix mixed", "-short, --long", &flag{
		shorthand: "short",
		longhand:  "long",
		nonposix:  false,
	})
}
