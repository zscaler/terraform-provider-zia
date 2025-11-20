---
subcategory: "Traffic Forwarding"
layout: "zscaler"
page_title: "ZIA: traffic_forwarding_static_ip"
description: |-
    Official documentation https://help.zscaler.com/zia/about-static-ip
    API documentation https://help.zscaler.com/zia/traffic-forwarding-0#/staticIP-get
    Creates and manages static IP addresses with automatic coordinate determination.
---

# zia_traffic_forwarding_static_ip (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-static-ip)
* [API documentation](https://help.zscaler.com/zia/traffic-forwarding-0#/staticIP-get)

The **zia_traffic_forwarding_static_ip** resource allows the creation and management of static IP addresses in the Zscaler Internet Access cloud. The resource can then be associated with other resources such as:

* VPN Credentials of type `IP`
* Location Management
* GRE Tunnel

## 🎯 Automatic Coordinate Determination (v4.6.2+)

Starting with **version 4.6.2**, the provider automatically determines latitude and longitude coordinates from the IP address, even when `geo_override = true`. This means:

* ✅ **No manual coordinate lookups** - Provider handles it automatically
* ✅ **No drift issues** - State always contains exact API values
* ✅ **Simpler configuration** - Omit `latitude` and `longitude` for automatic determination
* ✅ **Fully backward compatible** - Explicit coordinates still work if provided

**In short:** You can now use `geo_override = true` without specifying coordinates! See examples below.

## Example Usage

### Example 1: Auto-Determined Coordinates (Recommended)

```hcl
# ZIA Traffic Forwarding - Static IP
# The provider automatically determines latitude and longitude from the IP address
resource "zia_traffic_forwarding_static_ip" "example"{
    ip_address   = "122.164.82.249"
    routable_ip  = true
    comment      = "Static IP with auto-determined coordinates"
    geo_override = true
    # latitude and longitude are omitted - provider will auto-determine them
    # State will be populated with exact API values (e.g., latitude=13.0895, longitude=80.2739)
}
```

### Example 2: User-Specified Coordinates (Optional)

```hcl
# You can still explicitly provide coordinates if needed
resource "zia_traffic_forwarding_static_ip" "custom_location"{
    ip_address   = "1.1.1.1"
    routable_ip  = true
    comment      = "Static IP with custom coordinates"
    geo_override = true
    latitude     = -36.848461
    longitude    = 174.763336
}
```

### Example 3: Automatic Geolocation (geo_override = false)

```hcl
# When geo_override is false or omitted, all geo information is auto-determined
resource "zia_traffic_forwarding_static_ip" "auto_geo"{
    ip_address  = "8.8.8.8"
    routable_ip = true
    comment     = "Fully automatic geolocation"
    # geo_override defaults to false
    # latitude and longitude auto-determined and populated in state
}
```

## Argument Reference

The following arguments are supported:

### Required

* `ip_address` - (Required) The static IP address

### Optional

* `comment` - (Optional) Additional information about this static IP address
* `routable_ip` - (Optional) Indicates whether a non-RFC 1918 IP address is publicly routable. This attribute is ignored if there is no ZIA Private Service Edge associated to the organization.
* `geo_override` - (Optional) When set to `true`, allows custom geolocation settings. When `false` or omitted, geographic coordinates are automatically determined from the IP address. Default: `false`
* `latitude` - (Optional, Computed) Latitude coordinate with up to 7 decimal places precision (range: -90 to 90).
  * If **not provided**, the provider automatically determines the latitude from the IP address, even when `geo_override = true`
  * If **provided**, uses the specified value
  * Always populated in state with exact API values to prevent drift
* `longitude` - (Optional, Computed) Longitude coordinate with up to 7 decimal places precision (range: -180 to 180).
  * If **not provided**, the provider automatically determines the longitude from the IP address, even when `geo_override = true`
  * If **provided**, uses the specified value
  * Always populated in state with exact API values to prevent drift

### Computed

* `static_ip_id` - (Computed) The unique identifier for the static IP address
* `managed_by` - (Computed) Information about who manages this static IP

## How Latitude and Longitude Are Determined

The provider handles coordinates intelligently based on your configuration:

### When `geo_override = false` (or omitted)
* ✅ **Provider behavior**: Latitude and longitude are automatically determined by the ZIA API based on the IP address
* ✅ **State file**: Will contain the API-determined coordinates
* ✅ **User action**: None required - fully automatic

### When `geo_override = true` WITHOUT coordinates
* ✅ **Provider behavior**:
  1. Creates the static IP with `geo_override = false` first
  2. Retrieves the auto-determined coordinates from the API
  3. Updates the static IP with `geo_override = true` using those coordinates
* ✅ **State file**: Will contain the auto-determined coordinates
* ✅ **User action**: None required - provider handles it automatically
* ✅ **Result**: You get `geo_override = true` without manually looking up coordinates

### When `geo_override = true` WITH coordinates
* ✅ **Provider behavior**: Uses your specified coordinates
* ✅ **State file**: Will contain the exact values returned by the API (may have minor precision adjustments)
* ✅ **User action**: Provide `latitude` and `longitude` values
* ✅ **Result**: Your custom coordinates are used

### Key Benefits
* 🎯 **No drift issues** - State always contains exact API values
* 🎯 **No manual lookups** - API determines accurate coordinates from IP
* 🎯 **Flexible** - Can override coordinates when needed
* 🎯 **Always accurate** - Coordinates match the IP address geolocation

## Common Use Cases

### Use Case 1: GRE Tunnel with Auto-Determined Coordinates

```hcl
# Create static IP without specifying coordinates
resource "zia_traffic_forwarding_static_ip" "gre_endpoint" {
    ip_address   = "203.0.113.10"
    routable_ip  = true
    comment      = "GRE tunnel endpoint"
    geo_override = true
}

# Use the static IP with GRE VIP recommendation
data "zia_traffic_forwarding_gre_vip_recommended_list" "vips" {
    source_ip      = zia_traffic_forwarding_static_ip.gre_endpoint.ip_address
    required_count = 2
}

# Create GRE tunnel
resource "zia_traffic_forwarding_gre_tunnel" "main" {
    source_ip      = zia_traffic_forwarding_static_ip.gre_endpoint.ip_address
    comment        = "Main GRE tunnel"
    within_country = false
    ip_unnumbered  = false

    primary_dest_vip {
        datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[0].datacenter
        id         = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[0].id
        virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[0].virtual_ip
    }

    secondary_dest_vip {
        datacenter = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[1].datacenter
        id         = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[1].id
        virtual_ip = data.zia_traffic_forwarding_gre_vip_recommended_list.vips.list[1].virtual_ip
    }
}

# Output showing auto-determined coordinates
output "static_ip_coordinates" {
    value = {
        ip        = zia_traffic_forwarding_static_ip.gre_endpoint.ip_address
        latitude  = zia_traffic_forwarding_static_ip.gre_endpoint.latitude
        longitude = zia_traffic_forwarding_static_ip.gre_endpoint.longitude
    }
}
```

### Use Case 2: Multiple Static IPs with for_each

```hcl
locals {
    office_ips = {
        mumbai     = "103.21.244.1"
        chennai    = "122.164.82.249"
        singapore  = "203.0.113.50"
        tokyo      = "203.0.113.100"
    }
}

# Create multiple static IPs without specifying coordinates
resource "zia_traffic_forwarding_static_ip" "offices" {
    for_each = local.office_ips

    ip_address   = each.value
    routable_ip  = true
    comment      = "Office in ${each.key}"
    geo_override = true
    # No coordinates specified for any of them!
    # Provider auto-determines all coordinates
}

# Output all coordinates
output "office_coordinates" {
    value = {
        for name, ip in zia_traffic_forwarding_static_ip.offices :
        name => {
            ip_address = ip.ip_address
            latitude   = ip.latitude
            longitude  = ip.longitude
        }
    }
}
```

### Use Case 3: VPN Credentials Integration

```hcl
# Static IP with auto-determined coordinates
resource "zia_traffic_forwarding_static_ip" "vpn_endpoint" {
    ip_address   = "198.51.100.25"
    routable_ip  = true
    comment      = "VPN endpoint"
    geo_override = true
}

# VPN credentials using the static IP
resource "zia_traffic_forwarding_vpn_credentials" "branch_office" {
    type        = "IP"
    ip_address  = zia_traffic_forwarding_static_ip.vpn_endpoint.ip_address
    comments    = "Branch office VPN"
}
```

## Frequently Asked Questions (FAQ)

### Q: Do I need to specify latitude and longitude when using geo_override = true?

**A:** No! The provider will automatically determine coordinates from the IP address if you don't provide them. This is the **recommended approach** to avoid drift issues.

### Q: What if I want to use specific coordinates?

**A:** You can still provide `latitude` and `longitude` explicitly. The provider will use your values if provided.

### Q: Will there be drift if I don't specify coordinates?

**A:** No! The state file will contain the exact coordinates returned by the ZIA API. Subsequent `terraform plan` commands will show no changes.

### Q: What happens if I provide coordinates that don't match the IP location?

**A:** The API will accept your coordinates, but they may be adjusted for precision. The state file will always reflect the actual API response values.

### Q: Can I change from auto-determined to custom coordinates later?

**A:** Yes! Simply add `latitude` and `longitude` to your configuration and run `terraform apply`. The provider will update the static IP with your custom coordinates.

### Q: What precision does the API use for coordinates?

**A:** The API typically returns 4-7 decimal places depending on the IP location. The provider stores these exact values without rounding.

### Q: Why does my state show geo_override = true but I didn't set it?

**A:** The `geo_override` attribute has `Computed: true`, meaning it's populated from the API response. The API may set it based on other factors.

## Troubleshooting

### Error: "Missing geo Coordinates"

This error should no longer occur with the updated provider. If you still see it:

1. Ensure you're using provider version 4.6.2 or later
2. Check if coordinates are being populated: `terraform state show zia_traffic_forwarding_static_ip.<name>`
3. Enable debug logging: `export TF_LOG=DEBUG` and check for auto-population messages

### Unexpected Drift Detected

If `terraform plan` shows coordinate changes:

1. **Solution**: Remove explicit `latitude` and `longitude` from your configuration
2. **Reason**: API values may differ slightly from user-provided values due to precision
3. **After removal**: Run `terraform apply` once - state will sync with API values
4. **Future plans**: Will show no changes

### Coordinates Not in Expected Location

The coordinates reflect the IP address's actual geolocation as determined by Zscaler's geolocation database. If you need different coordinates:

1. Set `geo_override = true`
2. Provide your desired `latitude` and `longitude` explicitly
3. The API will use your values

## Best Practices

### ✅ Recommended: Let the Provider Auto-Determine Coordinates

```hcl
resource "zia_traffic_forwarding_static_ip" "best_practice" {
    ip_address   = "203.0.113.10"
    routable_ip  = true
    comment      = "Production endpoint"
    geo_override = true
    # Omit latitude and longitude
    # Provider will auto-determine accurate coordinates
    # No drift, no manual lookups, always accurate
}
```

**Why this is recommended:**

* ✅ No manual coordinate lookups required
* ✅ Zero drift - state always matches API
* ✅ Accurate - API knows the correct geolocation for each IP
* ✅ Maintainable - no hardcoded coordinates to update

### ⚠️ Use Custom Coordinates Only When Necessary

Only provide explicit coordinates if you have a specific requirement:

```hcl
resource "zia_traffic_forwarding_static_ip" "custom" {
    ip_address   = "203.0.113.10"
    routable_ip  = true
    comment      = "Custom location for testing"
    geo_override = true
    latitude     = 40.7128   # Only if you need specific coordinates
    longitude    = -74.0060  # Only if you need specific coordinates
}
```

**When to use custom coordinates:**

* Testing with specific geographic locations
* Compliance requirements for specific geo-coordinates
* Override API's geolocation database for special cases

## Migration Guide for Existing Users

If you're upgrading from an older provider version (< 4.6.2), you may have configurations like this:

### Old Configuration (Still Works, But Not Recommended)

```hcl
resource "zia_traffic_forwarding_static_ip" "old_style" {
    ip_address   = "122.164.82.249"
    routable_ip  = true
    comment      = "Old configuration"
    geo_override = true
    latitude     = 13.0895   # Manually specified
    longitude    = 80.2739   # Manually specified
}
```

### Migrating to New Approach (Recommended)

**Step 1:** Remove `latitude` and `longitude` from your configuration

```hcl
resource "zia_traffic_forwarding_static_ip" "old_style" {
    ip_address   = "122.164.82.249"
    routable_ip  = true
    comment      = "Migrated configuration"
    geo_override = true
    # Removed: latitude and longitude
}
```

**Step 2:** Run `terraform plan`

```bash
terraform plan
```

You'll see Terraform wants to update the resource (to remove explicitly set coordinates from state).

**Step 3:** Apply the changes

```bash
terraform apply
```

The provider will:

* Keep the same static IP (no destruction)
* Auto-determine coordinates from the IP
* Update state with API values
* No infrastructure change - just cleaner config!

**Step 4:** Verify no drift

```bash
terraform plan
# Expected: No changes. Your infrastructure matches the configuration.
```

### Migration Example: Full Before/After

**Before Migration:**

```hcl
# ❌ Old way - manual coordinates required
resource "zia_traffic_forwarding_static_ip" "chennai" {
    ip_address   = "122.164.82.249"
    routable_ip  = true
    comment      = "Chennai office"
    geo_override = true
    latitude     = 13.0895   # Had to look this up
    longitude    = 80.2739   # Had to look this up
}

resource "zia_traffic_forwarding_static_ip" "mumbai" {
    ip_address   = "103.21.244.1"
    routable_ip  = true
    comment      = "Mumbai office"
    geo_override = true
    latitude     = 19.0760   # Had to look this up
    longitude    = 72.8777   # Had to look this up
}
```

**After Migration:**

```hcl
# ✅ New way - auto-determined coordinates
resource "zia_traffic_forwarding_static_ip" "chennai" {
    ip_address   = "122.164.82.249"
    routable_ip  = true
    comment      = "Chennai office"
    geo_override = true
    # No coordinates needed!
}

resource "zia_traffic_forwarding_static_ip" "mumbai" {
    ip_address   = "103.21.244.1"
    routable_ip  = true
    comment      = "Mumbai office"
    geo_override = true
    # No coordinates needed!
}
```

**Migration Impact:**

* Configuration: 8 lines removed (cleaner)
* API calls: No additional overhead after migration
* Drift: Eliminated
* Maintenance: Easier

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

Static IP resources can be imported by using `<STATIC IP ID>` or `<IP ADDRESS>` as the import ID.

### Import by Static IP ID

```shell
terraform import zia_traffic_forwarding_static_ip.example <static_ip_id>
```

Example:

```shell
terraform import zia_traffic_forwarding_static_ip.chennai 3030759
```

### Import by IP Address

```shell
terraform import zia_traffic_forwarding_static_ip.example <ip_address>
```

Example:

```shell
terraform import zia_traffic_forwarding_static_ip.chennai 122.164.82.249
```

**After Import:**

* The state will include all attributes including latitude and longitude
* You can omit coordinates from your configuration - state will remain accurate
* Run `terraform plan` to see what configuration should look like
