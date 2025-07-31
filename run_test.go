package spec

import (
	_ "embed"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/sandbox"
	"github.com/carapace-sh/carapace/pkg/style"
)

//go:embed example/run.yaml
var runSpec string

func TestRunArray(t *testing.T) {
	if _, err := exec.LookPath("carapace"); err != nil {
		t.Skip(err.Error())
	}

	file, err := os.CreateTemp(os.TempDir(), "carapace-spec_TestRunArray")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := os.WriteFile(file.Name(), []byte("one\ntwo\n"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	runnableSpec(t, runSpec)(func(r runnable) {
		r.Run("alias", "array", file.Name()).
			Expect("two\n")
		r.Run("alias", "array", "-n1", file.Name()).
			Expect("two\n")
		r.Run("alias", "array", "--lines", "2", file.Name()).
			Expect("one\ntwo\n")
	})

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("alias", "array", file.Name()).
			Expect(carapace.ActionValues(filepath.Base(file.Name())).
				Prefix(os.TempDir() + "/").
				NoSpace('/').
				Tag("files"))
		s.Run("alias", "array", "--follow").
			Expect(carapace.ActionStyledValuesDescribed("--follow", "output appended data as the file grows", style.Carapace.FlagOptArg).
				NoSpace('.').
				Tag("longhand flags"))
	})
}

func TestRunString(t *testing.T) {
	if _, err := exec.LookPath("carapace"); err != nil {
		t.Skip(err.Error())
	}

	file, err := os.CreateTemp(os.TempDir(), "carapace-spec_TestRunString")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := os.WriteFile(file.Name(), []byte("one\ntwo\n"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	runnableSpec(t, runSpec)(func(r runnable) {
		r.Run("alias", "string", file.Name()).
			Expect("two\n")
		r.Run("alias", "string", "-n1", file.Name()).
			Expect("two\n")
		r.Run("alias", "string", "--lines", "2", file.Name()).
			Expect("one\ntwo\n")
	})

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("alias", "string", file.Name()).
			Expect(carapace.ActionValues(filepath.Base(file.Name())).
				Prefix(os.TempDir() + "/").
				NoSpace('/').
				Tag("files"))
		s.Run("alias", "string", "--follow").
			Expect(carapace.ActionStyledValuesDescribed("--follow", "output appended data as the file grows", style.Carapace.FlagOptArg).
				NoSpace('.').
				Tag("longhand flags"))
	})
}

func TestRunMacro(t *testing.T) {
	file, err := os.CreateTemp(os.TempDir(), "carapace-spec_TestRunMacro")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	if err := os.WriteFile(file.Name(), []byte("one\ntwo\n"), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	runnableSpec(t, runSpec)(func(r runnable) {
		r.Run("script", "macro", filepath.Base(file.Name())).
			Expect("two\n")
		r.Run("script", "macro", "-n1", filepath.Base(file.Name())).
			Expect("two\n")
		r.Run("script", "macro", "--lines", "2", filepath.Base(file.Name())).
			Expect("one\ntwo\n")
	})

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("script", "macro", filepath.Base(file.Name())).
			Expect(carapace.ActionValues(filepath.Base(file.Name())).
				NoSpace('/').
				Tag("files"))
	})
}

func TestRunShebang(t *testing.T) {
	runnableSpec(t, runSpec)(func(r runnable) {
		r.Run("script", "shebang", "one").
			Expect("one.suffix\n")
		r.Run("script", "shebang", "one", "two").
			Expect("one.suffix\ntwo.suffix\n")
		r.Run("script", "shebang", "one", "--suffix", ".backup", "two").
			Expect("one.backup\ntwo.backup\n")
	})

	sandboxSpec(t, runSpec)(func(s *sandbox.Sandbox) {
		s.Run("script", "shebang", "").
			Expect(carapace.ActionValues("one", "two"))

		s.Run("script", "shebang", "-").
			Expect(carapace.Batch(
				carapace.ActionStyledValuesDescribed("--suffix", "suffix to add", "blue").Tag("longhand flags"),
				carapace.ActionStyledValuesDescribed("-s", "suffix to add", "blue").Tag("shorthand flags"),
			).ToA().
				NoSpace('.'))

		s.Run("script", "shebang", "--suffix", "").
			Expect(carapace.ActionValues(".backup", ".copy").
				Usage("suffix to add"))
	})
}
