# Queue Examples

This directory contains examples for managing queues in LavinMQ.

## Usage

```bash
terraform init
terraform plan
terraform apply
```

## Examples

### Basic Queue

Creates a queue with default settings:

```hcl
resource "lavinmq_queue" "basic_example" {
  name  = "basic-queue"
  vhost = lavinmq_vhost.test.name
}
```

### Durable Queue

Creates a persistent queue that survives broker restarts:

```hcl
resource "lavinmq_queue" "durable_example" {
  name        = "durable-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}
```

### Temporary Queue

Creates a queue that is automatically deleted when no longer in use:

```hcl
resource "lavinmq_queue" "temp_example" {
  name        = "temporary-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = false
  auto_delete = true
}
```

### Queue with Arguments

Creates a queue with AMQP arguments for message TTL, max length, and dead letter routing:

```hcl
resource "lavinmq_queue" "queue_with_args" {
  name        = "ttl-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
  
  arguments = {
    "x-message-ttl"            = 60000
    "x-max-length"             = 1000
    "x-dead-letter-exchange"   = "dlx-exchange"
    "x-dead-letter-routing-key" = "dead.letters"
  }
}
```

Common queue arguments:
- `x-message-ttl`: Message time-to-live in milliseconds (integer)
- `x-max-length`: Maximum number of messages in queue (integer)
- `x-max-length-bytes`: Maximum queue size in bytes (integer)
- `x-dead-letter-exchange`: Exchange for expired/rejected messages (string)
- `x-dead-letter-routing-key`: Routing key for dead lettered messages (string)
- `x-expires`: Queue expiration time in milliseconds (integer)

## Cleanup

```bash
terraform destroy
```
