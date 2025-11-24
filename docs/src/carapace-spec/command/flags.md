# Flags

Flags of the command.

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

Additional formats when built with [carapace-pflag](https://github.com/carapace-sh/carapace-pflag).

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:nonposix}}
```

![](./nonposix.cast)

## Extended

There's also an extended notations for less common use cases.

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:extended}}
```

- `nargs` amount of arguments consumed
