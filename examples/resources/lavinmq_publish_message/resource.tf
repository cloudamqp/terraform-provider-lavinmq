resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_queue" "example" {
  name        = "example-queue"
  vhost       = lavinmq_vhost.example.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_publish_message" "example_message" {
  vhost       = lavinmq_vhost.example.name
  exchange    = "amq.default"
  routing_key = lavinmq_queue.example.name
  payload     = "{\"message\": \"Hello, World!\"}"
  properties = {
    content_type = "application/json"
  }
}
