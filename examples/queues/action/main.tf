# LavinMQ Queue Action Example
# This example demonstrates how to purge a queue

# Create a test vhost for our examples
resource "lavinmq_vhost" "test" {
  name = "test"
}

resource "lavinmq_exchange" "topic_exchange" {
  name        = "topic-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "topic"
  auto_delete = false
  durable     = true
}

resource "lavinmq_queue" "example_queue" {
  name  = "example-queue"
  vhost = lavinmq_vhost.test.name
}

resource "lavinmq_binding" "example_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.topic_exchange.name
  destination      = lavinmq_queue.example_queue.name
  destination_type = "queue"
  routing_key      = "example.*"
}

resource "lavinmq_publish_message" "example_message" {
  vhost       = lavinmq_vhost.test.name
  exchange    = lavinmq_exchange.topic_exchange.name
  routing_key = "example.rk"
  payload     = "Hello, LavinMQ!"

  depends_on = [
    lavinmq_binding.example_binding
  ]
}

resource "lavinmq_queue_action" "purge_example" {
  name   = lavinmq_queue.example_queue.name
  vhost  = lavinmq_vhost.test.name
  action = "purge"

  depends_on = [
    lavinmq_publish_message.example_message
  ]
}
