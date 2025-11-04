# LavinMQ Federation Upstreams Data Source Example
# This example demonstrates how to use the lavinmq_federation_upstreams data source

# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create multiple federation upstreams for demonstration
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "example-upstream-1"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream1.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name            = "example-upstream-2"
  vhost           = lavinmq_vhost.example.name
  uri             = "amqp://guest:guest@upstream2.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  queue           = "queue-1"
  prefetch_count  = 500
  ack_mode        = "on-publish"
  reconnect_delay = 10
}

resource "lavinmq_federation_upstream" "upstream3" {
  name     = "example-upstream-3"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream3.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-3"
  max_hops = 3
}

# List all federation upstreams across all vhosts
data "lavinmq_federation_upstreams" "all_upstreams" {
  depends_on = [
    lavinmq_federation_upstream.upstream1,
    lavinmq_federation_upstream.upstream2,
    lavinmq_federation_upstream.upstream3
  ]
}

# List federation upstreams for a specific vhost
data "lavinmq_federation_upstreams" "example_vhost_upstreams" {
  vhost = lavinmq_vhost.example.name
  depends_on = [
    lavinmq_federation_upstream.upstream1,
    lavinmq_federation_upstream.upstream2,
    lavinmq_federation_upstream.upstream3
  ]
}

# Output all federation upstreams information
output "all_upstreams_info" {
  value = {
    total_count = length(data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams)
    upstreams = [
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams : {
        name     = upstream.name
        vhost    = upstream.vhost
        exchange = upstream.exchange
        queue    = upstream.queue
        max_hops = upstream.max_hops
      }
    ]
  }
}

# Output vhost-specific federation upstreams with detailed configuration
output "vhost_upstreams_info" {
  value = {
    vhost       = data.lavinmq_federation_upstreams.example_vhost_upstreams.vhost
    total_count = length(data.lavinmq_federation_upstreams.example_vhost_upstreams.federation_upstreams)
    upstreams = [
      for upstream in data.lavinmq_federation_upstreams.example_vhost_upstreams.federation_upstreams : {
        name            = upstream.name
        uri             = upstream.uri
        exchange        = upstream.exchange
        queue           = upstream.queue
        prefetch_count  = upstream.prefetch_count
        ack_mode        = upstream.ack_mode
        reconnect_delay = upstream.reconnect_delay
        max_hops        = upstream.max_hops
      }
    ]
  }
}

# Filter upstreams by configuration
output "high_prefetch_upstreams" {
  description = "Federation upstreams with prefetch count > 100"
  value = [
    for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
    upstream.name if upstream.prefetch_count > 100
  ]
}

# Count upstreams by type
output "upstreams_by_type" {
  value = {
    exchange_federation = length([
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
      upstream if upstream.exchange != ""
    ])
    queue_federation = length([
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
      upstream if upstream.queue != ""
    ])
  }
}

# Group upstreams by ack mode
output "upstreams_by_ack_mode" {
  value = {
    on_confirm = [
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
      upstream.name if upstream.ack_mode == "on-confirm"
    ]
    on_publish = [
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
      upstream.name if upstream.ack_mode == "on-publish"
    ]
    no_ack = [
      for upstream in data.lavinmq_federation_upstreams.all_upstreams.federation_upstreams :
      upstream.name if upstream.ack_mode == "no-ack"
    ]
  }
}
