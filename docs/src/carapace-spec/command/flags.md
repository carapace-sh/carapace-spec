# Flags

Flags are defined as a map of name and description.

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:flags}}
```

![](./flags.cast)

## Modifiers:

Flags can have `0..n` modifier suffixes.

- `=` flag takes an argument
- `*` flag is repeatable
- `?` flag takes an optional argument
- `&` flag is hidden
- `!` flag is required

## Non-posix

With [carapace-pflag](https://github.com/carapace-sh/carapace-pflag) non-posix flags possible as well:

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:nonposix}}
```

![](./nonposix.cast)
