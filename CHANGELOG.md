# Changelog

All notable changes to beads will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.30.0] - 2025-11-01

### Added

- Complete namespace refactor: All `BD_*` environment variables renamed to `BEADS_*`
- Enhanced version bump script to include `.beads/config.json`
- Improved test error reporting with proper error channels

### Changed

- Updated `SetEnvPrefix("BD")` to `SetEnvPrefix("BEADS")` in config system
- Refactored test patterns to capture server.Start() errors
- Shortened Unix socket paths in tests to avoid macOS 104-character limit

### Fixed

- Fixed daemon discovery test timeouts on macOS
- Fixed typos: "subeadsirectory" → "subdirectory" throughout codebase
- Corrected test variable naming: `testIssueBD1` → `testIssueBeads1`
- Fixed environment variable references in documentation and examples

## [0.24.0] - 2025-11-01

### Added

- Initial public release preparation
- GoReleaser configuration for multi-platform builds
- GitHub Actions workflows for CI and releases
- Comprehensive install script with fallback methods

[0.30.0]: https://github.com/shaneholloman/beads/compare/v0.24.0...v0.30.0
[0.24.0]: https://github.com/shaneholloman/beads/releases/tag/v0.24.0
