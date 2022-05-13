# Traffic Forwarding  - Static IP Resource Example

This example will show you how to use Terraform to create a static IP address resource in the ZIA portal.
This example codifies [this API](https://help.zscaler.com/zia/api#/Traffic%20Forwarding/StaticIPResource_addStaticIP).

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
