resource "evident_external_account" "default" {
  name        = "EvidentProviderTest"
  external_id = var.external_id
  arn         = var.role_arn
  team_id     = var.team_id
}

variable "team_id" {
  type        = string
  description = "The ID of the team the external account belongs to"
}

variable "external_id" {
  type        = string
  description = "The ID of the team the external account belongs to"
}

variable "role_arn" {
  type        = string
  description = "Amazon Resource Name for the IAM role"
}

output "id" { value = evident_external_account.default.id }
