//! Shared helpers used by every builder and command.
//!
//! These are the Rust analogues of the primitives the bash pipeline relies on:
//! `sed -i` + `grep -rl` + `find -name ... | mv` + `jq | mv`. They avoid
//! macOS/Linux shell divergence entirely because the regex crate and
//! `fs_err` behave the same on every platform.

use std::path::{Path, PathBuf};

use anyhow::{Context, Result, anyhow};
use walkdir::WalkDir;

/// File names that are NEVER touched by placeholder replacement.
///
/// `template.yml` is skipped by every `builder/*.sh` script via
/// `grep -v "template.yml"`. This list captures that exact exclusion.
pub const PLACEHOLDER_SKIP_FILES: &[&str] = &["template.yml"];

/// Walk `root` in sorted order, yielding every file.
///
/// The sort is by full path using `LC_ALL=C` semantics (byte-wise), so the
/// order is stable across runs and platforms — required for canonical
/// equivalence with the bash baseline.
#[allow(clippy::unnecessary_wraps)] // shape reserved for future propagation
pub fn walk_files_sorted(root: &Path) -> Result<Vec<PathBuf>> {
    let mut files: Vec<PathBuf> = WalkDir::new(root)
        .into_iter()
        .filter_map(std::result::Result::ok)
        .filter(|e| e.file_type().is_file())
        .map(|e| e.path().to_path_buf())
        .collect();
    files.sort();
    Ok(files)
}

/// Replace every occurrence of `needle` with `replacement` in every file
/// under `root` that contains `needle`, **skipping** any file whose basename
/// is in [`PLACEHOLDER_SKIP_FILES`] and any file that is not valid UTF-8.
///
/// This is the Rust analogue of:
/// ```bash
/// grep -rl "$needle" "$root" | grep -v "template.yml" | while read -r file; do
///     sed -i "s/$needle/$replacement/g" "$file"
/// done
/// ```
///
/// The replacement is **literal** (no regex metacharacters) — bash's `sed`
/// replacement string is interpreted as a regex pattern, but the shell
/// scripts in this monorepo only use plain-text placeholders (port numbers,
/// template source names, `_CHANGE_ME_DESCRIPTION_`), so literal matching
/// preserves the bash semantics without surprises.
pub fn replace_in_files(root: &Path, needle: &str, replacement: &str) -> Result<usize> {
    let mut touched = 0_usize;
    for path in walk_files_sorted(root)? {
        if should_skip_placeholder_file(&path) {
            continue;
        }
        let Ok(original) = fs_err::read_to_string(&path) else {
            // Skip non-UTF-8 files silently, matching `grep -rl` behavior
            // which only emits text files.
            continue;
        };
        if !original.contains(needle) {
            continue;
        }
        let replaced = original.replace(needle, replacement);
        fs_err::write(&path, replaced)
            .with_context(|| format!("write back after replacement: {}", path.display()))?;
        touched += 1;
    }
    Ok(touched)
}

/// Replace `needle` with `replacement` inside a single file.
///
/// Analogue of `sed -i "s/$needle/$replacement/g" "$file"`.
pub fn replace_in_file(path: &Path, needle: &str, replacement: &str) -> Result<()> {
    if !path.exists() {
        return Ok(());
    }
    let original =
        fs_err::read_to_string(path).with_context(|| format!("read: {}", path.display()))?;
    let replaced = original.replace(needle, replacement);
    fs_err::write(path, replaced).with_context(|| format!("write: {}", path.display()))?;
    Ok(())
}

/// Returns true iff the file's basename is in `PLACEHOLDER_SKIP_FILES`.
fn should_skip_placeholder_file(path: &Path) -> bool {
    path.file_name()
        .and_then(|n| n.to_str())
        .is_some_and(|name| PLACEHOLDER_SKIP_FILES.contains(&name))
}

/// Rename every file under `root` whose extension matches `from_ext` so it
/// gets a `.raw.<from_ext>` extension instead.
///
/// Analogue of:
/// ```bash
/// find "$root" -type f -name "*.$from_ext" | while read -r file; do
///     mv "$file" "${file%.$from_ext}.raw.$from_ext"
/// done
/// ```
///
/// Idempotent: files already ending in `.raw.<from_ext>` are skipped.
pub fn rename_ext_to_raw(root: &Path, from_ext: &str) -> Result<usize> {
    let raw_marker = format!(".raw.{from_ext}");
    let mut renamed = 0_usize;
    for path in walk_files_sorted(root)? {
        let Some(name) = path.file_name().and_then(|n| n.to_str()) else {
            continue;
        };
        if !name.ends_with(&format!(".{from_ext}")) {
            continue;
        }
        if name.ends_with(&raw_marker) {
            continue; // already renamed
        }
        let stem = name.trim_end_matches(&format!(".{from_ext}"));
        let new_name = format!("{stem}{raw_marker}");
        let new_path = path.with_file_name(new_name);
        fs_err::rename(&path, &new_path)
            .with_context(|| format!("rename: {} -> {}", path.display(), new_path.display()))?;
        renamed += 1;
    }
    Ok(renamed)
}

