# LavinMQ Queues Data Source Example

This example demonstrates how to use the `lavinmq_queues` data source to list all queues in a specific vhost.

## Resources Created

- `lavinmq_vhost.example` - Creates a vhost named "example-vhost"
- `lavinmq_queue.queue1` - Creates a durable queue that doesn't auto-delete
- `lavinmq_queue.queue2` - Creates a non-durable queue that auto-deletes

## Data Source Usage

The `lavinmq_queues` data source lists all queues in the specified vhost:

```hcl
data "lavinmq_queues" "example_queues" {
  vhost = lavinmq_vhost.example.name
  depends_on = [lavinmq_queue.queue1, lavinmq_queue.queue2]
}
```

## Output

The example outputs information about all queues in the vhost, including:

- Queue name
- Durability setting
- Auto-delete setting

## Usage

1. Ensure LavinMQ is running on `localhost:15672`
2. Run `terraform init`
3. Run `terraform plan`
4. Run `terraform apply`

The output will show all queues found in the example vhost.
