---
subcategory: "Secure Browsing"
layout: "zscaler"
page_title: "ZIA: supported_browser_version"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-browser-control-policy
  API documentation https://help.zscaler.com/legacy-apis/browser-control-policy#/browserControlSettings/supportedBrowserVersions-get
  Get the list of all supported browsers and their versions
---

# zia_supported_browser_version (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-browser-control-policy)
* [API documentation](https://help.zscaler.com/legacy-apis/browser-control-policy#/browserControlSettings/supportedBrowserVersions-get)

Use the **zia_supported_browser_version** data source to retrieve the list of all supported browsers and their current and older version identifiers in the Zscaler Internet Access cloud. The returned values are the canonical version strings that the ZIA browser control policy uses on attributes such as `blocked_chrome_versions`, `blocked_firefox_versions`, `blocked_safari_versions`, `blocked_opera_versions`, and `blocked_internet_explorer_versions`.

## Example Usage - Retrieve All Browsers

```hcl
# Returns every supported browser entry. The `browsers` attribute is a list
# of objects, one per browser type.
data "zia_supported_browser_version" "all" {}

output "all_browsers" {
  value = data.zia_supported_browser_version.all.browsers
}
```

## Example Usage - Filter by Browser Type

```hcl
# Narrow the result set to a single browser type.
data "zia_supported_browser_version" "chrome" {
  browser_type = "CHROME"
}

# `browsers` is still a list; for a single-browser lookup use `one()`
# to flatten it into a scalar.
output "chrome_versions" {
  value = one([for b in data.zia_supported_browser_version.chrome.browsers : b.versions])
}

output "chrome_older_versions" {
  value = one([for b in data.zia_supported_browser_version.chrome.browsers : b.older_versions])
}
```

## Example Usage - Drive a `blocked_*_versions` Setting

The blocked-version attributes on `zia_browser_control_policy` (`blocked_chrome_versions`, `blocked_firefox_versions`, `blocked_safari_versions`, `blocked_opera_versions`, `blocked_internet_explorer_versions`) expect canonical ZIA version identifiers. This data source is the source of truth for those identifiers. Pick the pattern that matches what you're trying to express.

### Pattern A — Hardcode the exact versions you want to block

If you already know which versions to block and you don't need the catalogue lookup at plan time, just write the list. This is the most deterministic form and produces the most readable diffs.

```hcl
resource "zia_browser_control_policy" "this" {
  blocked_chrome_versions = ["CH147", "CH146"]
}
```

### Pattern B — Block every older Chrome version that ZIA still recognises

Useful when your policy is "always block whatever Zscaler has classified as outdated".

```hcl
data "zia_supported_browser_version" "chrome" {
  browser_type = "CHROME"
}

resource "zia_browser_control_policy" "this" {
  blocked_chrome_versions = data.zia_supported_browser_version.chrome.browsers[0].older_versions
}
```

### Pattern C — Block a specific subset, validated against the live catalogue

Asserts that every version you want to block is actually published by ZIA. If a version is missing from the upstream catalogue, the lifecycle precondition fails and the plan errors out instead of silently dropping it.

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
  blocked_chrome_versions = local.blocked_chrome

  lifecycle {
    precondition {
      condition     = length(local.blocked_chrome) == length(local.wanted_chrome)
      error_message = "One or more wanted Chrome versions are not present in the supported browser catalogue."
    }
  }
}
```

### Pattern D — Block every version matching a prefix

Useful for sweeping rules like "block every Chrome 1.x.x identifier". The `for` expression subsets the inner `versions` list, which is something the JMESPath `search` argument cannot do directly (see the note below).

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
  blocked_chrome_versions = local.blocked_chrome
}
```

## Example Usage - With JMESPath Search

The `search` argument applies a [JMESPath](https://jmespath.org/) expression to the response client-side after the data is retrieved from the API. Field names in expressions must use the API's camelCase names: `browserType`, `versions`, `olderVersions`.

~> **NOTE — `search` is a gate on the outer `browsers` array, not an inner-list subset.** A predicate like `[?contains(versions, 'CH147')]` returns the *entire* matching browser entry (with its full `versions` and `older_versions` lists) when the predicate is satisfied — it does **not** return only the matched version strings. If your goal is to block exactly the two versions in your predicate, see Patterns A and C in the previous section; if your goal is to block versions matching a prefix or other per-element rule, see Pattern D. Use `search` when you need to select **which browser entries** come back, not when you need to subset versions within an entry.

```hcl
# Browsers whose current `versions` list contains a specific identifier
data "zia_supported_browser_version" "has_c130x" {
  search = "[?contains(versions, 'C130X')]"
}
```

```hcl
# Browsers whose `olderVersions` list contains a specific identifier
data "zia_supported_browser_version" "had_c100x" {
  search = "[?contains(olderVersions, 'C100X')]"
}
```

```hcl
# Browsers with more than 100 older versions on record
data "zia_supported_browser_version" "deep_history" {
  search = "[?length(olderVersions) > `100`]"
}
```

```hcl
# Browsers that have at least one older version starting with a prefix
data "zia_supported_browser_version" "old_o8x" {
  search = "[?length(olderVersions[?starts_with(@, 'O8')]) > `0`]"
}
```

```hcl
# Reshape: return each browser with only the older versions that start
# with 'O5' (useful when you want to attach a narrowed slice to a policy).
data "zia_supported_browser_version" "trimmed" {
  search = "[*].{browserType: browserType, versions: versions, olderVersions: olderVersions[?starts_with(@, 'O5')]}"
}
```

```hcl
# Combine the enum filter with JMESPath. The JMESPath expression runs first,
# then `browser_type` narrows the remaining set. This returns Chrome only
# when Chrome's older versions contain 'C100X'.
data "zia_supported_browser_version" "chrome_with_c100x" {
  browser_type = "CHROME"
  search       = "[?contains(olderVersions, 'C100X')]"
}
```

~> **NOTE** JMESPath string literals are wrapped in single quotes (`'C100X'`); numeric and boolean literals are wrapped in backticks (`` `100` ``). Field names inside `search` are the API's **camelCase** names (`browserType`, `olderVersions`), not the Terraform snake_case attribute names.

## Argument Reference

The following arguments are supported:

* `browser_type` - (Optional, String) Return only the supported version entry for this browser type. Must be one of: `CHROME`, `FIREFOX`, `SAFARI`, `OPERA`, `MSCHREDGE`. Omit to return every browser type.

* `search` - (Optional, String) A [JMESPath](https://jmespath.org/) expression applied to the response client-side after all data has been retrieved from the API. Useful for content-based filtering (e.g. "browsers whose older versions contain X") and for reshaping the result set. Field names in expressions must use the API's camelCase names (`browserType`, `versions`, `olderVersions`). When combined with `browser_type`, the JMESPath expression runs first and the enum filter is applied afterwards.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) Constant identifier (`supported_browser_versions`). The endpoint is a tenant-scoped lookup and does not expose per-entry identifiers.

* `browsers` - (List of Object) The (optionally filtered) list of supported browser entries. Each element exposes the following attributes:

    * `browser_type` - (String) The browser type. One of `CHROME`, `FIREFOX`, `SAFARI`, `OPERA`, `MSCHREDGE`.

    * `versions` - (List of String) The currently supported version identifiers for this browser type.

    * `older_versions` - (List of String) The previously supported version identifiers (no longer current, but still recognised by the ZIA browser control policy) for this browser type.
