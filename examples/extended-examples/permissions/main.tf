resource "lavinmq_user" "app_user" {
  name     = "app-user"
  password = "app-password"
  tags     = ["management"]
}

resource "lavinmq_vhost" "app_vhost" {
  name = "app-vhost"
}

resource "lavinmq_permission" "app_user_full" {
  vhost     = "/"
  user      = lavinmq_user.app_user.name
  configure = ".*"
  read      = ".*"
  write     = ".*"
}

resource "lavinmq_permission" "app_user_limited" {
  vhost     = lavinmq_vhost.app_vhost.name
  user      = lavinmq_user.app_user.name
  configure = "^$"
  read      = ".*"
  write     = "^app-queue-.*"
}