/// Rename a single file by appending `.raw` before its final extension.
///
/// `foo/.mockery.yml` -> `foo/.mockery.raw.yml`
/// `foo/mailer_test.go` -> `foo/mailer_test.raw.go`
pub fn rename_one_to_raw(path: &Path) -> Result<()> {
    if !path.exists() {
        return Ok(());
    }
    let Some(name) = path.file_name().and_then(|n| n.to_str()) else {
        return Err(anyhow!("non-UTF-8 file name: {}", path.display()));
    };
    let new_name = match name.rfind('.') {
        Some(dot_idx) if dot_idx > 0 => {
            let (stem, ext) = name.split_at(dot_idx);
            format!("{stem}.raw{ext}")
        }
        _ => format!("{name}.raw"),
    };
    let new_path = path.with_file_name(&new_name);
    if new_path.exists() {
        return Ok(()); // idempotent
    }
    fs_err::rename(path, &new_path)
        .with_context(|| format!("rename: {} -> {}", path.display(), new_path.display()))?;
    Ok(())
}

/// Prepend a block of text to every file matching `ext` under `root`.
///
/// Used by the Astro builder, which injects `---\nforce: true\n---\n` at
/// the top of every `.astro` file before they get renamed to `.raw.astro`.
pub fn prepend_to_files_with_ext(root: &Path, ext: &str, prefix: &str) -> Result<()> {
    for path in walk_files_sorted(root)? {
        let Some(name) = path.file_name().and_then(|n| n.to_str()) else {
            continue;
        };
        if !name.ends_with(&format!(".{ext}")) || name.contains(".raw.") {
            continue;
        }
        let existing =
            fs_err::read_to_string(&path).with_context(|| format!("read: {}", path.display()))?;
        if existing.starts_with(prefix) {
            continue; // idempotent
        }
        let new_content = format!("{prefix}{existing}");
        fs_err::write(&path, new_content)
            .with_context(|| format!("prepend: {}", path.display()))?;
    }
    Ok(())
}

/// Remove every file under `dir` that is NOT listed in `keep` (relative
/// paths). Also removes every empty directory left behind.
///
/// Used by the go-modular builder, which preserves `docs/embed.go` while
/// deleting every other file in `templates/go-modular/docs/`, and similarly
/// for `web/` (keeping `embed.go` + `static/index.html`).
pub fn keep_only(dir: &Path, keep_relative: &[&str]) -> Result<()> {
    if !dir.exists() {
        return Ok(());
    }
    let keep_abs: Vec<PathBuf> = keep_relative.iter().map(|r| dir.join(r)).collect();
    let files = walk_files_sorted(dir)?;
    for path in files {
        if keep_abs.contains(&path) {
            continue;
        }
        fs_err::remove_file(&path)
            .with_context(|| format!("keep_only remove: {}", path.display()))?;
    }
    // Remove empty directories bottom-up.
    remove_empty_dirs(dir);
    Ok(())
}

fn remove_empty_dirs(dir: &Path) {
    let mut subdirs: Vec<PathBuf> = WalkDir::new(dir)
        .into_iter()
        .filter_map(std::result::Result::ok)
        .filter(|e| e.file_type().is_dir())
        .map(|e| e.path().to_path_buf())
        .collect();
    subdirs.sort_by_key(|p| std::cmp::Reverse(p.components().count()));
    for sub in subdirs {
        if sub == dir {
            continue;
        }
        if fs_err::read_dir(&sub).is_ok_and(|mut i| i.next().is_none()) {
            fs_err::remove_dir(&sub).ok();
        }
    }
}

/// Modify a JSON file in place by applying `f` to the parsed `Value`.
///
/// Reads, deserializes, calls `f`, reserializes with pretty printing
/// (2-space indent, trailing newline) to match what `jq` would produce.
/// If the file does not exist the operation is a no-op.
pub fn modify_json_file<F>(path: &Path, f: F) -> Result<()>
where
    F: FnOnce(&mut serde_json::Value) -> Result<()>,
{
    if !path.exists() {
        return Ok(());
    }
    let raw =
        fs_err::read_to_string(path).with_context(|| format!("read json: {}", path.display()))?;
    let mut value: serde_json::Value =
        serde_json::from_str(&raw).with_context(|| format!("parse json: {}", path.display()))?;
    f(&mut value)?;
    let mut pretty = serde_json::to_string_pretty(&value)
        .with_context(|| format!("reserialize json: {}", path.display()))?;
    if !pretty.ends_with('\n') {
        pretty.push('\n');
    }
    fs_err::write(path, pretty).with_context(|| format!("write json: {}", path.display()))?;
    Ok(())
}

