# Name

Name of the command.

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:name}}
```

It can also contain the one-line [usage](https://carapace-sh.github.io/carapace/carapace/action/usage.html) message.

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:usage}}
```

> [Recommended syntax is as follows](https://pkg.go.dev/github.com/spf13/cobra#Command):
> - `[ ]` identifies an **optional** argument. Arguments that are **not enclosed** in brackets are **required**.
> - `...` indicates that you can specify **multiple** values for the previous argument.
> - `|`   indicates **mutually exclusive** information. You can use the argument to the left of the separator or the argument to the right of the separator. You cannot use both arguments in a single use of the command.
> - `{ }` delimits a set of **mutually exclusive** arguments when one of the arguments is **required**. If the arguments are **optional**, they are enclosed in brackets (`[ ]`).

![](./name.cast)
