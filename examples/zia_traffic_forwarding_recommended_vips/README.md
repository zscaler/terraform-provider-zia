# Recommended GRE tunnel Virtual IP Addresses (VIPs) Example

This example will show you how to use Terraform to export/retrieve a list of recommended GRE tunnel virtual IP addresses (VIPs), based on source IP address or latitude/longitude coordinates.

This example codifies [this API](https://help.zscaler.com/zia/api#/Traffic%20Forwarding/CloudVipsResource_getRecommendedGreVips).

To run, configure your ZIA provider as described [Here](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```
