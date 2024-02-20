# Flag

Define [flag completion](https://rsteube.github.io/carapace/carapace/gen/flagCompletion.html).

```yaml
# yaml-language-server: $schema=https://carapace.sh/schemas/command.json
name: flag
flags:
  -e=: executables
  -f, --file=: file
completion:
  flag:
    e: ["$executables"]
    file: ["$files"]
```

![](./flag.cast)
