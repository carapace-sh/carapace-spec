# Parsing

Sets flag parsing mode. One of:

- `interspersed` mixed flags and positional arguments
- `non-interspersed` flag parsing stopped after first positional argument
- `disabled` flag parsing disabled

```yaml
{{#include ../../../../example/command.yaml:command}}
{{#include ../../../../example/command.yaml:parsing}}
```

![](./parsing.cast)
