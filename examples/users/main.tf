resource "lavinmq_user" "admin_user" {
  name             = "test-user"
  password         = "securepassword12345"
  password_version = 2
  # password_hash = {
  #   value     = "$6$rounds=656000$wHj3bX1bQz8JzE2G$y1r7Zk9h8jFzQxYv1K"
  #   algorithm = "sha512"
  # }
  tags = ["monitoring"]
}
