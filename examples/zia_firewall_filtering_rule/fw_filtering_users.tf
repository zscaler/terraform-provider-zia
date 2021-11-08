resource "zia_firewall_filtering_rule" "example" {
    name = "Example"
    description = "Example"
    action = "ALLOW"
    state = "ENABLED"
    order = 1
    enable_full_logging = true
    users {
        id = [ data.zia_user_management.user1.id,
               data.zia_user_management.user2.id,
               data.zia_user_management.user3.id,
               data.zia_user_management.user4.id
            ]
    }
}

data "zia_user_management" "user1" {
    name = "User1"
}

data "zia_user_management" "user2" {
    name = "User2"
}

data "zia_user_management" "user3" {
    name = "User3"
}

data "zia_user_management" "user4" {
    name = "User4"
}