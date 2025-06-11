resource "zia_virtual_service_edge_cluster" "this" {
  name  = "VSECluster01"
  status = "ENABLED"
  type = "VIP"
  ip_address = "10.0.0.2"
  subnet_mask = "255.255.255.0"
  default_gateway = "10.0.0.3"
  ip_sec_enabled = true
  virtual_zen_nodes {
    id = [9368]
  }
}