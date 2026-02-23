data "zia_datacenters" "this" {
  name = "SJC4"
}

# Using Unix timestamps
resource "zia_dc_exclusions" "this" {
  datacenter_id = data.zia_datacenters.this.datacenter_id
  start_time   = 1770422399
  end_time     = 1770508799
  description  = "Optional description"
}

# Using human-readable UTC date/time (same as zia_sub_cloud exclusions)
resource "zia_dc_exclusions" "utc" {
  datacenter_id  = data.zia_datacenters.this.datacenter_id
  start_time_utc = "02/05/2026 12:00 am"
  end_time_utc   = "02/19/2026 11:59 pm"
  description    = "Optional description"
}