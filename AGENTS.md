# AGENTS.md

Guide for AI agents working in the `carapace-spec` repository.

## Project Overview

`carapace-spec` defines shell completions using YAML spec files, built on top of [carapace](https://github.com/carapace-sh/carapace). A spec file describes a command's flags, subcommands, and completion actions; the library converts that into a `cobra.Command` tree that carapace bridges to all supported shells (bash, zsh, fish, elvish, nushell, oil, powershell, tcsh, xonsh).

The `carapace-spec` binary loads a spec YAML and outputs shell-specific completion scripts. [carapace-bin](https://github.com/carapace-sh/carapace-bin) is the recommended runtime for end users (it ships many custom macros), but this repo provides the library, spec format, codegen, and a standalone binary.

## Essential Commands

```bash
go build -v ./...                              # build all packages
go test -v ./...                               # run all tests
go test -v -coverprofile=profile.cov ./...     # tests with coverage (CI uses this)
gofmt -d -s .                                  # check formatting (CI fails on any diff)
go vet ./...                                   # vet
staticcheck ./...                              # static analysis (CI runs this)
go generate ./...                              # regenerate schema.json from struct tags
mdbook build docs                              # build the documentation site
```

### Schema generation

`schema.json` is generated from Go struct tags via `cmd/schema/main.go`:

```bash
go generate ./...          # runs: go run -C cmd/schema . ../../schema.json
```

### Releasing

Releases are handled by GoReleaser (`.goreleaser.yml`), triggered by git tags. Builds target linux/windows/darwin (CGO disabled) plus termux/android (CGO enabled). Packages are published to Homebrew, Scoop, AUR, and nFPM (apk/deb/rpm). Agents should not trigger releases.

## Repository Structure

This is a **Go workspace** (`go.work`) with two modules:

| Module | Path | Purpose |
|--------|------|---------|
| Library | `.` (`github.com/carapace-sh/carapace-spec`) | Spec types, YAML→cobra conversion, macro engine, codegen |
| CLI | `./cmd` (`github.com/carapace-sh/carapace-spec/cmd`) | The `carapace-spec` binary, schema generator |

The `cmd` module replaces `github.com/carapace-sh/carapace-spec` with the local library via a `replace` directive in `cmd/go.mod`.

### Library package layout (root module, package `spec`)

- `command.go` — `Command` type alias + `ToCobra()`/`ToCobraE()`: converts a spec `Command` into a `*cobra.Command`, wiring flags, completions, subcommands, run handlers
- `flag.go` — `addFlagTo()`: maps spec `Flag` to the correct `pflag.FlagSet` method (Bool/String/Count, P/N/S variants)
- `fork.go` — `flagSet` wrapper that detects whether the pflag fork (`carapace-pflag`) is in use via reflection (`IsFork()` checks for `BoolN` method)
- `action.go` — `action`/`value` types, `ActionMacro()`, `ActionSpec()`, and `action.Parse()` which resolves macro strings and `|||` modifier chains into `carapace.Action`s
- `modifier.go` — `modifier.Parse()`: applies modifier macros (`$filter`, `$list`, `$chdir`, etc.) to an action
- `macro.go` — `Macro` struct, `MacroN`/`MacroI`/`MacroV` constructors, `AddMacro`/`AddMacroI`/`AddMacroV` for custom macros, `macroName()` auto-naming logic
- `macromap.go` — `MacroMap` type + `Format()` method that generates Go source for a macro map (used by carapace-bin)
- `core.go` — `init()` registers core macros (`$files`, `$directories`, `$executables`, `$spec`, shell macros like `$bash`, `$sh`, etc.)
- `register.go` — `Register()`: adds the `_carapace macro` subcommand to a cobra command for macro introspection/invocation
- `run.go` — `run` type with three run modes: `parseMacro()`, `parseScript()`, `parseAlias()`
- `scan.go` — `ScanMacros()`: scans Go packages via `go doc` to discover `Action*` functions for macro registration (experimental)
- `codegen.go` — `Codegen()`: generates Go source files (cobra commands) from a spec, writing to a temp directory
- `schema.go` — embeds `schema.json` via `//go:embed`, exposes `Schema()`
- `internal/pflagfork/flag.go` — introspects pflag `Flag` fields (Nargs, Mode, OptargDelimiter) via reflection; classifies flags as Default/ShorthandOnly/NameAsShorthand
- `internal/shebang/shebang.go` — parses `#!` shebang lines for script run mode
- `pkg/command/` — `Command`, `Flag`, `FlagSet`, `Run`, `Parsing` types with YAML (de)serialization
- `pkg/macro/macro.go` — generic `Macro[T]`, `MacroMap[T]`, `MacroN`/`MacroI`/`MacroV` with YAML arg parsing and signature generation

### CLI package layout (`cmd/` module)

- `cmd/carapace-spec/main.go` — entry point, version handling
- `cmd/carapace-spec/cmd/root.go` — root command; loads spec YAML, bridges completion output. Contains `loadSpec()` and `bridgeCompletion()` (rewrites `_carapace` callbacks to reference the spec file path)
- `cmd/carapace-spec/cmd/codegen.go` — `codegen` subcommand
- `cmd/carapace-spec/cmd/run.go` — `run` subcommand (executes a spec as a real command)
- `cmd/carapace-spec/cmd/selfupdate.go` — self-update via `carapace-selfupdate`
- `cmd/schema/main.go` — standalone generator for `schema.json`; patches the FlagSet definition to support dual-type (string or extended object)

### Examples and docs

- `example/*.yaml` — spec files used both as documentation examples (included via `{{#include}}` in mdBook) and as test fixtures (embedded with `//go:embed`)
- `docs/` — mdBook documentation site; `docs/src/` contains markdown, `docs/src/carapace-spec/command/` has asciinema casts

## Architecture and Data Flow

### Spec → Cobra conversion

```
YAML spec file
    │  yaml.Unmarshal
    ▼
spec.Command (type alias of command.Command)
    │  Command.ToCobraE()
    ▼
*cobra.Command (with carapace completion registration)
    │  carapace bridge
    ▼
Shell completion script (bash/zsh/fish/...)
```

`Command.ToCobraE()` (`command.go:31`) runs a sequence of setup functions in order: `addFlags` → `addPersistentFlags` → `markFlagsExclusive` → `addRun` → `addFlagCompletion` → `addPositionalCompletion` → `addPositionalAnyCompletion` → `addDashCompletion` → `addDashAnyCompletion` → `addSubcommands` → `addAliasCompletion`. Each can return an error that aborts conversion. `ToCobra()` wraps this and, on error, creates a fallback command that surfaces the error as a completion message.

### Macro resolution

Completion values in a spec are strings that may be:
1. **Static values** — plain strings, tab-separated into value/description/style (`parseValue()` in `action.go:185`)
2. **Macros** — prefixed with `$`, e.g. `$files([.go])`, `$directories`, `$(echo one)`
3. **Modifiers** — applied via ` ||| ` delimiter, e.g. `["$files ||| $chdir(/tmp)"]`

`action.Parse()` (`action.go:107`) handles envsubst, splits on ` ||| `, and dispatches each element. Modifiers in a batch context wrap the entire batch; otherwise they wrap the preceding action.

Macro names are resolved through the global `macros` map (`macromap.go`). Core macros are registered with the `_.` prefix in `core.go`'s `init()`. Custom macros are added with `AddMacro`/`AddMacroI`/`AddMacroV` which also use the `_.` prefix. External macros (from other binaries like carapace-bin) are invoked by executing `<binary> _carapace macro <name> <args>`.

### Auto-naming (`macroName`)

`AddMacroI`/`AddMacroV` infer the macro name from the function's runtime name:
- Strips everything up to and including the last `/actions/` path segment (e.g. `.../pkg/actions/tools/git.ActionRefs` → `tools/git.Refs`)
- Strips the `Action` prefix from the function name
- Replaces `/` with `.` for nested package paths

### Run modes

The `run` field in a spec supports three types, detected by prefix:
- `$` → **macro**: executes a shell macro (e.g. `$sh(echo hello)`)
- `#!` → **script**: writes to a temp file and executes via the shebang interpreter
- `[` → **alias**: YAML array, first element is the command, rest are args (with envsubst)

All three set environment variables (`C_ARG<n>`, `C_FLAG_<NAME>`, `C_VALUE`) before execution.

### Flag syntax

Flags are defined as YAML map keys with a compact syntax (`pkg/command/flag.go:33` `format()` / `parseFlag()`):

| Suffix | Meaning |
|--------|---------|
| `=` | Takes a value |
| `?` | Optional argument (implies `=`) |
| `*` | Repeatable (slice/count) |
| `!` | Required |
| `&` | Hidden |

Examples: `--verbose`, `-v, --verbose`, `--output=`, `--count*`, `--optarg?=`, `--required!`, `--hidden&`

Non-POSIX flags (single `-` prefix for longhand, e.g. `-longflag`) are supported only with the carapace-pflag fork. See the pflag fork section below.

### Environment variables available in macros

Set by `action.Parse()` and `run.context()`:
- `${C_ARG<n>}` — positional arguments (0-indexed)
- `${C_FLAG_<FLAGNAME>}` — flag values (uppercase name, only if changed)
- `${C_VALUE}` — the word currently being completed
- `${C_PART<n>}` — parts of the current word during multipart completion

## Testing

### Test framework

Tests use the carapace `sandbox` package (`github.com/carapace-sh/carapace/pkg/sandbox`) for completion testing. The pattern in `example_test.go`:

```go
func sandboxSpec(t *testing.T, spec string) (f func(func(s *sandbox.Sandbox))) {
    var command Command
    if err := yaml.Unmarshal([]byte(spec), &command); err != nil {
        panic(err.Error())
    }
    return sandbox.Command(t, func() *cobra.Command {
        return command.ToCobra()
    })
}
```

Spec YAML is embedded with `//go:embed example/<name>.yaml` and passed to `sandboxSpec`. Tests then call `s.Run(args...).Expect(action)`.

For runnable specs (testing actual command execution), use `runnableSpec` + `runnable.Run(args...).Expect(output)`.

### Test files

| File | Covers |
|------|--------|
| `spec_test.go` | Non-POSIX flags (guarded by `skipNonFork`) |
| `command_test.go` | Commands, flags, parsing modes, exclusive flags, nargs |
| `flag_test.go` | Hidden/required/repeatable/optarg flags |
| `core_test.go` | Shell macros (bash, cmd, fish, etc.) — skips shells not in PATH |
| `modifier_test.go` | All modifier macros |
| `run_test.go` | Alias/macro/script run modes — some tests skip if `carapace` binary not in PATH |
| `macro_test.go` | Macro signatures and Default() interface |
| `pkg/command/flag_test.go` | Flag string parsing |
| `pkg/command/flagset_test.go` | FlagSet YAML marshal/unmarshal round-trip |

### Test gotchas

- `TestMain` in `spec_test.go` unsets `LS_COLORS` to ensure deterministic output.
- `skipNonFork(t)` skips tests that require the carapace-pflag fork (non-POSIX flags, longhand shorthands, nargs).
- Some tests (`TestRunArray`, `TestRunString`) skip if the `carapace` binary is not in PATH (they test alias completion which bridges through carapace-bin).
- Shell macro tests (`TestCore`) skip individual shells not installed locally.

## The pflag Fork

The project supports both `spf13/pflag` and `carapace-sh/carapace-pflag` (a fork). The fork adds features not in upstream pflag:

- **Non-POSIX flags** — longhand with single `-` prefix (e.g. `-longflag`)
- **Longhand shorthands** — multi-character shorthands
- **Shorthand-only flags** — no longhand at all
- **Nargs** — flags consuming a variable number of arguments

Detection is via reflection: `flagSet.IsFork()` (`fork.go:13`) checks whether `BoolN` method exists on the FlagSet. The `fork.go` wrapper uses reflection to call fork-only methods (`BoolN`, `BoolS`, `CountN`, `CountS`, `StringN`, `StringS`, `StringSliceN`, `StringSliceS`) so the code compiles with either pflag variant.

The `cmd/go.mod` replaces `spf13/pflag` with `carapace-sh/carapace-pflag`. In CI, the `nonposix` job explicitly runs `go work edit -replace github.com/spf13/pflag=github.com/carapace-sh/carapace-pflag@v1.1.0` to test fork features, while the `build` job tests without the replace.

## CI

`.github/workflows/go.yml` defines four jobs:

| Job | Trigger | Purpose |
|-----|---------|---------|
| `nonposix` | all | Build/test with carapace-pflag replace (non-POSIX flag features) |
| `build` | all | Build/test with standard pflag, gofmt check, staticcheck, coverage |
| `release` | tags only | GoReleaser publishes binaries + packages |
| `doc` | all (push to master deploys) | mdBook build, push to gh-pages |

Both test jobs install `carapace-bin` (v1.5.7) since some tests require the `carapace` binary in PATH.

## Code Conventions

- **Go 1.24** (set in `go.mod` and `go.work`)
- Standard `gofmt -s` formatting is enforced (CI checks with `gofmt -d -s .`)
- No external linter beyond `staticcheck` and `go vet`
- Package `spec` is the root library package; `command` and `macro` are sub-packages under `pkg/`
- Internal packages (`internal/pflagfork`, `internal/shebang`) are not importable externally
- YAML struct tags use `yaml:"name,omitempty"` and `json:"name,omitempty"` with `jsonschema_description` for schema generation
- The `Command` type in the root package is a type alias: `type Command command.Command` — methods are defined on the alias, not the underlying type

## Key Gotchas

1. **`go generate` changes `schema.json`** — The committed schema is outdated. Running `go generate ./...` will produce a different (more complete) `schema.json` with `FlagSet` `oneOf` support. Don't commit this unless your change is specifically about the schema.

2. **` ||| ` modifier delimiter is strict** — Not trimmed; must be exactly ` ||| ` with surrounding spaces (see `action.go:133` and modifier docs).

3. **Macro prefix convention** — Internal macros use `_.` prefix in the `macros` map. `ActionMacro()` converts `$<executable>.<name>` to `$_.<name>` for local resolution. The `$_` prefix (without dot) is deprecated in favor of `$carapace.` (`action.go:50-52`).

4. **`flagSet` reflection wrapper** — `fork.go` uses reflection to call fork-only pflag methods so the code compiles with both pflag versions. Don't call `BoolN`/`CountN` etc. directly on a `*pflag.FlagSet`; use the `flagSet` wrapper.

5. **`AddCommand` skips `_carapace`** — In codegen (`codegen.go:176`), the internal `_carapace` command is explicitly skipped when recursing into subcommands.

6. **Persistent flag completion** — `VisitAll` is used instead of `Visit()` in `action.Parse()` and `run.context()` because `Visit()` skips changed persistent flags from parent commands (`action.go:115`, `run.go:125`).

7. **Example YAML files serve triple duty** — They are test fixtures (`//go:embed`), documentation examples (mdBook `{{#include}}`), and live examples. Changes to example YAMLs affect all three. The `# ANCHOR:` / `# ANCHOR_END:` comments in `example/command.yaml` are used by mdBook to include specific sections in docs.

8. **`cmdVarName` special case** — The codegen function `cmdVarName()` (`codegen.go:189`) has a special case for a command literally named `root` (e.g. jj's `root` subcommand) to avoid a naming collision with the root command variable.
