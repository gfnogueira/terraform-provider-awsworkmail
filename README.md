# Terraform Provider AWS WorkMail

This provider allows you to manage AWS WorkMail resources via Terraform.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Installation

Add to your Terraform code:

```hcl
terraform {
  required_providers {
    awsworkmail = {
      source  = "gfnogueira/awsworkmail"
      version = ">= 0.1.0"
    }
  }
}
```

## Example Usage

```hcl
provider "awsworkmail" {}

resource "awsworkmail_organization" "example" {
  alias = "my-workmail-org"
}
```

## Features

- Manage AWS WorkMail Organizations
- Manage AWS WorkMail Users
- Manage AWS WorkMail Groups
- Domain resource stub (documented, not implemented in AWS SDK v2)

## Supported Resources

- `awsworkmail_organization`: Manages an AWS WorkMail organization
- `awsworkmail_user`: Manages a user in an AWS WorkMail organization
- `awsworkmail_group`: Manages a group in an AWS WorkMail organization
- `awsworkmail_domain`: Stub resource for documentation (manual step required)

## Development

To build the provider:

```sh
go install
```

To run acceptance tests:

```sh
make testacc
```

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines on how to contribute to this project.

## License

MPL-2.0
