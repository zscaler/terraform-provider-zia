# URL Filtering Policies Example

This example will show you how to use Terraform to create URL Filtering Policies in the ZIA portal.
This example codifies [this API](https://help.zscaler.com/zia/api#/URL%20Filtering%20Policies/addRule).

To run, configure your ZIA provider as described [Here](https://github.com/willguibr/terraform-provider-zia/blob/master/docs/index.html.markdown)

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
