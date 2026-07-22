# repo

A demo CLI that showcases [repometa](https://github.com/jedi-knights/repometa)
— a library for detecting the buildable components and monorepo layouts inside
an arbitrary source tree. `repo` scans a directory, renders what it found in
your chosen format, and shows how a downstream consumer would integrate the
library behind a clean hexagonal boundary.

## Status

v0. The CLI is intentionally small; it exists so `repometa`'s output has a
first-class human interface without dragging a CLI into the library itself
(which is an explicit non-goal upstream).

## Quickstart

```bash
go build -o repo ./cmd/repo

# Human-readable summary
./repo summary /path/to/some/repo

# Full report as text (default)
./repo scan /path/to/some/repo

# JSON for downstream tooling
./repo scan --format json /path/to/some/repo

# ASCII table for terminals or code review
./repo scan --format table /path/to/some/repo
```

## Commands

| Command       | Purpose                                           |
| ------------- | ------------------------------------------------- |
| `scan <path>` | Full component list in the selected format.       |
| `summary <path>` | Short ecosystem/workspace counts.              |
| `version`     | Print binary version.                             |

## Configuration

Global flags (bind to viper; `REPO_*` env vars and a config file also apply):

| Flag              | Env                     | Default | Purpose                                             |
| ----------------- | ----------------------- | ------- | --------------------------------------------------- |
| `--format`        | `REPO_FORMAT`           | `text`  | Output format: `text`, `json`, `table`.             |
| `--max-depth`     | `REPO_MAX_DEPTH`        | `0`     | Cap directory recursion depth. `0` = library default (20). |
| `--max-dirs`      | `REPO_MAX_DIRS`         | `0`     | Cap total directories visited. `0` = library default (50 000). |
| `--max-file-size` | `REPO_MAX_FILE_SIZE`    | `0`     | Cap bytes read per manifest file. `0` = library default (4 MiB). |
| `--config`        | —                       | —       | Path to a YAML/TOML/JSON config file loaded by viper. |

## Architecture

Hexagonal. Direction of imports points inward — `domain` never depends on
adapters or CLI.

```
cmd/repo/          # wiring only: build Deps, run cobra root
internal/
  domain/          # pure: Report/Component/Summary + Analyzer service
                   # + Scanner interface (defined at the consumer)
  ports/           # Renderer/RendererFactory (consumed by CLI)
  adapters/
    scanner/       # Adapter over github.com/jedi-knights/repometa
    renderer/      # Strategy + Factory for text/json/table
  cli/             # cobra + viper wiring; commands hold no logic
```

### Design patterns used

- **Adapter** — `adapters/scanner.Repometa` translates the repometa
  `Manifest` into the domain's `Report`. Swapping the upstream library
  means changing one file.
- **Strategy** — each renderer (`Text`, `JSON`, `Table`) satisfies
  `ports.Renderer`. The CLI holds no format-specific branches.
- **Factory** — `renderer.Factory` maps a format name to the right
  Strategy. Adding a new format is one map entry, no CLI change.
- **Dependency Injection** — `domain.NewAnalyzer(Scanner)` and
  `cli.NewRootCmd(Deps)` are constructor-injected; tests supply fakes.

## Development

```bash
make test    # go test ./...
make build   # binary at ./bin/repo
make lint    # go vet ./... + go build -o /dev/null ./...
make tidy    # go mod tidy
```

The `github.com/jedi-knights/repometa` dependency is currently a `replace`
directive to `../repometa` because the module is unpublished. Remove the
replace once repometa is tagged and published to the Go proxy.
