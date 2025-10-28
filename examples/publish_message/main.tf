resource "lavinmq_vhost" "test" {
  name = "test"
}

resource "lavinmq_exchange" "topic_exchange" {
  name        = "topic-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "topic"
  auto_delete = false
  durable     = true
}

resource "lavinmq_queue" "notifications_queue" {
  name        = "notifications-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_binding" "notifications_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.topic_exchange.name
  destination      = lavinmq_queue.notifications_queue.name
  destination_type = "queue"
  routing_key      = "notification.*"
}

resource "lavinmq_publish_message" "example_message" {
  vhost       = lavinmq_vhost.test.name
  exchange    = lavinmq_exchange.topic_exchange.name
  routing_key = "notification.test"
  payload     = "11"

  depends_on = [
    lavinmq_binding.notifications_binding
  ]
}

data "lavinmq_queues" "all_queues" {
  vhost = lavinmq_vhost.test.name

  depends_on = [
    lavinmq_publish_message.example_message
  ]
}
