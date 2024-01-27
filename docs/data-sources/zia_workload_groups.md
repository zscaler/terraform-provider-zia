---
subcategory: "Workload Groups"
layout: "zscaler"
page_title: "ZIA: workload_groups"
description: |-
  Get information about Workload Groups.
---

# Data Source: zia_workload_groups

Use the **zia_workload_groups** data source to get information about Workload Groups in the Zscaler Internet Access cloud or via the API. This data source can then be used as a criterion in ZIA policies such as, Firewall Filtering, URL Filtering, and Data Loss Prevention (DLP) to apply security policies to the workload traffic.

## Example Usage

```hcl
# ZIA Admin User Data Source
data "zia_workload_groups" "ios"{
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the workload group to be exported.
* `id` - (Optional) The unique identifer for the workload group.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `description` - (String) The description of the workload group.
* `expression_json` - (List) The workload group expression containing tag types, tags, and their relationships represented in a JSON format.
  * `expression_containers` - (List) Contains one or more tag types (and associated tags) combined using logical operators within a workload group
    * `tag_type` - (String) The tag type selected from a predefined list. Returned values are: ``ANY``, ``VPC``, ``SUBNET``, ``VM``, ``ENI``, ``ATTR``
    * `operator` - (String) The operator (either AND or OR) used to create logical relationships among tag types. Returned values are: ``AND``, ``OR``, ``OPEN_PARENTHESES``, ``CLOSE_PARENTHESES``
    * `tag_container` - (String) Contains one or more tags and the logical operator used to combine the tags within a tag type ``CLOSE_PARENTHESES``
      * `tags` - (String) One or more tags, each consisting of a key-value pair, selected within a tag type. If multiple tags are present within a tag type, they are combined using a logical operator. Note: A maximum of 8 tags can be added to a workload group, irrespective of the number of tag types present.
        * `key` - (String) The key component present in the key-value pair contained in a tag
        * `value` - (string) The value component present in the key-value pair contained in a tag
    * `operator` - (String) The operator (either AND or OR) used to create logical relationships among tag types. Returned values are: ``AND``, ``OR``, ``OPEN_PARENTHESES``, ``CLOSE_PARENTHESES``
* `expression` - (String) The workload group expression containing tag types, tags, and their relationships.
* `last_modified_time` - (Number) When the rule was last modified
* `last_modified_by`
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` -(String) The configured name of the entity
  * `extensions` - (Map of String)
