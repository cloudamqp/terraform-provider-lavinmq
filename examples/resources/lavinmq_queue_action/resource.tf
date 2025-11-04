resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_queue_action" "purge_example" {
  name   = "example-queue"
  vhost  = lavinmq_vhost.example.name
  action = "purge"
}
