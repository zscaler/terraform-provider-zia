# Retrieves and Filter Cloud Application by application category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class   = ["AI_ML"]
}

output "app_ids" {
  value = [for app in data.zia_cloud_applications.this.applications : app["app"]]
}

# Create a File Type Control rule and associated Cloud Application categories
resource "zia_file_type_control_rules" "this" {
    name               = "File_Type_Rule01"
    description        = "File Type Rule Created via Terraform"
    state              = "ENABLED"
    order              = 1
    rank               = 7
    filtering_action   = "BLOCK"
    operation          = "DOWNLOAD"
    active_content     = true
    unscannable        = false
    device_trust_levels = ["UNKNOWN_DEVICETRUSTLEVEL", "LOW_TRUST", "MEDIUM_TRUST", "HIGH_TRUST"]
    file_types         = ["FTCATEGORY_MS_WORD", "FTCATEGORY_MS_POWERPOINT", "FTCATEGORY_PDF_DOCUMENT", "FTCATEGORY_MS_EXCEL"]
    protocols          = ["FOHTTP_RULE", "FTP_RULE", "HTTPS_RULE", "HTTP_RULE"]
    cloud_applications = tolist([for app in data.zia_cloud_applications.this.applications : app["app"]])
    groups {
        id = [12006683]
    }
}
