# LavinMQ Federation Upstream Examples
# This example demonstrates how to configure federation upstreams
# for replicating exchanges and queues from remote brokers

# Basic Exchange Federation
# Federates an exchange from an upstream broker
resource "lavinmq_federation_upstream" "basic_exchange" {
  name     = "upstream-cluster"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange = "federated-exchange"
}

# Queue Federation
# Federates a queue from an upstream broker
resource "lavinmq_federation_upstream" "queue_federation" {
  name  = "upstream-queue-federation"
  vhost = "/"
  uri   = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  queue = "federated-queue"
}

# Federation with Custom Configuration
# Demonstrates all available configuration options
resource "lavinmq_federation_upstream" "custom_config" {
  name     = "custom-upstream"
  vhost    = "/"
  uri      = "amqp://user:password@upstream-broker.example.com:5672/%2f"
  exchange = "events-exchange"

  # Prefetch count controls message buffering
  prefetch_count = 500

  # Reconnection delay in seconds
  reconnect_delay = 10

  # Acknowledgement mode:
  # - "on-confirm": wait for destination confirmation (most reliable)
  # - "on-publish": acknowledge after publishing
  # - "no-ack": no acknowledgement (fastest, least reliable)
  ack_mode = "on-confirm"

  # Maximum federation hops to prevent loops
  max_hops = 3

  # Optional: Federation link expiry time in milliseconds
  expires = 3600000

  # Optional: Message TTL in milliseconds
  message_ttl = 600000
}

# High-Performance Federation
# Optimized for throughput
resource "lavinmq_federation_upstream" "high_performance" {
  name            = "fast-upstream"
  vhost           = "/"
  uri             = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange        = "high-volume-exchange"
  prefetch_count  = 5000
  ack_mode        = "no-ack"
  reconnect_delay = 5
}

# Reliable Federation
# Configured for maximum reliability
resource "lavinmq_federation_upstream" "reliable" {
  name            = "reliable-upstream"
  vhost           = "/"
  uri             = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange        = "critical-exchange"
  prefetch_count  = 10
  ack_mode        = "on-confirm"
  reconnect_delay = 30
  max_hops        = 1
}

# Cross-VHost Federation
# Federate from one vhost to another
resource "lavinmq_vhost" "source" {
  name = "source-vhost"
}

resource "lavinmq_vhost" "destination" {
  name = "destination-vhost"
}

resource "lavinmq_federation_upstream" "cross_vhost" {
  name     = "cross-vhost-upstream"
  vhost    = lavinmq_vhost.destination.name
  uri      = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.source.name)}"
  exchange = "shared-exchange"
}

# Multi-Hop Federation
# Allow messages to traverse multiple federation links
resource "lavinmq_federation_upstream" "multi_hop" {
  name     = "multi-hop-upstream"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange = "distributed-exchange"
  max_hops = 5
}

# Federation with Message TTL
# Ensure messages don't stay in the federation link indefinitely
resource "lavinmq_federation_upstream" "with_ttl" {
  name        = "ttl-upstream"
  vhost       = "/"
  uri         = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange    = "temp-exchange"
  message_ttl = 300000
  expires     = 3600000
}

# Federation with Consumer Tag
# Use a specific consumer tag for tracking
resource "lavinmq_federation_upstream" "with_tag" {
  name         = "tagged-upstream"
  vhost        = "/"
  uri          = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange     = "monitored-exchange"
  consumer_tag = "federation-monitor"
}

# Outputs
output "federation_info" {
  value = {
    basic_exchange = {
      name     = lavinmq_federation_upstream.basic_exchange.name
      exchange = lavinmq_federation_upstream.basic_exchange.exchange
    }
    custom_config = {
      name            = lavinmq_federation_upstream.custom_config.name
      prefetch_count  = lavinmq_federation_upstream.custom_config.prefetch_count
      max_hops        = lavinmq_federation_upstream.custom_config.max_hops
      reconnect_delay = lavinmq_federation_upstream.custom_config.reconnect_delay
    }
  }
}
