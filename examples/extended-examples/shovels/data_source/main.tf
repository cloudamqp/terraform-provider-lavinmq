# LavinMQ Shovels Data Source Example
# This example demonstrates how to use the lavinmq_shovels data source

# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create multiple shovels for demonstration
resource "lavinmq_shovel" "shovel1" {
  name     = "example-q2q-shovel"
  vhost    = lavinmq_vhost.example.name
  src_uri  = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.example.name)}"
  dest_uri = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.example.name)}"

  src_queue  = "source-queue-1"
  dest_queue = "destination-queue-1"
}

resource "lavinmq_shovel" "shovel2" {
  name     = "example-e2q-shovel"
  vhost    = lavinmq_vhost.example.name
  src_uri  = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.example.name)}"
  dest_uri = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.example.name)}"

  src_exchange     = "source-exchange"
  src_exchange_key = "events.#"
  dest_queue       = "destination-queue-2"

  src_prefetch_count = 500
  ack_mode           = "on-publish"
}

resource "lavinmq_shovel" "shovel3" {
  name     = "example-cross-vhost-shovel"
  vhost    = lavinmq_vhost.example.name
  src_uri  = "amqp://guest:guest@localhost:5672/${urlencode(lavinmq_vhost.example.name)}"
  dest_uri = "amqp://guest:guest@localhost:5672/%2f"

  src_queue  = "vhost-queue"
  dest_queue = "default-vhost-queue"
}

# List all shovels across all vhosts
data "lavinmq_shovels" "all_shovels" {
  depends_on = [
    lavinmq_shovel.shovel1,
    lavinmq_shovel.shovel2,
    lavinmq_shovel.shovel3
  ]
}

# List shovels for a specific vhost
data "lavinmq_shovels" "example_vhost_shovels" {
  vhost = lavinmq_vhost.example.name
  depends_on = [
    lavinmq_shovel.shovel1,
    lavinmq_shovel.shovel2,
    lavinmq_shovel.shovel3
  ]
}

# Output all shovels information
output "all_shovels_info" {
  value = {
    total_count = length(data.lavinmq_shovels.all_shovels.shovels)
    shovels = [
      for shovel in data.lavinmq_shovels.all_shovels.shovels : {
        name  = shovel.name
        vhost = shovel.vhost
        type = shovel.src_queue != "" ? (
          shovel.dest_queue != "" ? "queue-to-queue" : "queue-to-exchange"
          ) : (
          shovel.dest_queue != "" ? "exchange-to-queue" : "exchange-to-exchange"
        )
      }
    ]
  }
}

# Output vhost-specific shovels with detailed configuration
output "vhost_shovels_info" {
  value = {
    vhost       = data.lavinmq_shovels.example_vhost_shovels.vhost
    total_count = length(data.lavinmq_shovels.example_vhost_shovels.shovels)
    shovels = [
      for shovel in data.lavinmq_shovels.example_vhost_shovels.shovels : {
        name               = shovel.name
        src_queue          = shovel.src_queue
        src_exchange       = shovel.src_exchange
        dest_queue         = shovel.dest_queue
        dest_exchange      = shovel.dest_exchange
        src_prefetch_count = shovel.src_prefetch_count
        ack_mode           = shovel.ack_mode
        reconnect_delay    = shovel.reconnect_delay
      }
    ]
  }
}

# Filter shovels by configuration
output "high_prefetch_shovels" {
  description = "Shovels with prefetch count > 100"
  value = [
    for shovel in data.lavinmq_shovels.all_shovels.shovels :
    shovel.name if shovel.src_prefetch_count > 100
  ]
}

# Count shovels by type
output "shovels_by_type" {
  value = {
    queue_to_queue = length([
      for shovel in data.lavinmq_shovels.all_shovels.shovels :
      shovel if shovel.src_queue != "" && shovel.dest_queue != ""
    ])
    queue_to_exchange = length([
      for shovel in data.lavinmq_shovels.all_shovels.shovels :
      shovel if shovel.src_queue != "" && shovel.dest_exchange != ""
    ])
    exchange_to_queue = length([
      for shovel in data.lavinmq_shovels.all_shovels.shovels :
      shovel if shovel.src_exchange != "" && shovel.dest_queue != ""
    ])
    exchange_to_exchange = length([
      for shovel in data.lavinmq_shovels.all_shovels.shovels :
      shovel if shovel.src_exchange != "" && shovel.dest_exchange != ""
    ])
  }
}
