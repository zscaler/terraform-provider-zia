# Example Usage - Action - DECRYPT
resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Inspection_Rule_Decrypt"
  description                  = "SSL_Inspection_Rule_Decrypt"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                         = "DECRYPT"
    override_default_certificate = false

    ssl_interception_cert {
      id                  = 1
      name                = "Zscaler Intermediate CA Certificate"
      default_certificate = true
    }

    decrypt_sub_actions {
      server_certificates                 = "ALLOW"
      ocsp_check                          = true
      block_ssl_traffic_with_no_sni_enabled = true
      min_client_tls_version              = "CLIENT_TLS_1_0"
      min_server_tls_version              = "SERVER_TLS_1_0"
      block_undecrypt                    = true
      http2_enabled                       = false
    }
  }
}

// ## Example Usage - Action - DO_NOT_DECRYPT - Bypass Rule (False)
resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_Do_Not_Decrypt"
  description                  = "SSL_Rule_Do_Not_Decrypt"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "DO_NOT_DECRYPT"
    do_not_decrypt_sub_actions {
      bypass_other_policies                 = false
      server_certificates                   = "ALLOW"
      ocsp_check                            = true
      block_ssl_traffic_with_no_sni_enabled = true
      min_tls_version                       = "SERVER_TLS_1_0"
    }
  }
}


// ## Example Usage - Action - DO_NOT_DECRYPT - Bypass Rule (True)
resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_Do_Not_Decrypt02"
  description                  = "SSL_Rule_Do_Not_Decrypt02"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "DO_NOT_DECRYPT"
    do_not_decrypt_sub_actions {
      bypass_other_policies                 = true
      block_ssl_traffic_with_no_sni_enabled = true
    }
  }
}

// ## Example Usage - Action - BLOCK
resource "zia_ssl_inspection_rules" "this" {
  name                         = "SSL_Rule_BLOCK"
  description                  = "SSL_Rule_BLOCK"
  state                        = "ENABLED"
  order                        = 1
  rank                         = 7
  road_warrior_for_kerberos    = true
  cloud_applications           = ["CHATGPT_AI", "ANDI"]
  platforms                    = ["SCAN_IOS", "SCAN_ANDROID", "SCAN_MACOS", "SCAN_WINDOWS", "NO_CLIENT_CONNECTOR", "SCAN_LINUX"]

  action {
    type                                    = "BLOCK"
    ssl_interception_cert {
      id                                    = 1
    }
  }
}
