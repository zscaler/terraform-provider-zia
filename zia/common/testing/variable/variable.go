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
	GRETunnelComment       = "this is an acceptance test"
	GRETunnelWithinCountry = false
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
	DLPOCREnabled            = false
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
	AdminUserLoginName = "tf-acc-test-@bd-hashicorp.com"
	AdminUserName      = "tf-acc-test-admin"
	AdminUserEmail     = "tf-acc-test-@bd-hashicorp.com"
	AdminUserPassword  = "Password@123!"
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

// Activation
const (
	Status = "ACTIVE"
)
