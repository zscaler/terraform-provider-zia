---
subcategory: "IPS Control Policy"
layout: "zscaler"
page_title: "ZIA: ips_categories"
description: |-
  Official documentation https://help.zscaler.com/zia/configuring-custom-ips-signature-rules
  API documentation https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsCategories-get
  Retrieves an IPS category.
---

# zia_ips_categories (Data Source)

* [Official documentation](https://help.zscaler.com/zia/configuring-custom-ips-signature-rules)
* [API documentation](https://help.zscaler.com/legacy-apis/ips-control-policy#/ipsCategories-get)

Use the **zia_ips_categories** data source to retrieve IPS categories from the Zscaler Internet Access cloud. Supply `id` or `name` to resolve a single category, or omit both to return every category so you can discover which ones are available.

The category `name` returned by this data source is the exact identifier expected by the `res_categories` and `dest_ip_categories` attributes on firewall rules, so this data source lets you reference categories by a validated lookup instead of hard-coding literal strings.

## Example Usage

### Retrieve every category

When neither `id` nor `name` is set, the full list is returned in `categories`. This is the easiest way to discover the available category names:

```hcl
data "zia_ips_categories" "all" {}

output "category_names" {
  value = sort([for c in data.zia_ips_categories.all.categories : c.name])
}
```

You can filter that list in HCL to narrow it down:

```hcl
output "spyware_categories" {
  value = [
    for c in data.zia_ips_categories.all.categories : c
    if can(regex("SPYWARE", c.name))
  ]
}
```

### Look up by name

```hcl
data "zia_ips_categories" "adspyware" {
  name = "ADSPYWARE"
}
```

### Look up by ID

```hcl
data "zia_ips_categories" "this" {
  id = 74
}
```

### Reference from a Cloud Firewall DNS rule

Select the categories you need from a single unfiltered query in `locals`, then pass those locals to the rule attributes:

```hcl
data "zia_ips_categories" "all" {}

locals {
  dest_ip_categories = [
    for c in data.zia_ips_categories.all.categories : c.name
    if contains(["ADSPYWARE_SITES", "BOTNET", "DGA", "MALWARE_SITE", "PHISHING"], c.name)
  ]

  res_categories = [
    for c in data.zia_ips_categories.all.categories : c.name
    if contains(["ADSPYWARE_SITES", "BOTNET", "DGA", "MALWARE_SITE", "PHISHING"], c.name)
  ]
}

resource "zia_firewall_dns_rule" "this" {
  name               = "Block-Malicious-Domains"
  description        = "Block malicious domain categories"
  action             = "BLOCK"
  state              = "ENABLED"
  order              = 1
  rank               = 7
  protocols          = ["ANY_RULE"]
  dest_ip_categories = local.dest_ip_categories
  res_categories     = local.res_categories
}
```

The same pattern applies to `zia_firewall_ips_rule`, which exposes the same two attributes.

A name misspelled inside the `contains` list is silently skipped rather than reported, because the filter simply does not match. Individual lookups avoid that: they error at plan time when the category does not exist, which also makes them the better choice when you only need one or two categories.

```hcl
data "zia_ips_categories" "adspyware" {
  name = "ADSPYWARE_SITES"
}

data "zia_ips_categories" "botnet" {
  name = "BOTNET"
}

resource "zia_firewall_dns_rule" "this" {
  name      = "Block-Malicious-Domains"
  action    = "BLOCK"
  state     = "ENABLED"
  order     = 1
  rank      = 7
  protocols = ["ANY_RULE"]

  res_categories = [
    data.zia_ips_categories.adspyware.name,
    data.zia_ips_categories.botnet.name,
  ]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Optional, Integer) The unique identifier of the IPS category.
* `name` - (Optional, String) The name of the IPS category. Matching is case-insensitive.

~> **NOTE** Both arguments are optional. Omit them to return every category in `categories`. When both are set, `id` takes precedence. Looking up a category that does not exist is an error, which makes this a convenient way to validate category names at plan time.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `categories` - (List of Objects) Every available IPS category, populated when both `id` and `name` are omitted. Each entry exports `id`, `name`, `back_end_name`, `description`, `deleted`, `predefined`, and `ips_signature_rules_count`.

The following attributes describe the single matching category and are populated when `id` or `name` is set:

* `back_end_name` - (String) The descriptive name of the category as displayed in the admin portal, for example `Suspected Spyware or Adware`.
* `description` - (String) Additional information about the category.
* `deleted` - (Boolean) Indicates whether the category has been deleted.
* `predefined` - (Boolean) Indicates whether the category is predefined by the service.
* `ips_signature_rules_count` - (Integer) The number of IPS signature rules associated with the category.
