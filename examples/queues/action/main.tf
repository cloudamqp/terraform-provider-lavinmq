# LavinMQ Queue Action Example
# This example demonstrates how to purge a queue

# Create a test vhost for our examples
resource "lavinmq_vhost" "test" {
  name = "test"
}

# Create a basic durable queue with default settings
resource "lavinmq_queue" "basic_example" {
  name  = "basic-queue"
  vhost = lavinmq_vhost.test.name
}

# Purge the basic queue
resource "lavinmq_queue_action" "purge_example" {
  name   = lavinmq_queue.basic_example.name
  vhost  = lavinmq_vhost.test.name
  action = "purge"
}
