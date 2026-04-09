//! Top-level command implementations.
//!
//! Each module here corresponds to one entrypoint of the shell pipeline:
//! - [`build`] — `build-templates.sh`
//! - [`makezip`] — `makezip.sh`

pub mod build;
pub mod makezip;
