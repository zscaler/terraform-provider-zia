---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_idm_profile_lite"
description: |-
  Get information about ZIA DLP IDM Profile Lite.
---

# Data Source: zia_dlp_idm_profile_lite

Use the **zia_dlp_idm_profile_lite** data source to get summarized information about a ZIA DLP IDM Profile Lite in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP IDM Profile Lite by name
data "zia_dlp_idm_profile_lite" "example"{
    name = "Example"
}
```

```hcl
# Retrieve a DLP IDM Profile Lite by ID
data "zia_dlp_idm_profile_lite" "example"{
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `template_name` - (Required) The IDM template name.

### Optional

* `profile_id` - (Number) The unique identifier for the IDM template (i.e., IDM profile).
* `num_documents` - (Number) The number of documents associated to the IDM template.
* `client_vm` - (Number) This is an immutable reference to an entity. which mainly consists of id and name.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map) The configured name of the entity
* `last_modified_time` - (Number) The date and time the IDM template was last modified.
* `modified_by` - (Number) The admin that modified the IDM template last.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map) The configured name of the entity
