# AWS External Data Access

This example describes how to read an existing AWS external account from Evident. This can then be used to ensure that the evident IAM role maps as expected for the account.

## How to Run

You will need to have setup Evident API credentials, which you learn about in the [usage section](../../README.md).

```bash
terraform apply \
        -var="account=123456789"
````