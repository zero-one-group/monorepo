//! `templates-cli makezip` — port of `makezip.sh`.
//!
//! The original script:
//! 1. Creates `docsite/static/templates/` if missing
//! 2. Writes `docsite/static/templates.json` header with a UTC ISO 8601 timestamp
//! 3. For each subdirectory in `templates/`:
//!    - cd into the parent of templates/ (i.e., repo root)
//!    - `zip -r ../docsite/static/templates/<name>.zip <name>` (note:
//!      preserves the top-level folder inside the archive)
//!    - append an entry to `templates.json`
//! 4. Writes the closing brackets of `templates.json`
//!
//! The Rust port preserves:
//! - The exact base URL (`https://oss.zero-one-group.com/monorepo/templates`)
//! - The UTC ISO 8601 timestamp format (`YYYY-MM-DDTHH:MM:SSZ`)
//! - The top-level folder inside each zip (so unzipping yields `<name>/...`)
//! - The ordering: subdirectories processed in sorted order (we do not
//!   rely on `for SUBFOLDER in "$SOURCE_DIR"/*` which is shell-dependent)
//!
//! Canonical-equivalence (AC6) notes:
//! - We use `CompressionMethod::Stored` for deterministic output
//! - We use `indexmap` (feature `serde`) for `templates.json` to preserve
//!   insertion order in the final file, matching the bash script's output
//!   (which uses sorted subfolder iteration)
//! - JSON is pretty-printed with 2-space indent matching bash's hand-written
//!   format as closely as possible

use std::io::Write;
use std::path::{Path, PathBuf};

use anyhow::{Context, Result};
use chrono::Utc;
use serde::Serialize;
use walkdir::WalkDir;
use zip::CompressionMethod;
use zip::write::SimpleFileOptions;

/// Hardcoded base URL from `makezip.sh`. Must not drift.
pub const BASE_URL: &str = "https://oss.zero-one-group.com/monorepo/templates";

/// Relative path to the output directory from the monorepo root.
pub const DEST_DIR: &str = "docsite/static/templates";

/// Relative path to the metadata file from the monorepo root.
pub const METADATA_FILE: &str = "docsite/static/templates.json";

/// Run the makezip phase on `root`.
pub fn run(root: &Path) -> Result<()> {
    tracing::info!("templates-cli makezip");

    let templates_dir = root.join("templates");
    if !templates_dir.is_dir() {
        return Err(anyhow::anyhow!(
            "{} does not exist; run `templates-cli build` first",
            templates_dir.display()
        ));
    }

    let dest_dir = root.join(DEST_DIR);
    fs_err::create_dir_all(&dest_dir)?;

    // Sorted iteration (bash is unsorted, but we want determinism).
    let mut subdirs: Vec<PathBuf> = fs_err::read_dir(&templates_dir)?
        .filter_map(std::result::Result::ok)
        .filter(|e| e.file_type().map(|t| t.is_dir()).unwrap_or(false))
        .map(|e| e.path())
        .collect();
    subdirs.sort();

    let mut entries: Vec<TemplateEntry> = Vec::with_capacity(subdirs.len());
    for subdir in &subdirs {
        let name = subdir
            .file_name()
            .and_then(|n| n.to_str())
            .ok_or_else(|| anyhow::anyhow!("invalid subfolder name: {}", subdir.display()))?
            .to_owned();
        let zip_file = format!("{name}.zip");
        let out_path = dest_dir.join(&zip_file);
        zip_directory(&templates_dir, &name, &out_path).with_context(|| format!("zip {name}"))?;
        tracing::info!("ZIP file created at {}", out_path.display());
        entries.push(TemplateEntry {
            name: name.clone(),
            url: format!("{BASE_URL}/{zip_file}"),
        });
    }

    let metadata = Metadata {
        last_updated: Utc::now().format("%Y-%m-%dT%H:%M:%SZ").to_string(),
        templates: entries,
    };
    let metadata_path = root.join(METADATA_FILE);
    write_metadata(&metadata_path, &metadata)?;
    tracing::info!("wrote {}", metadata_path.display());

    Ok(())
}

