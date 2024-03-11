> [!IMPORTANT]
> In the process of moving to [github.com/carapace-sh](https://github.com/carapace-sh)

# carapace-spec

[![PkgGoDev](https://pkg.go.dev/badge/github.com/carapace-sh/carapace-spec/pkg/actions)](https://pkg.go.dev/github.com/carapace-sh/carapace-spec)
[![GoReportCard](https://goreportcard.com/badge/github.com/carapace-sh/carapace-spec)](https://goreportcard.com/report/github.com/carapace-sh/carapace-spec)
[![documentation](https://img.shields.io/badge/&zwnj;-documentation-blue?logo=gitbook)](https://carapace-sh.github.io/carapace-spec/)
[![Coverage Status](https://coveralls.io/repos/github/carapace-sh/carapace-spec/badge.svg?branch=master)](https://coveralls.io/github/carapace-sh/carapace-spec?branch=master)
[![Packaging status](https://repology.org/badge/tiny-repos/carapace-spec.svg)](https://repology.org/project/carapace-spec/versions)

Define simple completions using a spec file (based on [carapace](https://github.com/carapace-sh/carapace)).

The `carapace-spec` binary can be used to complete spec files, but [carapace-bin](https://github.com/rsteube/carapace-bin) is recommended as it supports a range of [custom macros](https://rsteube.github.io/carapace-bin/spec/macros.html).

```yaml
name: mycmd
description: my command
flags:
  --optarg?: optarg flag
  -r, --repeatable*: repeatable flag
  -v=: flag with value
persistentflags:
  --help: bool flag
completion:
  flag:
    optarg: ["one", "two\twith description", "three\twith style\tblue"]
    v: ["$files"]
commands:
- name: sub
  description: subcommand
  completion:
    positional:
      - ["$list(,)", "1", "2", "3"]
      - ["$directories"]
```

## Generators

- [carapace-spec-clap](https://github.com/carapace-sh/carapace-spec-clap) spec generation for clap-rs/clap
- [carapace-spec-kingpin](https://github.com/carapace-sh/carapace-spec-kingpin) spec generation for alecthomas/kingpin
- [carapace-spec-kong](https://github.com/carapace-sh/carapace-spec-kong) spec generation for alecthomas/kong
- [carapace-spec-man](https://github.com/carapace-sh/carapace-spec-man) spec generation for manpages
- [carapace-spec-urfavecli](https://github.com/carapace-sh/carapace-spec-urfavecli) spec generation for urfave/cli
