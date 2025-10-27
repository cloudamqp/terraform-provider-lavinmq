resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_exchange" "example" {
  name        = "example-exchange"
  vhost       = lavinmq_vhost.example.name
  type        = "topic"
  durable     = true
  auto_delete = false
}
