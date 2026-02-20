resource "zia_extranet" "this" {
    name        = "Extranet01"
    description = "Extranet01"

    extranet_dns_list {
        name                 = "DNS01"
        primary_dns_server   = "8.8.8.8"
        secondary_dns_server = "4.4.4.4"
        use_as_default       = true
    }

    extranet_dns_list {
        name                 = "DNS02"
        primary_dns_server   = "192.168.1.1"
        secondary_dns_server = "192.168.1.2"
        use_as_default       = false
    }

    extranet_ip_pool_list {
        name           = "TFS01"
        ip_start       = "10.0.0.1"
        ip_end         = "10.0.0.21"
        use_as_default = true
    }

    extranet_ip_pool_list {
        name           = "TFS02"
        ip_start       = "10.0.0.22"
        ip_end         = "10.0.0.43"
        use_as_default = false
    }
}