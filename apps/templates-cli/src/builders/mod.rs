//! Per-template builders.
//!
//! The bash pipeline has 11 `builder/*.sh` scripts that are 90% identical
//! boilerplate with small per-template tweaks. Instead of porting 11
//! near-duplicate Rust files, we express each builder as a [`BuilderSpec`]
//! and let a generic [`apply_spec`] function execute them all.
//!
//! When a template's logic doesn't fit the spec we attach a closure
//! via [`BuilderSpec::pre_custom`] (runs BEFORE the declarative pipeline)
//! or [`BuilderSpec::custom`] (runs AFTER). The two current users are:
//!
//! - **`go-modular`**: a `custom` pass that cleans `docs/` and `web/`
//!   after the placeholder/rename steps.
//! - **`shared-ui`**: a `pre_custom` pass that copies from
//!   `packages/shared-ui/` into `templates/shared-ui/` (instead of
//!   from another `templates/` dir) and strips runtime artifacts
//!   before the standard placeholder/rename flow runs.

use std::path::{Path, PathBuf};

use anyhow::{Context, Result, anyhow};

use crate::common::{
    copy_dir_contents, keep_only, modify_json_file, prepend_to_files_with_ext, rename_ext_to_raw,
    rename_one_to_raw, replace_in_file, replace_in_files,
};

/// Context passed to `pre_custom` and `custom` closures.
///
/// Holds every path a custom pass is likely to need, so callers don't
/// have to reconstruct them from a bare `&Path`. All four fields are
/// borrows of the same lifetime, so the struct is zero-cost.
pub struct BuildCtx<'a> {
    /// Monorepo root (the directory containing `apps/`, `templates/`,
    /// `packages/`, etc.).
    pub root: &'a Path,
    /// `<root>/templates` — the directory every template ends up in.
    pub templates_dir: &'a Path,
    /// `<root>/packages` — source for the `shared-ui` builder.
    pub packages_dir: &'a Path,
    /// `<root>/templates/<target_name>` — the template directory this
    /// builder is operating on.
    pub target: &'a Path,
}

/// Immutable spec for a single template build.
///
/// Every field here was derived from reading the corresponding
/// `builder/<name>.sh` shell script verbatim. Any divergence would
/// violate the canonical-equivalence contract (AC5/6 of the spec).
pub struct BuilderSpec {
    /// Name used for CLI lookup (`templates-cli builder <name>`). Matches
    /// the shell-script basename without the `.sh` suffix.
    pub name: &'static str,

    /// Source directory name under `templates/` before rename. Usually the
    /// same as the target, but differs for astro (`astro-web` → `astro`),
    /// nextjs (`nextjs-app` → `nextjs`), etc.
    pub source_name: &'static str,

    /// Target directory name under `templates/`. The final name users see
    /// in `moon generate`.
    pub target_name: &'static str,

    /// Human-readable banner printed when the builder runs. Matches the
    /// "Building X project templates..." echo at the top of each shell script.
    pub banner: &'static str,

    /// If set, every file in the target directory that contains this
    /// literal port string is rewritten with `{{ port_number }}`.
    pub default_port: Option<&'static str>,

