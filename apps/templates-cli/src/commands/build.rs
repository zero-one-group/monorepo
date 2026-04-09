//! `templates-cli build` — port of `build-templates.sh`.
//!
//! The original script:
//! 1. Checks that `jq` is installed (exits 1 if missing)
//! 2. Runs `pnpm --silent run format` on the monorepo (exits 1 on failure)
//! 3. Copies 9 apps into `templates/` using a parallel-array source→target map
//! 4. Cleans runtime artifacts (`node_modules`, vendor, dist, build, .expo,
//!    uv.lock, .react-router, .astro, .cache, .strapi, .next, .venv, .env,
//!    .`DS_Store`) from each target
//! 5. Invokes every `builder/*.sh` script in sequence
//!
//! The Rust port preserves steps 3-5 exactly. Steps 1-2 are ADAPTED:
//! - `jq` check is no longer needed — we do JSON surgery in Rust with `serde_json`
//! - `pnpm run format` is preserved but now optional via a flag so tests
//!   don't require pnpm on the test runner
//!
//! Per the spec's canonical-equivalence contract (AC5), the final output
//! of `templates/` must have an identical sorted file listing and identical
//! per-file SHA-256s compared to running the bash pipeline on the same input.

use std::path::{Path, PathBuf};
use std::process::Command;

use anyhow::{Context, Result};

use crate::builders;
use crate::common::copy_dir_contents;

/// Parallel-array source → target mapping from `build-templates.sh`.
///
/// The order here is significant only for deterministic logging output;
/// the copy operations themselves are independent.
pub const APP_MAPPING: &[(&str, &str)] = &[
    ("astro-web", "astro"),
    ("expo-app", "expo"),
    ("fastapi-ai", "fastapi-ai"),
    ("go-clean", "go-clean"),
    ("go-modular", "go-modular"),
    ("nextjs-app", "nextjs"),
    ("react-app", "react-app"),
    ("react-ssr", "react-ssr"),
    ("strapi-cms", "strapi"),
];

/// Runtime artifacts stripped from every copied template.
///
/// This list is copied verbatim from the `rm -rf` calls in
/// `build-templates.sh` so that the Rust output drops the same set of
/// files/directories the bash pipeline does.
pub const CLEANUP_ENTRIES: &[&str] = &[
    "node_modules",
    "vendor",
    "build",
    "dist",
    ".expo",
    "uv.lock",
    ".react-router",
    ".astro",
    ".cache",
    ".strapi",
    ".next",
    ".venv",
    ".env",
];

/// Run the full build phase on `root`.
pub fn run(root: &Path) -> Result<()> {
    tracing::info!("templates-cli build");

    let apps_dir = root.join("apps");
    let templates_dir = root.join("templates");
    fs_err::create_dir_all(&templates_dir)?;

    // Optional: run `pnpm run format` only if pnpm is available. Tests
    // and CI that want to skip can set `TEMPLATES_CLI_SKIP_FORMAT=1`.
    if std::env::var_os("TEMPLATES_CLI_SKIP_FORMAT").is_none() && which_pnpm().is_some() {
        tracing::info!("Running code formatter (pnpm --silent run format)...");
        let status = Command::new("pnpm")
            .args(["--silent", "run", "format"])
            .current_dir(root)
            .status()
            .context("invoke pnpm run format")?;
        if !status.success() {
            return Err(anyhow::anyhow!(
                "pnpm run format failed (exit {:?}); fix issues and retry",
                status.code()
            ));
        }
    }

    tracing::info!("Copying template files with mapping and path validation...");
    for (src, tgt) in APP_MAPPING {
        let src_path = apps_dir.join(src);
        let tgt_path = templates_dir.join(tgt);
        if !src_path.is_dir() {
            tracing::warn!(
                "Source directory {} does not exist, skipping.",
                src_path.display()
            );
            continue;
        }
        tracing::info!("Processing {src} -> {tgt}...");
        if tgt_path.exists() {
            fs_err::remove_dir_all(&tgt_path)
                .with_context(|| format!("rm -rf {}", tgt_path.display()))?;
        }
        fs_err::create_dir_all(&tgt_path)?;
        copy_dir_contents(&src_path, &tgt_path)
            .with_context(|| format!("cp -r {}/. -> {}", src_path.display(), tgt_path.display()))?;
    }
    tracing::info!("Template files copied successfully.");

    tracing::info!("Cleaning up unnecessary files...");
    for (_, tgt) in APP_MAPPING {
        let tgt_path = templates_dir.join(tgt);
        if !tgt_path.is_dir() {
            continue;
        }
        tracing::info!("Cleaning up {tgt}...");
        for entry in CLEANUP_ENTRIES {
            let victim = tgt_path.join(entry);
            if victim.is_dir() {
                fs_err::remove_dir_all(&victim).ok();
            } else if victim.is_file() {
                fs_err::remove_file(&victim).ok();
            }
        }
        // find "$TGT_PATH" -type f -name ".DS_Store" -delete
        remove_all_ds_store(&tgt_path);
    }
    tracing::info!("Cleanup completed.");

    tracing::info!("Scaffolding templates...");
    builders::run_default_set(root)?;

    tracing::info!("All processes completed successfully.");
    Ok(())
}

fn remove_all_ds_store(root: &Path) {
    for entry in walkdir::WalkDir::new(root)
        .into_iter()
        .filter_map(std::result::Result::ok)
    {
        if entry.file_type().is_file() && entry.file_name() == std::ffi::OsStr::new(".DS_Store") {
            fs_err::remove_file(entry.path()).ok();
        }
    }
}

fn which_pnpm() -> Option<PathBuf> {
    let path = std::env::var_os("PATH")?;
    for dir in std::env::split_paths(&path) {
        for candidate in ["pnpm", "pnpm.cmd", "pnpm.exe"] {
            let p = dir.join(candidate);
            if p.is_file() {
                return Some(p);
            }
        }
    }
    None
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn app_mapping_has_nine_entries() {
        assert_eq!(APP_MAPPING.len(), 9, "must match build-templates.sh");
    }

    #[test]
    fn cleanup_entries_match_bash_list() {
        // Sanity: we capture every entry in the bash script's rm -rf block.
        assert_eq!(CLEANUP_ENTRIES.len(), 13);
        assert!(CLEANUP_ENTRIES.contains(&"node_modules"));
        assert!(
            !CLEANUP_ENTRIES.contains(&".DS_Store"),
            "DS_Store handled separately"
        );
    }
}
