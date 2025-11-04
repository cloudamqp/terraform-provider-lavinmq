# LavinMQ Shovel Examples
# This example demonstrates how to configure shovels for message forwarding
# between queues and exchanges

# Basic Queue-to-Queue Shovel
# Moves messages from one queue to another on the same server
resource "lavinmq_shovel" "queue_to_queue" {
  name     = "q2q-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue  = "source-queue"
  dest_queue = "destination-queue"
}

# Queue-to-Exchange Shovel
# Reads from a queue and publishes to an exchange with a routing key
resource "lavinmq_shovel" "queue_to_exchange" {
  name     = "q2e-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue         = "source-queue"
  dest_exchange     = "target-exchange"
  dest_exchange_key = "routing.key"
}

# Exchange-to-Queue Shovel
# Subscribes to an exchange and forwards messages to a queue
resource "lavinmq_shovel" "exchange_to_queue" {
  name     = "e2q-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_exchange     = "source-exchange"
  src_exchange_key = "events.#"
  dest_queue       = "destination-queue"
}

# Exchange-to-Exchange Shovel
# Forwards messages from one exchange to another
resource "lavinmq_shovel" "exchange_to_exchange" {
  name     = "e2e-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_exchange      = "source-exchange"
  src_exchange_key  = "logs.*"
  dest_exchange     = "destination-exchange"
  dest_exchange_key = "forwarded.logs"
}

# Cross-VHost Shovel
# Moves messages between different virtual hosts on the same server
resource "lavinmq_vhost" "source_vhost" {
  name = "source-vhost"
}

resource "lavinmq_vhost" "dest_vhost" {
  name = "destination-vhost"
}

resource "lavinmq_shovel" "cross_vhost" {
  name     = "cross-vhost-shovel"
  vhost    = lavinmq_vhost.source_vhost.name
  src_uri  = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.source_vhost.name)}"
  dest_uri = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.dest_vhost.name)}"

  src_queue  = "source-queue"
  dest_queue = "destination-queue"
}

# Cross-Server Shovel
# Forwards messages from local server to a remote server
resource "lavinmq_shovel" "cross_server" {
  name     = "remote-backup-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://user:password@remote-server.example.com:5672/%2f"

  src_queue  = "local-queue"
  dest_queue = "remote-backup-queue"

  # Custom configuration for reliability
  src_prefetch_count = 500
  ack_mode           = "on-confirm"
  reconnect_delay    = 10
}

# Shovel with Custom Configuration
# Demonstrates all available configuration options
resource "lavinmq_shovel" "custom_config" {
  name     = "custom-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue  = "source-queue"
  dest_queue = "destination-queue"

  # Prefetch count controls how many messages to buffer
  src_prefetch_count = 1000

  # Delete mode: "never" (default) or "queue-length"
  src_delete_after = "never"

  # Reconnect delay in seconds
  reconnect_delay = 5

  # Acknowledgement mode:
  # - "on-confirm": acknowledge after destination confirms (most reliable)
  # - "on-publish": acknowledge after publishing to destination
  # - "no-ack": no acknowledgement (fastest, least reliable)
  ack_mode = "on-confirm"
}

# High-Performance Shovel
# Optimized for throughput with no-ack mode
resource "lavinmq_shovel" "high_performance" {
  name     = "fast-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue          = "high-volume-queue"
  dest_queue         = "processed-queue"
  src_prefetch_count = 5000
  ack_mode           = "no-ack"
}

# Reliable Shovel
# Configured for maximum reliability at the cost of some performance
resource "lavinmq_shovel" "reliable" {
  name     = "reliable-shovel"
  vhost    = "/"
  src_uri  = "amqp://guest:guest@localhost:5672/%2f"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue          = "critical-queue"
  dest_queue         = "processed-critical-queue"
  src_prefetch_count = 10
  ack_mode           = "on-confirm"
  reconnect_delay    = 30
}

# Outputs
output "shovel_info" {
  value = {
    queue_to_queue = {
      name       = lavinmq_shovel.queue_to_queue.name
      src_queue  = lavinmq_shovel.queue_to_queue.src_queue
      dest_queue = lavinmq_shovel.queue_to_queue.dest_queue
    }
    cross_server = {
      name            = lavinmq_shovel.cross_server.name
      prefetch_count  = lavinmq_shovel.cross_server.src_prefetch_count
      reconnect_delay = lavinmq_shovel.cross_server.reconnect_delay
    }
  }
}
