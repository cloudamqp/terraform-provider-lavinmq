terraform {
  required_providers {
    lavinmq = {
      source = "cloudamqp/lavinmq"
    }
  }
}

provider "lavinmq" {
  baseurl  = "http://localhost:15672"
  username = "guest"
  password = "guest"
}
