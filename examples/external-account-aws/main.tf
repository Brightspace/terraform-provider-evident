resource "random_uuid" "evident_external_id" {}

resource "evident_external_account" "default" {
  name        = "EvidentProviderTest"
  external_id = "${random_uuid.evident_external_id.result}"
  arn         = "${var.role_arn}"
  team_id     = "${var.team_id}"
}

variable "team_id" {
  type        = string
  description = "The ID of the team the external account belongs to"
}

variable "role_arn" {
  type        = string
  description = "Amazon Resource Name for the IAM role"
}

output "id" { value = evident_external_account.default.id }
