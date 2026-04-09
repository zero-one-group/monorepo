//! Configuration — env loading only (no separate logging/instrumentation
//! config module in the Rust port; those live in `lib.rs` and are small
//! enough to inline).

pub mod env;
