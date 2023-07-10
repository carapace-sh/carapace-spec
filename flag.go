package spec

import (
	"fmt"
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
	nonposix  bool
	hidden    bool
}

func parseFlag(s, usage string) (*flag, error) {
	r := regexp.MustCompile(`^(?P<shorthand>-[^-][^ =*?!]*)?(, )?(?P<longhand>-[-]?[^ =*?!]*)?(?P<modifier>[=*?!]*)$`)
	if !r.MatchString(s) {
		return nil, fmt.Errorf("flag syntax invalid: %v", s)
	}

	matches := findNamedMatches(r, s)

	f := &flag{}
	f.longhand = strings.TrimLeft(matches["longhand"], "-")
	f.nonposix = matches["longhand"] != "" && !strings.HasPrefix(matches["longhand"], "--")
	f.shorthand = strings.TrimPrefix(matches["shorthand"], "-")
	f.usage = usage
	f.slice = strings.Contains(matches["modifier"], "*")
	f.optarg = strings.Contains(matches["modifier"], "?")
	f.value = f.optarg || strings.Contains(matches["modifier"], "=")
	f.hidden = strings.Contains(matches["modifier"], "!")

	if f.longhand == "" && f.shorthand == "" {
		return nil, fmt.Errorf("malformed flag: '%v'", s)
	}
	return f, nil
}

func (f flag) addTo(fset *pflag.FlagSet) error {
	fs := flagSet{fset}
	if len(f.shorthand) > 1 && !fs.IsFork() {
		return fmt.Errorf("long shorthand only supported with rsteube/carapace-pflag: %v", f.shorthand)
	}
	if f.longhand == "" && !fs.IsFork() {
		return fmt.Errorf("shorthand-only only supported with rsteube/carapace-pflag: %v", f.shorthand)
	}

	if f.longhand != "" && f.shorthand != "" {
		if f.value {
			if f.slice {
				if !f.nonposix {
					fs.StringSliceP(f.longhand, f.shorthand, []string{}, f.usage)
				} else {
					fs.StringSliceN(f.longhand, f.shorthand, []string{}, f.usage)
				}
			} else {
				if !f.nonposix {
					fs.StringP(f.longhand, f.shorthand, "", f.usage)
				} else {
					fs.StringN(f.longhand, f.shorthand, "", f.usage)
				}
			}
		} else {
			if f.slice {
				if !f.nonposix {
					fs.CountP(f.longhand, f.shorthand, f.usage)
				} else {
					fs.CountN(f.longhand, f.shorthand, f.usage)
				}
			} else {
				if !f.nonposix {
					fs.BoolP(f.longhand, f.shorthand, false, f.usage)
				} else {
					fs.BoolN(f.longhand, f.shorthand, false, f.usage)
				}
			}
		}
	} else if f.longhand != "" {
		if f.value {
			if f.slice {
				if !f.nonposix {
					fs.StringSlice(f.longhand, []string{}, f.usage)
				} else {
					fs.StringSliceS(f.longhand, f.longhand, []string{}, f.usage)
				}
			} else {
				if !f.nonposix {
					fs.String(f.longhand, "", f.usage)
				} else {
					fs.StringS(f.longhand, f.longhand, "", f.usage)
				}
			}
		} else {
			if f.slice {
				if !f.nonposix {
					fs.Count(f.longhand, f.usage)
				} else {
					fs.CountS(f.longhand, f.longhand, f.usage)
				}
			} else {
				if !f.nonposix {
					fs.Bool(f.longhand, false, f.usage)
				} else {
					fs.BoolS(f.longhand, f.longhand, false, f.usage)
				}
			}
		}
	} else if f.shorthand != "" {
		if f.value {
			if f.slice {
				fs.StringSliceS(f.shorthand, f.shorthand, []string{}, f.usage)
			} else {
				fs.StringS(f.shorthand, f.shorthand, "", f.usage)
			}
		} else {
			if f.slice {
				fs.CountS(f.shorthand, f.shorthand, f.usage)
			} else {
				fs.BoolS(f.shorthand, f.shorthand, false, f.usage)
			}
		}
	}

	if f.optarg {
		if f.longhand != "" {
			fs.Lookup(f.longhand).NoOptDefVal = " "
		} else {
			fs.Lookup(f.shorthand).NoOptDefVal = " "
		}
	}

	if f.hidden {
		fs.Lookup(f.longhand).Hidden = f.hidden
	}

	return nil
}
