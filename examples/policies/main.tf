resource "lavinmq_vhost" "test_vhost" {
  name = "test_vhost"
}

resource "lavinmq_policy" "test_policy" {
  name    = "ha-all"
  vhost   = lavinmq_vhost.test_vhost.name
  pattern = ".*"
  definition = {
    "ha-mode"          = "all"
    "ha-sync-mode"     = "automatic"
    "max-length"       = 1000
    "max-length-bytes" = 10485760
  }
  priority = 1
  apply_to = "queues"
}
