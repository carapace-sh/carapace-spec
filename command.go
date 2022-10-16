package spec

import (
	"github.com/invopop/jsonschema"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

// static value or macro
type action string

func (action) JSONSchema() *jsonschema.Schema {
	enum := make([]interface{}, 0, len(macros))
	for macro := range macros {
		enum = append(enum, "$"+macro) // TODO full signature as in `carapace --macros XX`
	}
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Action",
		Description: "A static value or a macro",
		Enum:        enum,
	}
}

type Command struct {
	// Name of the command
	Name string `json:"name"`
	// Aliases of the command
	Aliases []string `json:"aliases,omitempty"`
	// Description of the command
	Description string `json:"description,omitempty"`
	// Flags of the command with their description
	Flags map[string]string `json:"flags,omitempty"`
	// Persistent flags of the command with their description
	PersistentFlags map[string]string `json:"persistentflags,omitempty"`
	// Completion definition
	Completion struct {
		// Flag completion
		Flag map[string][]action `json:"flag,omitempty"`
		// Positional completion
		Positional [][]action `json:"positional,omitempty"`
		// Positional completion for every other position
		PositionalAny []action `json:"positionalany,omitempty"`
		// Dash completion
		Dash [][]action `json:"dash,omitempty"`
		// Dash completion for every other position
		DashAny []action `json:"dashany,omitempty"`
	} `json:"completion,omitempty"`
	// Subcommands of the command
	Commands []Command `json:"commands,omitempty"`
}

func (c *Command) ToCobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:     c.Name,
		Aliases: c.Aliases,
		Short:   c.Description,
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

	for _, subcmd := range c.Commands {
		cmd.AddCommand(subcmd.ToCobra())
	}

	return cmd
}
