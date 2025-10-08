provider "lavinmq" {
  baseurl  = "http://localhost:15672/"
  username = "guest"
  password = "guest"
}

# Create a vhost for testing
resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

# Create some exchanges in the vhost
resource "lavinmq_exchange" "exchange1" {
  name    = "example-exchange-1"
  vhost   = lavinmq_vhost.example.name
  type    = "direct"
  durable = true
}

resource "lavinmq_exchange" "exchange2" {
  name        = "example-exchange-2"
  vhost       = lavinmq_vhost.example.name
  type        = "fanout"
  durable     = false
  auto_delete = true
}

# Use the data source to list all exchanges
data "lavinmq_exchanges" "all_exchanges" {
  depends_on = [lavinmq_exchange.exchange1, lavinmq_exchange.exchange2]
}

# Output the exchange information
output "exchanges_info" {
  value = [
    for exchange in data.lavinmq_exchanges.all_exchanges.exchanges : {
      name        = exchange.name
      vhost       = exchange.vhost
      type        = exchange.type
      durable     = exchange.durable
      auto_delete = exchange.auto_delete
    }
  ]
}
