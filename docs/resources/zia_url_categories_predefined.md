---
subcategory: "URL Categories"
layout: "zscaler"
page_title: "ZIA: url_categories_predefined"
description: |-
    Official documentation https://help.zscaler.com/zia/about-url-categories
    API documentation https://help.zscaler.com/zia/url-categories#/urlCategories-get
    Manages mutable fields of predefined URL categories such as custom URLs, keywords, and IP ranges.
---

# zia_url_categories_predefined (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-url-categories)
* [API documentation](https://help.zscaler.com/zia/url-categories#/urlCategories-get)

The **zia_url_categories_predefined** resource allows you to manage mutable fields of existing predefined URL categories. Predefined URL categories are built-in to the Zscaler platform and cannot be created or deleted — only specific fields can be updated.

This resource is designed for use cases where you need to add custom URLs, keywords, or IP ranges to a predefined category such as `EDUCATION`, `FINANCE`, `CORPORATE_MARKETING`, etc.

~> For managing **custom** URL categories (full CRUD lifecycle), use the [`zia_url_categories`](zia_url_categories) resource instead.

## How This Resource Works

Unlike standard Terraform resources, predefined URL categories have a unique lifecycle:

* **No creation** — The predefined category already exists on the Zscaler platform. The `Create` operation issues a PUT to update the category's mutable fields.
* **No deletion** — Predefined categories cannot be removed. Running `terraform destroy` simply removes the resource from state without making any API calls.
* **Incremental updates** — The provider uses incremental `ADD_TO_LIST` and `REMOVE_FROM_LIST` API operations for all list fields (`urls`, `ip_ranges`, `keywords`, `keywords_retaining_parent_category`, `ip_ranges_retaining_parent_category`). This means the provider compares the current API state against the desired Terraform configuration and issues targeted add/remove calls rather than a full replacement.

⚠️ **IMPORTANT**: The `description` attribute is **not** supported by this resource. Although the API accepts a description in the PUT payload, the value is not returned in the GET response for predefined categories, which would cause persistent state drift. If you need to set a description on a predefined category, use the ZIA Admin Portal directly.

## Example Usage - Adding Custom URLs

```hcl
resource "zia_url_categories_predefined" "education" {
  name = "EDUCATION"
  urls = [
    ".internal-learning.example.com",
    ".corporate-training.example.com",
  ]
}
```

## Example Usage - Adding Keywords and IP Ranges

```hcl
resource "zia_url_categories_predefined" "finance" {
  name = "FINANCE"
  keywords = [
    "internal-trading",
    "corporate-finance",
  ]
  ip_ranges = [
    "10.0.0.0/8",
    "172.16.0.0/12",
  ]
}
```

## Example Usage - URLs Retaining Parent Category

```hcl
resource "zia_url_categories_predefined" "corporate_marketing" {
  name = "CORPORATE_MARKETING"
  urls = [
    ".marketing-internal.example.com",
  ]
  urls_retaining_parent_category = [
    ".brand-portal.example.com",
  ]
  keywords_retaining_parent_category = [
    "brand-assets",
  ]
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The predefined URL category ID or display name (e.g., `FINANCE` or `Finance`). The provider performs a case-insensitive lookup and accepts either the category ID (e.g., `CORPORATE_MARKETING`) or the display name (e.g., `Corporate Marketing`). On the first apply, the provider resolves this to the canonical category ID and stores it in `category_id`.

-> The full list of predefined URL category IDs is available in the [ZIA API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/url-categories/add-custom-category). The provider does not maintain or validate this list locally — refer to the API documentation for the exact category IDs and names.

### Optional

* `urls` - (Optional) Custom URLs to add to the predefined URL category. Up to 25,000 custom URLs can be added per organization across all categories (including bandwidth classes).
* `urls_retaining_parent_category` - (Optional) URLs that are also retained under the original parent URL category.
* `keywords` - (Optional) Custom keywords associated to the URL category. Up to 2,048 custom keywords can be added per organization across all categories (including bandwidth classes).
* `keywords_retaining_parent_category` - (Optional) Retained custom keywords from the parent URL category. Up to 2,048 retained parent keywords can be added per organization across all categories.
* `ip_ranges` - (Optional) Custom IP address ranges associated to the URL category. Up to 2,000 custom IP address ranges and retaining parent custom IP address ranges can be added, per organization, across all categories.

⚠️ **NOTE**: This field is available only if the option to configure custom IP ranges is enabled for your organization. To enable this option, contact Zscaler Support.

* `ip_ranges_retaining_parent_category` - (Optional) The retaining parent custom IP address ranges associated to the URL category. Up to 2,000 custom IP ranges and retaining parent custom IP address ranges can be added, per organization, across all categories.

## Attribute Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the predefined URL category (same as `category_id`).
* `category_id` - The canonical predefined URL category identifier resolved by the provider.
* `configured_name` - The display name of the predefined URL category. Read-only.
* `super_category` - The super category of the predefined URL category. Read-only.
* `url_type` - The URL type (e.g., `EXACT`). Read-only.
* `type` - The type of the URL category. Read-only.
* `val` - (Number) The numeric identifier for the URL category.
* `editable` - (Boolean) Whether the category is editable.
* `db_categorized_urls` - URLs categorized by the Zscaler database. Read-only.
* `custom_urls_count` - The number of custom URLs associated to the URL category. Read-only.
* `urls_retaining_parent_category_count` - The number of URLs retaining the parent category. Read-only.
* `custom_ip_ranges_count` - The number of custom IP address ranges associated to the URL category. Read-only.
* `ip_ranges_retaining_parent_category_count` - The number of IP ranges retaining the parent category. Read-only.

* `url_keyword_counts` - URL and keyword counts for the URL category. Read-only.
  * `total_url_count` - Custom URL count for the category.
  * `retain_parent_url_count` - Count of URLs with retain parent category.
  * `total_keyword_count` - Total keyword count for the category.
  * `retain_parent_keyword_count` - Count of total keywords with retain parent category.

## Destroy Behavior

⚠️ **This resource does not support deletion.** Predefined URL categories are built-in to the Zscaler platform and cannot be removed.

When you run `terraform destroy` or remove this resource from your configuration:

1. The resource is removed from the Terraform state file.
2. **No API call is made** — the predefined category remains unchanged on the Zscaler platform.
3. Any custom URLs, keywords, or IP ranges that were added will **persist** on the predefined category.

If you need to remove custom URLs, keywords, or IP ranges from a predefined category, update the resource to set the fields to empty **before** destroying:

```hcl
resource "zia_url_categories_predefined" "education" {
  name      = "EDUCATION"
  urls      = []
  keywords  = []
  ip_ranges = []
}
```

Then run `terraform apply` to clear the items, followed by `terraform destroy` to remove from state.

## Import

Predefined URL categories can be imported by using the predefined category ID or display name as the import ID.

For example:

```shell
terraform import zia_url_categories_predefined.example EDUCATION
```

or

```shell
terraform import zia_url_categories_predefined.example FINANCE
```

⚠️ **NOTE**: This resource only supports importing **predefined** URL categories. For custom URL categories, use the [`zia_url_categories`](zia_url_categories) resource.

### Important Import Considerations

When importing a predefined category that already has custom URLs, keywords, or IP ranges configured (e.g., via the ZIA Admin Portal), those existing items will be captured in the Terraform state. If your HCL configuration does not include those items, the **next** `terraform apply` will attempt to **remove** them to match the desired state defined in your configuration.

To avoid unintended removals after import, ensure your HCL includes all existing custom items that should be retained, or review the plan output carefully before applying.
