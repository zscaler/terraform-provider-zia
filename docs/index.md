---
layout: "zscaler"
page_title: "Provider: Zscaler Internet Access (ZIA)"
description: |-
    The Zscaler Internet Access provider is used to interact with Zscaler Internet Access (ZIA) API
---

# Zscaler Internet Access (ZIA) Provider

The Zscaler Internet Access provider is used to interact with ZIA API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, Cloud Firewall Policies, DLP Dictionaries, Local Accounts etc. The provider is intended to save time and reducing configuration errors. With this ZIA provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZIA API credentials.

Use the navigation on the left to read about the available resources.

## Support Disclaimer

-> **Disclaimer:** Please refer to our [General Support Statement](guides/support.md) before proceeding with the use of this provider. You can also refer to our [troubleshooting guide](guides/troubleshooting.md) for guidance on typical problems.

## Feature Availability and API Parity

-> **Important:** The ZIA Terraform provider maintain parity with publicly available API endpoints. In some instances, certain features or attributes available via the Zscaler UI may not be immediately available through the API, and therefore cannot be included in the Terraform provider. This does not indicate that the provider is lagging behind; rather, it reflects that we implement only the features that are currently exposed by the public API.

If there is a feature or attribute you would like to see included in the provider, you are welcome to:

- Submit a feature request via [GitHub Issues](https://github.com/zscaler/terraform-provider-zia/issues)
- Contact Zscaler Global Support by opening a support ticket

Our team continuously works with product teams to expand API coverage and will incorporate new features into the provider as they become publicly available through the API.

## Zscaler OneAPI New Framework

The ZIA Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE**: Attention Government customers. OneAPI and Zidentity now support the government (FedRAMP) clouds via the unified `cloud=gov` and `cloud=govus` values. See the [OneAPI Government (FedRAMP) Cloud Environments](#oneapi-government-fedramp-cloud-environments) section below for details.

## Examples Usage - Client Secret Authentication

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "~> 4.0.0"
            source = "zscaler/zia"
        }
    }
}

# Configure the ZIA Provider (OneAPI Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zia" {
  client_id = "[ZSCALER_CLIENT_ID]"
  client_secret = "[ZSCALER_CLIENT_SECRET]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
}
```

## Examples Usage - Private Key Authentication

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "~> 4.0.0"
            source = "zscaler/zia"
        }
    }
}

# Configure the ZIA Provider (OneAPI Authentication) - Private Key
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zia" {
  client_id     = "[ZSCALER_CLIENT_ID]"
  private_key   = "[ZSCALER_PRIVATE_KEY]"
  vanity_domain = "[ZSCALER_VANITY_DOMAIN]"
  zscaler_cloud = "[ZSCALER_CLOUD]"
}
```

**NOTE**: The `zscaler_cloud` is optional for production comercial Clouds and ONLY required when authenticating to other environments. Currently the only supported value are:

- Production Commericial Clouds: The cloud parameter `IS NOT` required.
- Test (Beta) Commercial Clouds: `beta`
- FedRAMP Clouds: `gov` or `govus`

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For the resources and data sources examples, please check the [examples](https://github.com/zscaler/terraform-provider-zia/tree/master/examples) directory.

## Authentication - OneAPI New Framework

As of version v4.0.0, this provider supports authentication via the new Zscaler API framework [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)

Zscaler OneAPI uses the OAuth 2.0 authorization framework to provide secure access to Zscaler Internet Access (ZIA) APIs. OAuth 2.0 allows third-party applications to obtain controlled access to protected resources using access tokens. OneAPI uses the Client Credentials OAuth flow, in which client applications can exchange their credentials with the authorization server for an access token and obtain access to the API resources, without any user authentication involved in the process.

- [ZIA API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=managed%20using%20OneAPI.-,ZIA%20API,-Zscaler%20Internet%20Access)

### Default Environment variables

You can provide credentials via the `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET`, `ZSCALER_VANITY_DOMAIN`, `ZSCALER_CLOUD` environment variables, representing your Zidentity OneAPI credentials `clientId`, `clientSecret`, `vanityDomain` and `zscaler_cloud` respectively.

| Argument        | Description                                                                                         | Environment Variable     |
|-----------------|-----------------------------------------------------------------------------------------------------|--------------------------|
| `client_id`     | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.         | `ZSCALER_CLIENT_ID`      |
| `client_secret` | _(String)_ Secret key associated with the API Client ID for authentication.                         | `ZSCALER_CLIENT_SECRET`  |
| `privateKey`    | _(String)_ A string Private key value.                                                              | `ZSCALER_PRIVATE_KEY`    |
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization.                                     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud`         | _(String)_ The name of the Zidentity cloud, e.g., beta.                                             | `ZSCALER_CLOUD`          |

### Alternative OneAPI Cloud Environments

OneAPI supports authentication and can interact with alternative Zscaler enviornments i.e `beta`. To authenticate to these environments you must provide the following values:

| Argument        | Description                                                         | Environment Variable     |
|-----------------|---------------------------------------------------------------------|--------------------------|
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud` | _(String)_ The name of the Zidentity cloud i.e beta                | `ZSCALER_CLOUD`          |

For example: Authenticating to Zscaler Beta environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="beta"
```

## OneAPI Government (FedRAMP) Cloud Environments

OneAPI supports the Zscaler government (FedRAMP) clouds. These are FedRAMP-isolated environments served by a dedicated Zidentity identity provider and API gateway. To authenticate, set the `cloud` attribute (or `ZSCALER_CLOUD` environment variable) to one of the supported government values:

| Argument        | Description                                                         | Environment Variable     |
|-----------------|---------------------------------------------------------------------|--------------------------|
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization     | `ZSCALER_VANITY_DOMAIN`  |
| `zscaler_cloud` | _(String)_ Supported Zidentity Gov Cloud `gov` or `govus`                | `ZSCALER_CLOUD`          |

**NOTE** FedRAMP cloud is only supported when upgrading to the provider version `>=v4.7.25`. Earlier versions are NOT supported.

For example, authenticating to the GOV environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="gov"
```

### OneAPI (API Client Scope)

OneAPI Resources are automatically created within the ZIdentity Admin UI based on the RBAC Roles
applicable to APIs within the various products. For example, in ZIA, navigate to `Administration -> Role
Management` and select `Add API Role`.

Once this role has been saved, return to the ZIdentity Admin UI and from the Integration menu
select API Resources. Click the `View` icon to the right of Zscaler APIs and under the ZIA
dropdown you will see the newly created Role. In the event a newly created role is not seen in the
ZIdentity Admin UI a `Sync Now` button is provided in the API Resources menu which will initiate an
on-demand sync of newly created roles.

### Important Note - New Feature

- [API Session Timeout](https://help.zscaler.com/zia/release-upgrade-summary-2026#:~:text=Feature%20Available-,API%20Session%20Timeout,-When%20configuring%20advanced) - A new field, `api_session_timeout`, is available for the AdvancedSettings model in the /advancedSettings APIs. This configuration allows you to specify how long API-initiated sessions can be inactive before they are forced to reauthenticate. The timeout duration can range from 5 to 20 minutes. The attribute `api_session_timeout` is available via the resource `zia_advanced_settings`

  ~> **NOTE:** This setting directly affects long-running Terraform applies, because the platform activates pending changes when a session ends. See [API Session Timeout and Long-Running Applies](#api-session-timeout-and-long-running-applies).

## Legacy API Framework

### ZIA native authentication

- As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

### Examples Usage

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "~> 4.0.0"
            source = "zscaler/zia"
        }
    }
}

# Configure the ZIA Provider (Legacy Authentication)
#
# NOTE: Change place holder values denoted by brackets to real values, including
# the brackets.
#
# NOTE: If environment variables are utilized for provider settings the
# corresponding variable name does not need to be set in the provider config
# block.
provider "zia" {
  username            = "[ZIA_USERNAME]"
  password            = "[ZIA_PASSWORD]"
  api_key             = "[ZIA_API_KEY]"
  zia_cloud           = "[ZIA_CLOUD]"
  use_legacy_client   = "[ZSCALER_USE_LEGACY_CLIENT]"
}
```

The ZIA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

- `zscaler`
- `zscloud`
- `zscalerone`
- `zscalertwo`
- `zscalerthree`
- `zscalerbeta`
- `zscalergov`
- `zscalerten`
- `zspreview`

### Environment variables

You can provide credentials via the `ZIA_USERNAME`, `ZIA_PASSWORD`, `ZIA_API_KEY`, `ZIA_CLOUD`, `ZSCALER_USE_LEGACY_CLIENT` environment variables, representing your ZIA `username`, `password`, `api_key`,  `zia_cloud` and `use_legacy_client` respectively.

| Argument            | Description                                                                                                         | Environment variable        |
|---------------------|---------------------------------------------------------------------------------------------------------------------|------------------------------|
| `username`          | _(String)_ A string that contains the email ID of the API admin.                                                     | `ZIA_USERNAME`              |
| `password`          | _(String)_ A string that contains the password for the API admin.                                                    | `ZIA_PASSWORD`              |
| `api_key`           | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).  | `ZIA_API_KEY`               |
| `zia_cloud`         | _(String)_ The host and basePath for the cloud services API is `$zsapi.<Zscaler Cloud Name>/api/v1`.                | `ZIA_CLOUD`                 |
| `use_legacy_client` | _(Bool)_ Enable use of the legacy ZIA API Client.                                                                    | `ZSCALER_USE_LEGACY_CLIENT` |

```sh
# Change place holder values denoted by brackets to real values, including the
# brackets.

