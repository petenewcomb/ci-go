# Changelog

The `ci-go` project adheres to [Semantic Versioning]. This file documents all
notable changes to this project and generally follows the [Keep a Changelog]
format.

## [Unreleased]

### Added
- Development documentation files (`DEVELOPMENT.md`, `.github/copilot-instructions.md`)
- Working notes documentation (`WORKING_NOTES.md`)
- Pre-commit hook script (`.githooks/pre-commit`)
- LICENSE.txt file
- Permissions declaration to ci.yml workflow
- Multi-module support for all workflows (build, lint, test, vuln)
- Reusable `modules.yml` workflow for discovering Go modules in repositories
- Enhanced golangci-lint configuration with security, performance, and code quality linters
- License auditing with `go-licenses` to prevent restrictive dependencies
- Matrix-based test execution with granular timeout controls for each test configuration
- Comprehensive golangci-lint configuration validation in pre-commit hook

### Changed
- Moved `go generate` validation from build workflow to lint workflow for better separation of concerns
- Fixed syntax error in release.yml workflow (missing quote)
- Workflows now use matrix strategy to test each Go module independently
- Coverage reporting limited to root module only
- Separated repository-wide operations (formatting, linting) from per-module operations (tidy, vet, generate)
- Replaced generic `run:` job IDs with descriptive names (`build:`, `test:`, `vuln:`, `release:`)
- Test workflow unified with matrix-based parallel execution and automatic timeout calculation
- Pre-commit hook optimized to run repository-wide operations once, per-module operations separately

### Fixed
- Input defaulting behavior in workflows
- Various workflow configuration issues
- Workflows now handle repositories with multiple Go modules (matches pre-commit hook behavior)

## [0.0.8] - 2025-04-17

### Added

- Pass-through input parameters for CI workflow

## [0.0.7] - 2025-04-17

### Changed

- Moved triggers to omnibus "Check" (`ci.yml`) workflow
- Refactored defaulting now that component flows are always called

## [0.0.7] - 2025-04-17

### Changed

- Moved triggers to separate "Check" workflow
- Refactored defaulting now that component flows are always called

## [0.0.6] - 2025-04-17

### Added

- `jobTimeoutMinutes`, `testTimeout`, `golangciTimeout` input parameters
- Concurrency group configuration for all workflows
- Tests now run with `-cpu 1,2,4`
- Test execution on push and PR
- Weekly cronjob execution for `test` and `lint` workflows

### Changed

- Skip pinging `pkg.go.dev` unless new `publishDocs` input parameter is true
- Weekly cronjobs now scheduled for Saturdays at 2pm UTC 
- Removed `skipGenerate` parameter

## [0.0.5] - 2025-04-16

### Fixed

- Missing line in test workflow

## [0.0.4] - 2025-04-16

### Changed

- Moved testing out of build and merged into test workflow
- Moved linting steps out of build and merged into lint workflow
- Bumped test timeout to 30 min

## [0.0.3] - 2025-04-10

### Added
- Build status badge

### Fixed
- Parameter substitution in release action ping to pkg.go.dev

## [0.0.2] - 2025-04-10

### Added
- Made actions user-triggerable from the GitHub web interface

### Changed

- Use ncruces/go-coverage-report instead of codecov/codecov-action
- Made release action use what's in CHANGELOG.md to get the tag and release
  notes

### Fixed

- General operation of actions on the ci-go repository itself
- Checkout before installing Go to avoid [missing `go.sum`
  warning](https://github.com/actions/setup-go/issues/427#issuecomment-2273249463)
- Removed broken golangci configuration to use defaults instead

## [0.0.1] - 2025-04-10

### Added

- Initial codebase forked from https://github.com/cristalhq/.github
- Codebase trimmed and focused to serve petenewcomb's projects

[0.0.1]: https://github.com/petenewcomb/ci-go/releases/tag/v0.0.1
[Keep a Changelog]: https://keepachangelog.com/en/1.1.0/
[Semantic Versioning]: https://semver.org/spec/v2.0.0.html
