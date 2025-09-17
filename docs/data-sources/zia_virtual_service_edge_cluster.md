---
subcategory: "Virtual Service Edges"
layout: "zscaler"
page_title: "ZIA: virtual_service_edge_cluster"
description: |-
    Official documentation https://help.zscaler.com/zia/about-virtual-service-edge-clusters
    API documentation https://help.zscaler.com/zia/service-edges#/virtualZenClusters-get
   Retrieves a list of ZIA Virtual Service Edge clusters.
---

# zia_virtual_service_edge_cluster (Data Source)

* [Official documentation](https://help.zscaler.com/zia/about-virtual-service-edge-clusters)
* [API documentation](https://help.zscaler.com/zia/service-edges#/virtualZenClusters-get)

Use the **zia_virtual_service_edge_cluster** data source to get information about a Virtual Service Edge Cluster information for the specified `Name` or `ID`

## Example Usage

```hcl
data "zia_virtual_service_edge_cluster" "this"{
    name = "VSECluster01"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) Name of the Virtual Service Edge cluster
* `id` - (Optional) System-generated Virtual Service Edge cluster ID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - System-generated Virtual Service Edge cluster ID
* `cluster_id` - System-generated Virtual Service Edge cluster ID
* `name` - Name of the Virtual Service Edge cluster
* `status` - Specifies the status of the Virtual Service Edge cluster. The status is set to ENABLED by default
* `type` - The Virtual Service Edge cluster type. Supported values: `ANY`, `NONE`, `SME`, `SMSM`, `SMCA`, `SMUI`, `SMCDS`, `SMDNSD`, `SMAA`, `SMTP`, `SMQTN`, `VIP`, `UIZ`, `UIAE`, `SITEREVIEW`, `PAC`, `S_RELAY`, `M_RELAY`, `H_MON`, `SMIKE`, `NSS`, `SMEZA`, `SMLB`, `SMFCCLT`, `SMBA`, `SMBAC`, `SMESXI`, `SMBAUI`, `VZEN`, `ZSCMCLT`, `SMDLP`, `ZSQUERY`, `ADP`, `SMCDSDLP`, `SMSCIM`, `ZSAPI`, `ZSCMCDSSCLT`, `LOCAL_MTS`, `SVPN`, `SMCASB`, `SMFALCONUI`, `MOBILEAPP_REG`, `SMRESTSVR`, `FALCONCA`, `MOBILEAPP_NF`, `ZIRSVR`, `SMEDGEUI`, `ALERTEVAL`, `ALERTNOTIF`, `SMPARTNERUI`, `CQM`, `DATAKEEPER`, `SMBAM`, `ZWACLT`
* `ip_sec_enabled` - A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - The Virtual Service Edge cluster IP address
* `subnet_mask` - The Virtual Service Edge cluster subnet mask
* `default_gateway` - The IP address of the default gateway to the internet
* `virtual_zen_nodes` - List of ZPA Application Segments for which this rule is applicable. This field is applicable only for the ECZPA forwarding method (used for Zscaler Cloud Connector)
