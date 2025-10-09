# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create some queues in the vhost
resource "lavinmq_queue" "queue1" {
  name        = "example-queue-1"
  vhost       = lavinmq_vhost.example.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_queue" "queue2" {
  name        = "example-queue-2"
  vhost       = lavinmq_vhost.example.name
  durable     = false
  auto_delete = true
}

# Use the data source to list all queues in the vhost
data "lavinmq_queues" "example_queues" {
  vhost      = lavinmq_vhost.example.name
  depends_on = [lavinmq_queue.queue1, lavinmq_queue.queue2]
}

# Output the queue information
output "queues_info" {
  value = {
    vhost = data.lavinmq_queues.example_queues.vhost
    queues = [
      for queue in data.lavinmq_queues.example_queues.queues : {
        name        = queue.name
        durable     = queue.durable
        auto_delete = queue.auto_delete
      }
    ]
  }
}