$ export ZIA_USERNAME="[ZIA_USERNAME]"
$ export ZIA_PASSWORD="[ZIA_PASSWORD]"
$ export ZIA_API_KEY="[ZIA_API_KEY]"
$ export ZIA_CLOUD="[ZIA_CLOUD]"
$ export ZSCALER_USE_LEGACY_CLIENT=true
$ terraform plan
```

If you are on Windows, use PowerShell to set the environmenr variables using the following commands:

```pwsh
$env:username = 'xxxxxxxxxxxxxxxx'
$env:password = 'xxxxxxxxxxxxxxxx'
$env:api_key = 'xxxxxxxxxxxxxxxx'
$env:zia_cloud = '<zscaler_cloud_name>'
$env:use_legacy_client = true
```

```hcl
# provider settings established with values from environment variables
provider "zia" {}
```

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For details about how to retrieve your tenant Base URL and API key/token refer to the Zscaler help portal. <https://help.zscaler.com/zia/getting-started-zia-api>

## ZIA Configuration Activation

The ZIA platform saves configuration changes in a **pending** state. A successful `terraform apply` writes the configuration to the tenant but does not put it into effect for traffic — the changes must be **activated**. See [Saving and Activating Changes](https://help.zscaler.com/unified/saving-and-activating-changes-admin-console) for the platform's own description of this behaviour.

Activation is a tenant-wide operation, not a per-resource one. **A single activation call publishes every pending change in the tenant**, regardless of how many resources Terraform created, updated, or deleted.

The provider supports three activation methods, described below in order of preference.

### Method 1 — Out-of-band activation with the `ziaActivator` CLI (recommended)

Compile the dedicated activation binary and run it once, after Terraform has finished:

```bash
make build13 && sudo make ziaActivator
```

```bash
terraform apply && ziaActivator
```

This is the recommended method for CI/CD pipelines and for large configurations. It issues exactly one activation call for the entire run, keeps activation timing explicitly under your control, and adds no activation traffic to the apply itself. The activator reads the same credentials and environment variables as the provider, so no additional configuration is required.

For build and usage details, see the [ZIA Activator](guides/zia-activator-overview.md) guide.

### Method 2 — The `zia_activation_status` resource

Declare the resource in your configuration and use the `depends_on` meta-argument so that it runs after the resources it should activate:

```hcl
resource "zia_activation_status" "this" {
  status = "ACTIVE"

  depends_on = [
    zia_url_filtering_rules.block_gambling,
    zia_firewall_filtering_rule.allow_engineering,
  ]
}
```

This keeps activation inside the Terraform run, but it carries an important limitation: `depends_on` cannot be inferred, so **every resource whose changes must be activated has to be listed explicitly**. Any resource you forget may be applied *after* activation and therefore left pending. In configurations built from reusable modules this list becomes difficult to maintain and easy to get wrong, because you must depend on whole modules and keep that list in step with every future change. Prefer Method 1 for module-based or large configurations.

### Method 3 — The `ZIA_ACTIVATION` environment variable (discouraged)

Setting `ZIA_ACTIVATION=true` makes the provider activate changes in-flight, as resources are configured:

```bash
export ZIA_ACTIVATION=true
```

~> **NOTE:** This method is discouraged and may be removed in a future major release. It triggers an activation call for **every** resource created, updated, or deleted, even though a single call at the end of the run would publish the same changes. The activation endpoint is one of the most tightly rate-limited endpoints in the platform, at **10 POST requests per minute and 40 per hour**. A configuration of any meaningful size therefore exhausts that budget quickly and spends the rest of the run waiting on retries. Refer to the [API Rate Limit Summary](https://help.zscaler.com/legacy-apis/api-rate-limit-summary) for the published limits.

### API Session Timeout and Long-Running Applies

API-initiated sessions have a maximum lifetime controlled by the API session timeout in Advanced Settings. The value can range from **5 to 20 minutes and defaults to 5**. See [Configuring Advanced Settings](https://help.zscaler.com/zia/configuring-advanced-settings#session-timeout) for the platform's description of this setting.

This matters for Terraform because **the platform activates pending changes when a session ends**, including when it ends by reaching the session timeout. Two consequences follow for runs that last longer than the configured timeout:

- Pending changes may be activated part-way through the run, before the apply has finished writing the rest of the configuration.
- The provider must establish a new session and continue, so the run does not fail, but activation has already occurred at a point you did not choose.

A run of 30 minutes against a 5-minute session timeout will therefore cross this boundary several times. **Raise the API session timeout to its maximum of 20 minutes** before running large configurations. This can be done in either of two ways.

**In the ZIA Admin Portal** — navigate to **Administration > Advanced Settings** and set **API Session Timeout Duration (In Minutes)** under *Admin Portal Session Timeout*. Note that this is a separate field from the UI session timeout on the same page:

![ZIA Admin Portal — Advanced Settings, Admin Portal Session Timeout](https://raw.githubusercontent.com/zscaler/terraform-provider-zia/master/docs/guides/media/advanced_settings.png)

**With Terraform** — set the `api_session_timeout` attribute on the [`zia_advanced_settings`](resources/zia_advanced_settings.md) resource:

```hcl
resource "zia_advanced_settings" "this" {
  api_session_timeout = 20
  # … remaining advanced settings attributes
}
```

~> **NOTE:** Raising the timeout reduces how often a long run crosses a session boundary, but it does not eliminate the behaviour. Activation on session end is native to the platform and cannot be disabled or overridden by the provider. 20 minutes is a ceiling, not a guarantee — for very large configurations, split the work across smaller states so that no single run needs to outlive a session.

### Running Several Configurations Against One Tenant

Splitting a large deployment across several Terraform configurations is fully supported. Two operational rules apply, because a ZIA tenant has a single configuration-change lock and a single activation queue regardless of how state is split:

- **Apply one configuration at a time** against a given tenant. `terraform plan` can run in parallel — reads are unaffected. In CI, key the concurrency control on the tenant. Concurrent applies surface as `EDIT_LOCK_NOT_AVAILABLE` or `Failed during enter Org barrier`; the provider retries these, so the usual symptom is a slow run rather than an error.
- **Activate once, after the last apply** — not once per configuration. Activation publishes every pending change in the tenant, so a per-configuration call is redundant and will queue behind any other session that still has unactivated changes.

For the reasoning behind both rules, see the [ZIA Activator](guides/zia-activator-overview.md) guide.

## Rate Limiting

The ZIA platform enforces API rate limits on a per-endpoint basis. Different endpoints have different thresholds — for example, some endpoints allow multiple POST requests per second, while others (such as `/staticIP`) are limited to 1 POST request per second.

**You do not need to configure anything to stay within these limits.** When a limit is exceeded, the API returns an HTTP `429 Too Many Requests` response with a `Retry-After` header, and the provider waits the indicated interval and retries the request transparently. Exceeding a rate limit therefore costs a short delay rather than a failed apply, including on large bulk deployments.

Run Terraform with its default settings. Refer to the [Zscaler Rate Limiting Documentation](https://automate.zscaler.com/docs/api-reference-and-guides/guides/rate-limiting/zia) for details on per-endpoint limits.

Note that rate limiting is distinct from the tenant-wide write lock. An `HTTP 409` response reporting `EDIT_LOCK_NOT_AVAILABLE` or `Failed during enter Org barrier` indicates that another session is modifying the configuration, not that a rate limit was exceeded. See [Running Several Configurations Against One Tenant](#running-several-configurations-against-one-tenant).

## Zscaler Sandbox Authentication

As of version v4.0.0, the ZIA Terraform provider the legacy sandbox authentication environment variables `ZIA_CLOUD` and `ZIA_SANDBOX_TOKEN` are no longer supported.

Authentication to the Zscaler Sandbox service requires the following new environment variables the `ZSCALER_SANDBOX_CLOUD` and `ZSCALER_SANDBOX_TOKEN` or authentication attributes `sandbox_token` and `sandbox_cloud`. For details on how obtain the API Token visit the Zscaler help portal [About Sandbox API Token](https://help.zscaler.com/zia/about-sandbox-api-token)

## Argument Reference - OneAPI

Before starting with this Terraform provider you must create an API Client in the Zscaler Identity Service portal [Zidentity](https://help.zscaler.com/zidentity/what-zidentity) or have create an API key via the legacy method.

- `client_id` - (Required) This is the client ID for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_ID` environment variable.

