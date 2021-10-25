---
subcategory: "Firewall Filtering Time Window"
layout: "zia"
page_title: "ZIA: firewall_filtering_time_window"
description: |-
  Retrieve ZIA firewall rule time window.
  
---
# zia_firewall_filtering_time_window (Data Source)

The **zia_firewall_filtering_time_window** data source provides details about a specific time window option available in the Zscaler Internet Access cloud firewall. This data source can then be associated with a ZIA firewall filtering rule.

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

### Read-Only Attributes

* `start_time` - (String)
* `end_time` - (String)
* `day_of_week` - (String)
