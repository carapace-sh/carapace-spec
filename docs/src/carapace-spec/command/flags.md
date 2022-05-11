# Flags

Flags are defined as a map of name and description.
The name can contain shorthand, longhand and modifiers matching this [regex](https://regex101.com/r/to7O2W/1).

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
