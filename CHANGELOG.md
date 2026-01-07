# Changelog

All notable changes to this project will be documented in this file.

The format follows **Keep a Changelog**  
and this project adheres to **Semantic Versioning**.

---

## [0.1.0] - 2026-01-07

Initial open source release.

### Added
- Profile-based environment presets.
- Temporary application via `eval "$(gpx use <profile>)"`.
- Persistent application via `gpx apply` using a managed rc block.
- Safe shell quoting for exported values.
- Commands:
    - `list`, `status`, `use`, `set`, `unset`
    - `diff`, `apply`, `version`
- Config editing via CLI:
    - `gpx profile add|rm|rename|show|set|unset`
- Active profile marker in `gpx list` (`*` = last used/applied).
- Config validation on load.
- Atomic rc-file replacement to avoid corruption.

### Changed
- `gpx apply` does **not** create backup files by default.
  Use `--backup` to enable backups explicitly.

### Security
- rc files are updated atomically.
- No execution of shell code during config parsing.
