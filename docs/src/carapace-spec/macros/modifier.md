# Modifier

Modifiers change the completion of macros and values.
These can be set generic `["<macro>", "<value>", "<modifier>"]` specific `["<macro> ||| <modifier> ||| <modifier>"]`.

> The delimiter (` ||| `) is currently very strict and not trimmed.

## chdir

[`$chdir(<directory>)`](https://rsteube.github.io/carapace/carapace/action/chdir.html) changes the directory.

```yaml
["$files", "$chdir(/tmp)"]
```

The following macros can be passed as well instead of a static directory:

- [`$gitdir`] .git folder
- [`$gitworktree`] root of the working directory for a non-bare repository
- [`$parent(\[file, dir\])`] first parent directory containing any of the given names/directories
- [`$tempdir`] default directory to use for temporary files
- [`$usercachedir`] root directory to use for user-specific cached data
- [`$userconfigdir`] default root directory to use for user-specific configuration data
- [`$userhomedir`] current user's home directory
- [`$xdgcachehome`] cache directory (fallback to UserCacheDir)
- [`$xdgconfighome`] home directory (fallback to UserConfigDir)

```yaml
["$files", "$chdir($gitdir)"]
```

## filter

[`$filter([<value>])`](https://rsteube.github.io/carapace/carapace/action/filter.html) filters given values.

```yaml
["one", "two", "three", "$filter([two])"]
```

## filterargs

[`$filterargs`](https://rsteube.github.io/carapace/carapace/action/filterArgs.html) filters `Context.Args`.

```yaml
["$files", "$filterargs"]
```

## list

[`$list(<delimiter>)`](https://rsteube.github.io/carapace/carapace/action/list.html) creates a list with given divider.

```yaml
["one", "two", "three", "$list(,)"]
```

## multiparts

[`$multiparts([<delimiter>])`](https://rsteube.github.io/carapace/carapace/action/multiParts.html) completes values splitted by given delimiter(s) separately.

```yaml
["one/two/three", "$multiparts([/])"]
```

## nospace

[`$nospace(<characters>)`](https://rsteube.github.io/carapace/carapace/action/noSpace.html) disables space suffix for given character(s).

```yaml
["one", "two/", "three,", "$nospace(/,)"]
```

## prefix

[`$pefix(<prefix>)`](https://rsteube.github.io/carapace/carapace/action/prefix.html) adds a prefix to the inserted values.

```yaml
["$files", "$prefix(file://)"]
```

## retain

[`$retain([<value>])`](https://rsteube.github.io/carapace/carapace/action/retain.html) retains given values.

```yaml
["one", "two", "three", "$retain([two])"]
```

## shift

[`$shift(<n>)`](https://rsteube.github.io/carapace/carapace/action/shift.html) shifts positional arguments left n times.

```yaml
["one", "two", "three", "$filterargs", "$shift(1)"]
```

## split

[`$split`](https://rsteube.github.io/carapace/carapace/action/split.html) splits `Context.Value` lexicographically and replaces `Context.Args` with the tokens.

```yaml
["one", "two", "three", "$filterargs", "$split"]
```

## splitp

[`$splitp`](https://rsteube.github.io/carapace/carapace/action/splitP.html) is like Split but supports pipelines.

```yaml
["one", "two", "three", "$filterargs", "$splitp"]
```

## suffix

[`$suffix(<suffix>)`](https://rsteube.github.io/carapace/carapace/action/suffix.html) adds a suffix to the inserted values.

```yaml
["apple", "melon", "orange", "$suffix(juice)"]
```

## suppress

[`$suppress(<regex>)`](https://rsteube.github.io/carapace/carapace/action/suppress.html) suppresses specific error messages using a regular expression.
```yaml
["$message(fail)", "$suppress(fail)"]
```

## style

[`$style(<style>)`](https://rsteube.github.io/carapace/carapace/action/style.html) sets the style for all values.

```yaml
["one", "two", "three", "$style(underlined)"]
```

## tag

[`$tag(<tag>)`](https://rsteube.github.io/carapace/carapace/action/tag.html) sets the tag for all values.

```yaml
["one", "two", "three", "$tag(numbers)"]
```

## uniquelist

[`$uniquelist(<delimiter>)`](https://rsteube.github.io/carapace/carapace/action/uniqueList.html) creates a unique list with given divider.

```yaml
["one", "two", "three", "$uniquelist(,)"]
```

## usage

[`$usage(<usage>)`](https://rsteube.github.io/carapace/carapace/action/usage.html) sets the usage message.

```yaml
["$usage(custom)"]
```


[`$gitdir`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#GitDir
[`$gitworktree`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#GitWorkTree
[`$parent(\[file, dir\])`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#Parent
[`$tempdir`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#TempDir
[`$usercachedir`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#UserCacheDir
[`$userconfigdir`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#UserConfigDir
[`$userhomedir`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#UserHomeDir
[`$xdgcachehome`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#XdgCacheHome
[`$xdgconfighome`]:https://pkg.go.dev/github.com/rsteube/carapace/pkg/traverse#XdgConfigHome