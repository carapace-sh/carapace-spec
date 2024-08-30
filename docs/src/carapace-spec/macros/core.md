# Core

Core macros provided by [carapace-spec](https://github.com/carapace-sh/carapace-spec).

## directories

[`$directories`](https://carapace-sh.github.io/carapace/carapace/defaultActions/actionDirectories.html) completes directories.
```yaml
["$directories"]
```

## exec

[`$(<command>)`](https://carapace-sh.github.io/carapace/carapace/defaultActions/actionExecCommand.html) executes given command in a `sh` / `pwsh` shell.

```yaml
["$(echo -e 'a\nb\nc')"]
```

Any arguments or options or flags already parsed by carapace will be included in the executed command's environment variables, prefixed with `C_`

For example, with a spec like 
```yaml
# yaml-language-server: $schema=https://carapace.sh/schemas/command.json
name: context
persistentflags:
  -p, --persistent: persistent flag
commands:
  - name: sub
    flags:
      -s, --string=: string flag
      -b, --bool: bool flag
      --custom=: custom flag
    completion:
      flag:
        custom: ["$(env)"]
```
Typing `context sub --custom ` and hitting the **TAB** key will execute the unix `env` command and return all environment variables as completion options.
Typing `context --persistent sub --string one -b arg1 arg2 --custom C_` and hitting **TAB** will produce the following terminal completion options:
```console
C_ARG0=arg1                                                                                                                              
C_ARG1=arg2                                                                                                                              
C_FLAG_BOOL=true                                                                                                                         
C_FLAG_STRING=one                                                                                                                        
C_VALUE=C_
```

Every variable listed in [Variables](https://carapace-sh.github.io/carapace-spec/carapace-spec/variables.html) which has a value will be included in the executed command's environment

## executables

[`$executables`](https://carapace-sh.github.io/carapace/carapace/defaultActions/actionExecutables.html) completes executables either from [PATH] or given directories.
```yaml
["$executables", "$executables([~/.local/bin])"]
```

## files

[`$files([<suffixes>])`](https://carapace-sh.github.io/carapace/carapace/defaultActions/actionFiles.html) completes files with an optional list of suffixes to filter on.

```yaml
["$files([.go, go.mod, go.sum])"]
```

## message

[`$message(<message>)`](https://carapace-sh.github.io/carapace/carapace/defaultActions/actionMessage.html) adds given error message to completion.

```yaml
["$message(some error)"]
```

## spec

`$spec(<file>)` completes arguments using the given spec file.
This implicitly [disables flag parsing](https://pkg.go.dev/github.com/spf13/cobra#Command) for the corresponding (sub)command.

```yaml
["$spec(example.yaml)"]
```

[PATH]:https://en.wikipedia.org/wiki/PATH_(variable)
