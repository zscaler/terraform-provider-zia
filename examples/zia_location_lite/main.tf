data "zia_location_lite" "this" {
 name = "Road Warrior"
}

output "zia_location_lite"{
    value = data.zia_location_lite.this
}
