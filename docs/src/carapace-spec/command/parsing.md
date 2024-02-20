# Parsing

Sets flag parsing mode. One of:

- `interspersed` mixed flags and positional arguments
- `non-interspersed` flag parsing stopped after first positional argument
- `disabled` flag parsing disabled

```yaml
# yaml-language-server: $schema=https://carapace.sh/schemas/command.json
name: parsing
persistentflags:
  -h, --help: show help
commands:
  - name: disabled
    parsing: disabled
    completion:
      positionalany: [one, two, three]
  - name: interspersed
    parsing: interspersed
    completion:
      positionalany: [one, two, three]
  - name: non-interspersed
    parsing: non-interspersed
    completion:
      positionalany: [one, two, three]
```

![](./parsing.cast)
