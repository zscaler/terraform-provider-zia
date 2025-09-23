---
subcategory: "Workload Groups"
layout: "zscaler"
page_title: "ZIA: workload_groups"
description: |-
    Official documentation https://help.zscaler.com/zia/about-workload-groups
    API documentation https://help.zscaler.com/zia/workload-groups#/workloadGroups-get
    Get information about Workload Groups.
---

# zia_workload_groups (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-workload-groups)
* [API documentation](https://help.zscaler.com/zia/workload-groups#/workloadGroups-get)

Use the **zia_workload_groups** resource allows the creation and management of Workload Group objects in the Zscaler Internet Access. This resource can then be used as a criterion in ZIA policies such as, Firewall Filtering, URL Filtering, and Data Loss Prevention (DLP) to apply security policies to the workload traffic.

## Example Usage

```hcl
resource "zia_workload_groups" "example" {
  name = "Test Group"
  description = "Test Group"

  expression_json {
    expression_containers {
      tag_type = "ATTR"
      operator = "AND"

      tag_container {
        operator = "AND"

        tags {
          key   = "GroupName"
          value = "example"
        }
      }
    }

    expression_containers {
      tag_type = "ENI"
      operator = "AND"

      tag_container {
        operator = "AND"

        tags {
          key   = "GroupId"
          value = "123456789"
        }
      }
    }

    expression_containers {
      tag_type = "VPC"
      operator = "AND"

      tag_container {
        operator = "AND"

        tags {
          key   = "Vpc-id"
          value = "vpcid12344"
        }
      }
    }

    expression_containers {
      tag_type = "VM"
      operator = "AND"

      tag_container {
        operator = "AND"

        tags {
          key   = "IamInstanceProfile-Arn"
          value = "test01"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The name of the workload group.

### Optional

* `description` - (Optional) The description of the workload group.
* `expression_json` - (Optional) The workload group expression containing tag types, tags, and their relationships represented in a JSON format.
  * `expression_containers` - (Optional) Contains one or more tag types (and associated tags) combined using logical operators within a workload group.
    * `tag_type` - (Optional) The tag type selected from a predefined list. Supported values are: `ANY`, `VPC`, `SUBNET`, `VM`, `ENI`, `ATTR`.
    * `operator` - (Optional) The operator (either AND or OR) used to create logical relationships among tag types. Supported values are: `AND`, `OR`, `OPEN_PARENTHESES`, `CLOSE_PARENTHESES`.
    * `tag_container` - (Optional) Contains one or more tags and the logical operator used to combine the tags within a tag type.
      * `operator` - (Optional) The logical operator (either AND or OR) used to combine the tags within a tag type. Supported values are: `AND`, `OR`.
      * `tags` - (Optional) One or more tags, each consisting of a key-value pair, selected within a tag type. If multiple tags are present within a tag type, they are combined using a logical operator. Note: A maximum of 8 tags can be added to a workload group, irrespective of the number of tag types present.
        * `key` - (Optional) The key component present in the key-value pair contained in a tag.
        * `value` - (Optional) The value component present in the key-value pair contained in a tag.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - (String) A unique identifier assigned to the workload group.
* `group_id` - (Number) A unique identifier assigned to the workload group.
* `name` - (String) The name of the workload group.
* `description` - (String) The description of the workload group.
* `expression_json` - (List) The workload group expression containing tag types, tags, and their relationships represented in a JSON format.
  * `expression_containers` - (List) Contains one or more tag types (and associated tags) combined using logical operators within a workload group.
    * `tag_type` - (String) The tag type selected from a predefined list. Returned values are: `ANY`, `VPC`, `SUBNET`, `VM`, `ENI`, `ATTR`.
    * `operator` - (String) The operator (either AND or OR) used to create logical relationships among tag types. Returned values are: `AND`, `OR`, `OPEN_PARENTHESES`, `CLOSE_PARENTHESES`.
    * `tag_container` - (List) Contains one or more tags and the logical operator used to combine the tags within a tag type.
      * `operator` - (String) The logical operator (either AND or OR) used to combine the tags within a tag type. Returned values are: `AND`, `OR`.
      * `tags` - (List) One or more tags, each consisting of a key-value pair, selected within a tag type. If multiple tags are present within a tag type, they are combined using a logical operator. Note: A maximum of 8 tags can be added to a workload group, irrespective of the number of tag types present.
        * `key` - (String) The key component present in the key-value pair contained in a tag.
        * `value` - (String) The value component present in the key-value pair contained in a tag.

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_workload_groups** can be imported by using `<GROUP_ID>` or `<GROUP_NAME>` as the import ID.

For example:

```shell
terraform import zia_workload_groups.example <group_id>
```

or

```shell
terraform import zia_workload_groups.example <group_name>
```
