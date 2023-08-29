package spec

import (
	"fmt"
	"strings"

	"github.com/rsteube/carapace"
)

type modifier struct {
	carapace.Action
}

func (m modifier) Parse(s string) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		m.Action = updateEnv(m.Action) // TODO verify

		var err error
		if s, err = c.Envsubst(s); err != nil {
			return carapace.ActionMessage(err.Error())
		}

		modifiers := map[string]Macro{
			"$chdir":      MacroI(m.Action.Chdir),
			"$filter":     MacroV(m.Action.Filter),
			"$filterargs": MacroN(m.Action.FilterArgs),
			"$list":       MacroI(m.Action.List),
			"$multiparts": MacroV(m.Action.MultiParts),
			"$nospace":    MacroI(func(s string) carapace.Action { return m.Action.NoSpace([]rune(s)...) }),
			"$prefix":     MacroI(m.Action.Prefix),
			"$retain":     MacroV(m.Action.Retain),
			"$shift":      MacroI(m.Action.Shift),
			"$split":      MacroN(m.Action.Split),
			"$splitp":     MacroN(m.Action.SplitP),
			"$suffix":     MacroI(m.Action.Suffix),
			"$suppress":   MacroI(func(s string) carapace.Action { return m.Action.Suppress(s) }),
			"$style":      MacroI(m.Action.Style),
			"$tag":        MacroI(m.Action.Tag),
			"$uniquelist": MacroI(m.Action.UniqueList),
			"$usage":      MacroI(func(s string) carapace.Action { return m.Action.Usage(s) }),
		}

		if modifier, ok := modifiers[strings.SplitN(s, "(", 2)[0]]; ok {
			return modifier.parse(s)
		}
		return carapace.ActionMessage("unknown macro: %#v", s)
	})
}

func updateEnv(a carapace.Action) carapace.Action {
	return carapace.ActionCallback(func(c carapace.Context) carapace.Action {
		for index, arg := range c.Parts {
			c.Setenv(fmt.Sprintf("C_PART%v", index), arg)
		}
		c.Setenv("C_VALUE", c.Value)
		return a.Invoke(c).ToA()
	})
}
