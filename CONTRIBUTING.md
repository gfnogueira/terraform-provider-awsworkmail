# Contributing to terraform-provider-awsworkmail

Thank you for your interest in contributing to `terraform-provider-awsworkmail`!  
We welcome community contributions to improve the project.

## How to contribute

1. Fork this repository
2. Create a new branch (`git checkout -b feature/my-feature`)
3. Make your changes
4. Commit your changes (`git commit -am 'Add some feature'`)
5. Push the branch (`git push origin feature/my-feature`)
6. Open a Pull Request

## Code Style

We follow standard Go formatting and linting:

- `gofmt`
- [`golangci-lint`](https://golangci-lint.run)

### Run tests

```bash
make test
```
---

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

To release a new version:

1. Update the `CHANGELOG.md` with a summary of changes for the new version.
2. Commit your changes to the main branch.
3. Create and push a new semantic version tag (e.g., `v0.2.0`).
   ```sh
   git tag v0.2.0
   git push origin v0.2.0
   ```
4. The release workflow (CI) will:   
    - Build and sign the binaries  
    - Publish to GitHub Releases and the [Terraform Registry](https://registry.terraform.io/)    

## License

By contributing, you agree that your contributions will be licensed under the MPL-2.0 license.
