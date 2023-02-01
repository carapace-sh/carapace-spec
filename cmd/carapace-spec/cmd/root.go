package cmd

import (
	"bytes"
	"errors"
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
		switch args[0] {
		case "-h", "--help":
			cmd.Help()
		case "--scrape":
			if len(args) < 2 {
				return errors.New("flag needs an argument: --scrape")
			}
			command, err := loadSpec(args[1])
			if err != nil {
				return err
			}
			command.Scrape()
		case "--run":
			command, err := loadSpec(args[1])
			if err != nil {
				return err
			}
			cobraCmd := command.ToCobra()
			cobraCmd.SetArgs(args[2:])
			cobraCmd.Execute()
		default:
			abs, err := filepath.Abs(args[0])
			if err != nil {
				return err
			}
			specCmd, err := loadSpec(abs)
			if err != nil {
				return err
			}
			bridgeCompletion(specCmd.ToCobra(), abs, args[1:]...)
		}
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

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	rootCmd.Flags().Bool("scrape", false, "scrape to go code")
	rootCmd.Flags().Bool("run", false, "run with given args")

	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)

	carapace.Gen(rootCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			switch {
			case c.Args[0] == "--run":
				if len(c.Args) > 1 {
					path := c.Args[1]
					c.Args = c.Args[2:]
					return spec.ActionSpec(path).Invoke(c).ToA()
				}
				return carapace.ActionFiles(".yaml")

			case c.Args[0] == "--scrape" && len(c.Args) == 1:
				return carapace.ActionFiles(".yaml")

			case !strings.HasPrefix(c.Args[0], "-"):
				path := c.Args[0]
				if len(c.Args) < 3 {
					c.Args[0] = "_carapace"
				} else {
					c.Args = c.Args[3:]
				}
				return spec.ActionSpec(path).Invoke(c).ToA()

			default:
				return carapace.ActionValues()
			}
		}),
	)

	carapace.Gen(rootCmd).PreRun(func(cmd *cobra.Command, args []string) {
		switch len(args) {
		case 0, 1:
			cmd.DisableFlagParsing = false

		default:
			cmd.Flags().Parse(args[:1]) // TODO unnecessary
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
