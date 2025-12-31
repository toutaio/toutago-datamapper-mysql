# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed
- Updated minimum Go version to 1.22
- Removed local replace directive for independent module usage

### Added
- Comprehensive package documentation (doc.go)

## [0.1.0] - 2024-12-24

### Added
- Initial release
- MySQL adapter implementation for toutago-datamapper
- Full CRUD operations (Create, Read, Update, Delete)
- Bulk insert support for efficient batch operations
- Named parameter substitution ({param_name})
- Auto-generated ID handling (auto-increment)
- Optimistic locking support
- Connection pooling configuration
- Custom SQL execution and stored procedures
- CQRS pattern support via source configuration

[Unreleased]: https://github.com/toutaio/toutago-datamapper-mysql/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/toutaio/toutago-datamapper-mysql/releases/tag/v0.1.0
