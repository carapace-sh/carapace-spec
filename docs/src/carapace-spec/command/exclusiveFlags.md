# ExclusiveFlags

Mark flags [mutually exclusive](https://pkg.go.dev/github.com/spf13/cobra#Command.MarkFlagsMutuallyExclusive).

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:exclusiveflags}}
```

## TODO remove code and update asciicast

```yaml
# yaml-language-server: $schema=https://carapace.sh/schemas/command.json
name: exclusiveflags
flags:
  --add: add package
  --delete: delete package
exclusiveflags:
  - [add, delete]
```

![](./exclusiveflags.cast)
