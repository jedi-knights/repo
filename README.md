<div align="center">

# repo

**A demo CLI over the [repometa](https://github.com/jedi-knights/repometa) library — scan any source tree and report its buildable components and monorepo layouts.**

[![CI](https://github.com/jedi-knights/repo/actions/workflows/ci.yml/badge.svg)](https://github.com/jedi-knights/repo/actions/workflows/ci.yml)
[![Release](https://github.com/jedi-knights/repo/actions/workflows/release.yml/badge.svg)](https://github.com/jedi-knights/repo/actions/workflows/release.yml)
[![GoReleaser](https://github.com/jedi-knights/repo/actions/workflows/goreleaser.yml/badge.svg)](https://github.com/jedi-knights/repo/actions/workflows/goreleaser.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Badge](https://github.com/jedi-knights/repo/actions/workflows/badge.yaml/badge.svg)](https://github.com/jedi-knights/repo/actions/workflows/badge.yaml)
[![Coverage](https://img.shields.io/badge/Coverage-74.8%25-yellow)](https://jedi-knights.github.io/repo/?v=4)

[Overview](#overview) · [Features](#features) · [Requirements](#requirements) · [Installation](#installation) · [Usage](#usage) · [Configuration](#configuration) · [Development](#development) · [Contributing](#contributing)

</div>

---

## Overview

`repo` is a small, opinionated command-line consumer of the [repometa](https://github.com/jedi-knights/repometa) library. repometa detects the buildable components (Go modules, Rust crates, Node packages, Python packages, CMake / Make projects, loose C/asm trees) and monorepo workspace layouts (Go workspaces, Cargo workspaces, npm/yarn/pnpm workspaces, Nx, Turborepo, uv workspaces) inside an arbitrary directory. It intentionally does not ship a CLI — this repo is that CLI.

```
$ repo scan --format table /path/to/some/repo

Root: /path/to/some/repo   Components: 2   Monorepo: false

┌────────────────┬─────────────────┬────────────┬────────────┐
│      KIND      │      ROOT       │ CONFIDENCE │ WORKSPACES │
├────────────────┼─────────────────┼────────────┼────────────┤
│ make-project   │ .               │ 1.00       │ -          │
│ python-package │ lambdas/scanner │ 1.00       │ -          │
└────────────────┴─────────────────┴────────────┴────────────┘
```

## Features

- **Three output formats** — `text` for humans, `json` for downstream tooling, `table` for terminals and code review comments
- **Hexagonal architecture** — pure domain layer, adapters at the edges, dependency direction points inward
- **Layered configuration** — CLI flags override env vars (`REPO_*`) override config file (`--config file.yaml`) override built-in defaults, via [viper](https://github.com/spf13/viper)
- **Bounded scans** — the scanner's directory / depth / file-size caps are exposed as flags so operators can trade breadth for time
- **Small surface** — three commands (`scan`, `summary`, `version`), a handful of flags, no plugins
- **First-party design-pattern demonstration** — Adapter (scanner), Strategy + Factory (renderers), Constructor DI (analyzer, root command)
- **Pinned to a repometa pseudo-version** — swap to a semver tag once repometa cuts one

## Requirements

| Requirement | Version | Purpose |
|:---|:---|:---|
| Go | 1.26+ | Build from source |
| Homebrew | current | `brew install` distribution (macOS / Linux) |
| Git | any | Development only |

At runtime the binary is standalone — no Go toolchain required on the target host.

## Installation

### Homebrew (macOS and Linux)

```bash
brew install jedi-knights/tap/repo
```

### Pre-built binaries

Download the latest release for your platform from the [Releases page](https://github.com/jedi-knights/repo/releases).

```bash
# Linux x86_64
curl -fsSL https://github.com/jedi-knights/repo/releases/latest/download/repo_$(uname -s)_$(uname -m).tar.gz \
  | tar -xz -C /usr/local/bin repo

# macOS (Apple Silicon)
curl -fsSL https://github.com/jedi-knights/repo/releases/latest/download/repo_Darwin_arm64.tar.gz \
  | tar -xz -C /usr/local/bin repo
```

### Go install

```bash
go install github.com/jedi-knights/repo/cmd/repo@latest
```

### From source

```bash
git clone https://github.com/jedi-knights/repo
cd repo
make build   # binary at ./bin/repo
```

## Usage

### Scan a repository

```bash
# Human-readable text (default)
repo scan /path/to/some/repo

# JSON for downstream tooling
repo scan --format json /path/to/some/repo

# ASCII table for terminals and PR comments
repo scan --format table /path/to/some/repo
```

### Short summary

```bash
repo summary /path/to/some/repo
```

Prints the root, component count, monorepo flag, and per-ecosystem counts. Useful when scanning a directory of unrelated repos.

### Version

```bash
repo version
```

### Command reference

```
repo [command]

Commands:
  scan <path>     Scan a repository and render the detected components
  summary <path>  Print a short summary of ecosystems and workspaces in a repo
  version         Print binary version
  help            Help about any command

Global flags:
      --config string       path to a viper config file (yaml/toml/json)
      --format string       output format: text, json, table (default "text")
      --max-depth int       override scanner max directory depth (0 = library default 20)
      --max-dirs int        override scanner max directory count (0 = library default 50000)
      --max-file-size int   override scanner max file size in bytes (0 = library default 4 MiB)
  -h, --help                help for repo
```

## Configuration

Settings are resolved in this order (highest wins):

```
CLI flags  >  environment variables  >  config file  >  built-in defaults
```

### Environment variables

| Variable                | Description                                                    |
|:------------------------|:---------------------------------------------------------------|
| `REPO_FORMAT`           | Output format: `text`, `json`, `table`                         |
| `REPO_MAX_DEPTH`        | Cap directory recursion depth; `0` uses library default (20)   |
| `REPO_MAX_DIRS`         | Cap total directories visited; `0` uses library default (50000) |
| `REPO_MAX_FILE_SIZE`    | Cap bytes read per manifest file; `0` uses library default     |

### Config file

Any file viper can read (YAML / TOML / JSON). Example `repo.yaml`:

```yaml
format: table
max-depth: 15
max-dirs: 20000
max-file-size: 2097152
```

Load it with:

```bash
repo --config repo.yaml scan .
```

## Development

### Setup

```bash
git clone https://github.com/jedi-knights/repo
cd repo
go mod download
make build
make test
```

### Make targets

| Target       | Purpose                                              |
|:-------------|:-----------------------------------------------------|
| `make build` | Build the binary at `./bin/repo` with version ldflags |
| `make test`  | Run `go test ./...`                                  |
| `make vet`   | Run `go vet ./...`                                   |
| `make lint`  | Vet + full-tree build                                |
| `make tidy`  | Run `go mod tidy`                                    |
| `make clean` | Remove `./bin/`                                      |

### Architecture

Hexagonal. Import direction points inward — `domain` never depends on adapters or CLI.

```
cmd/repo/
  main.go                   Entry point, dependency wiring

internal/
  domain/                   Pure types and business logic (no I/O, no external deps)
    models.go               Report, Component, Summary
    analyzer.go             Analyzer service + Scanner interface (defined at consumer)
  ports/                    Consumer-defined interfaces
    ports.go                Renderer, RendererFactory
  adapters/
    scanner/                Adapter over github.com/jedi-knights/repometa
    renderer/               Strategy + Factory for text / json / table
  cli/                      cobra + viper wiring; commands hold no logic
```

The dependency rule: everything points inward toward `internal/domain`. Domain imports nothing beyond the standard library.

### Design patterns

- **Adapter** — `adapters/scanner.Repometa` translates the repometa `Manifest` into the domain `Report`. Swapping the upstream library means changing one file.
- **Strategy** — each renderer (`Text`, `JSON`, `Table`) satisfies `ports.Renderer`. The CLI holds no format-specific branches.
- **Factory** — `renderer.Factory` maps a format name to the right Strategy. Adding a format is one map entry, no CLI change.
- **Constructor injection** — `domain.NewAnalyzer(Scanner)` and `cli.NewRootCmd(Deps)` receive dependencies via constructor; tests supply fakes.

### Code style

- Standard `gofmt` formatting
- `go vet` and `golangci-lint` must pass with no warnings
- Functions ≤ 40 lines, cyclomatic complexity ≤ 7
- No globals; all dependencies injected via constructors

### Release process

Releases are driven end-to-end by conventional commits on `main`:

1. Commits merge to `main` → **CI** runs lint + tests + coverage gate (≥ 80%).
2. On CI success → **Release** workflow runs [go-semantic-release](https://github.com/jedi-knights/go-semantic-release), which analyzes commits since the last tag, computes the next version, writes `CHANGELOG.md` + `VERSION`, tags, and pushes back to `main`.
3. The new tag triggers **GoReleaser**, which cross-compiles for Linux / macOS / Windows on amd64 / arm64, creates the GitHub Release, and publishes a Homebrew formula to [`jedi-knights/homebrew-tap`](https://github.com/jedi-knights/homebrew-tap).
4. Also on release, the **Badge** workflow regenerates the coverage report, publishes it to GitHub Pages, and updates the coverage badge in this README.

## Contributing

Contributions are welcome. Please open an issue before starting significant work so the approach can be discussed.

- Follow [Conventional Commits](https://www.conventionalcommits.org/) — the release workflow depends on them
- One PR = one `type(scope)` pair
- New behavior needs a test; coverage stays ≥ 80%

## License

[MIT](./LICENSE)

---

<div align="center">
Made by <a href="https://github.com/jedi-knights">Jedi Knights</a>
</div>
