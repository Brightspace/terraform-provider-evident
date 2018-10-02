provider "evident" {
  access_key = "${var.access_key}"
  secret_key = "${var.secret_key}"
}

resource "random_uuid" "evident_external_id" {}

resource "evident_external_account" "evident" {
  name        = "ACE-Acme-Corporation"
  external_id = "${random_uuid.evident_external_id.result}"
  arn         = "${var.role_arn}"
  team_id     = "${var.team_id}"
}
