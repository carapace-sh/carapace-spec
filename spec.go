package spec

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

type Command struct {
	Name            string
	Description     string
	Flags           map[string]string
	PersistentFlags map[string]string
	Completion      struct {
		Flag          map[string][]string
		Positional    [][]string
		PositionalAny []string
		Dash          [][]string
		DashAny       []string
	}
	Commands []Command
}

func (c *Command) ToCobra() *cobra.Command {
	cmd := &cobra.Command{
		Use:   c.Name,
		Short: c.Description,
		Run:   func(cmd *cobra.Command, args []string) {},
	}
	carapace.Gen(cmd).Standalone()

	for id, description := range c.PersistentFlags {
		parseFlag(cmd.PersistentFlags(), id, description)
	}

	for id, description := range c.Flags {
		parseFlag(cmd.Flags(), id, description)
	}

	flagCompletions := make(carapace.ActionMap)
	for key, value := range c.Completion.Flag {
		flagCompletions[key] = parseAction(value)
	}
	carapace.Gen(cmd).FlagCompletion(flagCompletions)

	positionalCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Positional {
		positionalCompletions = append(positionalCompletions, parseAction(pos))
	}
	carapace.Gen(cmd).PositionalCompletion(positionalCompletions...)

	carapace.Gen(cmd).PositionalAnyCompletion(parseAction(c.Completion.PositionalAny))

	dashCompletions := make([]carapace.Action, 0)
	for _, pos := range c.Completion.Dash {
		dashCompletions = append(dashCompletions, parseAction(pos))
	}
	carapace.Gen(cmd).DashCompletion(dashCompletions...)

	carapace.Gen(cmd).DashAnyCompletion(parseAction(c.Completion.DashAny))

	for _, subcmd := range c.Commands {
		cmd.AddCommand(subcmd.ToCobra())
	}

	return cmd
}

func parseFlag(flagSet *pflag.FlagSet, id, description string) error {
	r := regexp.MustCompile(`^(?P<shorthand>-[^-])?(, *)?(?P<longhand>--[^ =*]*)?(?P<modifier>[=*]*)$`)
	matches := findNamedMatches(r, id)

	longhand := strings.TrimPrefix(matches["longhand"], "--")
	shorthand := strings.TrimPrefix(matches["shorthand"], "-")
	slice := strings.Contains(matches["modifier"], "*")
	value := strings.Contains(matches["modifier"], "=")

	if longhand != "" && shorthand != "" {
		if value {
			if slice {
				flagSet.StringSliceP(longhand, shorthand, []string{}, description)
			} else {
				flagSet.StringP(longhand, shorthand, "", description)
			}
		} else {
			if slice {
				flagSet.CountP(longhand, shorthand, description)
			} else {
				flagSet.BoolP(longhand, shorthand, false, description)
			}
		}
	} else if longhand != "" {
		if value {
			if slice {
				flagSet.StringSlice(longhand, []string{}, description)
			} else {
				flagSet.String(longhand, "", description)
			}
		} else {
			if slice {
				flagSet.Count(longhand, description)
			} else {
				flagSet.Bool(longhand, false, description)
			}
		}
	} else if shorthand != "" {
		if value {
			if slice {
				flagSet.StringSliceS(shorthand, shorthand, []string{}, description)
			} else {
				flagSet.StringS(shorthand, shorthand, "", description)
			}
		} else {
			if slice {
				flagSet.CountS(shorthand, shorthand, description)
			} else {
				flagSet.BoolS(shorthand, shorthand, false, description)
			}
		}
	} else {
		return fmt.Errorf("malformed flag: %v", id)
	}
	return nil
}

func parseAction(arr []string) carapace.Action {
	r := regexp.MustCompile(`^\$\((?P<cmd>.*)\)$`)

	batch := carapace.Batch()

	vals := make([]string, 0)
	for _, elem := range arr {
		if strings.HasPrefix(elem, "$_files") {
			return carapace.ActionFiles() // TODO params
		} else if strings.HasPrefix(elem, "$_directories") {
			return carapace.ActionDirectories()
		} else if r.MatchString(elem) {
			batch = append(batch, carapace.ActionExecCommand("sh", "-c", r.FindStringSubmatch(elem)[1])(func(output []byte) carapace.Action {
				// TODO parse like below
				lines := strings.Split(string(output), "\n")
				vals := make([]string, 0)
				for _, line := range lines {
					if line != "" {
						vals = append(vals, parseValue(line)...)
					}
				}
				return carapace.ActionStyledValuesDescribed(vals...)
			}))
		} else {
			vals = append(vals, parseValue(elem)...)
		}
	}
	// TODO $(execute) actions
	batch = append(batch, carapace.ActionStyledValuesDescribed(vals...))
	return batch.ToA()
}

func parseValue(s string) []string {
	if splitted := strings.SplitN(s, "\t", 3); len(splitted) > 2 {
		return splitted
	} else if len(splitted) > 1 {
		return []string{splitted[0], splitted[1], style.Default}
	} else {
		return []string{splitted[0], "", style.Default}
	}
}

func findNamedMatches(regex *regexp.Regexp, str string) map[string]string {
	match := regex.FindStringSubmatch(str)

	results := map[string]string{}
	for i, name := range match {
		results[regex.SubexpNames()[i]] = name
	}
	return results
}