- `client_secret` - (Optional) This is the client secret for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_SECRET` environment variable. `client_secret` conflicts with `private_key`.

- `private_key` - (Optional) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `ZSCALER_PRIVATE_KEY` environment variable. `private_key` conflicts with `client_secret`. The format of the PK is PKCS#1 unencrypted (header starts with `-----BEGIN RSA PRIVATE KEY-----` or PKCS#8 unencrypted (header starts with `-----BEGIN PRIVATE KEY-----`).

- `vanity_domain` - (Optional) This refers to the domain name used by your organization. It can also be sourced from the `ZSCALER_VANITY_DOMAIN`.

- `zscaler_cloud` - (Optional) This refers to Zscaler cloud name where API calls will be directed to i.e `beta`. It can also be sourced from the `ZSCALER_CLOUD`.

- `sandbox_token` - (Optional) This refers to the Zscaler Sandbox service API Token. It can also be sourced from the `ZSCALER_SANDBOX_TOKEN`.

- `sandbox_cloud` - (Optional) This refers to the Zscaler cloud name where API calls to the sandbox service will be directed. It can also be sourced from the `ZSCALER_SANDBOX_CLOUD`. Currently the following cloud names are supported:
  - `zscaler`
  - `zscalerone`
  - `zscalertwo`
  - `zscalerthree`
  - `zscloud`
  - `zscalerbeta`
  - `zscalergov`
  - `zscalerten`
  - `zspreview`

**NOTE**: Authentication to the Sandbox service is idependent from authentication to OneAPI or the Legacy API framework and can be set and used in standalone mode.

- `http_proxy` - (Optional) This is a custom URL endpoint that can be used for unit testing or local caching proxies. Can also be sourced from the `ZSCALER_HTTP_PROXY` environment variable.

- `max_retries` - (Optional) Maximum number of times a rate-limited request is retried before returning an error. The default is `100` and the maximum is `100`. Each retry waits out the interval reported by the API, so a high value costs nothing when rate limits are not being reached and allows large configurations to complete without manual intervention.

- `request_timeout` - (Optional) Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.

- `skip_credentials_validation` - (Optional) When set to `true`, the provider skips credential validation and does not initialize the API client. Can also be sourced from the `ZSCALER_SKIP_CREDENTIALS_VALIDATION` environment variable. This is intended for configurations where the ZIA provider is declared but every `zia_*` resource and data source is conditionally disabled (e.g., `count = 0`) — such as multi-environment deployments where Zscaler is not present in every environment. With this flag enabled, `terraform plan`/`apply` succeeds with a warning even when no credentials are supplied; any resource or data source that does attempt an API call fails with an explanatory error. Default: `false`.

- `username` - (Optional) Administrator account used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZIA_USERNAME` environment variable.

- `password` - (Optional) Administrator password used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZIA_PASSWORD` environment variable.

- `api_key` - (Optional) API key found in the Zscaler Internet Access portal `Administration > Cloud Service API Security > Cloud Service API Key`. Can also be sourced from the `ZIA_API_KEY` environment variable. Ensure you have the following SKU enabled `Z_API`

- `zia_cloud` - (Optional) This refers to the Zscaler cloud name where api calls will be forward to. Can also be sourced from the `ZIA_CLOUD` environment variable.
Currently the following cloud names are supported:
  - `zscaler`
  - `zscalerone`
  - `zscalertwo`
  - `zscalerthree`
  - `zscloud`
  - `zscalerbeta`
  - `zscalergov`
  - `zscalerten`
  - `zspreview`

- `use_legacy_client` - (Optional) This parameter is required when using the legacy API framework. Can also be sourced from the `ZSCALER_USE_LEGACY_CLIENT` environment variable.
