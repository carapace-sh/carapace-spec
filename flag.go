package spec

import (
	"fmt"

	"github.com/carapace-sh/carapace-spec/pkg/command"
	"github.com/spf13/pflag"
)

func addFlagTo(f command.Flag, fset *pflag.FlagSet) error {
	fs := flagSet{fset}
	if len(f.Shorthand) > 1 && !fs.IsFork() {
		return fmt.Errorf("long shorthand only supported with carapace-sh/carapace-pflag: %v", f.Shorthand)
	}
	if f.Longhand == "" && !fs.IsFork() {
		return fmt.Errorf("shorthand-only only supported with carapace-sh/carapace-pflag: %v", f.Shorthand)
	}
	if f.Nargs != 0 && !fs.IsFork() {
		return fmt.Errorf("nargs only supported with carapace-sh/carapace-pflag: %v", f.Shorthand)
	}

	if f.Longhand != "" && f.Shorthand != "" {
		if f.Value {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.StringSliceN(f.Longhand, f.Shorthand, []string{}, f.Description)
				} else {
					fs.StringSliceP(f.Longhand, f.Shorthand, []string{}, f.Description)
				}
			} else {
				if f.NameAsShorthand {
					fs.StringN(f.Longhand, f.Shorthand, "", f.Description)
				} else {
					fs.StringP(f.Longhand, f.Shorthand, "", f.Description)
				}
			}
		} else {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.CountN(f.Longhand, f.Shorthand, f.Description)
				} else {
					fs.CountP(f.Longhand, f.Shorthand, f.Description)
				}
			} else {
				if f.NameAsShorthand {
					fs.BoolN(f.Longhand, f.Shorthand, false, f.Description)
				} else {
					fs.BoolP(f.Longhand, f.Shorthand, false, f.Description)
				}
			}
		}
	} else if f.Longhand != "" {
		if f.Value {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.StringSliceS(f.Longhand, f.Longhand, []string{}, f.Description)
				} else {
					fs.StringSlice(f.Longhand, []string{}, f.Description)
				}
			} else {
				if f.NameAsShorthand {
					fs.StringS(f.Longhand, f.Longhand, "", f.Description)
				} else {
					fs.String(f.Longhand, "", f.Description)
				}
			}
		} else {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.CountS(f.Longhand, f.Longhand, f.Description)
				} else {
					fs.Count(f.Longhand, f.Description)
				}
			} else {
				if !f.NameAsShorthand {
					fs.Bool(f.Longhand, false, f.Description)
				} else {
					fs.BoolS(f.Longhand, f.Longhand, false, f.Description)
				}
			}
		}
	} else if f.Shorthand != "" {
		if f.Value {
			if f.Repeatable {
				fs.StringSliceS(f.Shorthand, f.Shorthand, []string{}, f.Description)
			} else {
				fs.StringS(f.Shorthand, f.Shorthand, "", f.Description)
			}
		} else {
			if f.Repeatable {
				fs.CountS(f.Shorthand, f.Shorthand, f.Description)
			} else {
				fs.BoolS(f.Shorthand, f.Shorthand, false, f.Description)
			}
		}
	}

	if f.Optarg {
		if f.Longhand != "" {
			fs.Lookup(f.Longhand).NoOptDefVal = " "
		} else {
			fs.Lookup(f.Shorthand).NoOptDefVal = " "
		}
	}

	if f.Hidden {
		fs.Lookup(f.Longhand).Hidden = f.Hidden
	}

	if f.Nargs != 0 {
		// TODO nargs only exists in fork
		fs.Lookup(f.Longhand).Nargs = f.Nargs
	}

	return nil
}
