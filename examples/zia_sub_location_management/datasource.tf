data "zia_location_management" "usa_sjc37" {
    name = "USA-SJC37"
    parent_name = "USA"
}

output "zia_location_management" {
    value = data.zia_location_management.usa_sjc37
}