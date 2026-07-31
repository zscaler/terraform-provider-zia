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
2. **Auto-activation by ZIA** — The platform activates pending changes automatically when a session ends, either because an admin logged out or because an API session reached its timeout. The API session timeout is set by `api_session_timeout` in Advanced Settings and can range from **5 to 20 minutes, defaulting to 5**. This applies to all API-based changes, including Terraform. See [API Session Timeout and Long-Running Applies](#api-session-timeout-and-long-running-applies) below.
3. **In-flight activation** — Using the provider option that activates as resources are applied (see below).

If you do not explicitly activate and do not use an in-flight option, your changes may still go live when ZIA auto-activates. That is expected ZIA behavior, not a Terraform bug. To avoid surprises, use one of the provider activation options below or push changes in a controlled fashion when you intend them to take effect.

For the official Zscaler explanation of saving and activating changes (including in the Admin UI), see [Saving and Activating Changes (Admin Console)](https://help.zscaler.com/unified/saving-and-activating-changes-admin-console).

## Activation Options with the Terraform Provider

Activation is a tenant-wide operation, not a per-resource one. **A single activation call publishes every pending change in the tenant**, regardless of how many resources Terraform created, updated, or deleted. The provider supports three ways to make that call:

| Method | Description |
|--------|-------------|
| **Out-of-band (recommended)** | Run the `ziaActivator` CLI after `terraform apply`, or as a dedicated stage in your pipeline. Issues exactly one activation call per run and keeps the timing explicitly under your control. |
| **`zia_activation_status` resource** | Declare the resource with Terraform `depends_on` listing the resources you configure, so activation runs after they are applied. Workable, but every resource must be listed explicitly — anything you forget may be applied after activation and left pending, which makes this hard to maintain in module-based configurations. |
| **Environment variable `ZIA_ACTIVATION`** | Set `ZIA_ACTIVATION=true` and the provider activates in-flight as resources are configured. **Discouraged**, and may be removed in a future major release: it fires one activation call per resource change against an endpoint limited to **10 POST requests per minute and 40 per hour** (see the [API Rate Limit Summary](https://help.zscaler.com/legacy-apis/api-rate-limit-summary)), so any configuration of meaningful size exhausts that budget and spends the rest of the run waiting on retries. |

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

!> **WARNING:** The ZIA platform has its own auto-activation behavior, independent of the Terraform provider. Pending changes are activated automatically when a session ends — either because an API session reached the `api_session_timeout` configured in Advanced Settings (5 to 20 minutes, default 5), or because an admin logged out. This applies to all API-based changes, including Terraform. If you do not want changes to go live until you decide, use one of the [activation options](#activation-options-with-the-terraform-provider) and activate explicitly, or push changes only when you intend them to take effect.

## API Session Timeout and Long-Running Applies

The lifetime of an API-initiated session is controlled by the API session timeout in Advanced Settings. It accepts **5 to 20 minutes and defaults to 5**. See [Configuring Advanced Settings](https://help.zscaler.com/zia/configuring-advanced-settings#session-timeout) for the platform's description of this setting.

Because the platform activates pending changes when a session ends, a Terraform run that lasts longer than the configured timeout will cross that boundary and have its pending changes activated part-way through — before the apply has finished writing the rest of the configuration. The provider then establishes a new session and continues, so the run does not fail, but activation has already happened at a point you did not choose. A 30-minute run against the default 5-minute timeout crosses this boundary several times.

Raise the timeout to its maximum of 20 minutes before running large configurations. This can be done in either of two ways.

### In the ZIA Admin Portal

Navigate to **Administration > Advanced Settings** and set **API Session Timeout Duration (In Minutes)** under *Admin Portal Session Timeout*. This is a separate field from the UI session timeout shown above it on the same page:

![ZIA Admin Portal — Advanced Settings, Admin Portal Session Timeout](https://raw.githubusercontent.com/zscaler/terraform-provider-zia/master/docs/guides/media/advanced_settings.png)

### With Terraform

Set the `api_session_timeout` attribute on the `zia_advanced_settings` resource:

```hcl
resource "zia_advanced_settings" "this" {
  api_session_timeout = 20
  # … remaining advanced settings attributes
}
```

~> **NOTE:** Raising the timeout reduces how often a long run crosses a session boundary, but it does not eliminate the behaviour — activation on session end is native to the platform and cannot be disabled or overridden by the provider. 20 minutes is a ceiling, not a guarantee, so for very large configurations also split the work across smaller states such that no single run needs to outlive a session.

## Several Terraform Configurations Against One Tenant

Splitting a large deployment across multiple Terraform configurations is a sound practice. Splitting the state, however, does not split the tenant, and activation is where that distinction matters most.

**Activation is not scoped to a configuration.** A single activation call publishes every pending change in the tenant. It also does not necessarily publish immediately: if another administrator or API session still has unactivated changes, the activation is **queued** until all editing administrators have activated, and a queued activation cannot be cancelled. See [Saving and Activating Changes](https://help.zscaler.com/legacy-zia/saving-and-activating-changes-admin-portal) for the platform's description of activation queuing.

For configurations that run concurrently, this has two consequences:

- Activating at the end of each configuration does not publish that configuration's changes independently. The activation waits for every other session, so publication timing is set by whichever run finishes last.
- When the queue clears, all pending tenant changes are published — including those written by a run still in progress, or one that failed part-way through. For rule-based resources, a partially created rule set can take effect with incomplete ordering.

Separately, the platform serialises configuration changes at the tenant level using a single write lock. Concurrent runs contend for it, and requests that cannot acquire it are rejected with `EDIT_LOCK_NOT_AVAILABLE`, `Resource Access Blocked`, or `Failed during enter Org barrier`. The provider retries these automatically, so contention typically appears as a progressively slower run rather than an immediate failure.

The recommended pipeline shape is therefore:

1. **Plan in parallel.** Read operations do not take the write lock, so planning several configurations concurrently is safe.
2. **Apply one configuration at a time** per tenant, using a CI concurrency control keyed on the **tenant** rather than the repository or workspace (`concurrency` in GitHub Actions, `resource_group` in GitLab CI, run queuing in HCP Terraform). Avoid cancelling an in-progress apply, which is itself a cause of partially applied configuration.
3. **Run `ziaActivator` once**, as a final pipeline stage after the last apply — not once per configuration.

~> **NOTE:** Do not lower Terraform's `-parallelism` in response to lock contention. It does not reduce the number of competing runs, and the longer run time increases the chance of crossing an API session boundary and triggering an unplanned activation. Reduce the number of concurrent runs instead.

If truly concurrent execution is required, a separate tenant is the only boundary that provides it, as each tenant has its own write lock and activation queue.

## FAQ: "We ran activation but nothing went live"

Check whether another session was editing the tenant at the time. When any administrator has unactivated changes, your activation is queued rather than published, and it stays queued until every editing administrator has activated. The ZIA Admin Portal lists the administrators concerned under **Queued Activations**, which is the quickest way to confirm this. Super administrators can use **Force Activate** to publish immediately, but note that this publishes *all* saved changes in the tenant, including any that are still mid-deployment.

The durable fix is to serialise applies and activate once at the end of the pipeline, as described in [Several Terraform Configurations Against One Tenant](#several-terraform-configurations-against-one-tenant).

## FAQ: "A change was pushed without us activating it"

If you see configuration take effect even though you did not run an explicit activation step, that is normal ZIA behavior. The platform can activate pending changes automatically (inactivity timeout or logout). The Terraform provider does not control when ZIA activates; it only writes configuration. To avoid unintended activation:

- Use the **out-of-band** `ziaActivator` (recommended) and run it only when you are ready, or
- Use **`zia_activation_status`** with `depends_on` so activation runs after your resources.

Keeping individual runs short, and raising `api_session_timeout` to 20 minutes, both reduce how often the platform activates on its own.

The provider cannot override ZIA’s native behavior. For the platform’s own description of save vs. activate behavior, see [Zscaler Help: Saving and Activating Changes](https://help.zscaler.com/unified/saving-and-activating-changes-admin-console).
