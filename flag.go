package spec

import (
	"fmt"

	"github.com/spf13/pflag"
)

type flag struct {
	longhand        string
	shorthand       string
	usage           string
	slice           bool
	optarg          bool
	value           bool
	nameAsShorthand bool
	hidden          bool
	required        bool
	nargs           int
}

func (f flag) addTo(fset *pflag.FlagSet) error {
	fs := flagSet{fset}
	if len(f.shorthand) > 1 && !fs.IsFork() {
		return fmt.Errorf("long shorthand only supported with carapace-sh/carapace-pflag: %v", f.shorthand)
	}
	if f.longhand == "" && !fs.IsFork() {
		return fmt.Errorf("shorthand-only only supported with carapace-sh/carapace-pflag: %v", f.shorthand)
	}
	if f.nargs != 0 && !fs.IsFork() {
		return fmt.Errorf("nargs only supported with carapace-sh/carapace-pflag: %v", f.shorthand)
	}

	if f.longhand != "" && f.shorthand != "" {
		if f.value {
			if f.slice {
				if f.nameAsShorthand {
					fs.StringSliceN(f.longhand, f.shorthand, []string{}, f.usage)
				} else {
					fs.StringSliceP(f.longhand, f.shorthand, []string{}, f.usage)
				}
			} else {
				if f.nameAsShorthand {
					fs.StringN(f.longhand, f.shorthand, "", f.usage)
				} else {
					fs.StringP(f.longhand, f.shorthand, "", f.usage)
				}
			}
		} else {
			if f.slice {
				if f.nameAsShorthand {
					fs.CountN(f.longhand, f.shorthand, f.usage)
				} else {
					fs.CountP(f.longhand, f.shorthand, f.usage)
				}
			} else {
				if f.nameAsShorthand {
					fs.BoolN(f.longhand, f.shorthand, false, f.usage)
				} else {
					fs.BoolP(f.longhand, f.shorthand, false, f.usage)
				}
			}
		}
	} else if f.longhand != "" {
		if f.value {
			if f.slice {
				if f.nameAsShorthand {
					fs.StringSliceS(f.longhand, f.longhand, []string{}, f.usage)
				} else {
					fs.StringSlice(f.longhand, []string{}, f.usage)
				}
			} else {
				if f.nameAsShorthand {
					fs.StringS(f.longhand, f.longhand, "", f.usage)
				} else {
					fs.String(f.longhand, "", f.usage)
				}
			}
		} else {
			if f.slice {
				if f.nameAsShorthand {
					fs.CountS(f.longhand, f.longhand, f.usage)
				} else {
					fs.Count(f.longhand, f.usage)
				}
			} else {
				if !f.nameAsShorthand {
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

	if f.nargs != 0 {
		// TODO nargs only exists in fork
		fs.Lookup(f.longhand).Nargs = f.nargs
	}

	return nil
}
