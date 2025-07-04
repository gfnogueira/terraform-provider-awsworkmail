# Changelog

All notable changes to this project will be documented in this file.

## [0.1.9] - 2025-07-03
### Changed
- Changed the `members` attribute of `awsworkmail_group` from a List to a Set to prevent configuration drift due to member order.
- Provider now ensures group members are always unique and sorted before saving to state, eliminating repeated updates when the order changes in configuration or AWS response.
- Updated documentation to reflect that `members` is now a Set.

### Fixed
- Resolved persistent drift when managing group membership, regardless of the order in the Terraform configuration.

## [0.1.8] - 2025-07-03
### Added
- `enabled` attribute for `awsworkmail_user` and `awsworkmail_group` resources, supporting enable/disable via AWS API.
- `first_name` and `last_name` attributes for `awsworkmail_user` resource.
- Improved documentation and usage examples for all new attributes.

### Changed
- Import and refresh for `awsworkmail_organization` now look up by OrganizationId, supporting correct import and state refresh.
- Provider now uses robust ImportState logic for all resources, including composite ID parsing.
- Provider never overwrites resource IDs with empty values and always falls back to state if needed.

### Fixed
- Improved error handling for AWS WorkMail `EntityStateException`, with clear diagnostics for unsupported state transitions.
- Cleaned up all unused imports and fixed linter warnings.
- Removed redundant `ignore_changes` lifecycle warnings in documentation and examples.
- All code comments, logs, and documentation are now in English for global use.

## [0.1.7] - 2025-07-03
### Fixed
- Improved import support for `awsworkmail_user`, `awsworkmail_group`, and `awsworkmail_domain` resources. Import now requires a composite ID in the format `<organization_id>,<resource_id>`, ensuring both required attributes are set in state. This fixes issues with importing resources that require both organization and resource IDs.

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
