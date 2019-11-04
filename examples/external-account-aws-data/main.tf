data "evident_external_account_aws" "default" {
  account = var.account
}

variable "account" {
  type        = string
  description = "The identifier for the AWS Account in Evident"
}

output "arn" { value = data.evident_external_account_aws.default.arn }
output "external_id" { value = data.evident_external_account_aws.default.external_id }
