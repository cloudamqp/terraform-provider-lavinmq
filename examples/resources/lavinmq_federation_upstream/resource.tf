resource "lavinmq_federation_upstream" "example" {
  name     = "upstream-cluster"
  vhost    = "/"
  uri      = "amqp://guest:guest@upstream-broker.example.com:5672/%2f"
  exchange = "federated-exchange"
}
