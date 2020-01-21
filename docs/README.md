# Evident Provider

The Evident provider is used to connect AWS accounts to Evident.io. This provider is not official, and is developed to cover a narrow slice of the Evident API.

The provider allows you to manage your Evident.io AWS connections. It needs to be configured with the proper credentials before it can be used.

## Example Usage

```hcl
# Configure the Evident Provider
provider "evident" {}

# Configure an AWS with evident
resource "evident_external_account" "aws" {
  name        = "ACME-AWSAccount"
  external_id = "cd45b8ec-df86-4698-9c80-6a63236fb6c6"
  arn         = "arn:aws:iam::123412341234:role/EvidentIORole"
  team_id     = "12345"
}
```

The evident provider is a [third party custom provider](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins). Third-party providers must be manually installed, since `terraform init` cannot automatically download them.

## Authentication

The Evident provider can be provided credentials for authentication using environment variables or static credentials.

### Static credentials

> Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file ever be committed to a public version control system

Static credentials can be provided by adding an `access_key` and `secret_key` in-line in the Evident provider block or by variables:

Usage:

```hcl
provider "evident" {
    access_key = "my-api-key"
    secret_key = "my-secret-key"
}
```
### Environment variables

You can provide your credentials via the `EVIDENT_ACCESS_KEY` and `EVIDENT_SECRET_KEY`, environment variables, representing your Evident public and private keys, respectively.

```hcl
provider "evident" {}
```

Usage:

```bash
export EVIDENT_ACCESS_KEY="my-api-key"
export EVIDENT_SECRET_KEY="my-secret-key"
terraform plan
```

## Argument Reference

In addition to [generic `provider` arguments](https://www.terraform.io/docs/configuration/providers.html), the following arguments are supported in the Evident provider block:

- `access_key` - (Optional) This is the Evident access key. It must be provided, but it can also be sourced from the `EVIDENT_ACCESS_KEY` environment variable.
- `secret_key` - (Optional) This is the Evident secret key. It must be provided, but it can also be sourced from the `EVIDENT_SECRET_KEY` environment variable.