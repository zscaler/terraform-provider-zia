---
subcategory: "Service Edge Cluster"
layout: "zscaler"
page_title: "ZIA: virtual_service_edge_cluster"
description: |-
    Official documentation https://help.zscaler.com/zia/about-virtual-service-edge-clusters
    API documentation https://help.zscaler.com/zia/service-edges#/virtualZenClusters-get
    Adds a new Virtual Service Edge cluster.
---

# zia_virtual_service_edge_cluster (Resource)

* [Official documentation](https://help.zscaler.com/zia/about-virtual-service-edge-clusters)
* [API documentation](https://help.zscaler.com/zia/service-edges#/virtualZenClusters-get)

Use the **zia_virtual_service_edge_cluster** resource allows the creation and management of Service Edge Cluster objects in the Zscaler Internet Access.

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The email address of the alert recipient

## Attribute Reference

In addition to all arguments above, the following attributes are supported:

* `status` - (Number) Specifies the status of the Virtual Service Edge cluster. The status is set to `ENABLED` by default.

* `ip_sec_enabled` - (String) A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - (String) The Virtual Service Edge cluster IP address. In a Virtual Service Edge cluster, the cluster IP address provides fault tolerance and is used to listen for user traffic. This interface doesn't explicitly get an IP address. The cluster IP address must be in the same VLAN as the proxy and load balancer IP addresses.
* `subnet_mask` - (String) The Virtual Service Edge cluster subnet mask
* `default_gateway` - (String) The IP address of the default gateway to the internet
* `last_modified_time` - (Number) When the cluster was last modified

* `virtual_zen_nodes` - (List of Object) The Virtual Service Edge instances you want to include in the cluster. A Virtual Service Edge cluster must contain at least two Virtual Service Edge instances.
  * `id` - (Number) Identifier that uniquely identifies an entity

* `type` - (String) The Virtual Service Edge cluster type
`ANY`, `NONE`, `SME`, `SMSM`, `SMCA`, `SMUI`, `SMCDS`, `SMDNSD`, `SMAA`, `SMTP`,`SMQTN`,`VIP`,
`UIZ`, `UIAE`, `SITEREVIEW`, `PAC`, `S_RELAY`, `M_RELAY`, `H_MON`, `SMIKE`, `NSS`, `SMEZA`, `SMLB`,
`SMFCCLT`, `SMBA`, `SMBAC`, `SMESXI`, `SMBAUI`, `VZEN`, `ZSCMCLT`, `SMDLP`, `ZSQUERY`, `ADP`, `SMCDSDLP`,
`SMSCIM`, `ZSAPI`, `ZSCMCDSSCLT`, `LOCAL_MTS`, `SVPN`, `SMCASB`, `SMFALCONUI`, `MOBILEAPP_REG`, `SMRESTSVR`, `FALCONCA`, `MOBILEAPP_NF`, `ZIRSVR`, `SMEDGEUI`, `ALERTEVAL`, `ALERTNOTIF`, `SMPARTNERUI`, `CQM`, `DATAKEEPER`,`SMBAM`, `ZWACLT`

## Import

Zscaler offers a dedicated tool called Zscaler-Terraformer to allow the automated import of ZIA configurations into Terraform-compliant HashiCorp Configuration Language.
[Visit](https://github.com/zscaler/zscaler-terraformer)

**zia_virtual_service_edge_cluster** can be imported by using `<CLUSTER_ID>` or `<CLUSTER_NAME>` as the import ID.

For example:

```shell
terraform import zia_virtual_service_edge_cluster.example <cluster_id>
```

or

```shell
terraform import zia_virtual_service_edge_cluster.example <cluster_name>
```
