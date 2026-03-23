resource "zia_virtual_service_edge_node" "this1" {
  name                              = "VSENode01"
  status                            = "ENABLED"
  type                              = "VZEN"
  ip_address                        = "10.0.0.10"
  subnet_mask                       = "255.255.255.0"
  default_gateway                   = "10.0.0.1"
  in_production                     = true
  load_balancer_ip_address          = "10.0.0.20"
  deployment_mode                   = "CLUSTER"
  vzen_sku_type                     = "LARGE"
}