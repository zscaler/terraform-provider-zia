---
subcategory: "Service Edge Cluster"
layout: "zscaler"
page_title: "ZIA: virtual_service_edge_cluster"
description: |-
   Retrieves a list of ZIA Virtual Service Edge clusters.
---
# Data Source: zia_virtual_service_edge_cluster

Use the **zia_virtual_service_edge_cluster** data source to get information about a Virtual Service Edge Cluster information for the specified `Name` or `ID`

```hcl
data "zia_virtual_service_edge_cluster" "this"{
    name = "VSECluster01"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (String) Name of the Virtual Service Edge cluster
* `id` - (String) USystem-generated Virtual Service Edge cluster ID

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `status` - (Number) Specifies the status of the Virtual Service Edge cluster. The status is set to `ENABLED` by default.

* `ip_sec_enabled` - (String) A Boolean value that specifies whether to terminate IPSec traffic from the client at selected Virtual Service Edge instances for the Virtual Service Edge cluster
* `ip_address` - (String) The Virtual Service Edge cluster IP address. In a Virtual Service Edge cluster, the cluster IP address provides fault tolerance and is used to listen for user traffic. This interface doesn't explicitly get an IP address. The cluster IP address must be in the same VLAN as the proxy and load balancer IP addresses.
* `subnet_mask` - (String) The Virtual Service Edge cluster subnet mask
* `default_gateway` - (String) The IP address of the default gateway to the internet
* `last_modified_time` - (Number) When the cluster was last modified

* `virtual_zen_nodes` - (List of Object) The Virtual Service Edge instances you want to include in the cluster. A Virtual Service Edge cluster must contain at least two Virtual Service Edge instances.
  * `id` - (Number) Identifier that uniquely identifies an entity
  * `name` - (String) The configured name of the entity
  * `external_id` (String) An external identifier used for an entity that is managed outside of ZIA.
  * `extensions` - (Map of String)

* `type` - (String) The Virtual Service Edge cluster type
`ANY`, `NONE`, `SME`, `SMSM`, `SMCA`, `SMUI`, `SMCDS`, `SMDNSD`, `SMAA`, `SMTP`,`SMQTN`,`VIP`,
`UIZ`, `UIAE`, `SITEREVIEW`, `PAC`, `S_RELAY`, `M_RELAY`, `H_MON`, `SMIKE`, `NSS`, `SMEZA`, `SMLB`,
`SMFCCLT`, `SMBA`, `SMBAC`, `SMESXI`, `SMBAUI`, `VZEN`, `ZSCMCLT`, `SMDLP`, `ZSQUERY`, `ADP`, `SMCDSDLP`,
`SMSCIM`, `ZSAPI`, `ZSCMCDSSCLT`, `LOCAL_MTS`, `SVPN`, `SMCASB`, `SMFALCONUI`, `MOBILEAPP_REG`, `SMRESTSVR`,
`FALCONCA`, `MOBILEAPP_NF`, `ZIRSVR`, `SMEDGEUI`, `ALERTEVAL`, `ALERTNOTIF`, `SMPARTNERUI`, `CQM`, `DATAKEEPER`,
`SMBAM`, `ZWACLT`
