package spec

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/invopop/jsonschema"
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

type (
	// static or dynamic value (macro)
	value  string
	action []value
)

func NewAction(s []string) action { // TODO rename
	a := make(action, len(s))
	for index, v := range s {
		a[index] = value(v)
	}
	return a
}

func (value) JSONSchema() *jsonschema.Schema {
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

func (a action) disableFlagParsing() bool {
	for _, value := range a {
		if strings.HasPrefix(string(value), "$") {
			macro := strings.SplitN(strings.TrimPrefix(string(value), "$"), "(", 2)[0]
			if m, ok := macros[macro]; ok && m.disableFlagParsing {
				return true
			}
		}
	}
	return false
}

func (a action) parse(cmd *cobra.Command) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		// TODO yuck - where to set thes best?
		for index, arg := range c.Args {
			c.Setenv(fmt.Sprintf("C_ARG%v", index), arg)
		}
		c.Setenv("C_VALUE", c.Value)

		cmd.Flags().Visit(func(f *pflag.Flag) {
			c.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
		})

		batch := carapace.Batch()
		batchAction := carapace.ActionCallback(func(c carapace.Context) carapace.Action {
			return batch.ToA()
		})

		for _, elem := range a {
			elemSubst, err := c.Envsubst(string(elem))
			if err != nil {
				batch = append(batch, carapace.ActionMessage("%v: %#v", err.Error(), elem))
				continue
			}

			splitted := strings.Split(elemSubst, " ||| ")

			if strings.HasPrefix(splitted[0], "$") { // macro
				switch strings.SplitN(splitted[0], "(", 2)[0] {
				case // generic modifier applied to batch
					"$chdir",
					"$filter",
					"$filterargs",
					"$list",
					"$multiparts",
					"$nospace",
					"$prefix",
					"$retain",
					"$shift",
					"$split",
					"$splitp",
					"$suffix",
					"$suppress",
					"$style",
					"$tag",
					"$uniquelist",
					"$usage":
					batchAction = modifier{batchAction}.Parse(splitted[0])
					if len(splitted) > 1 {
						for _, m := range splitted[1:] {
							batchAction = modifier{batchAction}.Parse(m)
						}
					}
				default:
					a := ActionMacro(splitted[0])
					if len(splitted) > 1 {
						for _, m := range splitted[1:] {
							a = modifier{a}.Parse(m)
						}
					}
					batch = append(batch, a)
				}
			} else {
				a := parseValue(splitted[0])
				if len(splitted) > 1 {
					for _, m := range splitted[1:] {
						a = modifier{a}.Parse(m)
					}
				}
				batch = append(batch, a)
			}
		}
		return batchAction.Invoke(c).ToA()
	})
}

func parseValue(s string) carapace.Action {
	splitted := strings.SplitN(s, "\t", 3)
	switch len(splitted) {
	case 1:
		return carapace.ActionValues(splitted...)
	case 2:
		return carapace.ActionValuesDescribed(splitted...)
	case 3:
		return carapace.ActionStyledValuesDescribed(splitted...)
	default:
		return carapace.ActionValues("invalid value: %#v", s)

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
