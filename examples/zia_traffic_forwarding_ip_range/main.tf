// It will return the first 10 available gre internal ip ranges
data "zia_gre_internal_ip_range_list" "example"{
}

output "zia_gre_internal_ip_range_list_example"{
    value = data.zia_gre_internal_ip_range_list.example
}

// Explicitly retrieves the first 10 available gre internal ip ranges
// You can select any number in the required_count
data "zia_gre_internal_ip_range_list" "example"{
    required_count = 10
}

output "zia_gre_internal_ip_range_list_example"{
    value = data.zia_gre_internal_ip_range_list.example
}