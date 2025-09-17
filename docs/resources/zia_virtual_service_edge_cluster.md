---
subcategory: "Virtual Service Edges"
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

* `name` - (Optional) Name of the Virtual Service Edge cluster
* `status` - (Optional) Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default
* `type` - (Optional) The Virtual Service Edge cluster type. Supported values: `ANY`, `NONE`, `SME`, `SMSM`, `SMCA`, `SMUI`, `SMCDS`, `SMDNSD`, `SMAA`, `SMTP`, `SMQTN`, `VIP`, `UIZ`, `UIAE`, `SITEREVIEW`, `PAC`, `S_RELAY`, `M_RELAY`, `H_MON`, `SMIKE`, `NSS`, `SMEZA`, `SMLB`, `SMFCCLT`, `SMBA`, `SMBAC`, `SMESXI`, `SMBAUI`, `VZEN`, `ZSCMCLT`, `SMDLP`, `ZSQUERY`, `ADP`, `SMCDSDLP`, `SMSCIM`, `ZSAPI`, `ZSCMCDSSCLT`, `LOCAL_MTS`, `SVPN`, `SMCASB`, `SMFALCONUI`, `MOBILEAPP_REG`, `SMRESTSVR`, `FALCONCA`, `MOBILEAPP_NF`, `ZIRSVR`, `SMEDGEUI`, `ALERTEVAL`, `ALERTNOTIF`, `SMPARTNERUI`, `CQM`, `DATAKEEPER`, `SMBAM`, `ZWACLT`
* `ip_sec_enabled` - (Optional) A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - (Optional) The Virtual Service Edge cluster IP address
* `subnet_mask` - (Optional) The Virtual Service Edge cluster subnet mask
* `default_gateway` - (Optional) The IP address of the default gateway to the internet
* `virtual_zen_nodes` - (Optional) List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)

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
