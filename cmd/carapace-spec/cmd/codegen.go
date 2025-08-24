package cmd

import (
	"github.com/carapace-sh/carapace"
	"github.com/spf13/cobra"
)

var codegenCmd = &cobra.Command{
	Use:   "codegen spec",
	Short: "generate code for spec file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		command, err := loadSpec(args[0])
		if err != nil {
			return err
		}
		return command.Codegen()
	},
}

func init() {
	rootCmd.AddCommand(codegenCmd)

	carapace.Gen(codegenCmd).PositionalCompletion(
		carapace.ActionFiles(".yaml"),
	)
}
