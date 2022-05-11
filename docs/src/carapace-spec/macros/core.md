# Core

Core macros provided by [carapace-spec](https://github.com/rsteube/carapace-spec).

## directories

`$directories` completes directories.
```yaml
["$directories"]
```

## command

`$(<command>)` executes given command in a `sh` shell.

```yaml
["$(echo -e 'a\nb\nc')"]
```

## files

`$files([<suffixes>])` completes files with an optional list of suffixes to filter on.

```yaml
["$files([.go, go.mod, go.sum])"]
```

## message

`$message(<message>)` adds given error message to completion.

```yaml
["$message(some error)"]
```

## spec

`$spec(<file>)` completes arguments using the given spec file.
This implicitly [disables flag parsing](https://pkg.go.dev/github.com/spf13/cobra#Command) for the corresponding (sub)command.

```yaml
["$spec(example.yaml)"]
```
