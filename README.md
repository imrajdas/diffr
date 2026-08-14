# diffr

Compare two directories or files. See what changed, what was added, and what was removed.

By default diffr opens a local web UI. You can also print a unified diff, write a patch, or emit JSON for scripts and CI.

[![Build](https://github.com/imrajdas/diffr/actions/workflows/build-check.yml/badge.svg?branch=main)](https://github.com/imrajdas/diffr/actions/workflows/build-check.yml)
[![Release](https://img.shields.io/github/v/release/imrajdas/diffr)](https://github.com/imrajdas/diffr/releases)
[![License](https://img.shields.io/github/license/imrajdas/diffr)](LICENSE)

<img src="./static/images/demo.png" alt="Diffr web UI">

## Install

### Homebrew (macOS / Linux)

```bash
brew install imrajdas/tap/diffr
```

```bash
diffr version
```

Upgrade later with:

```bash
brew upgrade imrajdas/tap/diffr
```

The fully qualified name (`imrajdas/tap/diffr`) taps this repo and trusts only this cask. Equivalent two-step form:

```bash
brew tap imrajdas/tap
brew install diffr
```

### GitHub Releases

Grab the archive for your OS from [Releases](https://github.com/imrajdas/diffr/releases).

**Linux / macOS**

```bash
tar -zxvf diffr_*_*.tar.gz
chmod +x diffr
sudo mv diffr /usr/local/bin/diffr
diffr version
```

**Windows** — unzip the archive and run `diffr.exe version`.

### From source

Needs [Go](https://go.dev/) 1.26.6 or later.

```bash
go install github.com/imrajdas/diffr@latest
```

Or clone and build:

```bash
git clone https://github.com/imrajdas/diffr
cd diffr
go build -o diffr .
```

## Quick start

```bash
# Open a side-by-side view in the browser
diffr ./old ./new

# Compare two files
diffr ./old.txt ./new.txt

# Print a unified diff (like diff -r)
diffr ./old ./new --stdout

# JSON for CI
diffr ./old ./new --json

# Write a patch
diffr ./old ./new --patch changes.patch
```

Without `--stdout`, `--json`, or `--patch`, diffr serves `http://localhost:8675` and opens it in your browser. Use `--no-open` if you only want the server.

Try the bundled demo trees from this repo:

```bash
diffr testdata/demo/left testdata/demo/right
```

## What it compares

diffr walks **both** trees:

| Status | Meaning |
| --- | --- |
| changed | Present on both sides, content differs |
| added | Only in the right tree |
| removed | Only in the left tree |
| identical | Counted, not listed |

File types are handled separately:

- **Text** — unified diffs, side-by-side in the UI
- **Images** (PNG, JPEG, GIF) — pixel comparison with a left/right slider
- **PDFs** — page count and extracted text
- **Binaries** — size/hash comparison, never dumped as text

The web UI also has a summary bar, filters (name, extension, status), unified vs side-by-side, and dark mode.

## Ignore rules

By default diffr skips `.git`, `node_modules`, `__pycache__`, `.svn`, `.hg`, `.DS_Store`, and `Thumbs.db`.

Rules are applied in this order:

1. Built-in excludes (unless `--no-default-exclude`)
2. `.diffrignore` in the current directory and at the root of each tree (gitignore syntax)
3. `.gitignore` from the git root of each tree (unless `--no-gitignore`)
4. `--exclude` / `-e` patterns (repeatable)
5. `--ignore-file` for an extra ignore file

Example `.diffrignore`:

```
*.log
dist/
vendor/
```

## Flags

```bash
diffr [dir1/file1] [dir2/file2] [flags]
diffr version
```

| Flag | Description |
| --- | --- |
| `-a, --address string` | Bind address for the web server. Default: `http://localhost` |
| `-p, --port int` | Port. Default: `8675` |
| `--stdout` | Print a unified diff instead of starting the web UI |
| `--json` | Print a JSON report instead of starting the web UI |
| `--patch string` | Write a unified patch (`-` for stdout) instead of starting the web UI |
| `--no-open` | Do not open a browser |
| `-C, --context int` | Unified-diff context lines. Default: `3` |
| `-w, --ignore-whitespace` | Ignore files that differ only in whitespace |
| `-i, --ignore-case` | Ignore files that differ only in case |
| `-e, --exclude stringArray` | Gitignore-style glob to exclude (repeatable) |
| `--ignore-file string` | Extra `.diffrignore` file |
| `--no-default-exclude` | Skip built-in excludes |
| `--no-gitignore` | Do not apply `.gitignore` |

`--stdout` is colored when stdout is a TTY. Set `NO_COLOR=1` for plain text. `--patch` files are always uncolored.

```bash
# UI on all interfaces, no browser
diffr ./old ./new -a 0.0.0.0 -p 9000 --no-open

# Ignore whitespace-only changes and extra paths
diffr ./old ./new --stdout -w -e '*.log' -e vendor
```

## Exit codes

Used with `--stdout`, `--json`, or `--patch`:

| Code | Meaning |
| --- | --- |
| `0` | No differences |
| `1` | Differences found |
| `2` | Error |

## Contributing

Bugs and ideas: open an [issue](https://github.com/imrajdas/diffr/issues). Code: fork and send a pull request.

## License

[Apache License 2.0](LICENSE)
