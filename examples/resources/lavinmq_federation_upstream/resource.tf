resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_federation_upstream" "example" {
  name     = "upstream-cluster"
  vhost    = lavinmq_vhost.example.name
  uri      = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange = "federated-exchange"
}
