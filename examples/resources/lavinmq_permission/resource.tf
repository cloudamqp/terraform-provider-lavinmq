resource "lavinmq_vhost" "example" {
  name = "example-vhost"
}

resource "lavinmq_user" "example" {
  name     = "example-user"
  password = "example-password"
  tags     = []
}

resource "lavinmq_permission" "example" {
  user      = lavinmq_user.example.name
  vhost     = lavinmq_vhost.example.name
  configure = ".*"
  write     = ".*"
  read      = ".*"
}
