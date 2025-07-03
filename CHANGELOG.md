# Changelog

All notable changes to this project will be documented in this file.

## [0.1.6] - 2025-07-03
### Fixed
- Fixed a bug where the `awsworkmail_organization` resource did not repopulate the `id` attribute in state after a refresh. The Read method now looks up the organization by alias and ensures the ID is always set, improving compatibility with outputs and dependent resources.

## [Unreleased]
- Work in progress

## [0.1.0] - 2025-07-01
### Added
- Initial release of the Terraform AWS WorkMail Provider.
- Support for managing AWS WorkMail Organizations, Users, and Groups.
- Import functionality for all real resources.
- Domain resource stub with clear documentation and warnings (manual step required).
- English-only, professional documentation and code comments.
- Complete meta files: CODEOWNERS, CONTRIBUTING, SECURITY, PR template, and CHANGELOG.
- CI/CD pipeline with GoReleaser and GPG signing for releases.
- Acceptance and unit test coverage with credential pre-checks.

### Notes
- This is the first public release, ready for HashiCorp Registry publication.
- Please see documentation for usage, import, and limitations.
