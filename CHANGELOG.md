# Changelog

The `ci-go` project adheres to [Semantic Versioning]. This file documents all
notable changes to this project and generally follows the [Keep a Changelog]
format.

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
