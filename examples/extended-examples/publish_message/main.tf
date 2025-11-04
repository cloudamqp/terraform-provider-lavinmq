locals {
  vhost_name = "/"
}

resource "lavinmq_queue" "publish_queue" {
  name        = "publish-queue"
  vhost       = local.vhost_name
  durable     = true
  auto_delete = false
}

resource "lavinmq_publish_message" "example_message" {
  vhost       = local.vhost_name
  exchange    = "amq.default"
  routing_key = lavinmq_queue.publish_queue.name
  payload     = "{\"message\": \"Hello, World!\"}"
  properties = {
    content_type = "application/json"
  }
  publish_message_counter = 1
}

data "lavinmq_queues" "all_queues" {
  vhost = local.vhost_name

  depends_on = [
    lavinmq_publish_message.example_message
  ]
}
