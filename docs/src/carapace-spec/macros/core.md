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
