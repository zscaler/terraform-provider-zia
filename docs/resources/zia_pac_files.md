---
subcategory: "PAC Files"
layout: "zscaler"
page_title: "ZIA: pac_files"
description: |-
  Official documentation https://help.zscaler.com/zia/about-hosted-pac-files
  API documentation https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/pac-files/pac-resource-add-pac-version
  Creates and manages ZIA hosted PAC files and their version lifecycle.
---

# zia_pac_files (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-hosted-pac-files)
* [API documentation](https://automate.zscaler.com/docs/api-reference-and-guides/api-reference/zia/pac-files/pac-resource-add-pac-version)

Use the **zia_pac_files** resource to create and manage hosted PAC files in the ZIA Admin Portal, including their full version lifecycle: deploying, staging, unstaging, and marking a version as last known good.

## How PAC File Versioning Works

Understanding the service's versioning model is essential to using this resource correctly:

* A PAC file is a container of **versions**. Each version carries its own content and commit message. A version's content can never be edited in place — changing content always means creating a **new version**.
* Each version holds at most one status: **`DEPLOYED`** (live), **`STAGE`** (staged for deployment), **`LKG`** (marked as last known good), or no status at all (**`UNSTAGED`**).
* Only one version can be deployed at a time, and only one can be staged at a time. Deploying a version automatically demotes the previously deployed version to unstaged.
* The **initial version is always created in the deployed state** — the service does not allow the first version of a PAC file to be staged.
* The currently **deployed version cannot be staged or unstaged** — there must always be a deployed version.
* PAC file content is validated by the service **before** the PAC file is created or a new version is saved. If the content contains syntax errors, the apply fails with the validation messages returned by the service.
* The service limits the number of versions a PAC file can hold (10). See `delete_version` below.

This resource manages **one PAC file and one of its versions at a time** — by default the latest version it created, recorded in the `pac_version` attribute. Your configuration declares three things that must always agree: the version being managed, its content, and its status.

## Managing Versions and Statuses — Quick Reference

| To achieve | Configure |
|---|---|
| Create the PAC file (version 1, deployed) | `pac_content` + `pac_version_status = "DEPLOYED"` (the only status allowed at creation) |
| Save a **new version** and deploy it | Change `pac_content`, keep `pac_version_status = "DEPLOYED"`, do **not** set `pac_version` |
| Save a **new version** and stage it | Change `pac_content`, set `pac_version_status = "STAGE"`, do **not** set `pac_version` |
| Deploy / unstage / mark the managed version | Change **only** `pac_version_status` (optionally `pac_commit_message`) — no new version is created |
| Act on a **specific existing version** | Set `pac_version = <N>` + the desired `pac_version_status`, with `pac_content` matching version N exactly |

Two rules tie the table together:

1. **A changed `pac_content` creates a new version; a changed `pac_version_status` transitions an existing one.** Changing both in one apply saves the new version first and then applies the status to it.
2. **`pac_version` is a selector, never a creator.** Declaring it targets an existing version, and your `pac_content` must match that version's content exactly — the plan describes version N, so the configuration must agree with it. To create a version, omit `pac_version`.

## Example Usage - Create the PAC File

The initial version is always deployed. Declaring any other status at creation returns an error.

```hcl
resource "zia_pac_files" "this" {
  name               = "example_pac_file"
  description        = "Example hosted PAC file"
  domain             = "acme.com"
  pac_commit_message = "Initial version"
  pac_version_status = "DEPLOYED"
  pac_content        = <<-EOT
function FindProxyForURL(url, host) {
    return "PROXY $${GATEWAY_FX}:80; PROXY $${SECONDARY_GATEWAY_FX}:80; DIRECT";
}
EOT
}
```

## Example Usage - Save a New Version in Staged Mode

Change the content and declare `STAGE`. The change is saved as a new version (e.g. version 2) and staged; the previously deployed version remains live.

```hcl
resource "zia_pac_files" "this" {
  name               = "example_pac_file"
  description        = "Example hosted PAC file"
  domain             = "acme.com"
  pac_commit_message = "Staged change"
  pac_version_status = "STAGE"
  pac_content        = <<-EOT
function FindProxyForURL(url, host) {
    /* updated logic */
    return "PROXY $${GATEWAY_FX}:80; DIRECT";
}
EOT
}
```

## Example Usage - Deploy the Staged Version

To promote the staged version, change **only** the status. No new version is created; the staged version becomes deployed and the previously deployed version becomes unstaged.

```hcl
resource "zia_pac_files" "this" {
  # ... name, domain, and pac_content unchanged from the staged version ...
  pac_commit_message = "Deploying staged change"
  pac_version_status = "DEPLOYED"
}
```

The same pattern applies to the other in-place transitions on the managed version:

* `pac_version_status = "UNSTAGED"` — removes the version's current status (removes staged status from a staged version, or the last-known-good marking from an LKG version) without deploying it.
* `pac_version_status = "LKG"` — marks the version as last known good.

## Example Usage - Promote a Specific Older Version

Declare `pac_version` to act on an existing version that is not the one currently managed — for example, re-staging version 1 after version 2 took over deployment. `pac_content` must match the declared version's content exactly.

```hcl
resource "zia_pac_files" "this" {
  name               = "example_pac_file"
  description        = "Example hosted PAC file"
  domain             = "acme.com"
  pac_version        = 1
  pac_commit_message = "Re-staging version 1"
  pac_version_status = "STAGE"
  pac_content        = <<-EOT
function FindProxyForURL(url, host) {
    return "PROXY $${GATEWAY_FX}:80; PROXY $${SECONDARY_GATEWAY_FX}:80; DIRECT";
}
EOT
}
```

A subsequent apply with `pac_version_status = "DEPLOYED"` (keeping `pac_version = 1`) deploys that version.

~> **NOTE:** `pac_version` cannot be combined with a content change, and cannot be set when creating the PAC file. If the declared content does not match the declared version, the apply fails with an error explaining the mismatch.

## Example Usage - Version Limit Reached

When the PAC file already holds the maximum number of versions, saving a new version fails. Name a version to remove with `delete_version`; it is used only for that save and can be removed from the configuration afterwards.

```hcl
resource "zia_pac_files" "this" {
  # ... changed pac_content ...
  pac_version_status = "DEPLOYED"
  delete_version     = 3
}
```

## Argument Reference

The following arguments are supported:

### Required

* `name` - (String) The name of the PAC file
* `pac_content` - (String) The content of the PAC file. The content is validated before the PAC file is created or a new version is saved. When `pac_version` is declared, the content must match that version's content exactly.

### Optional

* `description` - (String) The description of the PAC file
* `domain` - (String) The domain of your organization to which the PAC file applies
* `pac_commit_message` - (String) The commit message for the managed PAC file version. Changing only the commit message does not create a new version.
* `pac_url_obfuscated` - (Boolean) Indicates whether the PAC file URL is obfuscated. When true, the obfuscated URL is returned in the `pac_sub_url` attribute.
* `pac_version_status` - (String) The desired status of the managed PAC file version. The initial version is always deployed; `STAGE`, `LKG`, and `UNSTAGED` can be applied to subsequent versions. `UNSTAGED` removes the version's current status (staged or last known good) without deploying it. Supported values: `DEPLOYED`, `STAGE`, `LKG`, `UNSTAGED`. Defaults to `DEPLOYED`.
* `pac_version` - (Number) The PAC file version to apply `pac_version_status` to. When omitted, the resource manages the latest version it created. Declare it only for status transitions on an existing version — see "Managing Versions and Statuses" above.
* `delete_version` - (Number) The version to remove when saving a new version would exceed the maximum number of PAC file versions. Only used while saving a new version.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `pac_id` - (Number) The unique identifier for the PAC file
* `pac_version` - (Number) The PAC file version currently managed by this resource (also assignable — see Argument Reference)
* `pac_url` - (String) The URL location of the PAC file, auto-generated when the PAC file is first added
* `pac_sub_url` - (String) The obfuscated URL of the PAC file. Returned when `pac_url_obfuscated` is true.
* `pac_verification_status` - (String) The verification status of the PAC file. Supported values: `VERIFY_NOERR`, `VERIFY_ERR`, `NOVERIFY`
* `editable` - (Boolean) Indicates whether the PAC file is editable

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_pac_files** can be imported by using `<PAC_FILE_ID>` or `<PAC_FILE_NAME>` as the import ID.

For example:

```shell
terraform import zia_pac_files.example <pac_file_id>
```

or

```shell
terraform import zia_pac_files.example <pac_file_name>
```

Importing adopts the currently deployed version of the PAC file as the version managed by Terraform.
