# Location Management Lite - Example

This example will show you how to use Terraform to retrieve a lite version of a location management including "Road Warrior" location for use in several resource types such as ``zia_url_filtering_rules``, ``zia_firewall_filtering_rule`` and ``zia_dlp_web_rules``
This example codifies [this API](https://help.zscaler.com/zia/location-management#/locations/lite-get).

To run, configure your ZIA provider as described [Here](https://github.com/zscaler/terraform-provider-zia/blob/master/docs/index.md)

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
