package variable

// url filtering rules
const (
	URLFilteringRuleResourceName = "testAcc_url_filtering_rule"
	URLFilteringRuleDescription  = "testAcc_url_filtering_rule"
	URLFilteringRuleAction       = "ALLOW"
	URLFilteringRuleActionUpdate = "BLOCK"
	URLFilteringRuleState        = "ENABLED"
	URLFilteringRuleStateUpdate  = "DISABLED"
)

// Firewall IPS rules
const (
	FWIPSRuleResourceName = "testAcc_firewall_ips_rule"
	FWIPSRuleDescription  = "testAcc_firewall_ips_rule"
	FWIPSAction           = "ALLOW"
	FWIPSActionUpdate     = "BLOCK"
	FWIPSState            = "ENABLED"
	FWIPSUpdate           = "DISABLED"
)

// Firewall DNS rules
const (
	FWDNSRuleResourceName = "testAcc_firewall_dns_rule"
	FWDNSRuleDescription  = "testAcc_firewall_dns_rule"
	FWDNSAction           = "REDIR_RES"
	FWDNSActionUpdate     = "BLOCK"
	FWDNSState            = "ENABLED"
	FWDNSUpdate           = "DISABLED"
)

// File Type Control rules
const (
	FileTypeControlRuleResourceName = "testAcc_firewall_dns_rule"
	FileTypeControlRuleDescription  = "testAcc_firewall_dns_rule"
	FileTypeControlRuleAction       = "ALLOW"
	FileTypeControlRuleState        = "ENABLED"
)

// Sandbox rules
const (
	SandboxRuleResourceName = "testAcc_sandbox_rule"
	SandboxRuleDescription  = "testAcc_sandbox_rule"
	SandboxAction           = "ALLOW"
	SandboxActionUpdate     = "BLOCK"
	SandboxState            = "ENABLED"
	SandboxStateUpdate      = "DISABLED"
)

// SSL Inspection rules
const (
	SSLInspectionRuleName        = "testAcc_ssl_rule"
	SSLInspectionRuleDescription = "testAcc_ssl_rule"
	SSLInspectionRuleState       = "ENABLED"
	RoadWarriorKerberos          = true
)

// Custom URL Categories resource/datasource
const (
	CustomCategory = true
)

// Firewall Filtering Rule resource/datasource
const (
	FWRuleResourceName        = "this is an acceptance test"
	FWRuleResourceDescription = "this is an acceptance test"
	FWRuleResourceAction      = "ALLOW"
	FWRuleResourceState       = "ENABLED"
	FWRuleOrder               = "1"
	FWRuleResourceStateUpdate = "DISABLED"
	FWRuleEnableLogging       = false
)

// Firewall Filtering IP Destination Group resource/datasource
const (
	FWDSTGroupName         = "this is an acceptance test"
	FWDSTGroupDescription  = "this is an acceptance test"
	FWDSTGroupTypeDSTNFQDN = "DSTN_FQDN"
)

// Firewall Filtering IP Source Group resource/datasource
const (
	FWSRCGroupName        = "this is an acceptance test"
	FWSRCGroupDescription = "this is an acceptance test"
)

// Firewall network application groups resource/datasource
const (
	FWAppGroupName        = "this is an acceptance test"
	FWAppGroupDescription = "this is an acceptance test"
)

// Firewall network services groups resource/datasource
const (
	FWNetworkServicesGroupName        = "this is an acceptance test"
	FWNetworkServicesGroupDescription = "this is an acceptance test"
	// FWNetworkServices = "this is an acceptance test"
)

// Firewall network services resource/datasource
const (
	FWNetworkServicesName        = "this is an acceptance test"
	FWNetworkServicesDescription = "this is an acceptance test"
	FWNetworkServicesType        = "CUSTOM"
)

// Forwarding Control ZPA Gateway
const (
	FowardingControlDescription = "this is an acceptance test"
	FowardingControlType        = "FORWARDING"
	FowardingControlState       = "ENABLED"
)

// Forwarding Control Rule
const (
	FowardingControlUpdateDescription = "this is an updated acceptance test"
	FowardingControlUpdateState       = "ENABLED"
)

// Traffic Forwarding Static IP resource/datasource
const (
	StaticIPComment   = "this is an acceptance test"
	StaticRoutableIP  = true
	StaticGeoOverride = false
)

// Traffic Forwarding VPN Credentials resource/datasource
const (
	VPNCredentialTypeUFQDN = "UFQDN"
	VPNCredentialTypeIP    = "IP"
)

// Traffic Forwarding Location Management
const (
	LocName         = "this is an acceptance test"
	LocDesc         = "this is an acceptance test"
	LocAuthRequired = true
	LocSurrogateIP  = true
	LocXFF          = true
	LocOFW          = true
	LocOFWUpdate    = false
	LocIPS          = true
	LocIPSUpdate    = false
)

