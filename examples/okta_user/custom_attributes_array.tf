resource "okta_user_schema" "test" {
  index  = "customAttribute123"
  title  = "terraform acceptance test"
  type   = "string"
  master = "PROFILE_MASTER"
}

resource "okta_user_schema" "test_array" {
  index      = "array123"
  title      = "terraform acceptance test"
  type       = "array"
  array_type = "string"
  master     = "PROFILE_MASTER"
  depends_on = ["okta_user_schema.test_number", "okta_user_schema.test"]
}

resource "okta_user_schema" "test_number" {
  index      = "number123"
  title      = "terraform acceptance test"
  type       = "number"
  master     = "PROFILE_MASTER"
  depends_on = ["okta_user_schema.test"]
}

resource "okta_user" "testAcc_replace_with_uuid" {
  admin_roles = ["APP_ADMIN", "USER_ADMIN"]
  first_name  = "TestAcc"
  last_name   = "Smith"
  login       = "test-acc-replace_with_uuid@example.com"
  email       = "test-acc-replace_with_uuid@example.com"

  custom_profile_attributes = <<JSON
  {
    "array123": ["test"],
    "number123": 1
  }
JSON

  depends_on = ["okta_user_schema.test_array", "okta_user_schema.test_number"]
}
