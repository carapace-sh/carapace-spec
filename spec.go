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

	flagCompletions := make(carapace.ActionMap)

	for id, description := range c.PersistentFlags {
		fs := parseFlagSpec(id, description)
		fs.addTo(cmd.PersistentFlags())
		if a, ok := c.Completion.Flag[fs.name()]; ok {
			var action carapace.Action
			if fs.delimiter != 0 {
				action = carapace.ActionMultiParts(string(fs.delimiter), func(c carapace.Context) carapace.Action {
					return parseAction(cmd, a).Invoke(c).Filter(c.Parts).ToA()
				})
			} else {
				action = carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					return parseAction(cmd, a)
				})
			}

			if fs.nospace {
				action = action.NoSpace()
			}
			flagCompletions[fs.name()] = action
		}
	}

	for id, description := range c.Flags {
		fs := parseFlagSpec(id, description)
		fs.addTo(cmd.Flags())
		if a, ok := c.Completion.Flag[fs.name()]; ok {
			var action carapace.Action
			if fs.delimiter != 0 {
				action = carapace.ActionMultiParts(string(fs.delimiter), func(c carapace.Context) carapace.Action {
					return parseAction(cmd, a).Invoke(c).Filter(c.Parts).ToA()
				})
			} else {
				action = carapace.ActionCallback(func(c carapace.Context) carapace.Action {
					return parseAction(cmd, a)
				})
			}

			if fs.nospace {
				action = action.NoSpace()
			}
			flagCompletions[fs.name()] = action
		}
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

type flagSpec struct {
	longhand    string
	shorthand   string
	description string
	slice       bool
	optarg      bool
	value       bool
	nospace     bool
	delimiter   rune
}

func (fs flagSpec) name() string {
	if fs.longhand != "" {
		return fs.longhand
	}
	return fs.shorthand
}

func parseFlagSpec(spec, description string) flagSpec {
	r := regexp.MustCompile(`^(?P<shorthand>-[^-])?(, *)?(?P<longhand>--[^ =*?!\/:.,]*)?(?P<modifier>[=*?!\/:.,]*)(?P<delimiter>.*)$`)
	matches := findNamedMatches(r, spec)

	fs := flagSpec{
		longhand:    strings.TrimPrefix(matches["longhand"], "--"),
		shorthand:   strings.TrimPrefix(matches["shorthand"], "-"),
		description: description,
		slice:       strings.Contains(matches["modifier"], "*"),
		optarg:      strings.Contains(matches["modifier"], "?"),
		nospace:     strings.Contains(matches["modifier"], "!"),
	}
	fs.value = fs.optarg || strings.Contains(matches["modifier"], "=")

	for _, d := range "/:.," {
		if strings.ContainsRune(matches["modifier"], d) {
			fs.delimiter = d
			break
		}
	}
	return fs
}

func (fs flagSpec) addTo(flagSet *pflag.FlagSet) error {
	if fs.longhand != "" && fs.shorthand != "" {
		if fs.value {
			if fs.slice {
				flagSet.StringSliceP(fs.longhand, fs.shorthand, []string{}, fs.description)
			} else {
				flagSet.StringP(fs.longhand, fs.shorthand, "", fs.description)
			}
		} else {
			if fs.slice {
				flagSet.CountP(fs.longhand, fs.shorthand, fs.description)
			} else {
				flagSet.BoolP(fs.longhand, fs.shorthand, false, fs.description)
			}
		}
	} else if fs.longhand != "" {
		if fs.value {
			if fs.slice {
				flagSet.StringSlice(fs.longhand, []string{}, fs.description)
			} else {
				flagSet.String(fs.longhand, "", fs.description)
			}
		} else {
			if fs.slice {
				flagSet.Count(fs.longhand, fs.description)
			} else {
				flagSet.Bool(fs.longhand, false, fs.description)
			}
		}
	} else if fs.shorthand != "" {
		if fs.value {
			if fs.slice {
				flagSet.StringSliceS(fs.shorthand, fs.shorthand, []string{}, fs.description)
			} else {
				flagSet.StringS(fs.shorthand, fs.shorthand, "", fs.description)
			}
		} else {
			if fs.slice {
				flagSet.CountS(fs.shorthand, fs.shorthand, fs.description)
			} else {
				flagSet.BoolS(fs.shorthand, fs.shorthand, false, fs.description)
			}
		}
	} else {
		return fmt.Errorf("malformed flag") // TODO context info
	}

	if fs.optarg {
		if fs.longhand != "" {
			flagSet.Lookup(fs.longhand).NoOptDefVal = " "
		} else {
			flagSet.Lookup(fs.shorthand).NoOptDefVal = " "
		}
	}

	return nil
}

func parseAction(cmd *cobra.Command, arr []string) carapace.Action {
	r := regexp.MustCompile(`^\$\((?P<cmd>.*)\)$`)

	batch := carapace.Batch()

	vals := make([]string, 0)
	for _, elem := range arr {
		if strings.HasPrefix(elem, "$_files") {
			return carapace.ActionFiles() // TODO params
		} else if strings.HasPrefix(elem, "$_directories") {
			return carapace.ActionDirectories()
		} else if r.MatchString(elem) {
            elemCopy := elem
			batch = append(batch, carapace.ActionCallback(func(c carapace.Context) carapace.Action {
				for index, arg := range c.Args {
					c.Setenv(fmt.Sprintf("CARAPACE_ARG%v", index), arg)
				}
				c.Setenv("CARAPACE_CALLBACK", c.CallbackValue)

				cmd.Flags().Visit(func(f *pflag.Flag) {
					c.Setenv(fmt.Sprintf("CARAPACE_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
				})

				return carapace.ActionExecCommand("sh", "-c", r.FindStringSubmatch(elemCopy)[1])(func(output []byte) carapace.Action {
					// TODO parse like below
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
