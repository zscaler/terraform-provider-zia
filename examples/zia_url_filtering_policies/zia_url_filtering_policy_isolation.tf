data "zia_cloud_browser_isolation_profile" "this" {
    name = "BD_SA_Profile2_ZIA"
}

resource "zia_url_filtering_rules" "this" {
    name = "ExampleIsolation"
    description = "ExampleIsolation"
    state = "ENABLED"
    action = "ISOLATE"
    order = 1
    url_categories = ["ANY"]
    protocols = ["HTTPS_RULE", "HTTP_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
    cbi_profile {
        id = data.zia_cloud_browser_isolation_profile.this.id
        name = data.zia_cloud_browser_isolation_profile.this.name
        url = data.zia_cloud_browser_isolation_profile.this.url
    }
    user_agent_types = [
        "OPERA",
        "FIREFOX",
        "MSIE",
        "MSEDGE",
        "CHROME",
        "SAFARI",
        "MSCHREDGE"
    ]
}
