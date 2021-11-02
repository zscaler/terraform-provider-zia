---
subcategory: "URL Categories"
layout: "zia"
page_title: "ZIA: url_categories"
description: |-
        Adds a new custom URL category. If keywords are included within the request, they will be added to the new category.
---

# zia_url_categories (Resource)

The **zia_url_categories** - resource adds a new custom URL category.

```hcl
resource "zia_url_categories" "example" {
  super_category      = "USER_DEFINED"
  configured_name     = "MCAS Unsanctioned Apps"
  description         = "MCAS Unsanctioned Apps"
  keywords            = ["microsoft"]
  custom_category     = true
  type                = "URL_CATEGORY"
  scopes {
    type = "LOCATION"
    scope_entities {
      id = [ data.zpa_location_management.nyc_site.id ]
    }
    scope_group_member_entities {
      id = [ data.zia_group_management.engineering.id ]
    }
  }
  urls = [
    ".coupons.com",
    ".resource.alaskaair.net",
    ".techrepublic.com",
    ".dailymotion.com",
    ".osiriscomm.com",
    ".uefa.com",
    ".Logz.io",
    ".alexa.com",
    ".baidu.com",
    ".cnn.com",
    ".level3.com",
  ]
}

```

## Argument Reference

The following arguments are supported:

* `id` - (String) URL category
* `configured_name` - (String) Name of the URL category. This is only required for custom URL categories.
* `keywords` - (List of String) Custom keywords associated to a URL category. Up to 2048 custom keywords can be added per organization across all categories (including bandwidth classes).
* `super_category` - (String)
* `custom_category` - (Boolean) Set to true for custom URL category. Up to 48 custom URL categories can be added per organization.
* `custom_urls_count` - (Number)
* `db_categorized_urls` - (List of String) URLs added to a custom URL category are also retained under the original parent URL category (i.e., the predefined category the URL previously belonged to).
* `description` - (String) Description of the category.
* `editable` - (Boolean) Value is set to false for custom URL category when due to scope user does not have edit permission
* `type` - (String) Type of the custom categories. `URL_CATEGORY`, `TLD_CATEGORY`, `ALL`
* `urls` - (List of String) Custom URLs to add to a URL category. Up to 25,000 custom URLs can be added per organization across all categories (including bandwidth classes).
* `urls_retaining_parent_category_count` - (Number)
* `val` - (Number)

`scopes` - (List of Object) Scope of the custom categories.

* `scope_group_member_entities` - (List of Object) Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group. The attribute name is subject to change.
    **`id` - (String)
    **`extensions` - (Map of String)
    **`name` - (String)

* `type` - (String) The admin scope type. The attribute name is subject to change. `ORGANIZATION`, `DEPARTMENT`, `LOCATION`, `LOCATION_GROUP`

* `scope_entities` - (List of Object)
    **`id` - (String)
    **`extensions` - (Map of String)
    **`name` - (String)

`url_keyword_counts` - (List of Object) URL and keyword counts for the category.

* `total_url_count` - (Number) Custom URL count for the category.
* `retain_parent_url_count` - (Number) Count of URLs with retain parent category.
* `total_keyword_count` - (Number) Total keyword count for the category.
* `retain_parent_keyword_count` - (Number) Count of total keywords with retain parent category.
