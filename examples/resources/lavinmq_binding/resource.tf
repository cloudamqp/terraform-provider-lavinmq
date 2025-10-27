resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_exchange" "example" {
  name        = "example-exchange"
  vhost       = lavinmq_vhost.example.name
  type        = "direct"
  durable     = true
  auto_delete = false
}

resource "lavinmq_queue" "example" {
  name        = "example-queue"
  vhost       = lavinmq_vhost.example.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_binding" "example" {
  vhost            = lavinmq_vhost.example.name
  source           = lavinmq_exchange.example.name
  destination      = lavinmq_queue.example.name
  destination_type = "queue"
  routing_key      = "example.route"
}
