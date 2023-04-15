resource "zia_dlp_web_rules" "this" {
  name                      = "Terraform_Test"
  description               = "Terraform_Test"
  action                    = "BLOCK"
  ocr_enabled               =  true
  order                     = 1
  protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  cloud_applications        = ["WINDOWS_LIVE_HOTMAIL"]
  file_types                = [ "WINDOWS_META_FORMAT", "BITMAP", "JPEG", "PNG", "TIFF" ]
  rank                      = 7
  state                     = "ENABLED"
  zscaler_incident_reciever = true
  without_content_inspection = false
  url_categories {
    id = [data.zia_url_categories.marketing.val]
  }
  dlp_engines {
      id = [data.zia_dlp_engines.this.id]
    }
  locations {
    id = [data.zia_location_management.this.id]
  }
  labels {
      id = [data.zia_rule_labels.this.id]
    }
}

data "zia_url_categories" "marketing" {
    id = "CORPORATE_MARKETING"
}

data "zia_dlp_engines" "this" {
  name = "PCI"
}

data "zia_location_management" "this" {
  name = "Branch01"
}

data "zia_rule_labels" "this" {
  name = "Compliance"
}