---
title: Troubleshooting
slug: "troubleshooting"
---

## `moon generate` scaffolds an outdated template

**Symptom:** `moon generate <template>` produces files that don't match what's
currently in `templates/` on `main` (for example an older routing stack, or a
component fix you know has already landed upstream).

**Cause:** moon downloads remote template archives once and caches them. The
templates listed in `.moon/workspace.yml` are stored under:

```text
~/.moon/templates/archive/oss.zero-one-group.com/<template>/
```

Each directory contains an `.installed-at` stamp. As long as that directory
exists, moon reuses it and **does not re-download** the archive, even if the
remote zip has been republished since.

**Fix:** delete the cached archive before generating:

```bash
# one template
rm -rf ~/.moon/templates/archive/oss.zero-one-group.com/<template>

# or everything
rm -rf ~/.moon/templates/archive

moon generate <template>
```

To check whether your cache is stale without deleting it, compare the
`.installed-at` stamp against the `last_updated` field of the published
manifest:

```bash
curl -s https://oss.zero-one-group.com/monorepo/templates.json | head -3
cat ~/.moon/templates/archive/oss.zero-one-group.com/<template>/.installed-at
```

If `last_updated` is newer than your stamp, clear the cache and re-run
`moon generate`.
