package spec

import (
	"os"
	"path/filepath"
	"regexp"

	"github.com/carapace-sh/carapace"
	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Command command.Command

func (c Command) ToCobra() *cobra.Command {
	cmd, err := c.ToCobraE()
	if err != nil {
		cmd = &cobra.Command{
			Use:                c.Name,
			DisableFlagParsing: true,
			RunE:               func(cmd *cobra.Command, args []string) error { return err },
		}
		carapace.Gen(cmd).Standalone()
		carapace.Gen(cmd).PositionalAnyCompletion(
			carapace.ActionMessage(err.Error()),
		)
	}
	return cmd
}

func (c Command) ToCobraE() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:     c.Name,
		Aliases: c.Aliases,
		Short:   c.Description,
		GroupID: c.Group,
		Args:    cobra.ArbitraryArgs,
		Hidden:  c.Hidden,
		Run:     func(cmd *cobra.Command, args []string) {},
	}

	switch c.Parsing {
	case command.DISABLED:
		cmd.DisableFlagParsing = true
	case command.NON_INTERSPERSED:
		cmd.Flags().SetInterspersed(false)
	}

	carapace.Gen(cmd).Standalone()

	for _, f := range []func(cmd *cobra.Command) error{
		c.addFlags,
		c.addPersistentFlags,
		c.markFlagsExclusive,
		c.addRun,
		c.addFlagCompletion,
		c.addPositionalCompletion,
		c.addPositionalAnyCompletion,
		c.addDashCompletion,
		c.addDashAnyCompletion,
		c.addSubcommands,
		c.addAliasCompletion,
	} {
		if err := f(cmd); err != nil {
			return nil, err
		}
	}
	return cmd, nil
}

func (c Command) Codegen() error {
	cmd, err := c.ToCobraE()
	if err != nil {
		return err
	}
	return Codegen(cmd)
}

func (c Command) addPersistentFlags(cmd *cobra.Command) error {
	for id, description := range c.PersistentFlags {
		flag, err := parseFlag(id, description)
		if err != nil {
			return err
		}
		flag.addTo(cmd.PersistentFlags())
		if flag.required {
			cmd.MarkFlagRequired(flag.longhand)
		}
	}
	return nil
}

func (c Command) addFlags(cmd *cobra.Command) error {
	for id, description := range c.Flags {
		flag, err := parseFlag(id, description)
		if err != nil {
			return err
		}
		flag.addTo(cmd.Flags())
		if flag.required {
			cmd.MarkFlagRequired(flag.longhand)
		}
	}
	return nil
}

func (c Command) markFlagsExclusive(cmd *cobra.Command) error {
	for _, e := range c.ExclusiveFlags {
		cmd.MarkFlagsMutuallyExclusive(e...)
	}
	return nil
}

func (c Command) addFlagCompletion(cmd *cobra.Command) error {
	flagCompletions := make(carapace.ActionMap)
	for key, a := range c.Completion.Flag {
		flagCompletions[key] = NewAction(a).Parse(cmd)
	}
	carapace.Gen(cmd).FlagCompletion(flagCompletions)
	return nil
}

func (c Command) addRun(cmd *cobra.Command) error {
	if c.Run == "" {
		return nil
	}

	if len(c.Flags) == 0 && len(c.PersistentFlags) == 0 {
		cmd.DisableFlagParsing = true
	}

	cmd.RunE = run(c.Run).parse()
	return nil
}

func (c Command) addPositionalCompletion(cmd *cobra.Command) error {
	if len(c.Completion.Positional) == 0 {
		return nil
	}

	positionalCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Positional {
		positionalCompletions = append(positionalCompletions, NewAction(pos).Parse(cmd))
	}
	carapace.Gen(cmd).PositionalCompletion(positionalCompletions...)

	return nil
}

func (c Command) addPositionalAnyCompletion(cmd *cobra.Command) error {
	if len(c.Completion.PositionalAny) > 0 {
		carapace.Gen(cmd).PositionalAnyCompletion(NewAction(c.Completion.PositionalAny).Parse(cmd))
	}
	return nil
}

func (c Command) addDashCompletion(cmd *cobra.Command) error {
	if len(c.Completion.Dash) == 0 {
		return nil
	}

	dashCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Dash {
		dashCompletions = append(dashCompletions, NewAction(pos).Parse(cmd))
	}
	carapace.Gen(cmd).DashCompletion(dashCompletions...)
	return nil
}

func (c Command) addDashAnyCompletion(cmd *cobra.Command) error {
	if len(c.Completion.DashAny) > 0 {
		carapace.Gen(cmd).DashAnyCompletion(NewAction(c.Completion.DashAny).Parse(cmd))
	}
	return nil
}

func (c Command) addSubcommands(cmd *cobra.Command) error {
	groups := make(map[string]bool)
	for _, subcmd := range c.Commands {
		if subcmd.Group != "" {
			if _, exists := groups[subcmd.Group]; !exists {
				cmd.AddGroup(&cobra.Group{ID: subcmd.Group})
				groups[subcmd.Group] = true
			}
		}
		subcmdCobra, err := Command(subcmd).ToCobraE()
		if err != nil {
			return err
		}
		cmd.AddCommand(subcmdCobra)
	}
	return nil
}

func (c Command) addAliasCompletion(cmd *cobra.Command) error {
	if c.Run != "" &&
		len(c.Flags) == 0 &&
		len(c.PersistentFlags) == 0 &&
		len(c.Completion.Positional) == 0 &&
		len(c.Completion.PositionalAny) == 0 &&
		len(c.Completion.Dash) == 0 &&
		len(c.Completion.DashAny) == 0 {

		cmd.DisableFlagParsing = true
		carapace.Gen(cmd).PositionalAnyCompletion(
			carapace.ActionCallback(func(context carapace.Context) carapace.Action {
				switch {
				case regexp.MustCompile(`^\[.*\]$`).MatchString(string(c.Run)):
					var mArgs []string
					if err := yaml.Unmarshal([]byte(c.Run), &mArgs); err != nil {
						return carapace.ActionMessage(err.Error())
					}
					if len(mArgs) == 0 {
						return carapace.ActionMessage("empty alias: %#v", c.Run)
					}

					var err error
					for index, arg := range mArgs {
						if mArgs[index], err = context.Envsubst(arg); err != nil {
							return carapace.ActionMessage(err.Error())
						}
					}

					// TODO keep in sync with ActionCarapaceBin in carapace-bridge
					carapaceCmd := "carapace"
					if executable, err := os.Executable(); err == nil && filepath.Base(executable) == "carapace" {
						carapaceCmd = executable // workaround for sandbox tests: directly call executable which was built with "go run"
					}

					execArgs := []string{mArgs[0], "export", mArgs[0]}
					execArgs = append(execArgs, mArgs[1:]...)
					execArgs = append(execArgs, context.Args...)
					execArgs = append(execArgs, context.Value)
					return carapace.ActionExecCommand(carapaceCmd, execArgs...)(func(output []byte) carapace.Action {
						return carapace.ActionImport(output)
					})

				default:
					return carapace.ActionValues()
				}
			}),
		)
	}
	return nil
}
