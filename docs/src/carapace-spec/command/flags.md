# Flags

Flags of the command.

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:flags}}
```

![](./flags/flags.cast)

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

![](./flags/nonposix.cast)

## Extended

There's also an extended notations for less common use cases.

- `nargs` amount of arguments consumed
- `default` default value
- `optdefault` default value when an optional argument flag is used without a value
- `deprecated` deprecation message for the flag
- `shorthanddeprecated` deprecation message for the shorthand
- `delimiter` alternative delimiter for optional arguments

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:extended}}
```

![](./flags/extended.cast)
