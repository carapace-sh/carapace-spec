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

type flag struct {
	longhand  string
	shorthand string
	usage     string
	slice     bool
	optarg    bool
	value     bool
}

func parseFlag(s, usage string) (f flag) {
	r := regexp.MustCompile(`^(?P<shorthand>-[^-])?(, )?(?P<longhand>--[^ =*?]*)?(?P<modifier>[=*?]*)$`)
	matches := findNamedMatches(r, s)

	f.longhand = strings.TrimPrefix(matches["longhand"], "--")
	f.shorthand = strings.TrimPrefix(matches["shorthand"], "-")
	f.usage = usage
	f.slice = strings.Contains(matches["modifier"], "*")
	f.optarg = strings.Contains(matches["modifier"], "?")
	f.value = f.optarg || strings.Contains(matches["modifier"], "=")

	// TODO enable error check
	//if f.longhand == "" && f.shorthand == "" {
	//	err = fmt.Errorf("malformed flag: '%v'", s)
	//}
	return
}

func (f flag) addTo(flagSet *pflag.FlagSet) {
	if f.longhand != "" && f.shorthand != "" {
		if f.value {
			if f.slice {
				flagSet.StringSliceP(f.longhand, f.shorthand, []string{}, f.usage)
			} else {
				flagSet.StringP(f.longhand, f.shorthand, "", f.usage)
			}
		} else {
			if f.slice {
				flagSet.CountP(f.longhand, f.shorthand, f.usage)
			} else {
				flagSet.BoolP(f.longhand, f.shorthand, false, f.usage)
			}
		}
	} else if f.longhand != "" {
		if f.value {
			if f.slice {
				flagSet.StringSlice(f.longhand, []string{}, f.usage)
			} else {
				flagSet.String(f.longhand, "", f.usage)
			}
		} else {
			if f.slice {
				flagSet.Count(f.longhand, f.usage)
			} else {
				flagSet.Bool(f.longhand, false, f.usage)
			}
		}
	} else if f.shorthand != "" {
		if f.value {
			if f.slice {
				flagSet.StringSliceS(f.shorthand, f.shorthand, []string{}, f.usage)
			} else {
				flagSet.StringS(f.shorthand, f.shorthand, "", f.usage)
			}
		} else {
			if f.slice {
				flagSet.CountS(f.shorthand, f.shorthand, f.usage)
			} else {
				flagSet.BoolS(f.shorthand, f.shorthand, false, f.usage)
			}
		}
	}

	if f.optarg {
		if f.longhand != "" {
			flagSet.Lookup(f.longhand).NoOptDefVal = " "
		} else {
			flagSet.Lookup(f.shorthand).NoOptDefVal = " "
		}
	}
}

func parseAction(cmd *cobra.Command, arr []string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		rMacro := regexp.MustCompile(`^(?P<macro>\$[^(]*)(\((?P<arg>.*)\))?$`)
		listDelimiter := ""
		nospace := false

		batch := carapace.Batch()
		vals := make([]string, 0)
		for _, elem := range arr {
			if strings.HasPrefix(elem, "$") { // macro
				match := findNamedMatches(rMacro, elem) // TODO check if matches
				macro := match["macro"]
				arg := match["arg"]

				switch macro {
				case "$nospace":
					nospace = true
				case "$list":
					listDelimiter = arg
				case "$directories":
					return carapace.ActionDirectories()
				case "$files":
					if arg != "" {
						batch = append(batch, carapace.ActionFiles(strings.Fields(arg)...))
					}
					batch = append(batch, carapace.ActionFiles())
				case "$":
					batch = append(batch, carapace.ActionCallback(func(c carapace.Context) carapace.Action {
						for index, arg := range c.Args {
							c.Setenv(fmt.Sprintf("CARAPACE_ARG%v", index), arg)
						}
						for index, arg := range c.Parts {
							c.Setenv(fmt.Sprintf("CARAPACE_PART%v", index), arg)
						}
						c.Setenv("CARAPACE_CALLBACK", c.CallbackValue)

						cmd.Flags().Visit(func(f *pflag.Flag) {
							c.Setenv(fmt.Sprintf("CARAPACE_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
						})

						return carapace.ActionExecCommand("sh", "-c", arg)(func(output []byte) carapace.Action {
							lines := strings.Split(string(output), "\n")
							vals := make([]string, 0)
							for _, line := range lines {
								if line != "" {
									vals = append(vals, parseValue(line)...)
								}
							}
							return carapace.ActionStyledValuesDescribed(vals...)
						}).Invoke(c).ToA()
					}))
				default:
					batch = append(batch, carapace.ActionMessage(fmt.Sprintf("malformed macro: '%v'", elem)))
				}
			} else {
				vals = append(vals, parseValue(elem)...)
			}
		}
		batch = append(batch, carapace.ActionStyledValuesDescribed(vals...))

		if listDelimiter != "" {
			return carapace.ActionMultiParts(listDelimiter, func(c carapace.Context) carapace.Action {
				return batch.ToA().Invoke(c).Filter(c.Parts).ToA()
			})
		} else if nospace {
			return batch.ToA().NoSpace()
		}
		return batch.ToA()
	})
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
