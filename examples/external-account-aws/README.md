# AWS External Account

This example describes how to configure an AWS external account in Evident.

## How to Run

You will need to have setup Evident API credentials, which you learn about in the [usage section](../../README.md). You will also need an IAM role setup for this to work.

```bash
terraform apply \
        -var="team_id=12345" \
        -var="external_id=dfa8cc1d-6eb7-74f0-adb8-8af229cb6c83" \
        -var="role_arn=arn:aws:iam::123412341234:role/EvidentIO"
````