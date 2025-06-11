---
subcategory: "FTP Control Policy"
layout: "zscaler"
page_title: "ZIA: ftp_control_policy"
description: |-
  Retrieves the FTP Control status and the list of URL categories for which FTP is allowed.
---

# Data Source: zia_ftp_control_policy

Use the **zia_ftp_control_policy** data source to retrieves the FTP Control Policy configuration. To learn more see [Configuring the FTP Control Policy](https://help.zscaler.com/zia/configuring-ftp-control-policy)

## Example Usage

```hcl
data "zia_ftp_control_policy" "this" {}
```

## Argument Reference

This data source can be executed without the need of additional parameters.

## Attribute Reference

* `ftp_over_http_enabled` - (Boolean) Indicates whether to enable FTP over HTTP. By default, the Zscaler service doesn't allow users from a location to upload or download files from FTP sites that use FTP over HTTP. Select this to enable browsers to connect to FTP over HTTP sites and download files. If a remote user uses a dedicated port, then the service supports FTP over HTTP for them.
* `ftp_enabled` - (Boolean) Indicates whether to enable native FTP. When enabled, users can connect to native FTP sites and download files.

* `url_categories` - (List of Strings) List of URL categories that allow FTP traffic
* `urls` - (List of Strings) Domains or URLs included for the FTP Control settings