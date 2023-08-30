package spec

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

type Command struct {
	Name            string            `yaml:"name" jsonschema_description:"Name of the command"`
	Aliases         []string          `yaml:"aliases,omitempty" jsonschema_description:"Aliases of the command"`
	Description     string            `yaml:"description,omitempty" jsonschema_description:"Description of the command"`
	Group           string            `yaml:"group,omitempty" jsonschema_description:"Group of the command"`
	Hidden          bool              `yaml:"hidden,omitempty" jsonschema_description:"Hidden state of the command"`
	NonInterspersed bool              `yaml:"noninterspersed,omitempty" jsonschema_description:"Interspersed state of the command"`
	Flags           map[string]string `yaml:"flags,omitempty" jsonschema_description:"Flags of the command with their description"`
	PersistentFlags map[string]string `yaml:"persistentflags,omitempty" jsonschema_description:"Persistent flags of the command with their description"`
	ExclusiveFlags  [][]string        `yaml:"exclusiveflags,omitempty" jsonschema_description:"Flags that are mutually exclusive"`
	Run             run               `yaml:"run,omitempty" jsonschema_description:"Command or script to execute in runnable mode"`
	Completion      struct {
		Flag          map[string]action `yaml:"flag,omitempty" jsonschema_description:"Flag completion"`
		Positional    []action          `yaml:"positional,omitempty" jsonschema_description:"Positional completion"`
		PositionalAny action            `yaml:"positionalany,omitempty" jsonschema_description:"Positional completion for every other position"`
		Dash          []action          `yaml:"dash,omitempty" jsonschema_description:"Dash completion"`
		DashAny       action            `yaml:"dashany,omitempty" jsonschema_description:"Dash completion of every other position"`
	} `yaml:"completion,omitempty" jsonschema_description:"Completion definition"`
	Commands []Command `yaml:"commands,omitempty" jsonschema_description:"Subcommands of the command"`
}

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
	cmd.Flags().SetInterspersed(!c.NonInterspersed)

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
		c.disableFlagParsing,
	} {
		if err := f(cmd); err != nil {
			return nil, err
		}
	}
	return cmd, nil
}

func (c Command) Scrape() {
	cmd, err := c.ToCobraE()
	// TODO handle error
	if err == nil {
		Scrape(cmd)
	}
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
	for key, action := range c.Completion.Flag {
		flagCompletions[key] = action.parse(cmd)
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

	cmd.RunE = c.Run.parse()
	return nil
}

func (c Command) addPositionalCompletion(cmd *cobra.Command) error {
	if len(c.Completion.Positional) == 0 {
		return nil
	}

	positionalCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Positional {
		positionalCompletions = append(positionalCompletions, pos.parse(cmd))
	}
	carapace.Gen(cmd).PositionalCompletion(positionalCompletions...)

	return nil
}

func (c Command) addPositionalAnyCompletion(cmd *cobra.Command) error {
	if len(c.Completion.PositionalAny) > 0 {
		carapace.Gen(cmd).PositionalAnyCompletion(c.Completion.PositionalAny.parse(cmd))
	}
	return nil
}

func (c Command) addDashCompletion(cmd *cobra.Command) error {
	if len(c.Completion.Dash) == 0 {
		return nil
	}

	dashCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Dash {
		dashCompletions = append(dashCompletions, pos.parse(cmd))
	}
	carapace.Gen(cmd).DashCompletion(dashCompletions...)
	return nil
}

func (c Command) addDashAnyCompletion(cmd *cobra.Command) error {
	if len(c.Completion.DashAny) > 0 {
		carapace.Gen(cmd).DashAnyCompletion(c.Completion.DashAny.parse(cmd))
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
		subcmdCobra, err := subcmd.ToCobraE()
		if err != nil {
			return err
		}
		cmd.AddCommand(subcmdCobra)
	}
	return nil
}

func (c Command) disableFlagParsing(cmd *cobra.Command) error {
	for _, actions := range c.Completion.Flag {
		if actions.disableFlagParsing() {
			cmd.DisableFlagParsing = true
			return nil
		}
	}

	for _, actions := range c.Completion.Positional {
		if actions.disableFlagParsing() {
			cmd.DisableFlagParsing = true
			return nil
		}
	}

	if c.Completion.PositionalAny.disableFlagParsing() {
		cmd.DisableFlagParsing = true
		return nil
	}

	for _, actions := range c.Completion.Dash {
		if actions.disableFlagParsing() {
			cmd.DisableFlagParsing = true
			return nil
		}
	}

	if c.Completion.DashAny.disableFlagParsing() {
		cmd.DisableFlagParsing = true
		return nil
	}

	return nil
}
