resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_exchange" "direct_example" {
  name        = "direct-example-exchange"
  vhost       = lavinmq_vhost.example.name
  type        = "direct"
  durable     = true
  auto_delete = false
}
