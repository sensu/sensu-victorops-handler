# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic
Versioning](http://semver.org/spec/v2.0.0.html).

## Unreleased

## [0.3.1] - 2021-04-13

### Changed
- Updated Readme

## [0.3.0] - 2021-01-28

### Changed
- Updated SDK and other module dependencies
- Reorganized README
- Use Go 1.14
- Code cleanup for linter

### Added
- Template support for message and entity ID
- Log output of submission
- Lint GitHub action

### Removed
- Windows 386 build from goreleaser

## [0.2.0] - 2020-08-14

### Changed
- Updated SDK to 0.8.0
- Use secrets boolean to avoid exposing routing_key

## [0.1.2] - 2020-03-13

### Changed
- Fixed gorelaser config ldflags to properly include version

## [0.1.1] - 2020-03-13

### Changed
- Updated README
- Removed unnecessary test code
- Added comments to make golint happy

## [0.1.0] - 2020-03-12

### Added
- Initial release
