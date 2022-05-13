# GRE Internal IP Range Example

This example will show you how to use Terraform to export/retrieve a list containing the next available GRE tunnel internal IP address ranges from Zscaler.
Note: By default, Terraform will return the first 10 available internal IP ranges.
This example codifies [this API](https://help.zscaler.com/zia/api#/Traffic%20Forwarding/GreTunnelResource_validateAndGetAvailableInternalIpRanges).

To run, configure your ZIA provider as described [Here](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/index.html.markdown)

## Run the example

From inside of this directory:

```bash
terraform init
terraform plan -out theplan
terraform apply theplan
```
