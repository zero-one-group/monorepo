# templates-cli

Rust replacement for the monorepo's bash template pipeline. Ports the
semantics of the following shell scripts into a single workspace-managed
Cargo binary:

- `build-templates.sh` → `templates-cli build`
- `makezip.sh` → `templates-cli makezip`
- `builder/*.sh` (11 per-template scripts) → `templates-cli builder <name>`

## Why this exists

The original bash pipeline is portable across macOS and Linux but carries
known hazards: macOS/Linux `sed` syntax divergence, unquoted globbing,
non-deterministic `find` + `zip` output, JSON built via `echo >>` with
fragile comma placement, silent `sed` failures, and non-idempotent file
renames. This Rust port preserves the public contract (what files land
where, what placeholders get replaced, what the templates.json schema
looks like) while making the internals safer, deterministic, and
testable.

## Canonical-equivalence contract

The Rust output must
be **canonical-equivalent** to the bash output, defined as:

1. Sorted file listing of `templates/` and `docsite/static/templates/`
   is identical (via `find ... -print | LC_ALL=C sort`)
2. Per-file SHA-256 of every file is identical
3. `templates.json` matches after `jq -S '.'` canonicalization

Byte-for-byte zip equality is NOT required because bash's `zip` uses
non-deterministic Deflate compression. We use `zip::CompressionMethod::Stored`
so our output is deterministic, and we match the bash output by unzipping
and SHA-comparing the contents.

## Commands

```bash
# Full pipeline (equivalent to ./build-templates.sh followed by ./makezip.sh)
cargo run -p templates-cli -- all

# Just the build phase (copy apps -> templates, clean, run all builders)
cargo run -p templates-cli -- build

# Just the zip phase (zip each templates/ subdirectory + generate templates.json)
cargo run -p templates-cli -- makezip

# Run a single builder by name
cargo run -p templates-cli -- builder go-modular
```

## Invariants preserved from the bash pipeline

- Default port number per template (astro=4321, react=3000, etc.)
- Placeholder replacement: `_CHANGE_ME_DESCRIPTION_` → `{{ package_description }}`,
  source-name → `{{ package_name | kebab_case }}`, default-port → `{{ port_number }}`
- `.raw.*` file extension for language-specific raw template files
  (`.astro` → `.raw.astro`, `.py` → `.raw.py`, `.tsx` → `.raw.tsx`, etc.)
- `template.yml` is NEVER modified by the placeholder replacement pass
- `jq` removal of `@repo/shared-ui` from `package.json` and `tsconfig.json`
  for the 5 React/Next/Strapi/TanStack templates
- Astro frontmatter prepend: `---\nforce: true\n---` added to `.astro`
  files before the `.raw.astro` rename
- `go-modular/docs/` cleanup (keep `embed.go` only)
- `go-modular/web/` cleanup (keep `embed.go` and `static/index.html` only)
- Hardcoded zip base URL: `https://oss.zero-one-group.com/monorepo/templates`
- UTC ISO 8601 timestamp in `templates.json` (`Z` suffix, no fractional seconds)

## Not a drop-in for the bash scripts

This crate is **additive** during Phase A. The bash scripts stay in place
until the canonical-equivalence diff is green. Only after acceptance does
Task 18 of the ralplan (deferred to post-Phase-D cross-cutting cleanup)
remove the shell scripts.
