# AWS External Account

This example describes how to configure an AWS external account in Evident.

## How to Run

You will need to have setup Evident API credentials, which you learn about in the [usage section](../../README.md). You should also setup an IAM role for Evident to connect with, but this is not necessary.

```bash
terraform apply \
        -var="team_id=12345" \
        -var="role_arn=arn:aws:iam::123412341234:role/EvidentIO"
````