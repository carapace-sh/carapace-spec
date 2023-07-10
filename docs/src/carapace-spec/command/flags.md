# Flags

Flags are defined as a map of name and description.

```yaml
flags:
  -b: bool flag
  -v=: shorthand with value
  --repeatable*: longhand repeatable
  -o, --optarg?: shorthand and longhand with optional argument
```

## Modifiers:
- `=` flag takes an argument
- `*` flag is repeatable
- `?` flag takes an optional argument
- `!` flag is hidden

## Non-posix

With [carapace-pflag](https://github.com/rsteube/carapace-pflag) non-posix flags possible as well:

```yaml
  -np: non-posix shorthand
  -np, -nonposix:  non-posix shorthand and longhand
  -np, --nonposix: non-posix shorthand mixed with posix longhand
```
