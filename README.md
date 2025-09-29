# Terraform Provider AWS WorkMail
📦 Available at: [Terraform Registry](https://registry.terraform.io/providers/gfnogueira/awsworkmail/latest)

🔧 **Terraform provider to manage AWS WorkMail**, email users, and groups — with optional Route53 DNS setup.  
🔎 _Keywords: `aws`, `workmail`, `email`, `dns`, `route53`, `messaging`, `identity`_

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

### Basic Usage
```hcl
provider "awsworkmail" {
  region = "us-east-1"
}

resource "awsworkmail_organization" "example" {
  alias = "my-workmail-org"
}

# Register domain and get MX records
resource "awsworkmail_domain" "example" {
  organization_id = awsworkmail_organization.example.id
  domain          = "mycompany.com"
}

# Create users
resource "awsworkmail_user" "john" {
  organization_id = awsworkmail_organization.example.id
  name           = "john.doe"
  display_name   = "John Doe"
  password       = "TempPassword123!"
  email          = "john.doe@mycompany.com"
}

# Output MX records for DNS configuration
output "mx_records" {
  value = awsworkmail_domain.example.mx_records
}
```

### Multi-Account Setup with Assume Role
```hcl
provider "awsworkmail" {
  alias  = "target_account"
  region = "us-east-1"
  
  assume_role {
    role_arn     = "arn:aws:iam::123456789012:role/deployer-role"
    session_name = "terraform-workmail-session"
  }
}

resource "awsworkmail_organization" "cross_account" {
  provider = awsworkmail.target_account
  alias    = "my-cross-account-org"
}
```

## Features

- Manage AWS WorkMail Organizations, Users, and Groups
- Full Domain Management: Register, manage, and get MX records for WorkMail domains
- Multi-Account Support: Full assume_role functionality for cross-account deployments
- Import existing resources
- Comprehensive error handling and validation

## Supported Resources

- `awsworkmail_organization`: Manages an AWS WorkMail organization
- `awsworkmail_user`: Manages a user in an AWS WorkMail organization
- `awsworkmail_group`: Manages a group in an AWS WorkMail organization
- `awsworkmail_domain`: Manages domains in an AWS WorkMail organization

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
