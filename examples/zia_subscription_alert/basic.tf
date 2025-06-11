resource "zia_subscription_alert" "this" {
  email  = "alert@acme.com"
  description = "Terraform Alert"
  pt0_severities = ["CRITICAL"]
  secure_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  manage_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  comply_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
  system_severities = ["CRITICAL", "MAJOR", "MINOR", "INFO", "DEBUG"]
}