/// Recursive copy of every entry in `src` into `dst` — equivalent to
/// `cp -r "$src/." "$dst/"`.
///
/// `dst` must already exist; the caller creates it. Symlinks-to-files are
/// followed and copied as plain files (matches `cp -r`'s default; neither
/// `build-templates.sh` nor `builder/shared-ui.sh` use `-P` to preserve
/// symlinks).
pub fn copy_dir_contents(src: &Path, dst: &Path) -> Result<()> {
    for entry in fs_err::read_dir(src)? {
        let entry = entry?;
        let ty = entry.file_type()?;
        let from = entry.path();
        let to = dst.join(entry.file_name());
        if ty.is_dir() {
            fs_err::create_dir_all(&to)?;
            copy_dir_contents(&from, &to)?;
        } else if ty.is_file() {
            fs_err::copy(&from, &to)?;
        } else if ty.is_symlink() {
            let target = fs_err::read_link(&from)?;
            if target.is_file() {
                fs_err::copy(&from, &to)?;
            }
        }
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::tempdir;

    #[test]
    fn replace_in_files_skips_template_yml() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs_err::write(root.join("a.txt"), "port 8000").unwrap();
        fs_err::write(root.join("template.yml"), "port 8000").unwrap();
        replace_in_files(root, "8000", "{{ port_number }}").unwrap();
        assert_eq!(
            fs_err::read_to_string(root.join("a.txt")).unwrap(),
            "port {{ port_number }}"
        );
        assert_eq!(
            fs_err::read_to_string(root.join("template.yml")).unwrap(),
            "port 8000",
            "template.yml must be skipped"
        );
    }

    #[test]
    fn replace_in_files_is_literal_not_regex() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs_err::write(root.join("a.txt"), "value: $var").unwrap();
        replace_in_files(root, "$var", "{{ x }}").unwrap();
        assert_eq!(
            fs_err::read_to_string(root.join("a.txt")).unwrap(),
            "value: {{ x }}"
        );
    }

    #[test]
    fn rename_ext_to_raw_is_idempotent() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs_err::write(root.join("index.astro"), "hello").unwrap();
        rename_ext_to_raw(root, "astro").unwrap();
        assert!(root.join("index.raw.astro").exists());
        assert!(!root.join("index.astro").exists());
        rename_ext_to_raw(root, "astro").unwrap();
        assert!(root.join("index.raw.astro").exists());
    }

    #[test]
    fn rename_one_to_raw_handles_dotfile() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs_err::write(root.join(".mockery.yml"), "x").unwrap();
        rename_one_to_raw(&root.join(".mockery.yml")).unwrap();
        assert!(root.join(".mockery.raw.yml").exists());
        assert!(!root.join(".mockery.yml").exists());
    }

    #[test]
    fn prepend_to_files_is_idempotent() {
        let dir = tempdir().unwrap();
        let root = dir.path();
        fs_err::write(root.join("i.astro"), "<h1>hi</h1>").unwrap();
        // Matches what bash echo -e "---\nforce: true\n---\n" produces
        // (4 newlines: 3 from \n escapes + 1 from echo's own trailer).
        let prefix = "---\nforce: true\n---\n\n";
        prepend_to_files_with_ext(root, "astro", prefix).unwrap();
        let after = fs_err::read_to_string(root.join("i.astro")).unwrap();
        assert_eq!(after, "---\nforce: true\n---\n\n<h1>hi</h1>");
        prepend_to_files_with_ext(root, "astro", prefix).unwrap();
        let after2 = fs_err::read_to_string(root.join("i.astro")).unwrap();
        assert_eq!(after, after2, "prepend must be idempotent");
    }

    #[test]
    fn modify_json_file_preserves_shape_and_key_order() {
        let dir = tempdir().unwrap();
        let path = dir.path().join("package.json");
        // Use 3 keys so we can observe whether the key order is preserved
        // after the removal. With default `Map::remove` + `preserve_order`,
        // the behavior is `swap_remove`, which would move `z-last` into
        // position 0. We assert the original order is preserved via
        // `shift_remove`.
        fs_err::write(
            &path,
            r#"{"dependencies":{"@repo/shared-ui":"*","react":"18","z-last":"1"}}"#,
        )
        .unwrap();
        modify_json_file(&path, |v| {
            if let Some(deps) = v.get_mut("dependencies").and_then(|d| d.as_object_mut()) {
                deps.shift_remove("@repo/shared-ui");
            }
            Ok(())
        })
        .unwrap();
        let after_text = fs_err::read_to_string(&path).unwrap();
        let react_pos = after_text.find("\"react\"").unwrap();
        let zlast_pos = after_text.find("\"z-last\"").unwrap();
        assert!(
            react_pos < zlast_pos,
            "react must come BEFORE z-last (insertion order preserved): {after_text}"
        );
        let after: serde_json::Value = serde_json::from_str(&after_text).unwrap();
        assert!(after["dependencies"].get("@repo/shared-ui").is_none());
        assert_eq!(after["dependencies"]["react"], "18");
        assert_eq!(after["dependencies"]["z-last"], "1");
    }
}
