# LavinMQ Federation Upstream Set Examples
# Federation upstream sets group multiple upstreams for HA and failover scenarios

# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create multiple federation upstreams for different regions
resource "lavinmq_federation_upstream" "us_east" {
  name     = "us-east-upstream"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@us-east.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "orders"
}

resource "lavinmq_federation_upstream" "us_west" {
  name     = "us-west-upstream"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@us-west.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "orders"
}

resource "lavinmq_federation_upstream" "eu_central" {
  name     = "eu-central-upstream"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@eu-central.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "orders"
}

# Example 1: Basic HA set with multiple US regions
resource "lavinmq_federation_upstream_set" "us_ha" {
  name  = "us-ha-set"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.us_east.name,
    lavinmq_federation_upstream.us_west.name
  ]
}

# Example 2: Global HA set with all regions
resource "lavinmq_federation_upstream_set" "global_ha" {
  name  = "global-ha-set"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.us_east.name,
    lavinmq_federation_upstream.us_west.name,
    lavinmq_federation_upstream.eu_central.name
  ]
}

# Example 3: Create policy that uses the upstream set
resource "lavinmq_policy" "federate_orders" {
  name     = "federate-orders"
  vhost    = lavinmq_vhost.example.name
  pattern  = "^orders\\."
  priority = 10
  arguments = {
    federation-upstream-set = lavinmq_federation_upstream_set.global_ha.name
  }
  apply_to = "exchanges"
}

# Example 4: Queue federation with upstream set
resource "lavinmq_federation_upstream" "queue_upstream_1" {
  name     = "queue-upstream-1"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@queue1.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  queue    = "tasks"
  ack_mode = "on-confirm"
}

resource "lavinmq_federation_upstream" "queue_upstream_2" {
  name     = "queue-upstream-2"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@queue2.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  queue    = "tasks"
  ack_mode = "on-confirm"
}

resource "lavinmq_federation_upstream_set" "queue_ha" {
  name  = "queue-ha-set"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.queue_upstream_1.name,
    lavinmq_federation_upstream.queue_upstream_2.name
  ]
}

resource "lavinmq_policy" "federate_tasks" {
  name     = "federate-tasks"
  vhost    = lavinmq_vhost.example.name
  pattern  = "^tasks$"
  priority = 10
  arguments = {
    federation-upstream-set = lavinmq_federation_upstream_set.queue_ha.name
  }
  apply_to = "queues"
}

# Example 5: Multi-tier HA configuration
resource "lavinmq_federation_upstream" "primary" {
  name            = "primary-upstream"
  vhost           = lavinmq_vhost.example.name
  uri             = "amqp://guest:guest@primary.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange        = "events"
  reconnect_delay = 5
}

resource "lavinmq_federation_upstream" "secondary" {
  name            = "secondary-upstream"
  vhost           = lavinmq_vhost.example.name
  uri             = "amqp://guest:guest@secondary.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange        = "events"
  reconnect_delay = 10
}

resource "lavinmq_federation_upstream" "tertiary" {
  name            = "tertiary-upstream"
  vhost           = lavinmq_vhost.example.name
  uri             = "amqp://guest:guest@tertiary.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange        = "events"
  reconnect_delay = 15
}

resource "lavinmq_federation_upstream_set" "tiered_ha" {
  name  = "tiered-ha-set"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.primary.name,
    lavinmq_federation_upstream.secondary.name,
    lavinmq_federation_upstream.tertiary.name
  ]
}

# Output the upstream set names
output "us_ha_set" {
  value = {
    name      = lavinmq_federation_upstream_set.us_ha.name
    vhost     = lavinmq_federation_upstream_set.us_ha.vhost
    upstreams = lavinmq_federation_upstream_set.us_ha.upstreams
  }
}

output "global_ha_set" {
  value = {
    name      = lavinmq_federation_upstream_set.global_ha.name
    vhost     = lavinmq_federation_upstream_set.global_ha.vhost
    upstreams = lavinmq_federation_upstream_set.global_ha.upstreams
  }
}

output "tiered_ha_set" {
  value = {
    name      = lavinmq_federation_upstream_set.tiered_ha.name
    vhost     = lavinmq_federation_upstream_set.tiered_ha.vhost
    upstreams = lavinmq_federation_upstream_set.tiered_ha.upstreams
  }
}
