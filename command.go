package spec

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

type Command struct {
	Name            string            `json:"name" jsonschema_description:"Name of the command"`
	Aliases         []string          `json:"aliases,omitempty" jsonschema_description:"Aliases of the command"`
	Description     string            `json:"description,omitempty" jsonschema_description:"Description of the command"`
	Group           string            `json:"group,omitempty" jsonschema_description:"Group of the command"`
	Flags           map[string]string `json:"flags,omitempty" jsonschema_description:"Flags of the command with their description"`
	PersistentFlags map[string]string `json:"persistentflags,omitempty" jsonschema_description:"Persistent flags of the command with their description"`
	Completion      struct {
		Flag          map[string][]action `json:"flag,omitempty" jsonschema_description:"Flag completion"`
		Positional    [][]action          `json:"positional,omitempty" jsonschema_description:"Positional completion"`
		PositionalAny []action            `json:"positionalany,omitempty" jsonschema_description:"Positional completion for every other position"`
		Dash          [][]action          `json:"dash,omitempty" jsonschema_description:"Dash completion"`
		DashAny       []action            `json:"dashany,omitempty" jsonschema_description:"Dash completion of every other position"`
	} `json:"completion,omitempty" jsonschema_description:"Subcommands of the command"`
	Commands []Command `json:"commands,omitempty" jsonschema_description:"Completion definition"`
}

func (c *Command) ToCobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:     c.Name,
		Aliases: c.Aliases,
		Short:   c.Description,
		GroupID: c.Group,
		Run:     func(cmd *cobra.Command, args []string) {},
	}
	carapace.Gen(cmd).Standalone()

	for id, description := range c.PersistentFlags {
		parseFlag(id, description).addTo(cmd.PersistentFlags())
	}

	for id, description := range c.Flags {
		parseFlag(id, description).addTo(cmd.Flags())
	}

	flagCompletions := make(carapace.ActionMap)
	for key, value := range c.Completion.Flag {
		flagCompletions[key] = parseAction(cmd, value)
	}
	carapace.Gen(cmd).FlagCompletion(flagCompletions)

	positionalCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Positional {
		positionalCompletions = append(positionalCompletions, parseAction(cmd, pos))
	}
	carapace.Gen(cmd).PositionalCompletion(positionalCompletions...)

	carapace.Gen(cmd).PositionalAnyCompletion(parseAction(cmd, c.Completion.PositionalAny))

	dashCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Dash {
		dashCompletions = append(dashCompletions, parseAction(cmd, pos))
	}
	carapace.Gen(cmd).DashCompletion(dashCompletions...)

	carapace.Gen(cmd).DashAnyCompletion(parseAction(cmd, c.Completion.DashAny))

	groups := make(map[string]bool)
	for _, subcmd := range c.Commands {
		if _, exists := groups[subcmd.Group]; !exists {
			cmd.AddGroup(&cobra.Group{ID: subcmd.Group})
			groups[subcmd.Group] = true
		}
		cmd.AddCommand(subcmd.ToCobra())
	}

	return cmd
}

func (c *Command) Scrape() {
	Scrape(c.ToCobra())
}
