terraform {
  required_providers {
    zia = {
      version = "1.0.0"
      source  = "zscaler.com/zia/zia"
    }
  }
}

provider "zia" {}

resource "zia_firewall_filtering_rule" "example" {
  name                = "Example"
  description         = "Example"
  state               = "ENABLED"
  action              = "ALLOW"
  predefined          = false
  default_rule        = false
  access_control      = "READ_WRITE"
  enable_full_logging = false
  order               = 1
  dest_countries      = ["COUNTRY_CA", "COUNTRY_US", "COUNTRY_BR"]
  locations {
    id = [
      data.zia_location_management.sjc1.id
    ]
  }
}

data "zia_location_management" "sjc1" {
  name = "sjc-1"
}
