package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/rsteube/carapace"
	"github.com/rsteube/carapace/pkg/style"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

func ActionSpec(path string) {
	carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		abs, err := filepath.Abs(path)
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
		return carapace.ActionExecute(cmd.ToCobra())
	})
}

func parseAction(cmd *cobra.Command, arr []string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		// TODO yuck - where to set thes best?
		for index, arg := range c.Args {
			c.Setenv(fmt.Sprintf("C_ARG%v", index), arg)
		}
		c.Setenv("C_CALLBACK", c.CallbackValue)

		cmd.Flags().Visit(func(f *pflag.Flag) {
			c.Setenv(fmt.Sprintf("C_FLAG_%v", strings.ToUpper(f.Name)), f.Value.String())
		})

		listDelimiter := ""
		nospace := false
		chdir := ""
		multiparts := ""

		// TODO don't alter the map each time, solve this differently
		addCoreMacro("list", MacroI(func(s string) carapace.Action {
			listDelimiter = s
			return carapace.ActionValues()
		}))
		addCoreMacro("nospace", MacroI(func(s string) carapace.Action {
			nospace = true
			return carapace.ActionValues()
		}))
		addCoreMacro("chdir", MacroI(func(s string) carapace.Action {
			chdir = s
			return carapace.ActionValues()
		}))
		addCoreMacro("multiparts", MacroI(func(s string) carapace.Action {
			multiparts = s
			return carapace.ActionValues()
		}))

		addCoreMacro("files", MacroV(carapace.ActionFiles))
		addCoreMacro("directories", MacroN(carapace.ActionDirectories))
		addCoreMacro("message", MacroI(carapace.ActionMessage))
		addCoreMacro("", MacroI(func(s string) carapace.Action {
			return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
				return carapace.ActionExecCommand("sh", "-c", s)(func(output []byte) carapace.Action {
					lines := strings.Split(string(output), "\n")
					vals := make([]string, 0)
					for _, line := range lines {
						if line != "" {
							vals = append(vals, parseValue(line)...)
						}
					}
					return carapace.ActionStyledValuesDescribed(vals...)
				}).Invoke(c).ToA()
			})
		}))

		batch := carapace.Batch()
		vals := make([]string, 0)
		for _, elem := range arr {
			if elemSubst, err := c.Envsubst(elem); err != nil {
				batch = append(batch, carapace.ActionMessage(fmt.Sprintf("%v: %v", err.Error(), elem)))
			} else if strings.HasPrefix(elemSubst, "$") { // macro
				batch = append(batch, ActionMacro(elemSubst))
			} else {
				vals = append(vals, parseValue(elemSubst)...)
			}
		}
		batch = append(batch, carapace.ActionStyledValuesDescribed(vals...))

		action := batch.ToA()
		if chdir != "" {
			action = action.Chdir(chdir)
		}
		if multiparts != "" {
			actionCopy := action
			action = carapace.ActionCallback(func(c carapace.Context) carapace.Action {
				return actionCopy.Invoke(c).ToMultiPartsA(multiparts)
			})
		}

		if listDelimiter != "" {
			return carapace.ActionMultiParts(listDelimiter, func(c carapace.Context) carapace.Action {
				for index, arg := range c.Parts {
					c.Setenv(fmt.Sprintf("C_PART%v", index), arg)
				}
				c.Setenv("C_CALLBACK", c.CallbackValue)

				return action.Invoke(c).Filter(c.Parts).ToA()
			})
		} else if nospace {
			return action.NoSpace()
		}
		return action.Invoke(c).ToA()
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
