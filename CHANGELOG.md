# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [2.0.1] - 2025-08-23

### Fixed

- Renamed СtxManager to CtxManager without "С" cyrillic letter.

### Changes

- Bumped test up to go1.25.

## [2.0.0] - 2024-09-05

### Fixed

- CI/CD passed always even if some of them failed.

### Changes

- Bumped tests up to go1.23.

## [2.0.0-rc-9] - 2024-08-09

### Fixed

- transaction rollback on context cancel for PGX adapters.

## [2.0.0-rc-8] - 2024-01-04

### Fixed

- gouroutine leak.

## [2.0.0-rc-7] - 2024-01-04

### Changes

- Transferred drivers in separated modules.

## [1.4.0] - 2023-09-01

### Added

- pgx v4 adapter
- pgx v5 adapter

## [1.3.0] - 2023-06-16

### Added

- Ability to skip rollback due to an error

### Other

- Bumped library versions

## [1.2.2] - 2023-05-20

### Other

- Bumped library versions


## [1.2.1] - 2023-03-28

### Other

- Bumped go.mongodb.org/mongo-driver from 1.11.2 to 1.11.3
- Fixed lint issues

## [1.2.0] - 2023-03-10

### Added

- Redis adapter

## [1.1.0] - 2022-11-22

### Added

- GORM adapter

## [1.0.0] - 2022-11-13

### Fixed

- Documentation.

### Added

- Benchmark comparing case with and without trm.

## [1.0.0-beta] - 2022-11-05

### Added

- badges.
- examples.
- test with real mongo db.

## [1.0.0-beta] - 2022-10-24

### Added

- timeout and canceling of transaction.

## [1.0.0-beta] - 2022-10-20

### Added

- chained transaction.

## [1.0.0-beta] - 2022-10-17

### Added

- propagation.
- settings.

## [0.0.0] - 2022-09-08

### Added

- manager interface.
- transaction interface.
- mongo, sql implementations.

[0.0.0]: https://github.com/avito-tech/go-transaction-manager/
