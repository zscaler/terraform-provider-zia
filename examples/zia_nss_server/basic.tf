resource "zia_nss_server" "this" {
    name = "NSSServer01"
    status = "ENABLED"
    type = "NSS_FOR_WEB"
}

resource "zia_nss_server" "this" {
    name = "NSSServer01"
    status = "ENABLED"
    type = "NSS_FOR_FIREWALL"
}