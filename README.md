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

## Available Resources

- `awsworkmail_organization`: Manages an AWS WorkMail organization

## Development

To build the provider:

```sh
go install
```

To run acceptance tests:

```sh
make testacc
```

## License

MPL-2.0
