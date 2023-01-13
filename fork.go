package spec

import (
	"reflect"

	"github.com/spf13/pflag"
)

type flagSet struct {
	*pflag.FlagSet
}

func (f flagSet) IsFork() bool {
	return reflect.ValueOf(f.FlagSet).MethodByName("BoolN").IsValid()
}

func (f *flagSet) call(name string, args ...reflect.Value) {
	if v := reflect.ValueOf(f.FlagSet).MethodByName(name); v.IsValid() {
		v.Call(args)
	}
}

func (f *flagSet) BoolN(name, shorthand string, value bool, usage string) {
	f.call("BoolN", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}

func (f *flagSet) BoolS(name, shorthand string, value bool, usage string) {
	f.call("BoolS", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}

func (f *flagSet) CountN(name, shorthand, usage string) {
	f.call("CountN", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(usage))
}

func (f *flagSet) CountS(name, shorthand, usage string) {
	f.call("CountS", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(usage))
}

func (f *flagSet) StringSliceN(name, shorthand string, value []string, usage string) {
	f.call("StringSliceN", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}

func (f *flagSet) StringSliceS(name, shorthand string, value []string, usage string) {
	f.call("StringSliceS", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}

func (f *flagSet) StringN(name, shorthand, value, usage string) {
	f.call("StringN", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}

func (f *flagSet) StringS(name, shorthand, value, usage string) {
	f.call("StringS", reflect.ValueOf(name), reflect.ValueOf(shorthand), reflect.ValueOf(value), reflect.ValueOf(usage))
}
