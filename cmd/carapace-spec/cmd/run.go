package cmd

import (
	"github.com/carapace-sh/carapace"
	spec "github.com/carapace-sh/carapace-spec"
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run spec [arg]...",
	Short: "run spec",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command, err := loadSpec(args[0])
		if err != nil {
			return err
		}
		cobraCmd := command.ToCobra()
		cobraCmd.SetArgs(args[1:])
		return cobraCmd.Execute()
	},
}

func init() {
	runCmd.Flags().SetInterspersed(false)

	rootCmd.AddCommand(runCmd)

	carapace.Gen(runCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)

	carapace.Gen(runCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return spec.ActionSpec(c.Args[0]).Shift(1)
		}),
	)
}