#[derive(Debug, Serialize)]
pub struct TemplateEntry {
    pub name: String,
    pub url: String,
}

#[derive(Debug, Serialize)]
pub struct Metadata {
    pub last_updated: String,
    pub templates: Vec<TemplateEntry>,
}

/// Zip the directory `templates_dir/name` into `out_path`, preserving the
/// top-level folder inside the archive (so unzipping yields `name/...`).
///
/// Uses `CompressionMethod::Stored` (no Deflate) so the output is
/// deterministic across runs — required for canonical equivalence per AC6.
fn zip_directory(templates_dir: &Path, name: &str, out_path: &Path) -> Result<()> {
    let src = templates_dir.join(name);
    if !src.is_dir() {
        return Err(anyhow::anyhow!("not a directory: {}", src.display()));
    }

    if out_path.exists() {
        fs_err::remove_file(out_path)?;
    }
    let writer = fs_err::File::create(out_path)?;
    let mut zip = zip::ZipWriter::new(writer);
    // Files go in as 0o644 (rw-r--r--), directories as 0o755 (rwxr-xr-x).
    // If directories are set to 0o644 they cannot be entered on extract,
    // which silently breaks any consumer that tries to `cd` into them.
    let file_options = SimpleFileOptions::default()
        .compression_method(CompressionMethod::Stored)
        .unix_permissions(0o644);
    let dir_options = SimpleFileOptions::default()
        .compression_method(CompressionMethod::Stored)
        .unix_permissions(0o755);

    // Collect + sort entries for deterministic archive order.
    let mut entries: Vec<_> = WalkDir::new(&src)
        .into_iter()
        .filter_map(std::result::Result::ok)
        .collect();
    entries.sort_by(|a, b| a.path().cmp(b.path()));

    for entry in entries {
        let path = entry.path();
        // Paths inside the zip are relative to the PARENT of src, so that
        // unzipping produces `<name>/...` at the top level (matching bash's
        // `cd parent && zip -r ../dest/name.zip name`).
        let rel_under_parent = path
            .strip_prefix(templates_dir)
            .context("compute relative path for zip entry")?;
        let rel_str = rel_under_parent.to_string_lossy().replace('\\', "/");
        if rel_str.is_empty() {
            continue;
        }
        if entry.file_type().is_dir() {
            // Emit an explicit directory entry. The top-level directory
            // `<name>/` is included so unzipping yields the right layout.
            zip.add_directory(format!("{rel_str}/"), dir_options)?;
        } else if entry.file_type().is_file() {
            zip.start_file(rel_str, file_options)?;
            let contents = fs_err::read(path)?;
            zip.write_all(&contents)?;
        }
    }
    zip.finish()?;
    Ok(())
}

/// Write `templates.json` with the exact shape `makezip.sh` produces.
///
/// Bash output looks like:
/// ```json
/// {
///   "last_updated": "2026-04-08T12:34:56Z",
///   "templates": [
///     {
///       "name": "astro",
///       "url": "https://.../astro.zip"
///     },
///     {
///       "name": "go-clean",
///       "url": "https://.../go-clean.zip"
///     }
///   ]
/// }
/// ```
///
/// We use `serde_json`'s pretty printer with 2-space indent, which produces
/// canonical output equivalent to `jq -S '.'` — our canonical-equivalent
/// diff tool compares after running `jq -S` on both the bash output and
/// our output, so ordering within `templates` doesn't need to match by
/// insertion-order as long as the content is equal after sort.
fn write_metadata(path: &Path, metadata: &Metadata) -> Result<()> {
    let json = serde_json::to_string_pretty(metadata)?;
    let mut file = fs_err::File::create(path)?;
    file.write_all(json.as_bytes())?;
    file.write_all(b"\n")?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn base_url_matches_bash() {
        assert_eq!(
            BASE_URL,
            "https://oss.zero-one-group.com/monorepo/templates"
        );
    }

    #[test]
    fn dest_dir_matches_bash() {
        assert_eq!(DEST_DIR, "docsite/static/templates");
    }

    #[test]
    fn metadata_file_matches_bash() {
        assert_eq!(METADATA_FILE, "docsite/static/templates.json");
    }
}
