# Custom

Custom macros can be added with [`spec.AddMacro`](https://pkg.go.dev/github.com/rsteube/carapace-spec#AddMacro) (names are prefixed with `$_`).

```go
// `$_noarg` without argument
AddMacro("noarg", MacroN(func() carapace.Action { return carapace.ActionValues()}))

// `$_arg({user: example, enabled: true})` with argument (primitive or struct)
AddMacro("arg", MacroI(func(u User) carapace.Action { return carapace.ActionValues()}))

// `$_vararg([another, example])` with variable arguments (primitive or struct)
AddMacro("vararg", MacroV(func(s ...string) carapace.Action { return carapace.ActionValues()}))
```

Arguments are parsed as `yaml` so only struct keys deviating from the default need to be set.
