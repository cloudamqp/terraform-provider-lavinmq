# LavinMQ Binding Examples
# This example demonstrates how to create different types of bindings

# Create a test vhost for our examples
resource "lavinmq_vhost" "test" {
  name = "test"
}

# Create exchanges for binding
resource "lavinmq_exchange" "direct_exchange" {
  name        = "direct-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "direct"
  durable     = true
  auto_delete = false
}

resource "lavinmq_exchange" "topic_exchange" {
  name        = "topic-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "topic"
  durable     = true
  auto_delete = false
}

resource "lavinmq_exchange" "headers_exchange" {
  name        = "headers-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "headers"
  durable     = true
  auto_delete = false
}

# Create queues for binding
resource "lavinmq_queue" "orders_queue" {
  name        = "orders"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_queue" "notifications_queue" {
  name        = "notifications"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_queue" "logs_queue" {
  name        = "logs"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

# Basic binding with routing key
resource "lavinmq_binding" "orders_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.direct_exchange.name
  destination      = lavinmq_queue.orders_queue.name
  destination_type = "queue"
  routing_key      = "order.created"
}

# Topic exchange binding with wildcard routing key
resource "lavinmq_binding" "notifications_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.topic_exchange.name
  destination      = lavinmq_queue.notifications_queue.name
  destination_type = "queue"
  routing_key      = "notification.*"
}

# Multiple bindings to the same queue
resource "lavinmq_binding" "logs_binding_info" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.topic_exchange.name
  destination      = lavinmq_queue.logs_queue.name
  destination_type = "queue"
  routing_key      = "log.info"
}

resource "lavinmq_binding" "logs_binding_error" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.topic_exchange.name
  destination      = lavinmq_queue.logs_queue.name
  destination_type = "queue"
  routing_key      = "log.error"
}

# Headers exchange binding with arguments
resource "lavinmq_binding" "headers_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.headers_exchange.name
  destination      = lavinmq_queue.notifications_queue.name
  destination_type = "queue"

  arguments = {
    x-match  = "all"
    priority = "high"
    type     = "alert"
  }
}

# Exchange-to-exchange binding
resource "lavinmq_exchange" "backup_exchange" {
  name        = "backup-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "fanout"
  durable     = true
  auto_delete = false
}

resource "lavinmq_binding" "exchange_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.direct_exchange.name
  destination      = lavinmq_exchange.backup_exchange.name
  destination_type = "exchange"
  routing_key      = "backup.#"
}

# Outputs
output "bindings_info" {
  value = {
    orders = {
      source      = lavinmq_binding.orders_binding.source
      destination = lavinmq_binding.orders_binding.destination
      routing_key = lavinmq_binding.orders_binding.routing_key
    }
    notifications = {
      source      = lavinmq_binding.notifications_binding.source
      destination = lavinmq_binding.notifications_binding.destination
      routing_key = lavinmq_binding.notifications_binding.routing_key
    }
  }
}
