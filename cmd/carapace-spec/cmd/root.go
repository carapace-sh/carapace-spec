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
	Use:                "carapace-spec",
	Short:              "",
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

		content, err := os.ReadFile(abs)
		if err != nil {
			return err
		}

		var specCmd spec.Command
		if err := yaml.Unmarshal(content, &specCmd); err != nil {
			return err
		}
		bridgeCompletion(specCmd.ToCobra(), abs, args[1:]...)
		return nil
	},
}

func Execute() error {
	return rootCmd.Execute()
}
func init() {
	carapace.Gen(rootCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)

	carapace.Gen(rootCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
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
			return carapace.ActionExecute(specCmd.ToCobra()).Invoke(c).ToA()
		}),
	)
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
	patched := strings.Replace(string(out), fmt.Sprintf("%v _carapace", executableName), fmt.Sprintf("%v %v", executableName, spec), -1)      // general callback
	patched = strings.Replace(patched, fmt.Sprintf("'%v', '_carapace'", executableName), fmt.Sprintf("'%v', '%v'", executableName, spec), -1) // xonsh callback
	fmt.Print(patched)
}
