resource "lavinmq_user" "test-user-password" {
  name             = "test-user-password"
  password         = "test-password"
  password_version = 1
  tags             = ["monitoring"]
}

resource "lavinmq_user" "test-user-passwordhash" {
  name = "test-user-passwordhash"
  password_hash = {
    value     = "c6xQEdMpUle9NihE3SV8xcpXZtC6/z57IVlB22d/yEVw545L"
    algorithm = "sha256"
  }
  password_version = 1
  tags             = ["monitoring"]
}
