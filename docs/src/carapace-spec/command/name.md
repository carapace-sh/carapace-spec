# Name

Name of the command.

```yaml
name: add
```

It can also contain the one-line [usage](https://rsteube.github.io/carapace/carapace/action/usage.html) message.

```yaml
name: add [-F file | -D dir]... [-f format] profile
```

> Recommended syntax is as [follows](https://pkg.go.dev/github.com/spf13/cobra#Command):
> - `[ ]` identifies an optional argument. Arguments that are not enclosed in brackets are required.
> - `...` indicates that you can specify multiple values for the previous argument.
> - `|`   indicates mutually exclusive information. You can use the argument to the left of the separator or the argument to the right of the separator. You cannot use both arguments in a single use of the command.
> - `{ }` delimits a set of mutually exclusive arguments when one of the arguments is required. If the arguments are optional, they are enclosed in brackets ([ ]).
