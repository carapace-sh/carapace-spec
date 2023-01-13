package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/rsteube/carapace"
	spec "github.com/rsteube/carapace-spec"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var rootCmd = &cobra.Command{
	Use:  "carapace-spec",
	Long: "define simple completions using a spec file",
	Example: `  Spec completion:
    bash:       source <(carapace-spec example.yaml)
    elvish:     eval (carapace-spec example.yaml|slurp)
    fish:       carapace-spec example.yaml | source
    oil:        source <(carapace-spec example.yaml)
    nushell:    carapace-spec example.yaml | save example.nu ; nu -c 'source example.nu'
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
		if len(args) > 0 && args[0] == "-h" || args[0] == "--help" {
			cmd.Help()
			return nil
		}

		abs, err := filepath.Abs(args[0])
		if err != nil {
			return err
		}

		content, err := os.ReadFile(abs)
		if err != nil {
			return err
		}

		var specCmd spec.Command
		if err := yaml.Unmarshal(content, &specCmd); err != nil {
			return err
		}

		if flag := cmd.Flag("scrape"); flag != nil && flag.Changed {
			specCmd.Scrape()
		} else {
			specCmdCobra, err := specCmd.ToCobra()
			if err != nil {
				return err
			}
			bridgeCompletion(specCmdCobra, abs, args[1:]...)
		}
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.Flags().Bool("scrape", false, "scrape to go code")

	if len(os.Args) < 2 || os.Args[1] == "--scrape" {
		rootCmd.DisableFlagParsing = false
	}

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)

	carapace.Gen(rootCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if rootCmd.Flag("scrape").Changed {
				return carapace.ActionValues()
			}

			abs, err := filepath.Abs(c.Args[0])
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}

			content, err := os.ReadFile(abs)
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}

			var specCmd spec.Command
			if err := yaml.Unmarshal(content, &specCmd); err != nil {
				return carapace.ActionMessage(err.Error())
			}

			if len(c.Args) > 2 {
				c.Args = c.Args[3:]
			} else {
				c.Args[0] = "_carapace"
			}

			specCmdCobra, err := specCmd.ToCobra()
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return carapace.ActionExecute(specCmdCobra).Invoke(c).ToA()
		}),
	)

	carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
		if len(args) < 2 || args[0] == "--scrape" {
			cmd.DisableFlagParsing = false
		}
	})
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
	patched := strings.Replace(string(out), fmt.Sprintf("%v _carapace", executableName), fmt.Sprintf("%v '%v'", executableName, spec), -1)    // general callback
	patched = strings.Replace(patched, fmt.Sprintf("'%v', '_carapace'", executableName), fmt.Sprintf("'%v', '%v'", executableName, spec), -1) // xonsh callback
	fmt.Print(patched)
}
