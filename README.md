# carapace-spec

[![PkgGoDev](https://pkg.go.dev/badge/github.com/rsteube/carapace-spec/pkg/actions)](https://pkg.go.dev/github.com/rsteube/carapace-spec)
[![GoReportCard](https://goreportcard.com/badge/github.com/rsteube/carapace-spec)](https://goreportcard.com/report/github.com/rsteube/carapace-spec)
[![documentation](https://img.shields.io/badge/&zwnj;-documentation-blue?logo=gitbook)](https://rsteube.github.io/carapace-spec/)
[![Coverage Status](https://coveralls.io/repos/github/rsteube/carapace-spec/badge.svg?branch=master)](https://coveralls.io/github/rsteube/carapace-spec?branch=master)

Define simple completions using a spec file (based on [carapace](https://github.com/rsteube/carapace)).

The `carapace-spec` binary can be used to complete spec files, but [carapace-bin](https://github.com/rsteube/carapace-bin) is recommended as it supports a range of [custom macros](https://rsteube.github.io/carapace-bin/specs/macros.html).

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
