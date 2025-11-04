resource "lavinmq_queue_action" "purge_example" {
  name   = "example-queue"
  vhost  = "/"
  action = "purge"
}
