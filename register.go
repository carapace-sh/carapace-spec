package spec

import (
	"fmt"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace/pkg/uid"
	"github.com/spf13/cobra"
)

func Register(cmd *cobra.Command) {
	carapace.Gen(cmd)

	carapaceCmd, _, err := cmd.Find([]string{"_carapace"}) // TODO provide access to it using `carapace.Gen`
	if err != nil {
		carapace.LOG.Println(err.Error())
		return // should never happen
	}

	macroCmd := &cobra.Command{
		Use: "macro",
		RunE: func(cmd *cobra.Command, args []string) error {
			switch len(args) {
			case 0:
				keys := make([]string, 0, len(macros))
				for k := range macros {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				for _, key := range keys {
					fmt.Fprintln(cmd.OutOrStdout(), "$"+key)
				}
			case 1:
				m, ok := macros[args[0]]
				if !ok {
					return fmt.Errorf("unknown macro: %v", args[0])
				}
				fmt.Fprintln(cmd.OutOrStdout(), m.Signature())
			default:
				mCmd := &cobra.Command{
					DisableFlagParsing: true,
				}
				carapace.Gen(mCmd).Standalone()
				carapace.Gen(mCmd).PositionalAnyCompletion(
					ActionMacro(args[0]),
				)
				carapace.LOG.Printf("%#v", args)
				mCmd.SetArgs(append([]string{"_carapace", "export", ""}, args[1:]...))
				mCmd.SetOut(cmd.OutOrStdout())
				mCmd.SetErr(cmd.ErrOrStderr())
				return mCmd.Execute()
			}
			return nil
		},
	}

	macroCmd.Flags().SetInterspersed(false)

	carapaceCmd.AddCommand(macroCmd)

	carapace.Gen(macroCmd).PositionalCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			vals := make([]string, 0, len(macros))
			for key := range macros {
				if after, ok := strings.CutPrefix(key, "_."); ok {
					vals = append(vals, "$"+executable()+"."+after)
				}
			}
			return carapace.ActionValues(vals...).MultiParts(".")
		}),
	)

	carapace.Gen(macroCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			if strings.HasPrefix(c.Args[0], "$") {
				return ActionMacroM(c.Args[0]).Shift(1)
			}
			return ActionMacroM("$_." + c.Args[0]).Shift(1)
		}),
	)
}
