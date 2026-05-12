package spec

import (
	"fmt"
	"reflect"
	"strconv"
	"unicode/utf8"

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
	if f.Delimiter != "" && !fs.IsFork() {
		return fmt.Errorf("delimiter only supported with carapace-sh/carapace-pflag: %v", f.Shorthand)
	}

	defaultBool, err := boolDefault(f)
	if err != nil {
		return err
	}

	if f.Longhand != "" && f.Shorthand != "" {
		if f.Value {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.StringSliceN(f.Longhand, f.Shorthand, stringSliceDefault(f), f.Description)
				} else {
					fs.StringSliceP(f.Longhand, f.Shorthand, stringSliceDefault(f), f.Description)
				}
			} else {
				if f.NameAsShorthand {
					fs.StringN(f.Longhand, f.Shorthand, f.Default, f.Description)
				} else {
					fs.StringP(f.Longhand, f.Shorthand, f.Default, f.Description)
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
					fs.BoolN(f.Longhand, f.Shorthand, defaultBool, f.Description)
				} else {
					fs.BoolP(f.Longhand, f.Shorthand, defaultBool, f.Description)
				}
			}
		}
	} else if f.Longhand != "" {
		if f.Value {
			if f.Repeatable {
				if f.NameAsShorthand {
					fs.StringSliceS(f.Longhand, f.Longhand, stringSliceDefault(f), f.Description)
				} else {
					fs.StringSlice(f.Longhand, stringSliceDefault(f), f.Description)
				}
			} else {
				if f.NameAsShorthand {
					fs.StringS(f.Longhand, f.Longhand, f.Default, f.Description)
				} else {
					fs.String(f.Longhand, f.Default, f.Description)
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
					fs.Bool(f.Longhand, defaultBool, f.Description)
				} else {
					fs.BoolS(f.Longhand, f.Longhand, defaultBool, f.Description)
				}
			}
		}
	} else if f.Shorthand != "" {
		if f.Value {
			if f.Repeatable {
				fs.StringSliceS(f.Shorthand, f.Shorthand, stringSliceDefault(f), f.Description)
			} else {
				fs.StringS(f.Shorthand, f.Shorthand, f.Default, f.Description)
			}
		} else {
			if f.Repeatable {
				fs.CountS(f.Shorthand, f.Shorthand, f.Description)
			} else {
				fs.BoolS(f.Shorthand, f.Shorthand, defaultBool, f.Description)
			}
		}
	}

	flag := fs.Lookup(f.Name())
	if flag == nil {
		return fmt.Errorf("failed to add flag: %v", f.Name())
	}
	if f.Default != "" {
		flag.DefValue = f.Default
	}

	if f.Optarg {
		flag.NoOptDefVal = f.OptDefault
		if flag.NoOptDefVal == "" {
			flag.NoOptDefVal = " "
		}
	}

	if f.Hidden {
		flag.Hidden = f.Hidden
	}

	if f.Deprecated != "" {
		flag.Deprecated = f.Deprecated
	}

	if f.ShorthandDeprecated != "" {
		flag.ShorthandDeprecated = f.ShorthandDeprecated
	}

	if f.Nargs != 0 {
		// TODO move to carapace (pflagfork)
		if field := reflect.ValueOf(flag).Elem().FieldByName("Nargs"); field.IsValid() && field.Kind() == reflect.Int {
			field.SetInt(int64(f.Nargs))
		}
	}

	if f.Delimiter != "" {
		delimiter, size := utf8.DecodeRuneInString(f.Delimiter)
		if delimiter == utf8.RuneError || size != len(f.Delimiter) {
			return fmt.Errorf("invalid delimiter: %v", f.Delimiter)
		}
		if field := reflect.ValueOf(flag).Elem().FieldByName("OptargDelimiter"); field.IsValid() && field.Kind() == reflect.Int32 {
			field.SetInt(int64(delimiter))
		}
	}

	return nil
}

func stringSliceDefault(f command.Flag) []string {
	if f.Default == "" {
		return []string{}
	}
	return []string{f.Default}
}

func boolDefault(f command.Flag) (bool, error) {
	if f.Value || f.Repeatable || f.Default == "" {
		return false, nil
	}

	value, err := strconv.ParseBool(f.Default)
	if err != nil {
		return false, fmt.Errorf("invalid bool default for %v: %w", f.Name(), err)
	}
	return value, nil
}
