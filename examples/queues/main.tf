# LavinMQ Queue Examples
# This example demonstrates how to create different types of queues

# Create a test vhost for our examples
resource "lavinmq_vhost" "test" {
  name = "test"
}

# Create a basic durable queue with default settings
resource "lavinmq_queue" "basic_example" {
  name  = "basic-queue"
  vhost = lavinmq_vhost.test.name
}

# Create a durable queue
resource "lavinmq_queue" "durable_example" {
  name        = "durable-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

# Create a temporary queue that auto-deletes
resource "lavinmq_queue" "temp_example" {
  name        = "temporary-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = false
  auto_delete = true
}

# Create a queue on a custom vhost
resource "lavinmq_vhost" "custom" {
  name = "my-app"
}

resource "lavinmq_queue" "custom_vhost_example" {
  name        = "app-queue"
  vhost       = lavinmq_vhost.custom.name
  durable     = true
  auto_delete = false
}

# Outputs
output "basic_queue_id" {
  value = lavinmq_queue.basic_example.id
}

output "custom_queue_info" {
  value = {
    name        = lavinmq_queue.custom_vhost_example.name
    vhost       = lavinmq_queue.custom_vhost_example.vhost
    durable     = lavinmq_queue.custom_vhost_example.durable
    auto_delete = lavinmq_queue.custom_vhost_example.auto_delete
  }
}
