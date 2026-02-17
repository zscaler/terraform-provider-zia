---
subcategory: "Activation"
layout: "zscaler"
page_title: "ZIA Config Activation"
subcategory: "Activation"
---

# Activation Overview

## Who Controls When Changes Are Activated?

**Activation timing is controlled by the ZIA platform, not by the Terraform provider.** When Terraform (or any API client) creates or updates configuration, ZIA saves those changes into a pending state. Those changes do not take effect for traffic until they are *activated*. The provider cannot override this native platform behavior.

Activation can happen in three ways:

1. **Explicit activation** — You activate after your Terraform run (via the options below or via the ZIA Admin UI).
2. **Auto-activation by ZIA** — The platform may activate pending changes automatically after a period of inactivity (e.g. ~15–30 minutes, configurable in Advanced Settings) or when an admin logs out. This applies to all API-based changes, including Terraform.
3. **In-flight activation** — Using the provider option that activates as resources are applied (see below).

If you do not explicitly activate and do not use an in-flight option, your changes may still go live when ZIA auto-activates. That is expected ZIA behavior, not a Terraform bug. To avoid surprises, use one of the provider activation options below or push changes in a controlled fashion when you intend them to take effect.

For the official Zscaler explanation of saving and activating changes (including in the Admin UI), see [Saving and Activating Changes (Admin Console)](https://help.zscaler.com/unified/saving-and-activating-changes-admin-console).

## Activation Options with the Terraform Provider

The provider supports three ways to activate changes:

| Method | Description |
|--------|-------------|
| **`zia_activation_status` resource** | Use the `zia_activation_status` data source or resource with Terraform `depends_on` listing the resources you configure. This ensures activation runs after those resources are applied. Not the most optimal but doable. |
| **Environment variable `ZIA_ACTIVATION`** | Set `ZIA_ACTIVATION` in the environment when running Terraform. The provider will activate changes as resources are configured in-flight. |
| **Out-of-band (recommended)** | Run the `ziaActivator` CLI after `terraform apply` (or as part of your pipeline). This gives you explicit control over when activation happens and is the recommended approach. |

The rest of this guide describes the **recommended** out-of-band method using the `ziaActivator` binary.

## Building and Using the Out-of-Band Activator (ziaActivator)

Terraform does not provide built-in support for commits or post-activation configuration, so many users handle activation out-of-band with the `ziaActivator` CLI.

The activation CLI source is built alongside the provider. You can install the binary into your `$PATH` (e.g. `$HOME/bin`). By default it is installed to `/usr/local/bin/`. Build and install with:

```bash
make build13 && sudo make ziaActivator
```

~> You may or may not need `sudo` to install `ziaActivator` into your path.

Example: run activation after apply in one shot:

```bash
terraform init && terraform apply && ziaActivator
```

The authentication credentials can be given multiple ways, and if all are present then this is the order, from highest to lowest priority:

!> **WARNING:** Providing authentication credentials via CLI argument is insecure and
is not recommended.

1. CLI arguments
2. Environment variables

Refer to the ZIA provider argument reference documentation for more information on the JSON config file and the environment variables that are used.

!> **WARNING:** The ZIA platform has its own auto-activation behavior, independent of the Terraform provider. Pending changes may be activated automatically when: (1) the session has been inactive for a configurable period (e.g. 30 minutes, see Advanced Settings), or (2) an admin logs out. This applies to all API-based changes, including Terraform. If you do not want changes to go live until you decide, use one of the [activation options](#activation-options-with-the-terraform-provider) and activate explicitly, or push changes only when you intend them to take effect.

## FAQ: "A change was pushed without us activating it"

If you see configuration take effect even though you did not run an explicit activation step, that is normal ZIA behavior. The platform can activate pending changes automatically (inactivity timeout or logout). The Terraform provider does not control when ZIA activates; it only writes configuration. To avoid unintended activation:

- Use the **out-of-band** `ziaActivator` (recommended) and run it only when you are ready, or
- Use **`ZIA_ACTIVATION`** if you want activation to happen during the same Terraform run, or
- Use **`zia_activation_status`** with `depends_on` so activation runs after your resources.

The provider cannot override ZIA’s native behavior. For the platform’s own description of save vs. activate behavior, see [Zscaler Help: Saving and Activating Changes](https://help.zscaler.com/unified/saving-and-activating-changes-admin-console).
