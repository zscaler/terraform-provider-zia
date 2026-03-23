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

resource "zia_virtual_service_edge_node" "this2" {
  name                              = "VSENode02"
  status                            = "ENABLED"
  type                              = "VZEN"
  ip_address                        = "10.0.0.11"
  subnet_mask                       = "255.255.255.0"
  default_gateway                   = "10.0.0.1"
  in_production                     = true
  load_balancer_ip_address          = "10.0.0.20"
  deployment_mode                   = "CLUSTER"
  vzen_sku_type                     = "LARGE"
}

resource "zia_virtual_service_edge_cluster" "this" {
  name  = "VSECluster01"
  status = "ENABLED"
  type = "VIP"
  ip_address = "10.0.0.2"
  subnet_mask = "255.255.255.0"
  default_gateway = "10.0.0.3"
  ip_sec_enabled = true
  virtual_zen_nodes {
    id = [zia_virtual_service_edge_node.this1.id, zia_virtual_service_edge_node.this2.id]
  }
}