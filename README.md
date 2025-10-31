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

- `lavinmq_binding` - Manage bindings between exchanges and queues/exchanges
- `lavinmq_exchange` - Manage exchanges
- `lavinmq_permission` - Manage user permissions on vhosts
- `lavinmq_policy` - Manage policies
- `lavinmq_publish_message` - Publish messages to an exchange
- `lavinmq_queue` - Manage queues
- `lavinmq_queue_action` - Perform actions on queues (pause/resume/purge)
- `lavinmq_shovel` - Manage shovels
- `lavinmq_user` - Manage users
- `lavinmq_vhost` - Manage virtual hosts

## Data Sources

- `lavinmq_bindings` - List all bindings
- `lavinmq_exchanges` - List all exchanges
- `lavinmq_permissions` - List all permissions
- `lavinmq_policies` - List all policies
- `lavinmq_queues` - List all queues
- `lavinmq_shovels` - List all shovels
- `lavinmq_users` - List all users
- `lavinmq_vhosts` - List all vhosts
