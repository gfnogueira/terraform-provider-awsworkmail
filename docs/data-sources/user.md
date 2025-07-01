# awsworkmail_user Data Source

Provides details about an AWS WorkMail user.

## Example Usage

```hcl
data "awsworkmail_user" "example" {
  organization_id = "m-1234567890"
  user_id        = "S-1-1-12-1234567890-1234567890-1234567890-1234"
}

output "user_name" {
  value = data.awsworkmail_user.example.name
}
output "user_email" {
  value = data.awsworkmail_user.example.email
}
output "user_state" {
  value = data.awsworkmail_user.example.state
}
```

## Argument Reference

- `organization_id` (Required) - The WorkMail Organization ID.
- `user_id` (Optional) - The WorkMail User ID. If not provided, the data source will return an error.

## Attributes Reference

- `name` - The name of the user.
- `email` - The primary email of the user.
- `state` - The state of the user.

## Import

This data source does not support import as it is read-only.

## Limitations

- The `user_id` must be known in advance. Listing all users is not currently supported.
- The organization must exist and the user must be active in WorkMail.
