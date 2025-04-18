resource "zia_admin_roles" "this" {
  name               = "AdminRoleTerraform"
  rank               = 7
  alerting_access    = "READ_WRITE"
  dashboard_access   = "READ_WRITE"
  report_access      = "READ_WRITE"
  analysis_access    = "READ_ONLY"
  username_access    = "READ_ONLY"
  device_info_access = "READ_ONLY"
  admin_acct_access  = "READ_WRITE"
  policy_access      = "READ_WRITE"
  permissions = [
 "NSS_CONFIGURATION", "LOCATIONS", "HOSTED_PAC_FILES", "EZ_AGENT_CONFIGURATIONS",
 "SECURE_AGENT_NOTIFICATIONS", "VPN_CREDENTIALS", "AUTHENTICATION_SETTINGS", "STATIC_IPS",
 "GRE_TUNNELS", "CLIENT_CONNECTOR_PORTAL", "SECURE", "POLICY_RESOURCE_MANAGEMENT",
 "CUSTOM_URL_CAT", "OVERRIDE_EXISTING_CAT", "TENANT_PROFILE_MANAGEMENT", "COMPLY",
 "SSL_POLICY", "ADVANCED_SETTINGS", "PROXY_GATEWAY", "SUBCLOUDS", "IDENTITY_PROXY_SETTINGS",
 "USER_MANAGEMENT", "APIKEY_MANAGEMENT", "FIREWALL_DNS", "VZEN_CONFIGURATION",
 "PARTNER_INTEGRATION", "USER_ACCESS", "CUSTOMER_ACCT_INFO", "CUSTOMER_SUBSCRIPTION",
 "CUSTOMER_ORG_SETTINGS", "ZIA_TRAFFIC_CAPTURE", "REMOTE_ASSISTANCE_MANAGEMENT"
  ]
}


resource "zia_admin_roles" "this" {
  name               = "SDWANAdminRoleTerraform"
  rank               = 7
  policy_access      = "READ_WRITE"
  alerting_access    = "NONE"
  role_type          = "SDWAN"
}