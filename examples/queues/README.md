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

## Cleanup

```bash
terraform destroy
```
