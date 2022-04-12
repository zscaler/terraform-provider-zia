---
layout: "zscaler"
page_title: "Provider: Zscaler Internet Access (ZIA)"
sidebar_current: "docs-zia-index"
description: |-
        The Zscaler Internet Access provider is used to interact with ZIA API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, cloud firewall policies, dlp dictionaries, local accounts etc. The provider is intended to save time and reducing configuration errors. With this ZIA provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZIA API credentials.
---

# Zscaler Internet Access (ZIA) Provider

The Zscaler Internet Access provider is used to interact with ZIA API, to automate the provisioning of new locations, IPSec and GRE tunnels, URL filtering policies, cloud firewall policies, dlp dictionaries, local accounts etc. The provider is intended to save time and reducing configuration errors. With this ZIA provider, DevOps teams can automate their security and transform it into DevSecOps workflows. To use this  provider, you must create ZIA API credentials.

Use the navigation on the left to read about the available resources.

## Examples Usage

```hcl
# Configure the Zscaler Internet Access Provider
terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}
```

```hcl
# Create a location management resource
resource "zia_location_management" "gre_canada_toronto_branch01"{
    name = "GRE Canada - Toronto - Branch01"
    description = "GRE Canada - Toronto - Branch01"
    country = "CANADA"
    tz = "CANADA_AMERICA_TORONTO"
    ip_addresses = [ zia_traffic_forwarding_static_ip.gre_ca_toronto_branch01.ip_address ]
    depends_on = [ zia_traffic_forwarding_static_ip.gre_ca_toronto_branch01 ]
}
```

```hcl
# Create GRE Tunnel
resource "zia_traffic_forwarding_gre_tunnel" "gre_ca_toronto_branch01" {
  source_ip = zia_traffic_forwarding_static_ip.gre_ca_toronto_branch01.ip_address
  comment   = "GRE Canada - Toronto - Branch01"
  within_country = true
  country_code = "CA"
  ip_unnumbered = false
  depends_on = [ zia_traffic_forwarding_static_ip.gre_ca_toronto_branch01 ]
}
```

```hcl
# Create Static IP Address
resource "zia_traffic_forwarding_static_ip" "gre_ca_toronto_branch01"{
    ip_address =  "xx.xxx.xxx.xxx"
    routable_ip = true
    comment = "GRE Canada - Toronto - Branch01"
    geo_override = false
}
```

## Authentication

The ZIA provider offers various means of providing credentials for authentication. The following methods are supported:

* Static credentials
* Environment variables

### Static credentials

⚠️ **WARNING:** Hard-coding credentials into any Terraform configuration is not recommended, and risks secret leakage should this file be committed to public version control

Static credentials can be provided by specifying the `username`, `password`, `api_key`, `zia_cloud` arguments in-line in the ZIA provider block:

**Usage:**

```hcl
provider "zia" {
  username      = 'xxxxxxxxxxxxxxxx'
  password      = 'xxxxxxxxxxxxxxxx'
  api_key       = 'xxxxxxxxxxxxxxxx'
  zia_cloud  = 'https://zsapi.<zscaler_cloud_name>.net'
}
```

### Environment variables

You can provide credentials via the `ZIA_USERNAME`, `ZIA_PASSWORD`, `ZIA_API_KEY`, `zia_cloud` environment variables, representing your ZIA username, password, API Key credentials and tenant base URL, respectively.

```hcl
provider "zia" {}
```

**Usage:**

```sh
export ZIA_USERNAME = "xxxxxxxxxxxxxxxx"
export ZIA_PASSWORD = "xxxxxxxxxxxxxxxx"
export ZIA_API_KEY  = "xxxxxxxxxxxxxxxx"
export zia_cloud = "xxxxxxxxxxxxxxxx"
terraform plan
```

If you are on Windows, use PowerShell to set the environmenr variables using the following commands:

```sh
$env:username = 'xxxxxxxxxxxxxxxx'
$env:password = 'xxxxxxxxxxxxxxxx'
$env:api_key = 'xxxxxxxxxxxxxxxx'
$env:zia_cloud = 'https://zsapi.<zscaler_cloud_name>.net'
```

For details about how to retrieve your tenant Base URL and API key/token refer to the Zscaler help portal. <https://help.zscaler.com/zia/getting-started-zia-api>
