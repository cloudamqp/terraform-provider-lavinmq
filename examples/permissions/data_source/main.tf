resource "lavinmq_vhost" "app_vhost" {
  name = "app-vhost"
}

resource "lavinmq_user" "admin_user" {
  name     = "admin-user"
  password = "admin-password"
  tags     = ["administrator"]
}

resource "lavinmq_user" "app_user" {
  name     = "app-user"
  password = "app-password"
  tags     = ["management"]
}

resource "lavinmq_permission" "admin_default" {
  vhost     = "/"
  user      = lavinmq_user.admin_user.name
  configure = ".*"
  read      = ".*"
  write     = ".*"
}

resource "lavinmq_permission" "app_default" {
  vhost     = "/"
  user      = lavinmq_user.app_user.name
  configure = "^$"
  read      = ".*"
  write     = "^app-.*"
}

resource "lavinmq_permission" "app_vhost_perm" {
  vhost     = lavinmq_vhost.app_vhost.name
  user      = lavinmq_user.app_user.name
  configure = ".*"
  read      = ".*"
  write     = ".*"
}

# List all permissions
data "lavinmq_permissions" "all" {
  depends_on = [
    lavinmq_permission.admin_default,
    lavinmq_permission.app_default,
    lavinmq_permission.app_vhost_perm
  ]
}

# Filter permissions by vhost
data "lavinmq_permissions" "app_vhost_only" {
  vhost = lavinmq_vhost.app_vhost.name
  depends_on = [
    lavinmq_permission.app_vhost_perm
  ]
}

# Filter permissions by user
data "lavinmq_permissions" "app_user_only" {
  user = lavinmq_user.app_user.name
  depends_on = [
    lavinmq_permission.app_default,
    lavinmq_permission.app_vhost_perm
  ]
}

# Filter by both vhost and user (returns single permission)
data "lavinmq_permissions" "app_user_on_app_vhost" {
  vhost = lavinmq_vhost.app_vhost.name
  user  = lavinmq_user.app_user.name
  depends_on = [
    lavinmq_permission.app_vhost_perm
  ]
}

output "all_permissions" {
  value = data.lavinmq_permissions.all.permissions
}

output "app_vhost_permissions" {
  value = data.lavinmq_permissions.app_vhost_only.permissions
}

output "app_user_permissions" {
  value = data.lavinmq_permissions.app_user_only.permissions
}

output "specific_permission" {
  value = data.lavinmq_permissions.app_user_on_app_vhost.permissions
}
