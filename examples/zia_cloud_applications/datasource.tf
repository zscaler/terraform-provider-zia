# Retrieves a list of Predefined and User Defined Cloud Applications associated with the DLP rules, Cloud App Control rules, Advanced Settings, Bandwidth Classes, and File Type Control rules.
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
}

output "zia_cloud_applications" {
  value = data.zia_cloud_applications.this
}

# Retrieves and Filter Cloud Application by application category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_policy"
  app_class   = ["AI_ML"]
}

output "app_ids" {
  value = [for app in data.zia_cloud_applications.this.applications : app["app"]]
}

# Retrieves specific application by name and category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class = ["SOCIAL_NETWORKING"]
  app_name = "Nebenan"
}

output "zia_cloud_applications" {
    value = data.zia_cloud_applications.this
}


# Retrieves a list of Predefined and User Defined Cloud Applications associated with the SSL Inspection rules.
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
}

output "zia_cloud_applications" {
  value = data.zia_cloud_applications.this
}

#Retrieves and Filter Cloud Application associated with a SSL inspection rule by application category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class   = ["AI_ML"]
}

output "app_ids" {
  value = [for app in data.zia_cloud_applications.this.applications : app["app"]]
}

# Retrieves specific application associated with a SSL inspection rule by name and category
data "zia_cloud_applications" "this" {
  policy_type = "cloud_application_ssl_policy"
  app_class = ["SOCIAL_NETWORKING"]
  app_name = "Nebenan"
}

output "zia_cloud_applications" {
    value = data.zia_cloud_applications.this
}
