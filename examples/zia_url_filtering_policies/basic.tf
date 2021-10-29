resource "zia_url_filtering_rules" "block_streaming" {
    name = "Block Streaming"
    description = "Block Video Streaming."
    state = "ENABLED"
    action = "BLOCK"
    order = 2
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}

resource "zia_url_filtering_rules" "caution_for_gambling" {
    name = "Caution for Gambling"
    description = "Caution Marketing users going to gambling sites, but allow them after the caution. This is to support marketing ad campaigns on Gambling websites."
    state = "ENABLED"
    action = "CAUTION"
    order = 3
    url_categories = ["ANY"]
    protocols = ["ANY_RULE"]
    request_methods = [ "CONNECT", "GET", "HEAD" ]
}

resource "zia_url_filtering_rules" "allow_all_other_traffic" {
    name = "Allow All Other Traffic"
    description = "Allow All Other Traffic"
    state = "ENABLED"
    action = "ALLOW"
    order = 4
    url_categories = ["ANY"]
    protocols = ["HTTPS_RULE", "HTTP_RULE"]
    request_methods = [ "CONNECT", "DELETE", "GET", "HEAD", "OPTIONS", "OTHER", "POST", "PUT", "TRACE"]
}