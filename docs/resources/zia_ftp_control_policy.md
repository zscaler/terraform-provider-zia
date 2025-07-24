---
subcategory: "FTP Control Policy"
layout: "zscaler"
page_title: "ZIA: ftp_control_policy"
description: |-
  Official documentation https://help.zscaler.com/zia/about-ftp-control
  API documentation https://help.zscaler.com/zia/ftp-control-policy#/ftpSettings-get
  Updates the FTP Control settings.
---

# zia_ftp_control_policy (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-ftp-control)
* [API documentation](https://help.zscaler.com/zia/ftp-control-policy#/ftpSettings-get)

The **zia_ftp_control_policy** resource allows you to update FTP control Policy. To learn more see [Configuring the FTP Control Policy](https://help.zscaler.com/zia/configuring-ftp-control-policy)

## Example Usage

```hcl
resource "zia_ftp_control_policy" "this" {
    ftp_enabled = true
    ftp_over_http_enabled = true
    url_categories = ["HOBBIES_AND_LEISURE", "HEALTH","HISTORY","INSURANCE","IMAGE_HOST","INTERNET_SERVICES","GOVERNMENT"]
    urls = ["test1.acme.com", "test10.acme.com"]
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `ftp_over_http_enabled` - (Boolean) Indicates whether to enable FTP over HTTP. By default, the Zscaler service doesn't allow users from a location to upload or download files from FTP sites that use FTP over HTTP. Select this to enable browsers to connect to FTP over HTTP sites and download files. If a remote user uses a dedicated port, then the service supports FTP over HTTP for them.
* `ftp_enabled` - (Boolean) Indicates whether to enable native FTP. When enabled, users can connect to native FTP sites and download files.

* `url_categories` - (List of Strings) List of URL categories that allow FTP traffic. See the [URL Categories API](https://help.zscaler.com/zia/url-categories#/urlCategories-get) for the list of available categories or use the data source `zia_url_categories` to retrieve the list of URL categories.
* `urls` - (List of Strings) Domains or URLs included for the FTP Control settings

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_ftp_control_policy** can be imported by using `ftp_control` as the import ID.

For example:

```shell
terraform import zia_ftp_control_policy.this "ftp_control"
```
