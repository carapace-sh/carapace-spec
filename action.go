package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/carapace-sh/carapace"
	"github.com/invopop/jsonschema"
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
		examples = append(examples, fmt.Sprintf("$%v(%v)", name, macros[name].macro.Signature()))
	}
	return &jsonschema.Schema{
		Type:        "string",
		Title:       "Action",
		Description: "A static value or a macro",
		Examples:    examples,
	}
}

func executable() string {
	s, err := os.Executable()
	if err != nil {
		panic(err.Error()) // TODO handle error, eval symlink, how to handle "go test"
	}

	return filepath.Base(s)
}

// ActionMacro completes given macro
func ActionMacro(s string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		r := regexp.MustCompile(`^\$(?P<macro>[^(]*)(\((?P<arg>.*)\))?$`)
		matches := r.FindStringSubmatch(s)
		if matches == nil {
			return carapace.ActionMessage("malformed macro: %#v", s)
		}
		if strings.HasPrefix(matches[1], "_") && !strings.HasPrefix(matches[1], "_.") {
			return carapace.ActionMessage(`"$_" deprecated: replace %#v with %#v`, "$"+matches[1], "$carapace."+strings.TrimPrefix(matches[1], "_"))
		}
		prefix := fmt.Sprintf("$%v.", executable())

		switch {
		case !strings.HasPrefix(matches[1], "_.") && strings.Contains(matches[1], ".") && !strings.HasPrefix(s, prefix):
			splitted := strings.SplitN(strings.TrimPrefix(s, "$"), ".", 2)
			args := []string{"_carapace", "macro"}
			args = append(args, splitted[1])
			args = append(args, c.Args...)
			args = append(args, c.Value)
			carapace.LOG.Printf("%#v", args)
			return carapace.ActionExecCommand(splitted[0], args...)(func(output []byte) carapace.Action {
				return carapace.ActionImport(output)
			})

		default:
			if strings.HasPrefix(s, prefix) {
				s = "$_." + strings.TrimPrefix(s, prefix)
			}
			m, err := macros.Lookup(s)
			if err != nil {
				return carapace.ActionMessage(err.Error())
			}
			return m.Parse(s)
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

// TODO experimentally public
func (a action) Parse(cmd *cobra.Command) carapace.Action {
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
