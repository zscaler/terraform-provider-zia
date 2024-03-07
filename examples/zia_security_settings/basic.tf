# Add URLs to ZIA Whitelist and BlackList - Advanced Threat Protection & Malware Protection respectively
resource "zia_security_settings" "this" {
  whitelist_urls = [
    "resource5.acme.net",
    "resource6.acme.net",
    "resource7.acme.net",
    "resource8.acme.net",
  ]
  blacklist_urls = [
    "resource1.acme.net",
    "resource2.acme.net",
    "resource3.acme.net",
    "resource4.acme.net",
  ]
}

# Add URLs to ZIA Blacklist - Advanced Threat Protection
resource "zia_security_settings" "this" {
  blacklist_urls = [
    "resource1.acme.net",
    "resource2.acme.net",
    "resource3.acme.net",
    "resource4.acme.net",
  ]
}

# Add URLs to ZIA Whitelist - Malware Protection
resource "zia_security_settings" "this" {
  whitelist_urls = [
    "resource5.acme.net",
    "resource6.acme.net",
    "resource7.acme.net",
    "resource8.acme.net",
  ]
}