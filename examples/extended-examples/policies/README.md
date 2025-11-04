# LavinMQ Policy Resource

This example demonstrates how to use the `lavinmq_policy` resource to manage policies in LavinMQ.

## Policy Resource

The `lavinmq_policy` resource allows you to create and manage policies that define configuration rules for exchanges and queues matching specific patterns.

### Basic Usage

```terraform
resource "lavinmq_policy" "example" {
  name    = "my-policy"
  vhost   = "/"
  pattern = "^example-"

  definition = {
    "message-ttl" = "3600000"  # 1 hour in milliseconds
  }

  priority = 1
  apply_to = "queues"
}
```

### Arguments

- `name` - (Required) The name of the policy.
- `vhost` - (Required) The virtual host where the policy applies.
- `pattern` - (Required) Regular expression pattern matching queue/exchange names.
- `definition` - (Required) Map of policy definition key-value pairs.
- `priority` - (Optional) Policy priority. Higher numbers = higher priority. Default: 0.
- `apply_to` - (Optional) What the policy applies to: "all", "exchanges", or "queues". Default: "all".

### Policy Definition Keys

Common policy definition keys include:

#### Message TTL and Expiry
- `message-ttl` - Message time-to-live in milliseconds
- `expires` - Queue expiry time in milliseconds

#### Queue Length Limits
- `max-length` - Maximum number of messages in a queue
- `max-length-bytes` - Maximum queue size in bytes

#### Dead Letter Configuration
- `dead-letter-exchange` - Dead letter exchange name
- `dead-letter-routing-key` - Dead letter routing key

#### High Availability (if supported)
- `ha-mode` - High availability mode ("all", "exactly", "nodes")
- `ha-params` - High availability parameters
- `ha-sync-mode` - Synchronization mode ("automatic", "manual")

### Import

Policies can be imported using the format `vhost@policy_name`:

```bash
terraform import lavinmq_policy.example "/@my-policy"
```

For non-default vhosts:

```bash
terraform import lavinmq_policy.example "my-vhost@my-policy"
```

## Testing

The policy resource includes comprehensive tests covering:

- Basic policy creation and updates
- Policy definition changes
- Priority modifications
- Apply-to field changes
- Dead letter exchange configuration
- Queue length limits
- Import functionality

Run tests with:

```bash
TF_ACC=1 go test ./lavinmq -run TestAccPolicy -v
```
