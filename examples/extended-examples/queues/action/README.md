# Queue Action Example

This example demonstrates how to perform actions on LavinMQ queues using the `lavinmq_queue_action` resource.

## What This Example Does

1. **Creates a test vhost** - Sets up a virtual host named "test" for organizing resources
2. **Creates a basic queue** - Sets up a durable queue named "basic-queue"
3. **Purges the queue** - Removes all messages from the queue using the queue action resource

## Queue Actions Available

- **purge**: Removes all messages from the queue without deleting the queue itself

## Usage

1. **Initialize Terraform**:

   ```bash
   terraform init
   ```

2. **Plan the deployment**:

   ```bash
   terraform plan
   ```

3. **Apply the configuration**:

   ```bash
   terraform apply
   ```

4. **View the results**:
   After applying, the queue will be created and then purged of any messages.

## Understanding Queue Actions

The `lavinmq_queue_action` resource is designed for one-time operations on queues:

- **Action-based**: Performs a specific action when the resource is created
- **Requires replacement**: Any change to the action or target queue creates a new resource
- **No server state**: The action itself isn't stored on the server, only in Terraform state

## Resource Behavior

```hcl
resource "lavinmq_queue_action" "purge_example" {
  name   = lavinmq_queue.basic_example.name  # Target queue name
  vhost  = lavinmq_vhost.test.name           # Virtual host containing the queue
  action = "purge"                           # Action to perform
}
```

### Key Properties:

- **name**: The name of the queue to act upon
- **vhost**: The virtual host containing the target queue
- **action**: The action to perform (currently only "purge" is supported)

## When to Use Queue Actions

Queue actions are useful for:

- **Cleanup operations**: Purging messages during deployments
- **Testing scenarios**: Clearing queues between test runs
- **Maintenance tasks**: Automated queue management

## Important Notes

1. **Destructive operation**: Purging removes ALL messages from the queue permanently
2. **Queue must exist**: The target queue must exist before the action can be performed
3. **One-time execution**: Actions are performed when the resource is created
4. **Replacement on changes**: Modifying any attribute recreates the resource and re-executes the action

## Clean Up

To destroy the created resources:

```bash
terraform destroy
```

This will remove the queue action from state, delete the queue, and delete the vhost.

## Example Variations

### Purge Multiple Queues

```hcl
resource "lavinmq_queue_action" "purge_queue1" {
  name   = "queue-1"
  vhost  = "/"
  action = "purge"
}

resource "lavinmq_queue_action" "purge_queue2" {
  name   = "queue-2" 
  vhost  = "/"
  action = "purge"
}
```

### Conditional Purging

```hcl
resource "lavinmq_queue_action" "conditional_purge" {
  count  = var.purge_on_deploy ? 1 : 0
  name   = lavinmq_queue.example.name
  vhost  = lavinmq_queue.example.vhost
  action = "purge"
}
```

## Error Handling

The queue action resource includes built-in error handling:

- Verifies the target queue exists before performing the action
- Handles resource drift (queue deleted externally)
- Provides warning messages if queue not found

If the target queue doesn't exist, the action will no-op and provider warning message.
