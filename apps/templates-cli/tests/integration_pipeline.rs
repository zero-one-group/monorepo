//! End-to-end integration tests for the templates-cli pipeline.
//!
//! These tests would have caught the drifts the canonical-equivalence
//! diff caught during Phase A, without requiring a full bash pipeline
//! baseline. Specifically:
//!
//! - **Drift 5 regression guard**: `zip::FileOptions` for directory
//!   entries must use 0o755 perms (or `add_directory` with matching
//!   permissions), otherwise extracted directories can't be entered.
//!   The canonical-equivalence diff only caught this via an explicit
//!   `unzip -q -o ...` on a real archive; this test catches it with
//!   a synthetic fixture in under a second.
//!
//! - **Shape assertions on `templates.json`**: schema, URL base, and
//!   insertion order of entries.
//!
//! - **End-to-end canonical shape of a zip**: listing, per-file
//!   contents, and file/dir permissions.

use std::fs::Permissions;
use std::io::Read as _;
use std::os::unix::fs::PermissionsExt as _;
use std::path::{Path, PathBuf};

use anyhow::Result;
use pretty_assertions::assert_eq;
use tempfile::tempdir;
use templates_cli::commands;

/// Build a synthetic `templates/` subtree inside `root` that exercises:
/// - A single-file template (smallest possible archive)
/// - A nested-directory template (exercises the `add_directory` path)
/// - An empty-directory edge case (shouldn't cause failures)
///
/// Returns the absolute root path so tests can operate on it.
fn build_fixture(root: &Path) -> Result<()> {
    let templates = root.join("templates");
    fs_err::create_dir_all(&templates)?;

    // 1) smallest: one file at the top
    let alpha = templates.join("alpha");
    fs_err::create_dir_all(&alpha)?;
    fs_err::write(alpha.join("README.md"), "# alpha\n")?;

    // 2) nested: directories-within-directories, multiple files, a file
    //    with known non-ASCII content, and a dotfile
    let beta = templates.join("beta");
    fs_err::create_dir_all(beta.join("src/lib"))?;
    fs_err::create_dir_all(beta.join("tests"))?;
    fs_err::write(beta.join("README.md"), "beta template\n")?;
    fs_err::write(beta.join(".gitignore"), "target/\n")?;
    fs_err::write(
        beta.join("src/main.rs"),
        "fn main() { println!(\"β 界\"); }\n",
    )?;
    fs_err::write(beta.join("src/lib/util.rs"), "pub fn hi() {}\n")?;
    fs_err::write(
        beta.join("tests/smoke.rs"),
        "#[test] fn smoke() { assert!(true); }\n",
    )?;

    // 3) empty-subdirectory edge case: gamma has an empty sub-dir
    let gamma = templates.join("gamma");
    fs_err::create_dir_all(gamma.join("empty"))?;
    fs_err::write(gamma.join("only.txt"), "only file\n")?;

    Ok(())
}

/// Extract a zip archive into `dest` and return a list of
/// `(relative_path, file_or_dir, content_bytes_if_file, unix_mode)`
/// entries, sorted by path for deterministic assertions.
fn read_zip_entries(zip_path: &Path) -> Result<Vec<ZipEntry>> {
    let bytes = fs_err::read(zip_path)?;
    let reader = std::io::Cursor::new(bytes);
    let mut archive = zip::ZipArchive::new(reader)?;
    let mut out = Vec::with_capacity(archive.len());
    for i in 0..archive.len() {
        let mut entry = archive.by_index(i)?;
        let name = entry.name().to_owned();
        let is_dir = entry.is_dir();
        let mode = entry.unix_mode();
        let mut contents = Vec::new();
        if !is_dir {
            entry.read_to_end(&mut contents)?;
        }
        out.push(ZipEntry {
            name,
            is_dir,
            mode,
            contents,
        });
    }
    out.sort_by(|a, b| a.name.cmp(&b.name));
    Ok(out)
}

#[derive(Debug, Clone)]
struct ZipEntry {
    name: String,
    is_dir: bool,
    mode: Option<u32>,
    contents: Vec<u8>,
}

