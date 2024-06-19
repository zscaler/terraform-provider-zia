resource "zia_url_filtering_rules" "allow_rule" {
    name = "ExampleAllow"
    description = "ExampleAllow"
    state = "ENABLED"
    action = "ALLOW"
    order = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    url_categories = ["ANY"]
    protocols = ["HTTPS_RULE", "HTTP_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    time_quota = 15
    size_quota = 10
}

resource "zia_url_filtering_rules" "block_rule" {
    name = "ExampleBlock"
    description = "ExampleBlock"
    state = "ENABLED"
    action = "BLOCK"
    order = 1
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    user_agent_types = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    time_quota = 15
    size_quota = 10
    block_override = true
    override_users {
        id = [ 45513075 ]
    }
    override_groups {
        id = [ 76662385 ]
    }
}

resource "zia_url_filtering_rules" "caution_rule" {
    name = "ExampleCaution"
    description = "ExampleCaution"
    state = "ENABLED"
    action = "CAUTION"
    order = 3
    enforce_time_validity = true
    validity_start_time   = "Mon, 17 Jun 2024 23:30:00 UTC"
    validity_end_time     = "Tue, 17 Jun 2025 23:00:00 UTC"
    validity_time_zone_id = "US/Pacific"
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    end_user_notification_url = "https://caution.acme.com"
    user_agent_types = ["OPERA", "FIREFOX", "MSIE", "MSEDGE", "CHROME", "SAFARI", "MSCHREDGE"]
    request_methods = [ "CONNECT", "GET", "HEAD" ]
    time_quota = 15
    size_quota = 10
}
