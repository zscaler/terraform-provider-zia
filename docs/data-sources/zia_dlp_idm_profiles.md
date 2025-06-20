---
subcategory: "Data Loss Prevention"
layout: "zscaler"
page_title: "ZIA: dlp_idm_profile"
description: |-
  Official documentation https://help.zscaler.com/zia/data-loss-prevention#/idmprofile-get
  API documentation https://help.zscaler.com/zia/about-indexed-document-match
  Get information about ZIA DLP IDM Profile.
---

# zia_dlp_idm_profile (Data Source)

* [Official documentation](https://help.zscaler.com/zia/data-loss-prevention#/idmprofile-get)
* [API documentation](https://help.zscaler.com/zia/about-indexed-document-match)

Use the **zia_dlp_idm_profile** data source to get information about a ZIA DLP IDM Profile in the Zscaler Internet Access cloud or via the API.

## Example Usage

```hcl
# Retrieve a DLP IDM Profile by name
data "zia_dlp_idm_profile" "example"{
    name = "Example"
}
```

```hcl
# Retrieve a DLP IDM Profile by ID
data "zia_dlp_idm_profile" "example"{
    name = "Example"
}
```

## Argument Reference

The following arguments are supported:

### Required

* `profile_name` - (Required) The IDM template name, which is unique per Index Tool.

### Optional

* `profile_id` - (String) The identifier (1-64) for the IDM template (i.e., IDM profile) that is unique within the organization.
* `profile_desc` - (String) The IDM template's description.
* `profile_type` - (String) The IDM template's type. The returned values are:
  * `LOCAL`
  * `REMOTECRON`
  * `REMOTE`
* `host` - (String) The fully qualified domain name (FQDN) of the IDM template's host machine.
* `port` - (Number) The port number of the IDM template's host machine.
* `profile_dir_path` - (String) The IDM template's directory file path, where all files are present.
* `schedule_type` - (String) The schedule type for the IDM template's schedule (i.e., Monthly, Weekly, Daily, or None). This attribute is required by PUT and POST requests.
  * `NONE`
  * `MONTHLY`
  * `WEEKLY`
  * `DAILY`
* `schedule_day` - (Number) The day the IDM template is scheduled for. This attribute is required by PUT and POST requests.
* `schedule_day_of_month` - (Number) The day of the month that the IDM template is scheduled for. This attribute is required by PUT and POST requests, and when scheduleType is set to MONTHLY.
* `schedule_time` - (Number) The time of the day (in minutes) that the IDM template is scheduled for. For example: at 3am= 180 mins. This attribute is required by PUT and POST requests.
* `schedule_disabled` - (Bool) If set to true, the schedule for the IDM template is temporarily in a disabled state. This attribute is required by PUT requests in order to disable or enable a schedule.
* `upload_status` - (Bool) The status of the file uploaded to the Index Tool for the IDM template.
  * `IDM_PROF_CREATED`
  * `IDM_PROF_PROCESSING`
  * `IDM_PROF_UPLOAD_COMPLETED`
  * `IDM_PROF_VALIDATION_FAILURE`
  * `IDM_PROF_DEL_FAILURE`
  * `IDM_PROF_MD5_CONVERSION_FAILURE`
  * `IDM_PROF_UPLOAD_ERROR`
  * `IDM_PROF_RESV_CRED_ERROR`
  * `IDM_PROF_GENERIC_ERROR`
  * `IDM_PROF_CDSSDLP_COMM_FAIL`
  * `IDM_PROF_CDSSDLP_REJECT`
  * `IDM_PROF_UPLOAD_ABORTED`
  * `IDM_PROF_UPLOAD_RESOURCE_FAIL`

* `username` - (String) The username to be used on the IDM template's host machine.
* `version` - (Number) The version number for the IDM template.
* `idm_client` - (String ) The unique identifer for the Index Tool that was used to create the IDM template. This attribute is required by POST requests, but ignored if provided in PUT requests.
* `volume_of_documents` - (Number) The total volume of all the documents associated to the IDM template.
* `num_documents` - (Number) The number of documents associated to the IDM template.
* `last_modified_time` - (Number) The date and time the IDM template was last modified.
* `modified_by` - (Number) The admin that modified the IDM template last.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `extensions` - (Map) The configured name of the entity
