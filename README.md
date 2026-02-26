# xk6-subcommand-docs

A [k6 extension](https://grafana.com/docs/k6/latest/extensions/) that provides offline k6 documentation in the terminal. Built as an [xk6 subcommand](https://grafana.com/docs/k6/latest/extensions/create/subcommand-extensions/), it registers under `k6 x docs`.

Designed for AI agents and developers who need quick access to k6 docs without leaving the terminal or opening a browser.

## Build

```bash
xk6 build --with github.com/grafana/xk6-subcommand-docs
```

## Usage

```
k6 x docs                                  # Table of contents
k6 x docs http                             # JS API: k6/http module
k6 x docs http get                         # JS API: HTTP get function
k6 x docs browser page click              # JS API: browser page click
k6 x docs using-k6 scenarios              # Using k6: scenarios
k6 x docs jslib                            # JS API: jslib (no k6- prefix needed)
k6 x docs crypto                           # JS API: crypto
k6 x docs javascript-api/k6-http/get      # Full slug (always works)
k6 x docs --list                           # List top-level categories
k6 x docs http --list                      # List subtopics under http
k6 x docs search threshold                # Search docs
k6 x docs search "close context"          # Fuzzy: matches closecontext
k6 x docs search "http-debugging"         # Fuzzy: matches http debugging
k6 x docs best-practices                  # Best practices guide
k6 x docs --all                            # Dump all docs
```

### Non-k6-prefixed modules

Most JS API modules use the `k6-` prefix (`k6-http`, `k6-browser`, etc.), but some don't. Modules like `jslib`, `crypto`, `init-context`, and `error-codes` can be accessed directly by name without any prefix:

```
k6 x docs jslib
k6 x docs crypto
k6 x docs init-context
```

### Fuzzy search

Search ignores spaces and dashes, so you don't need to remember exact formatting:

- `k6 x docs search "close context"` matches `closecontext`
- `k6 x docs search "http-debugging"` matches `http debugging`
- `k6 x docs search "closecontext"` matches `close-context`

Search results are grouped by topic/module with hierarchical output rather than a flat list.

### Markdown renderer

You can configure a markdown renderer (e.g. [glow](https://github.com/charmbracelet/glow)) so docs output is automatically rendered in the terminal.

Create `~/.config/k6/docs.yaml`:

```yaml
renderer: glow -p 200
```

When configured and stdout is a TTY (interactive use), output is piped through the renderer. When piped (`k6 x docs http | grep something`), the renderer is bypassed and raw markdown is emitted. If the renderer command fails, raw output is printed as a fallback.

## How doc bundles work

Docs are **not** embedded in the binary. They are fetched at runtime and cached locally.

On first run, the extension:

1. Detects the k6 version from build info (e.g. `v1.5.0` maps to `v1.5.x`)
2. Downloads the matching `.tar.zst` bundle from this repository's GitHub releases
3. Extracts it to `~/.local/share/k6/docs/{version}/`

Subsequent runs read straight from cache. No network needed.

Why not embed? xk6 has no build-time version detection mechanism, and embedding would bloat every k6 binary regardless of whether the user wants docs.

## How CI auto-publishes bundles

Docs are published as GitHub releases through two triggers:

- **Primary:** `repository_dispatch` from k6's release workflow
- **Fallback:** Scheduled poll (disabled by default)

The pipeline clones [k6-docs](https://github.com/grafana/k6-docs), transforms Hugo markdown into a compact format, compresses with zstd, and publishes as an asset to the single `doc-bundles` GitHub release.

## Claude Code Skill

This extension includes a [Claude Code](https://docs.anthropic.com/en/docs/claude-code) skill for AI-assisted k6 documentation lookup.

### Installation

```bash
npx @anthropic-ai/claude-code skill install --url https://github.com/grafana/xk6-subcommand-docs
```

### Usage

Once installed, Claude Code can look up k6 documentation during conversations using `k6 x docs`.

## Development

```
make test                                               # Run tests
make lint                                               # Run linter
make build                                              # Build k6 with this extension
make prepare K6_VERSION=v1.5.x K6_DOCS_PATH=~/k6-docs  # Prepare docs bundle locally
```
