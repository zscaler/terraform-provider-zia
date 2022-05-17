---
subcategory: "Firewall Policies"
layout: "zscaler"
page_title: "ZIA): firewall_filtering_time_window"
description: |-
  Get information about firewall rule time window.

---
# Data Source: zia_firewall_filtering_time_window

Use the **zia_firewall_filtering_time_window** data source to get information about a time window option available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

## Example Usage

```hcl
# ZIA Time Window - Work Hours
data "zia_firewall_filtering_time_window" "work_hours"{
    name = "Work hours"
}
```

```hcl
# ZIA Time Window - Weekends
data "zia_firewall_filtering_time_window" "weekends"{
    name = "Weekends"
}
```

```hcl
# ZIA Time Window - Off Hours
data "zia_firewall_filtering_time_window" "off_hours"{
    name = "Off hours"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the time window to be exported.
* `id` - (Optional) The ID of the time window resource.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `start_time` - (String)
* `end_time` - (String)
* `day_of_week` - (String). The supported values are:
  * `ANY` - (String)
  * `NONE` - (String)
  * `EVERYDAY` - (String)
  * `SUN` - (String)
  * `MON` - (String)
  * `TUE` - (String)
  * `WED` - (String)
  * `THU` - (String)
  * `FRI` - (String)
  * `SAT` - (String)
