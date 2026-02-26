**RULES:**
1. **Update this file concisely** whenever features are added, removed, or changed.
2. **TDD**: Always use red/green/refactor. Tests must compile and fail on assertions before writing implementation.
3. **Plans**: Store plans in `.claude/plans/` with incrementing numbers (next: `5-<name>.md`).

---

`k6 x docs` — offline k6 documentation in the terminal. For humans and AI agents. Docs are not embedded in the binary. On first run, the extension detects the k6 version from build info, downloads a matching compressed doc bundle (`.tar.zst`) from GitHub releases, and caches it locally (`~/.local/share/k6/docs/{version}/`). Subsequent runs serve from cache with no network. A separate standalone prepare tool (`cmd/prepare/`) builds these bundles by cloning the k6-docs Hugo repository, transforming markdown into CLI-friendly format, building a searchable index (`sections.json`), and compressing everything. CI auto-publishes a bundle per k6 release.

## Browsing
- `k6 x docs` shows categories with children and truncated descriptions (80 char max). Each category has a usage hint footer.
- `k6 x docs http get` resolves args to a slug (case-insensitive), reads the cached markdown, and prints it. If the topic has children, a subtopics footer is appended with comma-separated child names (redundant parent prefix stripped, e.g. `cookiejar-clear` → `clear`) and a usage hint showing the full CLI path via `slugToArgs`.
- `k6 x docs best-practices` prints a curated guide (embedded in the prepare tool via `//go:embed`).
- `k6 x docs search <query>` fuzzy searches for the query (case-insensitive, ignores punctuation and spaces).

## Slug resolution
- `k6 x docs http get` → `javascript-api/k6-http/get`
- `k6 x docs javascript-api/k6-http/get` → `javascript-api/k6-http/get`
- `k6 x docs using-k6 scenarios` → `using-k6/scenarios`
- Parent-prefix fallback: `k6 x docs http cookiejar clear` → tries `.../cookiejar/clear` (miss) → `.../cookiejar/cookiejar-clear` (hit). Handled by `withParentFallback` in `resolve.go`.

### Rendering
- Optional configurable renderer (e.g. `glow`) for pretty terminal output in `~/.config/k6/docs.yaml`.
- Links to the current version's online docs are stripped: `[text](https://grafana.com/docs/k6/v1.6.1/foo)` → `text`.
- Stripped: Shared shortcodes (`{{< docs/shared >}}`), code tags (`{{< code >}}`), section tags (`{{< section >}}`), React/MDX component tags (`<Glossary>`), `<br/>`, internal doc links, image links, remaining markdown links, HTML comments, YAML frontmatter.
- Converted: Admonitions (`{{< admonition type="warning" >}}`) → `> **Warning:** ...` blockquotes.
- Placeholders replaced: `<K6_VERSION>` → actual version.
- Internal doc links to included categories are converted to plain text (URL stripped). Links to excluded categories keep the URL.
- List topic and descriptions (also truncated to 80 chars with `...`) are aligned in columns for better readability.
- Duplicate child names are deduplicated (e.g. `javascript-api/k6-http/get` and `using-k6/http-requests/get` both have `get` child, but only one is shown in the parent list).

### Documentation version handling
- Auto-detects k6 version from Go build info.
- Maps to wildcard: `v1.5.0` → `v1.5.x`, `v1.6.0-rc.1` → `v1.6.x`.
- Override via `--version` flag or `K6_DOCS_VERSION` env var.
- Cache dir override via `--cache-dir` flag or `K6_DOCS_CACHE_DIR` env var.

### Bundle preparation (standalone `cmd/prepare/`)
- Clones k6-docs if not present, checks out matching tag.
- Builds shared content map from `docs/sources/shared/`.
- Walks markdown files, parses YAML frontmatter (deduplicates duplicate keys by keeping first occurrence), derives slugs, filters to included categories.
- Handles slug collisions: prefers `_index.md` over leaf `.md` (it has children).
- Populates parent→child relationships.
- Outputs: `dist/sections.json`, `dist/markdown/**/*.md`, `dist/best_practices.md`.

### CI/CD
- **CI** — lint + test + build on push/PR to main.
- **Release bundle** — triggered by k6 release dispatch or manual. Clones k6-docs, runs prepare, compresses with `zstd --ultra -22`, publishes GitHub release tagged `docs-v{major}.{minor}.x`.
- **Release poll** — manual fallback (schedule disabled). Polls k6 releases, builds if bundle missing.

### Categories
`javascript-api`, `using-k6`, `using-k6-browser`, `examples`, `extensions`, `results-output`, `testing-guides`, `set-up`
