resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_queue" "example" {
  name        = "example-queue"
  vhost       = lavinmq_vhost.example.name
  durable     = true
  auto_delete = false
  arguments = {
    "x-message-ttl" = 60000
  }
}
