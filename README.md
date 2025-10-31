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

## VCR Testing

The provider can be tested with Terraform Acceptance Test together with [Go-VCR] package. When using
the Go-VCR package all HTTP interactions to the API backend can be recorded or replayed and used
while testing the provider.

### Record

To record VCR cassettes, you need a running LavinMQ instance (preferably installed locally). 

Set the following environment variables in a `.env` file:

```
LAVINMQ_API_BASEURL="http://localhost:15672/"
LAVINMQ_API_USERNAME="guest"
LAVINMQ_API_PASSWORD="guest"
# VCR-TEST VARIABLES
TEST_AMQP_URI="amqp://guest:guest@localhost:5672//"
```

**Record all tests:**

```sh
LAVINMQ_RECORD=1 TF_ACC=1 dotenv -f .env go test ./lavinmq/ -v -timeout 5s
```

**Record a single test:**

```sh
LAVINMQ_RECORD=1 TF_ACC=1 dotenv -f .env go test ./lavinmq/ -v -run {TestName} -timeout 5s
```

### Replay

**Replay single test:**

```sh
TF_ACC=1 go test ./lavinmq/ -v -run {TestName}
```

**Replay all tests:**

```sh
TF_ACC=1 go test ./lavinmq/ -v
```

[Go-VCR]: https://github.com/dnaeon/go-vcr