    /// Per-template file-extension renames (e.g. `.py` → `.raw.py`).
    pub rename_exts_to_raw: &'static [&'static str],

    /// Single-file renames (relative to target dir). Used for things like
    /// `.mockery.yml` → `.mockery.raw.yml` and specific test files.
    pub rename_single_files: &'static [&'static str],

    /// Whether to remove `@repo/shared-ui` from `package.json` via JSON
    /// surgery. Set for every frontend builder that currently `jq`-deletes
    /// it from package.json (react-app, react-ssr, nextjs-app, strapi-cms,
    /// tanstack-start — five builders).
    pub strip_shared_ui_from_package_json: bool,

    /// Whether to remove the `../../packages/shared-ui` entry from
    /// `tsconfig.json`'s `references` array. Set for every frontend builder
    /// that currently `jq`-filters this out. **Not** set for `strapi-cms`
    /// — strapi's tsconfig.json is JSONC (has comments) and the bash
    /// script intentionally never touches it, so we match that exactly
    /// (four builders: react-app, react-ssr, nextjs-app, tanstack-start).
    pub strip_shared_ui_from_tsconfig_json: bool,

    /// Whether to prepend `---\nforce: true\n---\n` to `.astro` files
    /// before renaming them. Astro-specific.
    pub astro_frontmatter_prepend: bool,

    /// Optional PRE-custom pass, run BEFORE step 1 of `apply_spec`.
    ///
    /// Used by `shared-ui` to copy from `packages/shared-ui/` into
    /// `templates/shared-ui/` (a different source location than the
    /// standard `templates/<source_name>` lookup).
    pub pre_custom: Option<fn(ctx: &BuildCtx<'_>) -> Result<()>>,

    /// Optional post-custom pass, run AFTER all the declarative steps.
    /// Used by `go-modular` for its `docs/` and `web/` directory pruning.
    pub custom: Option<fn(ctx: &BuildCtx<'_>) -> Result<()>>,
}

/// The full registry. Order matches the order in `build-templates.sh`.
pub const REGISTRY: &[&BuilderSpec] = &[
    &ASTRO,
    &EXPO_APP,
    &FASTAPI_AI,
    &GO_CLEAN,
    &GO_MODULAR,
    &NEXTJS_APP,
    &REACT_APP,
    &REACT_SSR,
    &RUST_AI,
    &RUST_CLEAN,
    &RUST_MODULAR,
    &STRAPI_CMS,
    &TANSTACK_START,
    &SHARED_UI,
];

/// Look up a builder by name and run it against `root`.
pub fn run_builder(name: &str, root: &Path) -> Result<()> {
    let spec = REGISTRY
        .iter()
        .find(|s| s.name == name)
        .copied()
        .ok_or_else(|| anyhow!("unknown builder: {name}"))?;
    apply_spec(spec, root)
}

/// The 10 builders invoked by the original `build-templates.sh` orchestrator,
/// in the exact order it called them.
///
/// **`tanstack-start` is intentionally absent.** The bash orchestrator never
/// invoked `builder/tanstack-start.sh` from `build-templates.sh` — it could
/// only be run on demand via `bash builder/tanstack-start.sh`. Including
/// tanstack-start in the default `templates-cli build` set would silently
/// reformat `templates/tanstack-start/` on every full rebuild, drifting from
/// the bash baseline. Users who want to scaffold tanstack-start use
/// `templates-cli builder tanstack-start` explicitly.
const DEFAULT_BUILD_SET: &[&BuilderSpec] = &[
    &ASTRO,
    &EXPO_APP,
    &FASTAPI_AI,
    &GO_CLEAN,
    &GO_MODULAR,
    &NEXTJS_APP,
    &REACT_APP,
    &REACT_SSR,
    &RUST_AI,
    &RUST_CLEAN,
    &RUST_MODULAR,
    &STRAPI_CMS,
    &SHARED_UI,
];

/// Run the default set of 10 builders that the bash `build-templates.sh`
/// orchestrator invoked, in its declared order.
///
/// This is what `templates-cli build` calls. To run a single builder
/// (including `tanstack-start`), use `templates-cli builder <name>`.
pub fn run_default_set(root: &Path) -> Result<()> {
    for spec in DEFAULT_BUILD_SET {
        apply_spec(spec, root)?;
    }
    Ok(())
}

/// Execute a [`BuilderSpec`] against a monorepo root.
///
/// The order of operations mirrors the shell-script order EXACTLY:
/// 0. Pre-custom pass (shared-ui copies from packages/ here)
/// 1. If source directory exists AND source != target, rename source → target
/// 2. Replace placeholders in `moon.yml` (source name, description)
/// 3. Replace port number across all files (except `template.yml`)
/// 4. Replace source name across all files (except `template.yml`)
/// 5. Replace `_CHANGE_ME_DESCRIPTION_` across all files (except `template.yml`)
/// 6. Astro frontmatter prepend (astro only)
/// 7. File-extension renames (.py → .raw.py, .astro → .raw.astro, etc.)
/// 8. Single-file renames (.mockery.yml → .mockery.raw.yml, etc.)
/// 9. Strip `@repo/shared-ui` from package.json + tsconfig.json (if enabled)
/// 10. Post-custom pass (go-modular's docs/web cleanup runs here)
pub fn apply_spec(spec: &BuilderSpec, root: &Path) -> Result<()> {
    tracing::info!("{}", spec.banner);

    let templates_dir: PathBuf = root.join("templates");
    let packages_dir: PathBuf = root.join("packages");
    let source_path = templates_dir.join(spec.source_name);
    let target_path = templates_dir.join(spec.target_name);

    // Step 0: pre-custom pass. Must build a ctx with `target = target_path`
    // even though target may not exist yet — pre_custom is allowed to
    // create it (that's the whole point for shared-ui).
    if let Some(pre_custom) = spec.pre_custom {
        let ctx = BuildCtx {
            root,
            templates_dir: &templates_dir,
            packages_dir: &packages_dir,
            target: &target_path,
        };
        pre_custom(&ctx)?;
    }

    // Step 1: source → target rename when they differ.
    if source_path.exists() && spec.source_name != spec.target_name {
        fs_err::rename(&source_path, &target_path).with_context(|| {
            format!(
                "rename source -> target: {} -> {}",
                source_path.display(),
                target_path.display()
            )
        })?;
    }

    if !target_path.exists() {
        tracing::warn!(
            "target dir does not exist, skipping: {}",
            target_path.display()
        );
        return Ok(());
    }

    // Step 2: moon.yml placeholder replacement (strict order — moon.yml first).
    let moon_yml = target_path.join("moon.yml");
    if moon_yml.exists() {
        replace_in_file(
            &moon_yml,
            spec.source_name,
            "{{ package_name | kebab_case }}",
        )?;
        replace_in_file(
            &moon_yml,
            "_CHANGE_ME_DESCRIPTION_",
            "{{ package_description }}",
        )?;
    }

    // Step 3: port number replacement across all files (bash order: port first,
    // then source name, then description — we preserve that exactly).
    if let Some(port) = spec.default_port {
        replace_in_files(&target_path, port, "{{ port_number }}")?;
    }

    // Step 4: source name replacement across all files.
    replace_in_files(
        &target_path,
        spec.source_name,
        "{{ package_name | kebab_case }}",
    )?;

    // Step 5: description placeholder replacement across all files.
    // For shared-ui this is a no-op because packages/shared-ui has no
    // `_CHANGE_ME_DESCRIPTION_` literals (verified 2026-04-08).
    replace_in_files(
        &target_path,
        "_CHANGE_ME_DESCRIPTION_",
        "{{ package_description }}",
    )?;

    // Step 6: astro frontmatter prepend (must run BEFORE the .astro rename).
    // The bash script uses `echo -e "---\nforce: true\n---\n"` which emits
    // `---\nforce: true\n---\n` plus an extra trailing newline from echo
    // itself, giving a total of 4 newlines (blank line after the second
    // `---`). We must match that exactly.
    if spec.astro_frontmatter_prepend {
        prepend_to_files_with_ext(&target_path, "astro", "---\nforce: true\n---\n\n")?;
    }

    // Step 7: file-extension renames.
    for ext in spec.rename_exts_to_raw {
        rename_ext_to_raw(&target_path, ext)?;
    }

    // Step 8: single-file renames.
    for rel in spec.rename_single_files {
        rename_one_to_raw(&target_path.join(rel))?;
    }

    // Step 9: strip @repo/shared-ui references (package.json and/or tsconfig.json).
    if spec.strip_shared_ui_from_package_json {
        strip_shared_ui_from_package_json(&target_path)?;
    }
    if spec.strip_shared_ui_from_tsconfig_json {
        strip_shared_ui_from_tsconfig_json(&target_path)?;
    }

    // Step 10: post-custom pass.
    if let Some(custom) = spec.custom {
        let ctx = BuildCtx {
            root,
            templates_dir: &templates_dir,
            packages_dir: &packages_dir,
            target: &target_path,
        };
        custom(&ctx)?;
    }

    Ok(())
}

/// Remove `@repo/shared-ui` from `package.json`'s `.dependencies` (and
/// `.devDependencies` as a safety net — some frontends keep it there).
///
/// Mirrors the `jq` invocation present in every frontend builder:
/// `jq 'if .dependencies then .dependencies |= del(.["@repo/shared-ui"]) else . end'`.
///
/// IMPORTANT: uses `shift_remove` instead of `remove` so the insertion
/// order of the remaining keys is preserved. With the `preserve_order`
/// feature enabled on `serde_json`, the default `Map::remove` uses
/// `swap_remove` which swaps the LAST element into the removed slot,
/// silently reordering the keys. `jq del(.[...])` preserves order, so
/// we must match that behavior to stay canonical-equivalent.
fn strip_shared_ui_from_package_json(target: &Path) -> Result<()> {
    let package_json = target.join("package.json");
    modify_json_file(&package_json, |value| {
        if let Some(deps) = value
            .get_mut("dependencies")
            .and_then(|d| d.as_object_mut())
        {
            deps.shift_remove("@repo/shared-ui");
        }
        if let Some(deps) = value
            .get_mut("devDependencies")
            .and_then(|d| d.as_object_mut())
        {
            deps.shift_remove("@repo/shared-ui");
        }
        Ok(())
    })?;
    Ok(())
}

/// Remove `../../packages/shared-ui` from `tsconfig.json`'s `.references`.
///
/// Mirrors `jq 'if .references then .references |= map(select(.path != "../../packages/shared-ui")) else . end'`.
/// Strapi-cms intentionally does NOT use this because its tsconfig.json is
/// JSONC (has `//` comments) and the bash script never touches it.
fn strip_shared_ui_from_tsconfig_json(target: &Path) -> Result<()> {
    let tsconfig_json = target.join("tsconfig.json");
    modify_json_file(&tsconfig_json, |value| {
        if let Some(refs) = value.get_mut("references").and_then(|r| r.as_array_mut()) {
            refs.retain(|r| {
                r.get("path").and_then(|p| p.as_str()) != Some("../../packages/shared-ui")
            });
        }
        Ok(())
    })?;
    Ok(())
}

// -----------------------------------------------------------------------
// Per-template specs below. Each one is a straight translation of the
// corresponding shell script. The comments cite the source script and any
// non-obvious ports.
// -----------------------------------------------------------------------

/// `builder/astro.sh`: astro-web → astro, port 4321, frontmatter prepend,
/// .astro → .raw.astro.
pub const ASTRO: BuilderSpec = BuilderSpec {
    name: "astro",
    source_name: "astro-web",
    target_name: "astro",
    banner: "Building Astro project templates...",
    default_port: Some("4321"),
    rename_exts_to_raw: &["astro"],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: true,
    pre_custom: None,
    custom: None,
};

/// `builder/expo-app.sh`: expo-app → expo, no port replacement.
///
/// Note: the original script does NOT do a port replacement because expo's
/// default dev port (19000/19006) isn't in the source files as a hardcoded
/// literal. We preserve that by setting `default_port = None`.
pub const EXPO_APP: BuilderSpec = BuilderSpec {
    name: "expo-app",
    source_name: "expo-app",
    target_name: "expo",
    banner: "Building Expo project templates...",
    default_port: None,
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/fastapi-ai.sh`: fastapi-ai (same name), port 8080, .py → .raw.py.
pub const FASTAPI_AI: BuilderSpec = BuilderSpec {
    name: "fastapi-ai",
    source_name: "fastapi-ai",
    target_name: "fastapi-ai",
    banner: "Building FastAPI-AI project templates...",
    default_port: Some("8080"),
    rename_exts_to_raw: &["py"],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/go-clean.sh`: go-clean (same name), port 8000,
/// .mockery.yml → .mockery.raw.yml.
pub const GO_CLEAN: BuilderSpec = BuilderSpec {
    name: "go-clean",
    source_name: "go-clean",
    target_name: "go-clean",
    banner: "Building Go Clean Architecture project templates...",
    default_port: Some("8000"),
    rename_exts_to_raw: &[],
    rename_single_files: &[".mockery.yml"],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/go-modular.sh`: go-modular (same name), port 8000,
/// multiple single-file renames, docs/ and web/ cleanups via custom pass.
pub const GO_MODULAR: BuilderSpec = BuilderSpec {
    name: "go-modular",
    source_name: "go-modular",
    target_name: "go-modular",
    banner: "Building Go Modular project templates...",
    default_port: Some("8000"),
    rename_exts_to_raw: &[],
    rename_single_files: &[".mockery.yml", "internal/notification/mailer_test.go"],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: Some(go_modular_custom),
};

fn go_modular_custom(ctx: &BuildCtx<'_>) -> Result<()> {
    // Rename templates/emails/*.html to *.raw.html (only that subdirectory).
    let emails = ctx.target.join("templates/emails");
    if emails.exists() {
        rename_ext_to_raw(&emails, "html")?;
    }
    // Keep only docs/embed.go
    keep_only(&ctx.target.join("docs"), &["embed.go"])?;
    // Keep only web/embed.go and web/static/index.html
    keep_only(&ctx.target.join("web"), &["embed.go", "static/index.html"])?;
    Ok(())
}

/// Rust port of fastapi-ai: rust-ai, port 8080. Custom pass replaces
/// the Rust underscore crate name (`rust_ai`) with the `snake_case`
/// template variable.
pub const RUST_AI: BuilderSpec = BuilderSpec {
    name: "rust-ai",
    source_name: "rust-ai",
    target_name: "rust-ai",
    banner: "Building Rust AI project templates...",
    default_port: Some("8080"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: Some(rust_ai_custom),
};

fn rust_ai_custom(ctx: &BuildCtx<'_>) -> Result<()> {
    replace_in_files(ctx.target, "rust_ai", "{{ package_name | snake_case }}")?;
    Ok(())
}

/// Rust port of go-clean: rust-clean, port 8000. Custom pass replaces
/// the Rust underscore crate name (`rust_clean`).
pub const RUST_CLEAN: BuilderSpec = BuilderSpec {
    name: "rust-clean",
    source_name: "rust-clean",
    target_name: "rust-clean",
    banner: "Building Rust Clean Architecture project templates...",
    default_port: Some("8000"),
    rename_exts_to_raw: &[],
    rename_single_files: &[".mockery.yml"],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: Some(rust_clean_custom),
};

fn rust_clean_custom(ctx: &BuildCtx<'_>) -> Result<()> {
    replace_in_files(ctx.target, "rust_clean", "{{ package_name | snake_case }}")?;
    Ok(())
}

/// Rust port of go-modular: rust-modular, port 8000. Custom pass
/// replaces the Rust underscore crate name (`rust_modular`).
pub const RUST_MODULAR: BuilderSpec = BuilderSpec {
    name: "rust-modular",
    source_name: "rust-modular",
    target_name: "rust-modular",
    banner: "Building Rust Modular project templates...",
    default_port: Some("8000"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: Some(rust_modular_custom),
};

fn rust_modular_custom(ctx: &BuildCtx<'_>) -> Result<()> {
    replace_in_files(
        ctx.target,
        "rust_modular",
        "{{ package_name | snake_case }}",
    )?;
    Ok(())
}

/// `builder/nextjs-app.sh`: nextjs-app → nextjs, port 3200, strip shared-ui.
pub const NEXTJS_APP: BuilderSpec = BuilderSpec {
    name: "nextjs-app",
    source_name: "nextjs-app",
    target_name: "nextjs",
    banner: "Building Next.js project templates...",
    default_port: Some("3200"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: true,
    strip_shared_ui_from_tsconfig_json: true,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/react-app.sh`: react-app (same name), port 3000, strip shared-ui.
pub const REACT_APP: BuilderSpec = BuilderSpec {
    name: "react-app",
    source_name: "react-app",
    target_name: "react-app",
    banner: "Building React SPA project templates...",
    default_port: Some("3000"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: true,
    strip_shared_ui_from_tsconfig_json: true,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/react-ssr.sh`: react-ssr (same name), port 3100, strip shared-ui.
pub const REACT_SSR: BuilderSpec = BuilderSpec {
    name: "react-ssr",
    source_name: "react-ssr",
    target_name: "react-ssr",
    banner: "Building React SSR project templates...",
    default_port: Some("3100"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: true,
    strip_shared_ui_from_tsconfig_json: true,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/strapi-cms.sh`: strapi-cms → strapi, port 1337, strip shared-ui
/// from `package.json` ONLY. Strapi's `tsconfig.json` is JSONC (has `//`
/// comments) and the bash script intentionally never touches it.
pub const STRAPI_CMS: BuilderSpec = BuilderSpec {
    name: "strapi-cms",
    source_name: "strapi-cms",
    target_name: "strapi",
    banner: "Building Strapi CMS project templates...",
    default_port: Some("1337"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: true,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/tanstack-start.sh`: tanstack-start (same name), port 3300,
/// strip shared-ui.
pub const TANSTACK_START: BuilderSpec = BuilderSpec {
    name: "tanstack-start",
    source_name: "tanstack-start",
    target_name: "tanstack-start",
    banner: "Building TanStack Start project templates...",
    default_port: Some("3300"),
    rename_exts_to_raw: &[],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: true,
    strip_shared_ui_from_tsconfig_json: true,
    astro_frontmatter_prepend: false,
    pre_custom: None,
    custom: None,
};

/// `builder/shared-ui.sh`: special — copies from `packages/shared-ui` into
/// `templates/shared-ui`, strips runtime artifacts, then runs the standard
/// placeholder/rename pipeline (which handles `.tsx`/`.mdx` → `.raw.*`).
///
/// The copy step is wired through `pre_custom` so it runs BEFORE step 1
/// of `apply_spec`. The standard pipeline handles everything else.
pub const SHARED_UI: BuilderSpec = BuilderSpec {
    name: "shared-ui",
    source_name: "shared-ui",
    target_name: "shared-ui",
    banner: "Building Shared UI project templates...",
    default_port: None,
    rename_exts_to_raw: &["tsx", "mdx"],
    rename_single_files: &[],
    strip_shared_ui_from_package_json: false,
    strip_shared_ui_from_tsconfig_json: false,
    astro_frontmatter_prepend: false,
    pre_custom: Some(shared_ui_pre_custom),
    custom: None,
};

/// Cleanup entries removed from `templates/shared-ui/` after copy.
///
/// Mirrors the `rm -rf` calls in `builder/shared-ui.sh` exactly. Note
/// that `storybook-static` is specific to shared-ui and is NOT in the
/// general `CLEANUP_ENTRIES` list in `commands/build.rs`.
const SHARED_UI_CLEANUP_ENTRIES: &[&str] = &["node_modules", "storybook-static", "dist", "build"];

/// `pre_custom` pass for shared-ui.
///
/// Mirrors the first half of `builder/shared-ui.sh`:
/// 1. Remove any existing `templates/shared-ui/`
/// 2. Recursively copy `packages/shared-ui/` → `templates/shared-ui/`
/// 3. Strip runtime artifacts (`node_modules`, `storybook-static`, `dist`,
///    `build`, plus the `.DS_Store` find-and-delete pass)
///
/// After this runs, `apply_spec` continues with its standard pipeline
/// (step 1 source/target rename is a no-op because both names equal
/// "shared-ui", step 2+ do the placeholder replacement and .tsx/.mdx
/// renames).
fn shared_ui_pre_custom(ctx: &BuildCtx<'_>) -> Result<()> {
    let src = ctx.packages_dir.join("shared-ui");
    if !src.is_dir() {
        tracing::warn!(
            "shared-ui source {} does not exist, skipping",
            src.display()
        );
        return Ok(());
    }

    if ctx.target.exists() {
        fs_err::remove_dir_all(ctx.target)
            .with_context(|| format!("rm -rf {}", ctx.target.display()))?;
    }
    fs_err::create_dir_all(ctx.target)?;
    copy_dir_contents(&src, ctx.target)
        .with_context(|| format!("cp -R {} -> {}", src.display(), ctx.target.display()))?;

    // Strip runtime artifacts from the freshly-copied target.
    for entry in SHARED_UI_CLEANUP_ENTRIES {
        let victim = ctx.target.join(entry);
        if victim.is_dir() {
            fs_err::remove_dir_all(&victim).ok();
        } else if victim.is_file() {
            fs_err::remove_file(&victim).ok();
        }
    }
    // `find "$TARGET_PATH" -type f -name ".DS_Store" -delete`
    for entry in walkdir::WalkDir::new(ctx.target)
        .into_iter()
        .filter_map(std::result::Result::ok)
    {
        if entry.file_type().is_file() && entry.file_name() == std::ffi::OsStr::new(".DS_Store") {
            fs_err::remove_file(entry.path()).ok();
        }
    }
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn registry_has_all_fourteen_builders() {
        let names: Vec<&str> = REGISTRY.iter().map(|s| s.name).collect();
        assert_eq!(names.len(), 14);
        assert!(names.contains(&"astro"));
        assert!(names.contains(&"expo-app"));
        assert!(names.contains(&"fastapi-ai"));
        assert!(names.contains(&"go-clean"));
        assert!(names.contains(&"go-modular"));
        assert!(names.contains(&"nextjs-app"));
        assert!(names.contains(&"react-app"));
        assert!(names.contains(&"react-ssr"));
        assert!(names.contains(&"rust-ai"));
        assert!(names.contains(&"rust-clean"));
        assert!(names.contains(&"rust-modular"));
        assert!(names.contains(&"strapi-cms"));
        assert!(names.contains(&"tanstack-start"));
        assert!(names.contains(&"shared-ui"));
    }

    #[test]
    fn registry_names_are_unique() {
        let mut names: Vec<&str> = REGISTRY.iter().map(|s| s.name).collect();
        names.sort_unstable();
        let unique_count = names.iter().collect::<std::collections::HashSet<_>>().len();
        assert_eq!(unique_count, names.len(), "builder names must be unique");
    }

    #[test]
    fn default_build_set_excludes_tanstack_start() {
        // The bash `build-templates.sh` orchestrator NEVER invoked
        // `builder/tanstack-start.sh`. The Rust default-build-set must
        // match that exclusion exactly, or `templates-cli build` would
        // silently reformat `templates/tanstack-start/` on every run and
        // drift from the bash baseline.
        let names: Vec<&str> = DEFAULT_BUILD_SET.iter().map(|s| s.name).collect();
        assert!(
            !names.contains(&"tanstack-start"),
            "tanstack-start must NOT be in the default build set; \
             it is invoked on demand only via `templates-cli builder tanstack-start`"
        );
        assert_eq!(
            names.len(),
            13,
            "default build set must have exactly 13 entries"
        );
    }

    #[test]
    fn default_build_set_order_matches_bash_orchestrator() {
        let expected = [
            "astro",
            "expo-app",
            "fastapi-ai",
            "go-clean",
            "go-modular",
            "nextjs-app",
            "react-app",
            "react-ssr",
            "rust-ai",
            "rust-clean",
            "rust-modular",
            "strapi-cms",
            "shared-ui",
        ];
        let actual: Vec<&str> = DEFAULT_BUILD_SET.iter().map(|s| s.name).collect();
        assert_eq!(
            actual, expected,
            "default build set order must match `build-templates.sh` exactly"
        );
    }

    #[test]
    fn frontend_builders_strip_shared_ui_from_package_json() {
        // All 5 frontend builders strip @repo/shared-ui from package.json.
        for spec in REGISTRY {
            if matches!(
                spec.name,
                "react-app" | "react-ssr" | "nextjs-app" | "strapi-cms" | "tanstack-start"
            ) {
                assert!(
                    spec.strip_shared_ui_from_package_json,
                    "{} must strip @repo/shared-ui from package.json",
                    spec.name
                );
            }
        }
    }

    #[test]
    fn strapi_does_not_touch_tsconfig_json() {
        // Strapi's tsconfig.json is JSONC (has comments) and the bash
        // script intentionally never touches it. Only react-app, react-ssr,
        // nextjs-app, and tanstack-start strip @repo/shared-ui from
        // tsconfig.json — NOT strapi-cms.
        let frontends_that_touch_tsconfig =
            ["react-app", "react-ssr", "nextjs-app", "tanstack-start"];
        for spec in REGISTRY {
            if frontends_that_touch_tsconfig.contains(&spec.name) {
                assert!(
                    spec.strip_shared_ui_from_tsconfig_json,
                    "{} must strip @repo/shared-ui from tsconfig.json",
                    spec.name
                );
            }
            if spec.name == "strapi-cms" {
                assert!(
                    !spec.strip_shared_ui_from_tsconfig_json,
                    "strapi-cms must NOT touch tsconfig.json (JSONC, bash skips it)"
                );
            }
        }
    }

    #[test]
    fn only_astro_has_frontmatter_prepend() {
        for spec in REGISTRY {
            if spec.name == "astro" {
                assert!(spec.astro_frontmatter_prepend);
            } else {
                assert!(!spec.astro_frontmatter_prepend);
            }
        }
    }
}
