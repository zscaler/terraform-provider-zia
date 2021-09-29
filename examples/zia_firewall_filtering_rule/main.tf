terraform {
    required_providers {
        zia = {
            version = "1.0.0"
            source = "zscaler.com/zia/zia"
        }
    }
}

provider "zia" {}

data "zia_firewall_filtering_rule" "example"{
    name = "Example"
}

output "zia_firewall_filtering_rule_example"{
    value = data.zia_firewall_filtering_rule.example
}

