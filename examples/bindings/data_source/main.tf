# LavinMQ Bindings Data Source Examples
# This example demonstrates how to query bindings

# Query all bindings in the default vhost
data "lavinmq_bindings" "default_vhost_bindings" {
  vhost = "/"
}

# Create a test vhost and bindings
resource "lavinmq_vhost" "test" {
  name = "test"
}

resource "lavinmq_exchange" "test_exchange" {
  name        = "test-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "direct"
  durable     = true
  auto_delete = false
}

resource "lavinmq_queue" "test_queue" {
  name        = "test-queue"
  vhost       = lavinmq_vhost.test.name
  durable     = true
  auto_delete = false
}

resource "lavinmq_binding" "test_binding" {
  vhost            = lavinmq_vhost.test.name
  source           = lavinmq_exchange.test_exchange.name
  destination      = lavinmq_queue.test_queue.name
  destination_type = "queue"
  routing_key      = "test.key"
}

# Query all bindings in the test vhost
data "lavinmq_bindings" "test_vhost_bindings" {
  vhost      = lavinmq_vhost.test.name
  depends_on = [lavinmq_binding.test_binding]
}

# Output bindings information
output "default_vhost_bindings_count" {
  value       = length(data.lavinmq_bindings.default_vhost_bindings.bindings)
  description = "Number of bindings in the default vhost"
}

output "test_vhost_bindings" {
  value = [
    for binding in data.lavinmq_bindings.test_vhost_bindings.bindings : {
      source           = binding.source
      destination      = binding.destination
      destination_type = binding.destination_type
      routing_key      = binding.routing_key
      properties_key   = binding.properties_key
    }
  ]
  description = "All bindings in the test vhost"
}

output "test_binding_details" {
  value = {
    for binding in data.lavinmq_bindings.test_vhost_bindings.bindings :
    "${binding.source}->${binding.destination}" => {
      routing_key    = binding.routing_key
      properties_key = binding.properties_key
    }
  }
  description = "Detailed binding information mapped by source->destination"
}