#[test]
fn makezip_emits_one_archive_per_template_subdir() -> Result<()> {
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let zip_dir = dir.path().join("docsite/static/templates");
    let mut zips: Vec<PathBuf> = fs_err::read_dir(&zip_dir)?
        .filter_map(std::result::Result::ok)
        .filter(|e| e.path().extension().and_then(|ext| ext.to_str()) == Some("zip"))
        .map(|e| e.path())
        .collect();
    zips.sort();

    let zip_names: Vec<&str> = zips
        .iter()
        .filter_map(|p| p.file_name().and_then(|n| n.to_str()))
        .collect();
    assert_eq!(zip_names, vec!["alpha.zip", "beta.zip", "gamma.zip"]);
    Ok(())
}

#[test]
fn makezip_archive_contents_match_source() -> Result<()> {
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let zip_path = dir.path().join("docsite/static/templates/beta.zip");
    assert!(zip_path.exists(), "beta.zip should exist");

    let entries = read_zip_entries(&zip_path)?;

    // Every entry must be prefixed with "beta/" (the top-level folder)
    for e in &entries {
        assert!(
            e.name.starts_with("beta/"),
            "entry {:?} is not prefixed with beta/",
            e.name
        );
    }

    // Pick out the file entries (skip directory entries) and compare
    // their contents against the source.
    let source_root = dir.path().join("templates/beta");
    let file_entries: Vec<&ZipEntry> = entries.iter().filter(|e| !e.is_dir).collect();
    assert!(
        file_entries.len() >= 5,
        "expected at least 5 files in beta.zip, got {}",
        file_entries.len()
    );
    for e in &file_entries {
        // Strip the top-level "beta/" prefix to locate the source file.
        let rel = e
            .name
            .strip_prefix("beta/")
            .expect("already checked prefix");
        let source_file = source_root.join(rel);
        let source_contents = fs_err::read(&source_file).expect("source file present");
        assert_eq!(
            e.contents, source_contents,
            "contents mismatch for {}",
            e.name
        );
    }
    Ok(())
}

/// Regression guard for **Drift 5** caught during Phase A canonical-diff:
/// directory entries in the zip must have execute-permission bits set so
/// that extracted trees are enterable. Before the fix, every directory
/// came out at 0o644, triggering EACCES on `cd` and `ls` — the first Rust
/// makezip run produced a flood of "Permission denied" errors.
///
/// This test would have caught that in under a second.
#[test]
fn makezip_directory_entries_have_x_bit_set_drift5_regression() -> Result<()> {
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let zip_path = dir.path().join("docsite/static/templates/beta.zip");
    let entries = read_zip_entries(&zip_path)?;

    // Count directory entries — there should be at least 4 (beta/, src/,
    // src/lib/, tests/). Any directory entry that lacks the x-bit is a
    // Drift-5 regression.
    let mut dir_count = 0usize;
    for e in &entries {
        if e.is_dir {
            dir_count += 1;
            let mode = e.mode.unwrap_or(0);
            // Check owner execute bit (0o100) is set.
            assert!(
                mode & 0o100 != 0,
                "directory entry {} has mode {:o}, missing owner execute bit",
                e.name,
                mode
            );
            // For defense-in-depth, also check group+other x bits so
            // `ls -R` from any user can descend.
            assert!(
                mode & 0o755 == 0o755,
                "directory entry {} has mode {:o}, expected at least 0o755",
                e.name,
                mode
            );
        }
    }
    assert!(
        dir_count >= 4,
        "expected at least 4 directory entries in beta.zip, got {dir_count}"
    );
    Ok(())
}

#[test]
fn makezip_file_entries_have_regular_file_mode() -> Result<()> {
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let zip_path = dir.path().join("docsite/static/templates/beta.zip");
    let entries = read_zip_entries(&zip_path)?;

    for e in &entries {
        if !e.is_dir {
            let mode = e.mode.unwrap_or(0);
            // Owner read bit must be set.
            assert!(
                mode & 0o400 != 0,
                "file entry {} has mode {:o}, missing owner read bit",
                e.name,
                mode
            );
            // Owner execute bit must NOT be set — these are plain
            // content files, not scripts.
            assert!(
                mode & 0o100 == 0,
                "file entry {} has mode {:o}, should not have execute bit",
                e.name,
                mode
            );
        }
    }
    Ok(())
}

