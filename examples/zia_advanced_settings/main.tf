resource "zia_advanced_settings" "this" {
  auth_bypass_apps                                            = []
  auth_bypass_urls                                            = [".newexample1.com", ".newexample2.com"]
  dns_resolution_on_transparent_proxy_apps                    = ["CHATGPT_AI"]
  basic_bypass_url_categories                                 = ["NONE"]
  http_range_header_remove_url_categories                     = ["NONE"]
  kerberos_bypass_urls                                        = ["test1.com"]
  kerberos_bypass_apps                                        = []
  dns_resolution_on_transparent_proxy_urls                    = ["test1.com", "test2.com"]
  ui_session_timeout                                          = 300
  enable_dns_resolution_on_transparent_proxy                  = true
  enable_evaluate_policy_on_global_ssl_bypass                 = true
  enable_office365                                            = true
  log_internal_ip                                             = true
  enforce_surrogate_ip_for_windows_app                        = true
  track_http_tunnel_on_http_ports                             = true
  block_http_tunnel_on_non_http_ports                         = false
  block_domain_fronting_on_host_header                        = false
  zscaler_client_connector_1_and_pac_road_warrior_in_firewall = true
  cascade_url_filtering                                       = true
  enable_policy_for_unauthenticated_traffic                   = true
  block_non_compliant_http_request_on_http_ports              = true
  enable_admin_rank_access                                    = true
  http2_nonbrowser_traffic_enabled                            = true
  ecs_for_all_enabled                                         = false
  dynamic_user_risk_enabled                                   = false
  block_connect_host_sni_mismatch                             = false
  prefer_sni_over_conn_host                                   = false
  sipa_xff_header_enabled                                     = false
  block_non_http_on_http_port_enabled                         = true
}
