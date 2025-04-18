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

resource "zia_dlp_web_rules" "this" {
  name                      = "Terraform_Test"
  description               = "Terraform_Test"
  action                    = "BLOCK"
  order                     = 1
  protocols                 = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  cloud_applications        = ["WINDOWS_LIVE_HOTMAIL"]
  file_types                = [ "WINDOWS_META_FORMAT", "BITMAP", "JPEG", "PNG", "TIFF" ]
  rank                      = 7
  state                     = "ENABLED"
  zscaler_incident_receiver = true
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

resource "zia_dlp_web_sub_rules" "subrule1" {
  name                       = "Terraform_Test_subrule_prod_tf"
  description                = "Terraform_Test_subrule_prod_tf"
  action                     = "BLOCK"
  state                      = "ENABLED"
  order                      = 1
  rank                       = 0
  protocols                  = ["FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
  cloud_applications         = ["WINDOWS_LIVE_HOTMAIL"]
  without_content_inspection = false
  file_types                 = ["FTCATEGORY_MS_WORD", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_MS_EXCEL"]
  match_only                 = false
  min_size                   = 10
  parent_rule                = zia_dlp_web_rules.this.id # This attribuite is required and indicates the parent rule to which it is associated with.
  icap_server {
    id = data.zia_dlp_incident_receiver_servers.this.id
  }
  groups {
    id = data.zia_group_management.this[*].id
  }
  notification_template {
    id = data.zia_dlp_notification_templates.this.id
  }
  url_categories {
    id = [data.zia_url_categories.this.val]
  }
  dlp_engines {
    id = data.zia_dlp_engines.this[*].id
  }
}