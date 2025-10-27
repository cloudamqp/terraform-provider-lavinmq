# Terraform provider for LavinMQ HTTP API

## Prerequisites

1. **Golang**: Install Golang, follow [Golang's installation guide](https://go.dev/doc/install)
2. **Terraform**: Install Terraform (>= 0.12). Follow [Terraform's installation guide](https://developer.hashicorp.com/terraform/downloads).

## Example usage

Create new user to gain access HTTP API, management interface and AMQP broker.

```hcl
provider lavinmq {
  baseurl = "<http-api-url>"
  username = "<username>"
  password = "<password>"
}

resource "lavinmq_user" "this" {
  name     = "<username>"
  password = "<password>"
  tags     = ["administrator"]
}

resource "lavinmq_permission" "this" {
  vhost     = "/"
  user      = lavinmq_user.this.name
  configure = ".*"
  read      = ".*"
  write     = ".*"
}
```

## Resources

- `lavinmq_user` - Manage users
- `lavinmq_vhost` - Manage virtual hosts
- `lavinmq_queue` - Manage queues
- `lavinmq_exchange` - Manage exchanges
- `lavinmq_binding` - Manage bindings between exchanges and queues/exchanges
- `lavinmq_policy` - Manage policies
- `lavinmq_permission` - Manage user permissions on vhosts

## Data Sources

- `lavinmq_users` - List all users
- `lavinmq_vhosts` - List all vhosts
- `lavinmq_queues` - List all queues
- `lavinmq_exchanges` - List all exchanges
- `lavinmq_bindings` - List all bindings
- `lavinmq_policies` - List all policies
- `lavinmq_permissions` - List all permissions

## Documentation

Documentation is automatically generated using [terraform-plugin-docs](https://github.com/hashicorp/terraform-plugin-docs).

To generate or update documentation:

```bash
make docs
```

This will:
- Format example configurations
- Generate resource and data source documentation from schemas
- Update all files in the `docs/` directory
