package spec

import (
	"regexp"
	"strings"

	"github.com/spf13/pflag"
)

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
