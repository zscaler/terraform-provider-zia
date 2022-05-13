# GRE Tunnel Information Example

This example will show you how to use Terraform to export/retrieve a list of IP addresses from a specific GRE Tunnel with its respective details.

This example codifies [this API](https://help.zscaler.com/zia/api#/Traffic%20Forwarding/getIPGWDetails).

To run, configure your ZIA provider as described [Here](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```

## Destroy 💥

```bash
terraform destroy
```
