resource "lavinmq_user" "example" {
  name     = "example-user"
  password = "example-password"
  tags     = ["administrator"]
}
