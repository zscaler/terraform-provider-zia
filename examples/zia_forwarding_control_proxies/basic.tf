resource "zia_forwarding_control_proxies" "this" {
  name  = "Proxy01_Terraform"
  description = "Proxy01_Terraform"
  type = "PROXYCHAIN"
  address = "192.168.1.150"
  port = 5000
  insert_xau_header = true
  base64_encode_xau_header = true
  cert {
    id = 18492369
  }
}