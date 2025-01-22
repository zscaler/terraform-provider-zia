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

## Zscaler OneAPI New Framework

The ZIA Terraform Provider now offers support for [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi) Oauth2 authentication through [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

**NOTE** Notice that OneAPI and Zidentity is not currently supported for the following clouds: `zscalergov` and `zscalerten`. Refer to the [Legacy API Framework](#legacy-api-framework) for more information on how authenticate to these environments

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

**NOTE**: The `zscaler_cloud` is optional and only required when authenticating to other environments i.e `beta`

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

For the resources and data sources examples, please check the [examples](https://github.com/zscaler/terraform-provider-zia/tree/master/examples) directory.

## Authentication - OneAPI New Framework

As of version v4.0.0, this provider supports authentication via the new Zscaler API framework [OneAPI](https://help.zscaler.com/oneapi/understanding-oneapi)

Zscaler OneAPI uses the OAuth 2.0 authorization framework to provide secure access to Zscaler Internet Access (ZIA) APIs. OAuth 2.0 allows third-party applications to obtain controlled access to protected resources using access tokens. OneAPI uses the Client Credentials OAuth flow, in which client applications can exchange their credentials with the authorization server for an access token and obtain access to the API resources, without any user authentication involved in the process.

* [ZIA API](https://help.zscaler.com/oneapi/understanding-oneapi#:~:text=managed%20using%20OneAPI.-,ZIA%20API,-Zscaler%20Internet%20Access)

### Default Environment variables

You can provide credentials via the `ZSCALER_CLIENT_ID`, `ZSCALER_CLIENT_SECRET`, `ZSCALER_VANITY_DOMAIN`, `ZSCALER_CLOUD` environment variables, representing your Zidentity OneAPI credentials `clientId`, `clientSecret`, `vanityDomain` and `cloud` respectively.

| Argument        | Description                                                                                         | Environment Variable     |
|-----------------|-----------------------------------------------------------------------------------------------------|--------------------------|
| `client_id`     | _(String)_ Zscaler API Client ID, used with `clientSecret` or `PrivateKey` OAuth auth mode.         | `ZSCALER_CLIENT_ID`      |
| `client_secret` | _(String)_ Secret key associated with the API Client ID for authentication.                         | `ZSCALER_CLIENT_SECRET`  |
| `privateKey`    | _(String)_ A string Private key value.                                                              | `ZSCALER_PRIVATE_KEY`    |
| `vanity_domain` | _(String)_ Refers to the domain name used by your organization.                                     | `ZSCALER_VANITY_DOMAIN`  |
| `cloud`         | _(String)_ The name of the Zidentity cloud, e.g., beta.                                             | `ZSCALER_CLOUD`          |

### Alternative OneAPI Cloud Environments

OneAPI supports authentication and can interact with alternative Zscaler enviornments i.e `beta`. To authenticate to these environments you must provide the following values:

| Argument         | Description                                                                                         |   | Environment Variable     |
|------------------|-----------------------------------------------------------------------------------------------------|---|--------------------------|
| `vanity_domain`   | _(String)_ Refers to the domain name used by your organization |   | `ZSCALER_VANITY_DOMAIN`  |
| `cloud`          | _(String)_ The name of the Zidentity cloud i.e beta      |   | `ZSCALER_CLOUD`          |

For example: Authenticating to Zscaler Beta environment:

```sh
export ZSCALER_VANITY_DOMAIN="acme"
export ZSCALER_CLOUD="beta"
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

### OneAPI Framework - Configuration file

You can use a configuration file to specify your credentials. The
file location must be `$HOME/.zscaler/zscaler.yaml` on Linux and OS X, or
`"%USERPROFILE%\.zscaler/zscaler.yaml"` for Windows users.
If we fail to detect credentials inline, or in the environment variable, Terraform will check
this location.

Usage:

```terraform
provider "zia" {}
```

```yaml
zscaler:
  client:
    client_id: "{yourClientId}"
    client_secret: "{yourClientSecret}"
    vanity_domain: "{youVanityDomain}"
    cloud: "{yourZidentityCloud}"
```

## Legacy API Framework

### ZIA native authentication

* As of version v4.0.0, this Terraform provider offers backwards compatibility to the Zscaler legacy API framework. This is the recommended authentication method for organizations whose tenants are still not migrated to [Zidentity](https://help.zscaler.com/zidentity/what-zidentity).

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
  cloud               = "[ZIA_CLOUD]"
  use_legacy_client   = "[ZSCALER_USE_LEGACY_CLIENT]"
}
```

The ZIA Cloud is identified by several cloud name prefixes, which determines which API endpoint the requests should be sent to. The following cloud environments are supported:

* `zscaler`
* `zscloud`
* `zscalerone`
* `zscalertwo`
* `zscalerthree`
* `zscalerbeta`
* `zscalergov`
* `zscalerten`
* `zspreview`

### Environment variables

You can provide credentials via the `ZIA_USERNAME`, `ZIA_PASSWORD`, `ZIA_API_KEY`, `ZIA_CLOUD`, `ZSCALER_USE_LEGACY_CLIENT` environment variables, representing your ZIA `username`, `password`, `api_key`,  `cloud` and `use_legacy_client` respectively.

| Argument     | Description | Environment variable |
|--------------|-------------|-------------------|
| `username`       | _(String)_ A string that contains the email ID of the API admin.| `ZIA_USERNAME` |
| `password`       | _(String)_ A string that contains the password for the API admin.| `ZIA_PASSWORD` |
| `api_key`       | _(String)_ A string that contains the obfuscated API key (i.e., the return value of the obfuscateApiKey() method).| `ZIA_API_KEY` |
| `cloud`       | _(String)_ The host and basePath for the cloud services API is `$zsapi.<Zscaler Cloud Name>/api/v1`.| `ZIA_CLOUD` |
| `use_legacy_client`       | _(Bool)_ Enable use of the legacy ZIA API Client.| `ZSCALER_USE_LEGACY_CLIENT` |

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

### Legacy API Framework - Configuration file

You can use a configuration file to specify your credentials. The
file location must be `$HOME/.zscaler/zscaler.yaml` on Linux and OS X, or
`"%USERPROFILE%\.zscaler/zscaler.yaml"` for Windows users.
If we fail to detect credentials inline, or in the environment variable, Terraform will check
this location.

Usage:

```terraform
provider "zia" {}
```

```yaml
zia:
  client:
    username: "{yourUsername}"
    password: "{yourPassword}"
    api_key: "{youApiKey}"
    zia_cloud: "{yourZiaCloud}"
```

### ZIA Configuration Activation

The ZIA platform requires every configuration to be activated. As of version [v2.8.0](https://github.com/zscaler/terraform-provider-zia/releases/tag/v2.8.0) the provider supports implicit activation. In order to make this process more flexible, we have implemented a dedicated environment variable `ZIA_ACTIVATION`, which when set to `true` will implicitly activate the changes as resources are configured.
If the environment variable `ZIA_ACTIVATION` is not set, you must then use the out of band activation method described here [zia activator](guides/zia-activator-overview.md) or leverage the dedicated activation resource `zia_activation_status`.

### Zscaler Sandbox Authentication

As of version v4.0.0, the ZIA Terraform provider the legacy sandbox authentication environment variables `ZIA_CLOUD` and `ZIA_SANDBOX_TOKEN` are no longer supported.

Authentication to the Zscaler Sandbox service requires the following new environment variables the `ZSCALER_SANDBOX_CLOUD` and `ZSCALER_SANDBOX_TOKEN` or authentication attributes `sandbox_token` and `sandbox_cloud`. For details on how obtain the API Token visit the Zscaler help portal [About Sandbox API Token](https://help.zscaler.com/zia/about-sandbox-api-token)

## Argument Reference - OneAPI

Before starting with this Terraform provider you must create an API Client in the Zscaler Identity Service portal [Zidentity](https://help.zscaler.com/zidentity/what-zidentity) or have create an API key via the legacy method.

* `client_id` - (Required) This is the client ID for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_ID` environment variable.

* `client_secret` - (Optional) This is the client secret for obtaining the API token. It can also be sourced from the `ZSCALER_CLIENT_SECRET` environment variable. `client_secret` conflicts with `private_key`.

* `private_key` - (Optional) This is the private key for obtaining the API token (can be represented by a filepath, or the key itself). It can also be sourced from the `ZSCALER_PRIVATE_KEY` environment variable. `private_key` conflicts with `client_secret`. The format of the PK is PKCS#1 unencrypted (header starts with `-----BEGIN RSA PRIVATE KEY-----` or PKCS#8 unencrypted (header starts with `-----BEGIN PRIVATE KEY-----`).

* `vanity_domain` - (Optional) This refers to the domain name used by your organization.. It can also be sourced from the `ZSCALER_VANITY_DOMAIN`.

* `zscaler_cloud` - (Optional) This refers to Zscaler cloud name where API calls will be directed to i.e `beta`. It can also be sourced from the `ZSCALER_CLOUD`.

* `sandbox_token` - (Optional) This refers to the Zscaler Sandbox service API Token. It can also be sourced from the `ZSCALER_SANDBOX_TOKEN`.

* `sandbox_cloud` - (Optional) This refers to the Zscaler cloud name where API calls to the sandbox service will be directed. It can also be sourced from the `ZSCALER_SANDBOX_CLOUD`. Currently the following cloud names are supported:
  * `zscaler`
  * `zscalerone`
  * `zscalertwo`
  * `zscalerthree`
  * `zscloud`
  * `zscalerbeta`
  * `zscalergov`
  * `zscalerten`
  * `zspreview`

**NOTE**: Authentication to the Sandbox service is idependent from authentication to OneAPI or the Legacy API framework and can be set and used in standalone mode.

* `http_proxy` - (Optional) This is a custom URL endpoint that can be used for unit testing or local caching proxies. Can also be sourced from the `ZSCALER_HTTP_PROXY` environment variable.

* `parallelism` - (Optional) Number of concurrent requests to make within a resource where bulk operations are not possible. [Learn More](https://help.zscaler.com/oneapi/understanding-rate-limiting)

* `max_retries` - (Optional) Maximum number of retries to attempt before returning an error, the default is `5`.

* `request_timeout` - (Optional) Timeout for single request (in seconds) which is made to Zscaler, the default is `0` (means no limit is set). The maximum value can be `300`.

* `username` - (Optional) Administrator account used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZIA_USERNAME` environment variable.

* `password` - (Optional) Administrator password used when authenticating to the legacy Zscaler API framework. Can also be sourced from the `ZIA_PASSWORD` environment variable.

* `api_key` - (Optional) API key found in the Zscaler Internet Access portal `Administration > Cloud Service API Security > Cloud Service API Key`. Can also be sourced from the `ZIA_API_KEY` environment variable. Ensure you have the following SKU enabled `Z_API`

* `cloud` - (Optional) This refers to the Zscaler cloud name where api calls will be forward to. Can also be sourced from the `ZIA_CLOUD` environment variable.
Currently the following cloud names are supported:
  * `zscaler`
  * `zscalerone`
  * `zscalertwo`
  * `zscalerthree`
  * `zscloud`
  * `zscalerbeta`
  * `zscalergov`
  * `zscalerten`
  * `zspreview`

* `use_legacy_client` - (Optional) This parameter is required when using the legacy API framework. Can also be sourced from the `ZSCALER_USE_LEGACY_CLIENT` environment variable.
