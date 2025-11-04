# LavinMQ Exchanges Data Source Example

This example demonstrates how to use the `lavinmq_exchanges` data source to list all exchanges.

## Resources Created

- `lavinmq_vhost.example` - Creates a vhost named "example-vhost"
- `lavinmq_exchange.exchange1` - Creates a durable direct exchange
- `lavinmq_exchange.exchange2` - Creates a non-durable fanout exchange that auto-deletes

## Data Source Usage

The `lavinmq_exchanges` data source lists all exchanges across all vhosts:

```hcl
data "lavinmq_exchanges" "all_exchanges" {
  depends_on = [lavinmq_exchange.exchange1, lavinmq_exchange.exchange2]
}
```

## Output

The example outputs information about all exchanges, including:

- Exchange name
- Virtual host
- Exchange type
- Durability setting
- Auto-delete setting

## Usage

1. Ensure LavinMQ is running on `localhost:15672`
2. Run `terraform init`
3. Run `terraform plan`
4. Run `terraform apply`

The output will show all exchanges found in the system.
