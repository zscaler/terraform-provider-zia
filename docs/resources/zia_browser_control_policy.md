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

data "zia_group_management" "this" {
 name = "Finance"
}

data "zia_user_management" "this" {
 email = "adam.ashcroft@acme.com"
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

    smart_isolation_groups {
        id = [ data.zia_group_management.this.id ]
    }
    smart_isolation_users = {
    id = [ data.zia_user_management.this.id ]
  }
}
```

## Example Usage - Drive `blocked_*_versions` from the Browser Catalogue

The `blocked_chrome_versions`, `blocked_firefox_versions`, `blocked_safari_versions`, `blocked_opera_versions`, and `blocked_internet_explorer_versions` attributes expect the canonical version identifiers that ZIA publishes through the [`zia_supported_browser_version`](../data-sources/zia_supported_browser_version.md) data source. The patterns below cover the common ways to wire those two together. The data source page documents each pattern in full.

### Block a specific, hand-picked list

Most deterministic. Use this when you know exactly which versions you want to block.

```hcl
resource "zia_browser_control_policy" "this" {
  plugin_check_frequency  = "DAILY"
  enable_warnings         = true
  blocked_chrome_versions = ["CH147", "CH146"]
}
```

### Block every older Chrome version that ZIA recognises

```hcl
data "zia_supported_browser_version" "chrome" {
  browser_type = "CHROME"
}

resource "zia_browser_control_policy" "this" {
  plugin_check_frequency  = "DAILY"
  enable_warnings         = true
  blocked_chrome_versions = data.zia_supported_browser_version.chrome.browsers[0].older_versions
}
```

### Block a specific subset, validated against the live catalogue

Fails the plan if a wanted version is missing from the upstream catalogue, rather than silently dropping it from the set.

```hcl
data "zia_supported_browser_version" "chrome" {
  browser_type = "CHROME"
}

locals {
  wanted_chrome = ["CH147", "CH146"]

  available_chrome = toset(concat(
    data.zia_supported_browser_version.chrome.browsers[0].versions,
    data.zia_supported_browser_version.chrome.browsers[0].older_versions,
  ))

  blocked_chrome = setintersection(toset(local.wanted_chrome), local.available_chrome)
}

resource "zia_browser_control_policy" "this" {
  plugin_check_frequency  = "DAILY"
  enable_warnings         = true
  blocked_chrome_versions = local.blocked_chrome

  lifecycle {
    precondition {
      condition     = length(local.blocked_chrome) == length(local.wanted_chrome)
      error_message = "One or more wanted Chrome versions are not present in the supported browser catalogue."
    }
  }
}
```

### Block every version matching a prefix

```hcl
data "zia_supported_browser_version" "chrome" {
  browser_type = "CHROME"
}

locals {
  blocked_chrome = [
    for v in data.zia_supported_browser_version.chrome.browsers[0].versions :
    v if startswith(v, "CH1")
  ]
}

resource "zia_browser_control_policy" "this" {
  plugin_check_frequency  = "DAILY"
  enable_warnings         = true
  blocked_chrome_versions = local.blocked_chrome
}
```

~> **NOTE — Do not use the data source `search` argument to subset version strings.** The `search` argument on `zia_supported_browser_version` filters **which browser entries** are returned, not **which versions within an entry**. A predicate like `[?contains(versions, 'CH147')]` returns the entire CHROME entry (with its complete `versions` list) when satisfied — so wiring `browsers[0].versions` from such a query into `blocked_chrome_versions` blocks every Chrome version, not just `CH147`. For "block exactly these versions" use the hand-picked list or `setintersection` pattern above. For "block versions matching a prefix" use the `for` expression pattern.

~> **NOTE — `ANY` is accepted but normalized to `NONE` on the wire.** The API treats `ANY`, `NONE`, and an empty list as equivalent. Writing `blocked_chrome_versions = ["ANY"]` (and similarly for `bypass_plugins`, `bypass_applications`, and the other `blocked_*_versions` attributes) is accepted and the provider preserves `["ANY"]` in state so subsequent plans do not show drift. Switching from `["ANY"]` to a concrete list is a real change and will be applied.

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
