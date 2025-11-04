resource "lavinmq_shovel" "example" {
  name     = "example-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue  = "source-queue"
  dest_queue = "destination-queue"
}
