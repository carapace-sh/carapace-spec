package spec

import (
	"encoding/json"
	"fmt"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace"
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
				exe := executable()
				keys := make([]string, 0, len(macros))
				for k := range macros {
					keys = append(keys, k)
				}
				sort.Strings(keys)

				type macroEntry struct {
					Name        string `json:"name"`
					Signature   string `json:"signature"`
					Description string `json:"description"`
					Version     string `json:"version"`
					Function    string `json:"function"`
				}

				type macroList struct {
					Version string       `json:"version"`
					Macros  []macroEntry `json:"macros"`
				}

				mainVersion := mainModuleVersion()

				entries := make([]macroEntry, 0, len(keys))
				for _, key := range keys {
					if strings.HasPrefix(key, "_") {
						m := macros[key]
						sig := m.Signature()
						if sig == "" {
							sig = "—"
						}
						pkgPath, _, _ := strings.Cut(m.Function, "#")
						entries = append(entries, macroEntry{
							Name:        exe + strings.TrimPrefix(key, "_"),
							Signature:   sig,
							Description: m.Description,
							Version:     resolveVersion(pkgPath, mainVersion),
							Function:    m.Function,
						})
					}
				}

				output, _ := json.MarshalIndent(macroList{
					Version: mainVersion,
					Macros:  entries,
				}, "", "  ")
				fmt.Fprintln(cmd.OutOrStdout(), string(output))
			case 1:
				m, ok := macros["_."+args[0]]
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
					ActionMacro("$_." + args[0]),
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
					vals = append(vals, after)
				}
			}
			return carapace.ActionValues(vals...).MultiParts(".")
		}),
	)

	carapace.Gen(macroCmd).PositionalAnyCompletion(
		carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return ActionMacro("$_." + c.Args[0]).Shift(1)
		}),
	)
}

func mainModuleVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		v := info.Main.Version
		if v != "" {
			return v
		}
	}
	return "unknown"
}

func resolveVersion(pkgPath, mainVersion string) string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if strings.HasPrefix(pkgPath, info.Main.Path+"/") {
			return mainVersion
		}
		for _, dep := range info.Deps {
			if pkgPath == dep.Path || strings.HasPrefix(pkgPath, dep.Path+"/") {
				return dep.Version
			}
		}
	}
	return ""
}
