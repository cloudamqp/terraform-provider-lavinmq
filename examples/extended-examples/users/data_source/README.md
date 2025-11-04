# Users Data Source Examples

This directory contains examples for retrieving users from LavinMQ using the data source.

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Examples

### List All Users

Retrieves all users with their tags:

```hcl
data "lavinmq_users" "all" {}

output "all_users" {
  value = data.lavinmq_users.all.users
}
```

### Filter Users by Tag

After retrieving all users, you can filter them using Terraform expressions:

```hcl
data "lavinmq_users" "all" {}

locals {
  admin_users = [
    for user in data.lavinmq_users.all.users :
    user if contains(user.tags, "administrator")
  ]
}

output "admin_users" {
  value = local.admin_users
}
```

### Get User Names

Extract only the user names:

```hcl
data "lavinmq_users" "all" {}

output "user_names" {
  value = [for user in data.lavinmq_users.all.users : user.name]
}
```

## Cleanup

```bash
terraform destroy
```
