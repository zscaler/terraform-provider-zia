# Firewall Filtering - IP Source Groups Example

This example will show you how to use Terraform to create IP source groups for use in a Cloud Firewall rule in the ZIA portal.
This example codifies [this API](https://help.zscaler.com/zia/api#/Firewall%20Policies/IpDestinationGroupResource_addDestinationIpGroup).

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
