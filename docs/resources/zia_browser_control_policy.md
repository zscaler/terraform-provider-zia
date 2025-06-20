---
subcategory: "Browser Control Policy"
layout: "zscaler"
page_title: "ZIA: browser_control_policy"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-browser-control-policy
  API documentation https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get
  Updates the malware protection policy configuration details
---

# Resource: zia_browser_control_policy

* [Official documentation](https://help.zscaler.com/zia/configuring-browser-control-policy)
* [API documentation](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

The **zia_browser_control_policy** resource allows you to update the malware protection policy configuration details. To learn more see [Configuring the Browser Control Policy](https://help.zscaler.com/zia/configuring-browser-control-policy)

## Example Usage

```hcl
resource "zia_browser_control_policy" "this" {
    plugin_check_frequency = "DAILY"
    bypass_plugins = ["ACROBAT", "FLASH", "SHOCKWAVE"]
    bypass_applications = ["OUTLOOKEXP", "MSOFFICE"]
    blocked_internet_explorer_versions = ["IE10", "MSE81", "MSE92"]
    blocked_chrome_versions = ["CH143", "CH142"]
    blocked_firefox_versions = ["MF145", "MF144"]
    blocked_safari_versions = ["AS19", "AS18"]
    blocked_opera_versions = ["O129X", "O130X"]
    bypass_all_browsers = true
    allow_all_browsers = true
    enable_warnings = true
}
```

## Example Usage  - Enable Smart Isolation

```hcl
data "zia_cloud_browser_isolation_profile" "this" {
    name = "ZS_CBI_Profile1"
}

resource "zia_browser_control_policy" "this" {
    plugin_check_frequency = "DAILY"
    bypass_plugins = ["ACROBAT", "FLASH", "SHOCKWAVE"]
    bypass_applications = ["OUTLOOKEXP", "MSOFFICE"]
    blocked_internet_explorer_versions = ["IE10", "MSE81", "MSE92"]
    blocked_chrome_versions = ["CH143", "CH142"]
    blocked_firefox_versions = ["MF145", "MF144"]
    blocked_safari_versions = ["AS19", "AS18"]
    blocked_opera_versions = ["O129X", "O130X"]
    bypass_all_browsers = true
    allow_all_browsers = true
    enable_warnings = true
    enable_smart_browser_isolation = true
    smart_isolation_profile {
      id = data.zia_cloud_browser_isolation_profile.this.id
    }
}
```

## Argument Reference

The following arguments are supported:

### Optional

* `plugin_check_frequency` - (String) Specifies how frequently the service checks browsers and relevant applications to warn users regarding outdated or vulnerable browsers, plugins, and applications. If not set, the warnings are disabled. Supported Values:
  * `DAILY`
  * `WEEKLY`
  * `MONTHLY`,
  * `EVERY_2_HOURS`
  * `EVERY_4_HOURS`
  * `EVERY_6_HOURS`
  * `EVERY_8_HOURS`
  * `EVERY_12_HOURS`

* `bypass_plugins` - (List) List of plugins that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable plugins are warned.Supported Values:
  * `ANY`
  * `NONE`
  * `ACROBAT`
  * `FLASH`
  * `SHOCKWAVE`
  * `QUICKTIME`
  * `DIVX`
  * `GOOGLEGEARS`
  * `DOTNET`
  * `SILVERLIGHT`
  * `REALPLAYER`
  * `JAVA`
  * `TOTEM`
  * `WMP`

* `bypass_applications` - (List) List of applications that need to be bypassed for warnings. This attribute has effect only if the 'enableWarnings' attribute is set to true. If not set, all vulnerable applications are warned. Supported Values:
  * `ANY`
  * `NONE`
  * `OUTLOOKEXP`
  * `MSOFFICE`

* `blocked_internet_explorer_versions` - (List) Versions of Microsoft browser that need to be blocked. If not set, all Microsoft browser versions are allowed. See all [Supported values](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

* `blocked_chrome_versions` - (List) Versions of Google Chrome browser that need to be blocked. If not set, all Google Chrome versions are allowed. See all [Supported values](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

* `blocked_firefox_versions` - (List) Versions of Mozilla Firefox browser that need to be blocked. If not set, all Mozilla Firefox versions are allowed. See all [Supported values](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

* `blocked_safari_versions` - (List) Versions of Apple Safari browser that need to be blocked. If not set, all Apple Safari versions are allowed. See all [Supported values](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

* `blocked_opera_versions` - (List) Versions of Opera browser that need to be blocked. If not set, all Opera versions are allowed. See all [Supported values](https://help.zscaler.com/zia/browser-control-policy#/browserControlSettings-get)

* `bypass_all_browsers` - (Boolean) If set to true, all the browsers are bypassed for warnings

* `allow_all_browsers` - (Boolean) A Boolean value that specifies whether or not to allow all the browsers and their respective versions access to the internet

* `enable_warnings` - (Boolean) A Boolean value that specifies if the warnings are enabled

* `enable_smart_browser_isolation` - (Boolean) A Boolean value that specifies if Smart Browser Isolation is enabled

* `smart_isolation_profile` - (Block, Max: 1) The isolation profile ID used for DLP email alerts sent to the auditor.
  * `id` - (int) A unique identifier for an entity.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_browser_control_policy** can be imported by using `browser_settings` as the import ID.

For example:

```shell
terraform import zia_browser_control_policy.this "browser_settings"
```
