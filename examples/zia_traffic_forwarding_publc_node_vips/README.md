# List of Public Virtual IP Addresses (VIPs) Example

This example will show you how to use Terraform to export/retrieve a list of the virtual IP addresses (VIPs) available in the Zscaler cloud, including region and data center information.

This example codifies [this API](https://help.zscaler.com/zia/api#/Traffic%20Forwarding/getZscalerNodesDetails).

To run, configure your ZIA provider as described [Here](https://github.com/willguibr/terraform-provider-zia/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```
