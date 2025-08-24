package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/carapace-sh/carapace"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:  "carapace-spec spec [arg]...",
	Long: "define simple completions using a spec file",
	Example: `  Spec completion:
    bash:       source <(carapace-spec example.yaml)
    elvish:     eval (carapace-spec example.yaml|slurp)
    fish:       carapace-spec example.yaml | source
    oil:        source <(carapace-spec example.yaml)
    nushell:    carapace-spec example.yaml
    powershell: carapace-spec example.yaml | Out-String | Invoke-Expression
    tcsh:       eval ` + "`" + `carapace-spec example.yaml` + "`" + `
    xonsh:      exec($(carapace-spec example.yaml))
    zsh:        source <(carapace-spec example.yaml)
    `,
	Args:               cobra.MinimumNArgs(1),
	DisableFlagParsing: true,
	CompletionOptions: cobra.CompletionOptions{
		DisableDefaultCmd: true,
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}
		specCmd, err := loadSpec(abs)
		if err != nil {
			return err
		}
		bridgeCompletion(specCmd.ToCobra(), abs, args[1:]...)
		return nil
	},
}

func loadSpec(path string) (*spec.Command, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	content, err := os.ReadFile(abs)
	if err != nil {
		return nil, err
	}

	var specCmd spec.Command
	if err := yaml.Unmarshal(content, &specCmd); err != nil {
		return nil, err
	}
	return &specCmd, nil
}

func Execute(version string) error {
	rootCmd.Version = version
	return rootCmd.Execute()
}
func init() {
	rootCmd.Flags().SetInterspersed(false)

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)

	carapace.Gen(rootCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return spec.ActionSpec(c.Args[0]).Shift(1)
		}),
	)

	spec.AddMacro("Spec", spec.MacroI(spec.ActionSpec))
	spec.Register(rootCmd)
}

func bridgeCompletion(cmd *cobra.Command, spec string, args ...string) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	a := []string{"_carapace"}
	a = append(a, args...)
	cmd.SetArgs(a)
	cmd.Execute()

	w.Close()
	out := <-outC
	os.Stdout = old

	executable, err := os.Executable()
	if err != nil {
		panic(err.Error()) // TODO exit with error message
	}

	executableName := filepath.Base(executable)
	patched := strings.ReplaceAll(string(out), fmt.Sprintf("%v _carapace", executableName), fmt.Sprintf("%v '%v'", executableName, spec))    // general callback
	patched = strings.ReplaceAll(patched, fmt.Sprintf("'%v', '_carapace'", executableName), fmt.Sprintf("'%v', '%v'", executableName, spec)) // xonsh callback
	fmt.Print(patched)
}
