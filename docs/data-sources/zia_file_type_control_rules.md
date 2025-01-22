---
subcategory: "File Type Control Policy"
layout: "zscaler"
page_title: "ZIA: file_type_control_rules"
description: |-
  Retrieves all the rules in the File Type Control policy.
---

# Data Source: zia_file_type_control_rules

Use the **zia_file_type_control_rules** data source to retrieves File Type Control rules.

## Example Usage

```hcl
# Retrieve a File Type Control Rule by name
data "zia_file_type_control_rules" "this" {
    name = "Example"
}
```

```hcl
# Retrieve a File Type Control Rule by ID
data "zia_file_type_control_rules" "this" {
    name = "12134558"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (Required) The DLP policy rule name.
* `id` - (Required) Optional rules.

### Optional

* `description` - (String) Additional information about the File Type rule.
* `state` - (String) Rule State.
* `order` - (Integer) Order of policy execution with respect to other file-type policies.
* `filtering_action` - (String) Action taken when traffic matches policy. Supported values: "BLOCK", "CAUTION", "ALLOW".
* `time_quota` - (Integer) Time quota in minutes, after which the policy must be applied. Ignored if action is BLOCK.
* `size_quota` - (Integer) Size quota in KB, beyond which the policy must be applied. Ignored if action is BLOCK.
* `access_control` - (String) The access privilege for this DLP policy rule based on the admin's state.
* `rank` - (Integer) Admin rank of the admin who creates this rule.
* `capture_pcap` - (Boolean) A Boolean value that indicates whether packet capture (PCAP) is enabled.
* `operation` - (String) File operation performed.
* `active_content` - (Boolean) Flag to check whether a file has active content.
* `unscannable` - (Boolean) Flag to check whether a file is scannable.
* `cloud_applications` - (List of Strings) The list of cloud applications to which the File Type Control policy rule must be applied.
* `file_types` - (List of Strings) The list of file types to which the Sandbox Rule must be applied.
* `min_size` - (Integer) Minimum file size (in KB) used for evaluation of the rule.
* `max_size` - (Integer) Maximum file size (in KB) used for evaluation of the rule.
* `protocols` - (List of Strings) Protocol for the given rule.
* `url_categories` - (List of Strings) The list of URL categories to which the DLP policy rule must be applied.
* `last_modified_time` - (Integer) When the rule was last modified.
* `last_modified_by` - (Object) Who modified the rule last. Contains Name-ID pairs.
* `locations` - (List of Objects) Name-ID pairs of locations for which rule must be applied.
* `location_groups` - (List of Objects) Name-ID pairs of the location groups to which the rule must be applied.
* `groups` - (List of Objects) Name-ID pairs of groups for which rule must be applied.
* `departments` - (List of Objects) Name-ID pairs of departments for which rule must be applied.
* `users` - (List of Objects) Name-ID pairs of users for which rule must be applied.
* `time_windows` - (List of Objects) Name-ID pairs of time intervals during which rule must be enforced.
* `labels` - (List of Objects) The URL Filtering rule's label.
* `device_groups` - (List of Objects) Applicable for devices managed using Zscaler Client Connector.
* `devices` - (List of Objects) Name-ID pairs of devices for which rule must be applied.
* `device_trust_levels` - (List of Strings) List of device trust levels for which the rule must be applied.
* `zpa_app_segments` - (List of Objects) The list of ZPA Application Segments for which this rule is applicable.
