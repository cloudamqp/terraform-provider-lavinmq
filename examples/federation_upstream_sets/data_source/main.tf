# LavinMQ Federation Upstream Sets Data Source Example
# This example demonstrates how to use the lavinmq_federation_upstream_sets data source

# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create several federation upstreams
resource "lavinmq_federation_upstream" "upstream1" {
  name     = "example-upstream-1"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream1.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-1"
}

resource "lavinmq_federation_upstream" "upstream2" {
  name     = "example-upstream-2"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream2.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-2"
}

resource "lavinmq_federation_upstream" "upstream3" {
  name     = "example-upstream-3"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream3.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-3"
}

resource "lavinmq_federation_upstream" "upstream4" {
  name     = "example-upstream-4"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream4.example.com:5672/${urlencode(lavinmq_vhost.example.name)}"
  exchange = "exchange-4"
}

# Create multiple federation upstream sets
resource "lavinmq_federation_upstream_set" "ha_set_1" {
  name  = "ha-set-1"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name
  ]
}

resource "lavinmq_federation_upstream_set" "ha_set_2" {
  name  = "ha-set-2"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.upstream3.name,
    lavinmq_federation_upstream.upstream4.name
  ]
}

resource "lavinmq_federation_upstream_set" "ha_set_all" {
  name  = "ha-set-all"
  vhost = lavinmq_vhost.example.name
  upstreams = [
    lavinmq_federation_upstream.upstream1.name,
    lavinmq_federation_upstream.upstream2.name,
    lavinmq_federation_upstream.upstream3.name,
    lavinmq_federation_upstream.upstream4.name
  ]
}

# List all federation upstream sets across all vhosts
data "lavinmq_federation_upstream_sets" "all_sets" {
  depends_on = [
    lavinmq_federation_upstream_set.ha_set_1,
    lavinmq_federation_upstream_set.ha_set_2,
    lavinmq_federation_upstream_set.ha_set_all
  ]
}

# List federation upstream sets for a specific vhost
data "lavinmq_federation_upstream_sets" "example_vhost_sets" {
  vhost = lavinmq_vhost.example.name
  depends_on = [
    lavinmq_federation_upstream_set.ha_set_1,
    lavinmq_federation_upstream_set.ha_set_2,
    lavinmq_federation_upstream_set.ha_set_all
  ]
}

# Output all federation upstream sets information
output "all_sets_info" {
  value = {
    total_count = length(data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets)
    sets = [
      for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets : {
        name           = set.name
        vhost          = set.vhost
        upstream_count = length(set.upstreams)
        upstream_names = set.upstreams
      }
    ]
  }
}

# Output vhost-specific federation upstream sets
output "vhost_sets_info" {
  value = {
    vhost       = data.lavinmq_federation_upstream_sets.example_vhost_sets.vhost
    total_count = length(data.lavinmq_federation_upstream_sets.example_vhost_sets.federation_upstream_sets)
    sets = [
      for set in data.lavinmq_federation_upstream_sets.example_vhost_sets.federation_upstream_sets : {
        name           = set.name
        upstreams      = set.upstreams
        upstream_count = length(set.upstreams)
      }
    ]
  }
}

# Filter sets by upstream count
output "large_sets" {
  description = "Federation upstream sets with more than 2 upstreams"
  value = [
    for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets :
    set.name if length(set.upstreams) > 2
  ]
}

# Find sets containing a specific upstream
output "sets_with_upstream1" {
  description = "Sets containing upstream-1"
  value = [
    for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets :
    set.name if contains(set.upstreams, "example-upstream-1")
  ]
}

# Group sets by upstream count
output "sets_by_size" {
  value = {
    small = [
      for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets :
      set.name if length(set.upstreams) <= 2
    ]
    medium = [
      for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets :
      set.name if length(set.upstreams) > 2 && length(set.upstreams) <= 4
    ]
    large = [
      for set in data.lavinmq_federation_upstream_sets.all_sets.federation_upstream_sets :
      set.name if length(set.upstreams) > 4
    ]
  }
}

# Get details of specific set
locals {
  ha_set_1_details = [
    for set in data.lavinmq_federation_upstream_sets.example_vhost_sets.federation_upstream_sets :
    set if set.name == "ha-set-1"
  ][0]
}

output "ha_set_1_details" {
  value = {
    name           = local.ha_set_1_details.name
    vhost          = local.ha_set_1_details.vhost
    upstreams      = local.ha_set_1_details.upstreams
    upstream_count = length(local.ha_set_1_details.upstreams)
  }
}
