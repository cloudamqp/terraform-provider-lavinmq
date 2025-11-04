resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_policy" "example" {
  name     = "example-policy"
  vhost    = lavinmq_vhost.example.name
  pattern  = ".*"
  priority = 0
  apply_to = "queues"
  definition = {
    "max-length" = 1000
  }
}
