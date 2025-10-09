resource "lavinmq_user" "admin_user" {
  name     = "admin-user"
  password = "admin-password"
  tags     = ["administrator"]
}

resource "lavinmq_user" "monitoring_user" {
  name     = "monitoring-user"
  password = "monitoring-password"
  tags     = ["monitoring"]
}

resource "lavinmq_user" "app_user" {
  name     = "app-user"
  password = "app-password"
  tags     = ["management", "policymaker"]
}

data "lavinmq_users" "all" {
  depends_on = [
    lavinmq_user.admin_user,
    lavinmq_user.monitoring_user,
    lavinmq_user.app_user
  ]
}

output "all_users" {
  value = data.lavinmq_users.all.users
}

output "user_names" {
  value = [for user in data.lavinmq_users.all.users : user.name]
}
