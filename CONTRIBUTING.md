# Contributing Guide

Thank you for considering contributing to the `terraform-provider-awsworkmail`!

## How to contribute

- Fork the repository
- Create a new branch (`git checkout -b feature/my-feature`)
- Make your changes
- Commit your changes (`git commit -am 'Add some feature'`)
- Push to the branch (`git push origin feature/my-feature`)
- Create a new Pull Request

## Code Style

We follow standard Go formatting and linting tools:
- `gofmt`
- `golangci-lint`

Run tests with:
```bash
make test
```

## Requirements

- Go >= 1.20
- Terraform Plugin Framework

## Reporting Bugs

Use the [bug report template](.github/ISSUE_TEMPLATE/bug_report.md).

## Suggesting Features

Use the [feature request template](.github/ISSUE_TEMPLATE/feature_request.md).

## Acceptance Tests and AWS Credentials

Acceptance tests (those that interact with real AWS resources) require valid AWS credentials. These tests will be skipped automatically if the following environment variables are not set:
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`

To run acceptance tests, export your credentials before running tests:

```bash
export AWS_ACCESS_KEY_ID=your-access-key-id
export AWS_SECRET_ACCESS_KEY=your-secret-access-key
make test
```

Unit tests and static analysis do not require AWS credentials.

## Release Process

For each new release:
1. Update the `CHANGELOG.md` with a summary of changes for the new version.
2. Commit your changes to the main branch.
3. Create and push a new semantic version tag (e.g., `v0.2.0`).
   ```sh
   git tag v0.2.0
   git push origin v0.2.0
   ```
4. The release workflow will run automatically, sign the release, and publish it to GitHub and the Terraform Registry.

## License

By contributing, you agree that your contributions will be licensed under the MPL-2.0 license.
