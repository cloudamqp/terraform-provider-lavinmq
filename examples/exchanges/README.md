# LavinMQ Exchange Examples

This directory contains examples of how to use the `lavinmq_exchange` resource to create and manage exchanges in LavinMQ.

## Prerequisites

1. LavinMQ server running locally on `http://127.0.0.1:15672`
2. Valid credentials (default: guest/guest)
3. Terraform installed

## Running the Examples

1. Initialize Terraform:
   ```bash
   terraform init
   ```

2. Plan the deployment:
   ```bash
   terraform plan
   ```

3. Apply the configuration:
   ```bash
   terraform apply
   ```

## What This Example Creates

- **Test Vhost**: `test` - A dedicated vhost for all example exchanges
- **Direct Exchange**: `direct-exchange` - Routes messages with exact routing key matches
- **Fanout Exchange**: `fanout-exchange` - Routes messages to all bound queues
- **Topic Exchange**: `topic-exchange` - Routes messages using routing key patterns
- **Headers Exchange**: `headers-exchange` - Routes messages based on message headers
- **Temporary Exchange**: `temporary-exchange` - Auto-delete exchange for short-lived use
- **Custom Vhost Exchange**: Creates an additional custom vhost and exchange within it

## Exchange Types

### Direct Exchange
Routes messages to queues where the routing key exactly matches the binding key.

### Fanout Exchange  
Routes messages to all queues bound to the exchange, ignoring routing keys.

### Topic Exchange
Routes messages to queues based on wildcard matching between routing key and binding pattern.

### Headers Exchange
Routes messages based on message header attributes instead of routing keys.

## Configuration Options

- `name`: Unique name for the exchange
- `vhost`: Virtual host (examples use "test" vhost)
- `type`: Exchange type (direct, fanout, topic, headers)
- `auto_delete`: Delete when no longer used (default: false)
- `durable`: Survive broker restarts (default: false)

## Cleanup

To remove all created resources:
```bash
terraform destroy
```