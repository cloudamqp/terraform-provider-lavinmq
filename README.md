# Terraform provider for LavinMQ HTTP API

## Prerequisites

1. **Golang**: Install Golang, follow [Golang's installation guide](https://go.dev/doc/install)
2. **Terraform**: Install Terraform (>= 0.12). Follow [Terraform's installation guide](https://developer.hashicorp.com/terraform/downloads).

## Example usage

Create new user to gain access HTTP API, manegment interface and AMQP broker.

```hcl
provider lavinmq {
  baseurl = "<http-api-url>"
  username = "<username>"
  password = "<password>"
}


resource "lavinmq_user" "this" {
  name              = "<username>"
  password          = "<password>"
  tags              = ["administrator"]
}
```