// Traffic Forwarding GRE resource/datasource
const (
	GRETunnelComment       = "GRE Tunnel Created with Terraform"
	GRETunnelWithinCountry = true
	GRETunnelIPUnnumbered  = false
)

// DLP Dictionaries resource/datasource
const (
	DLPWebRuleName           = "this is an acceptance test"
	DLPWebRuleDesc           = "this is an acceptance test"
	DLPRuleResourceAction    = "ALLOW"
	DLPRuleResourceState     = "ENABLED"
	DLPRuleContentInspection = false
	DLPMatchOnly             = false
)

// DLP Dictionaries resource/datasource
const (
	DLPDictionaryResourceName = "this is an acceptance test"
	DLPDictionaryDescription  = "this is an acceptance test"
)

// DLP Engines resource/datasource
const (
	DLPCustomEngine = true
)

// DLP Dictionaries resource/datasource
const (
	DLPNoticationTemplateAttachContent = true
	DLPNoticationTemplateTLSEnabled    = true
)

// Admin Users
const (
	AdminUserLoginName = "tf-acc-test-@securitygeek.io"
	AdminUserName      = "tf-acc-test-admin"
	AdminUserEmail     = "tf-acc-test-@securitygeek.io"
)

// User Management
const (
	UserName = "tf-acc-test-user"
)

// Rule Labels
const (
	RuleLabelName        = "testAcc_rule_label"
	RuleLabelDescription = "testAcc_rule_label"
)

// Admin Roles
const (
	RoleName = "testAcc_admin_role"
)

// Activation
const (
	Status = "ACTIVE"
)

// Cloud Application Control Rules
const (
	CloudAppControlRuleResourceName = "tf-acc-test-cloud-app-control"
	CloudAppControlRuleDescription  = "tf-acc-test-cloud-app-control"
	CloudAppControlRuleType         = "STREAMING_MEDIA"
	CloudAppControlRuleState        = "ENABLED"
)

const (
	AlertEmail       = "alert@acme.com"
	AlertDescription = "testAcc_Alert_Subscription"
)

// Forwarding Control Proxies
const (
	ProxyDescription           = "testAcc_proxy"
	ProxyType                  = "PROXYCHAIN"
	ProxyAddress               = "192.168.1.150"
	ProxyPort                  = 5000
	ProxyInsertXauHeader       = true
	ProxyBase64EncodeXauHeader = true
)

// NAT Control rules
const (
	NATControlRuleName        = "testAcc_nat_control_rule"
	NATControlRuleDescription = "testAcc_nat_control_rule"
	NATControlRuleState       = "ENABLED"
	NATControlRuleUpdate      = "DISABLED"
	NATControlRuleLogging     = true
	NATControlRedirectPort    = 5000
	NATControlRedirectIP      = "192.168.100.150"
)

// Service Edge Cluster
const (
	VzenType           = "VIP"
	VzenStatus         = "ENABLED"
	VzenIPAddress      = "10.0.0.2"
	VzenSubnetMask     = "255.255.255.0"
	VzenDefaultGateway = "10.0.0.3"
	VzenIpSecEnabled   = true
)

// Service Edge Node
const (
	VzenNodeType              = "VZEN"
	VzenNodeStatus            = "ENABLED"
	VzenNodeIPAddress         = "10.0.0.10"
	VzenNodeSubnetMask        = "255.255.255.0"
	VzenNodeDefaultGateway    = "10.0.0.20"
	VzenNodeLoadBalancer      = "10.0.0.30"
	VzenNodeDeploymentMode    = "STANDALONE"
	VZenSKUType               = "LARGE"
	VzenNodeInProduction      = true
	VzenOnDemandSupportTunnel = true
)

// Rule Labels
const (
	NSSStatus = "ENABLED"
	NSSType   = "NSS_FOR_FIREWALL"
)

const (
	BandwdithControlRuleDescription = "tf-acc-test-cloud-app-control"
	BandwdithControlRulestate       = "ENABLED"
)

const (
	UrlCatReviewEnabled                  = true
	EunUrlCatReviewSubmitToSecurityCloud = false
	EunSecurityReviewEnabled             = true
	EunWebDlpReviewEnabled               = true
	EunWebDlpReviewSubmitToSecurityCloud = false
	DisplayCompReason                    = true
	DisplayCompName                      = true
	DisplayCompLogo                      = true
)

// DC Exclusions
const (
	DCExclusionsDescription = "tf-acc-dc-exclusions"
)

// Rule Labels
const (
	ExtranetName               = "testAcc_extranet"
	ExtranetDescription        = "testAcc_extranet"
	ExtranetDNSName            = "testAcc_extranet_dns"
	ExtranetDNSServer          = "192.168.1.1"
	ExtranetDNSServer2         = "192.168.1.2"
	ExtranetDNSUseAsDefault    = true
	ExtranetIPPoolName         = "testAcc_extranet_ip_pool"
	ExtranetIPPoolStart        = "192.168.1.1"
	ExtranetIPPoolEnd          = "192.168.1.2"
	ExtranetIPPoolUseAsDefault = true
)