// `ends_with(".zip")` is fine in tests: we generate the .zip filenames
// ourselves in `makezip::run`, they are always lowercase, and case-mixed
// filenames would be a bug the test wants to catch.
#[allow(clippy::case_sensitive_file_extension_comparisons)]
#[test]
fn makezip_writes_templates_json_with_expected_shape() -> Result<()> {
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let json_path = dir.path().join("docsite/static/templates.json");
    assert!(json_path.exists(), "templates.json should exist");

    let raw = fs_err::read_to_string(&json_path)?;
    let parsed: serde_json::Value = serde_json::from_str(&raw)?;

    // Must have last_updated and templates fields
    assert!(parsed.get("last_updated").is_some(), "missing last_updated");
    let templates = parsed
        .get("templates")
        .and_then(|t| t.as_array())
        .expect("templates must be an array");
    assert_eq!(templates.len(), 3, "expected 3 template entries");

    // Names should be sorted and match our fixture
    let names: Vec<&str> = templates
        .iter()
        .filter_map(|t| t.get("name").and_then(|n| n.as_str()))
        .collect();
    assert_eq!(names, vec!["alpha", "beta", "gamma"]);

    // Each URL must start with the hardcoded BASE_URL
    let expected_base = "https://oss.zero-one-group.com/monorepo/templates";
    for t in templates {
        let url = t
            .get("url")
            .and_then(|u| u.as_str())
            .expect("url string required");
        assert!(
            url.starts_with(expected_base),
            "url {url} does not start with {expected_base}"
        );
        assert!(url.ends_with(".zip"), "url {url} does not end with .zip");
    }

    // last_updated must be ISO 8601 UTC (ends with Z)
    let last_updated = parsed["last_updated"]
        .as_str()
        .expect("last_updated must be string");
    assert!(
        last_updated.ends_with('Z'),
        "last_updated {last_updated} must be UTC ISO 8601 (ends with Z)"
    );
    Ok(())
}

#[test]
fn makezip_unzipped_files_are_readable_after_cp_equivalent_extract() -> Result<()> {
    // This test simulates what `unzip -o <archive> -d <target>` would do
    // and then tries to READ every extracted file + LIST every extracted
    // directory. Before Drift 5 was fixed, this would fail with EACCES
    // on the directories.
    let dir = tempdir()?;
    build_fixture(dir.path())?;

    commands::makezip::run(dir.path())?;

    let zip_path = dir.path().join("docsite/static/templates/beta.zip");
    let extract_root = dir.path().join("_extract");
    fs_err::create_dir_all(&extract_root)?;

    // Extract using the zip crate's reader API with explicit permission
    // application (mimicking what /usr/bin/unzip does on Linux).
    let bytes = fs_err::read(&zip_path)?;
    let mut archive = zip::ZipArchive::new(std::io::Cursor::new(bytes))?;
    for i in 0..archive.len() {
        let mut entry = archive.by_index(i)?;
        let out_path = extract_root.join(entry.name());
        if entry.is_dir() {
            fs_err::create_dir_all(&out_path)?;
            if let Some(mode) = entry.unix_mode() {
                fs_err::set_permissions(&out_path, Permissions::from_mode(mode))?;
            }
        } else {
            if let Some(parent) = out_path.parent() {
                fs_err::create_dir_all(parent)?;
            }
            let mut file = fs_err::File::create(&out_path)?;
            std::io::copy(&mut entry, &mut file)?;
            if let Some(mode) = entry.unix_mode() {
                fs_err::set_permissions(&out_path, Permissions::from_mode(mode))?;
            }
        }
    }

    // Now walk the extracted tree. Every file must be readable; every
    // directory must be listable. If a directory was extracted with
    // 0o644 (Drift 5), the `read_dir` call here would fail with EACCES.
    for entry in walkdir::WalkDir::new(extract_root.join("beta"))
        .into_iter()
        .filter_map(std::result::Result::ok)
    {
        if entry.file_type().is_file() {
            let _ = fs_err::read(entry.path())?;
        } else if entry.file_type().is_dir() {
            // Just listing the directory is enough to prove the x-bit
            // is set for the current user.
            let _: Vec<_> = fs_err::read_dir(entry.path())?.collect();
        }
    }
    Ok(())
}
