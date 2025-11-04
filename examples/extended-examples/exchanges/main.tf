# LavinMQ Exchange Examples
# This example demonstrates how to create different types of exchanges

terraform {
  required_providers {
    lavinmq = {
      source = "hashicorp/lavinmq"
    }
  }
}

provider "lavinmq" {
  baseurl  = "http://127.0.0.1:15672"
  username = "guest"
  password = "guest"
}

# Create a test vhost for our examples
resource "lavinmq_vhost" "test" {
  name = "test"
}

# Create a direct exchange
resource "lavinmq_exchange" "direct_example" {
  name  = "direct-exchange"
  vhost = lavinmq_vhost.test.name
  type  = "direct"
  # auto_delete = false
  # durable     = true
}

# Create a fanout exchange
resource "lavinmq_exchange" "fanout_example" {
  name        = "fanout-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "fanout"
  auto_delete = false
  durable     = true
}

# Create a topic exchange
resource "lavinmq_exchange" "topic_example" {
  name        = "topic-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "topic"
  auto_delete = false
  durable     = true
}

# Create a headers exchange
resource "lavinmq_exchange" "headers_example" {
  name        = "headers-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "headers"
  auto_delete = false
  durable     = true
}

# Create an auto-delete exchange
resource "lavinmq_exchange" "temp_exchange" {
  name        = "temporary-exchange"
  vhost       = lavinmq_vhost.test.name
  type        = "direct"
  auto_delete = true
  durable     = false
}

# Create an exchange with alternate exchange for unroutable messages
resource "lavinmq_exchange" "backup_exchange" {
  name    = "backup-exchange"
  vhost   = lavinmq_vhost.test.name
  type    = "fanout"
  durable = true
}

resource "lavinmq_exchange" "main_exchange_with_args" {
  name    = "main-exchange"
  vhost   = lavinmq_vhost.test.name
  type    = "direct"
  durable = true

  arguments = {
    "alternate-exchange" = lavinmq_exchange.backup_exchange.name
  }
}

# Create a custom vhost and exchange
resource "lavinmq_vhost" "custom" {
  name = "my-app"
}

resource "lavinmq_exchange" "custom_vhost_example" {
  name        = "app-exchange"
  vhost       = lavinmq_vhost.custom.name
  type        = "topic"
  auto_delete = false
  durable     = true
}

# Outputs
output "direct_exchange_id" {
  value = lavinmq_exchange.direct_example.id
}

output "custom_exchange_info" {
  value = {
    name  = lavinmq_exchange.custom_vhost_example.name
    vhost = lavinmq_exchange.custom_vhost_example.vhost
    type  = lavinmq_exchange.custom_vhost_example.type
  }
}
