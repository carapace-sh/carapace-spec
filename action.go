package spec

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

// static value or macro
type action string

func (action) JSONSchema() *jsonschema.Schema {
	sortedNames := make([]string, 0, len(macros))
	for name := range macros {
		sortedNames = append(sortedNames, name)
	}
	sort.Strings(sortedNames)

	examples := make([]interface{}, 0, len(macros))
	for _, name := range sortedNames {
		examples = append(examples, fmt.Sprintf("$%v(%v)", name, macros[name].Signature()))
	}
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Action",
		Description: "A static value or a macro",
		Examples:    examples,
	}
}

// ActionMacro completes given macro
func ActionMacro(s string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		r := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`)
		if !r.MatchString(s) {
			return carapace.ActionMessage("malformed macro: %#v", s)
		}

		matches := findNamedMatches(r, s)
		if m, ok := macros[matches["macro"]]; !ok {
			return carapace.ActionMessage("unknown macro: %#v", s)
		} else {
			return m.f(matches["arg"])
		}
	})
}

// ActionSpec completes a spec
func ActionSpec(path string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		abs, err := c.Abs(path)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		content, err := os.ReadFile(abs)
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}

		var cmd Command
		if err := yaml.Unmarshal(content, &cmd); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		cmdCobra, err := cmd.ToCobraE()
		if err != nil {
			return carapace.ActionMessage(err.Error())
		}
		return carapace.ActionExecute(cmdCobra)
	})
}

func parseAction(cmd *cobra.Command, arr []action) carapace.Action {
	if !cmd.DisableFlagParsing {
		for _, entry := range arr {
			if strings.HasPrefix(string(entry), "$") {
				macro := strings.SplitN(strings.TrimPrefix(string(entry), "$"), "(", 2)[0]
				if m, ok := macros[macro]; ok && m.disableFlagParsing {
					cmd.DisableFlagParsing = true // implicitly disable flag parsing
					break
				}
			}
		}
	}

	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		// TODO yuck - where to set thes best?
		for index, arg := range c.Args {
			c.Setenv(fmt.Sprintf("C_ARG%v", index), arg)
		}
		c.Setenv("C_CALLBACK", c.CallbackValue)

		cmd.Flags().Visit(func(f *pflag.Flag) {
			c.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
		})

		batch := carapace.Batch()
		action := carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return batch.ToA()
		})

		vals := make([]string, 0)
		for _, elem := range arr {
			if elemSubst, err := c.Envsubst(string(elem)); err != nil {
				batch = append(batch, carapace.ActionMessage("%v: %#v", err.Error(), elem))
			} else if strings.HasPrefix(elemSubst, "$") { // macro
				switch strings.SplitN(elemSubst, "(", 2)[0] {
				case "$chdir":
					action = MacroI(action.Chdir).parse(elemSubst)
				case "$list":
					action = MacroI(updateEnv(action).List).parse(elemSubst)
				case "$multiparts":
					action = MacroV(action.MultiParts).parse(elemSubst)
				case "$nospace":
					localAction := action
					action = MacroI(func(s string) carapace.Action {
						return localAction.NoSpace([]rune(s)...)
					}).parse(elemSubst)
				case "$uniquelist":
					action = MacroI(updateEnv(action).UniqueList).parse(elemSubst)
				default:
					batch = append(batch, ActionMacro(elemSubst))
				}
			} else {
				vals = append(vals, parseValue(elemSubst)...)
			}
		}
		batch = append(batch, carapace.ActionStyledValuesDescribed(vals...))
		return action.Invoke(c).ToA()
	})
}

func updateEnv(a carapace.Action) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		for index, arg := range c.Parts {
			c.Setenv(fmt.Sprintf("C_PART%v", index), arg)
		}
		c.Setenv("C_CALLBACK", c.CallbackValue)
		return a.Invoke(c).ToA()
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
