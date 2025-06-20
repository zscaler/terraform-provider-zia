---
subcategory: "URL Categories"
layout: "zscaler"
page_title: "ZIA: url_categories"
description: |-
    Official documentation https://help.zscaler.com/zia/about-url-categories
    API documentation https://help.zscaler.com/zia/url-categories#/urlCategories-get
    Gets information about all or custom URL categories. By default, the response includes keywords.
---

# zia_url_categories (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-url-categories)
* [API documentation](https://help.zscaler.com/zia/url-categories#/urlCategories-get)

Use the **zia_url_categories** data source to get information about all or custom URL categories. By default, the response includes keywords.

```hcl
// Query a URL Category by Name
data "zia_url_categories" "this"{
    configured_name = "Example"
}
```

```hcl
// Query a URL Category by its Custom ID
data "zia_url_categories" "this"{
    id = "CUSTOM_08"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (String) URL category

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `configured_name` - (String) Name of the URL category. This is only required for custom URL categories.
* `keywords` - (List of String) Custom keywords associated to a URL category. Up to 2048 custom keywords can be added per organization across all categories (including bandwidth classes).
* `super_category` - (String)
* `custom_category` - (Boolean) Set to true for custom URL category. Up to 48 custom URL categories can be added per organization.
* `custom_urls_count` - (Number) The number of custom URLs associated to the URL category.
* `db_categorized_urls` - (List of String) URLs added to a custom URL category are also retained under the original parent URL category (i.e., the predefined category the URL previously belonged to).
* `description` - (String) Description of the category.
* `editable` - (Boolean) Value is set to false for custom URL category when due to scope user does not have edit permission
* `type` - (String) Type of the custom categories. `URL_CATEGORY`, `TLD_CATEGORY`, `ALL`
* `urls` - (List of String) Custom URLs to add to a URL category. Up to 25,000 custom URLs can be added per organization across all categories (including bandwidth classes).
* `urls_retaining_parent_category_count` - (Number) The number of custom URLs associated to the URL category, that also need to be retained under the original parent category.
* `val` - (Number) The unique ID for the URL category.

* `scopes` - (List of Object) Scope of the custom categories.
  * `scope_group_member_entities` - (List of Object) Only applicable for the LOCATION_GROUP admin scope type, in which case this attribute gives the list of ID/name pairs of locations within the location group. The attribute name is subject to change.
        - `id` - (String) Identifier that uniquely identifies an entity
        - `name` - (String) The configured name of the entity
        - `extensions` - (Map of String)

  * `type` - (String) The admin scope type. The attribute name is subject to change. `ORGANIZATION`, `DEPARTMENT`, `LOCATION`, `LOCATION_GROUP`

  * `scope_entities` - (List of Object)
        - `id` - (String) Identifier that uniquely identifies an entity
        - `name` - (String) The configured name of the entity
        - `extensions` - (Map of String)

* `url_keyword_counts` - (List of Object) URL and keyword counts for the category.
  * `total_url_count` - (Number) Custom URL count for the category.
  * `retain_parent_url_count` - (Number) Count of URLs with retain parent category.
  * `total_keyword_count` - (Number) Total keyword count for the category.
  * `retain_parent_keyword_count` - (Number) Count of total keywords with retain parent category.

* `custom_ip_ranges_count` - (Number) The number of custom IP address ranges associated to the URL category.
* `ip_ranges_retaining_parent_category_count` - (Number) The number of custom IP address ranges associated to the URL category, that also need to be retained under the original parent category.
