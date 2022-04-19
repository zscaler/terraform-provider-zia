package variable

// url filtering rules
const (
	URLFilteringRuleResourceName = "testAcc_url_filtering_rule"
	URLFilteringRuleDescription  = "testAcc_url_filtering_rule"
	URLFilteringRuleAction       = "ALLOW"
	URLFilteringRuleState        = "ENABLED"
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

// Traffic Forwarding Static IP resource/datasource
const (
	StaticIPComment   = "this is an acceptance test"
	StaticRoutableIP  = true
	StaticGeoOverride = true
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
	LocIPS          = true
)

// DLP Dictionaries resource/datasource
const (
	DLPWebRuleName           = "this is an acceptance test"
	DLPWebRuleDesc           = "this is an acceptance test"
	DLPRuleResourceAction    = "ALLOW"
	DLPRuleResourceState     = "ENABLED"
	DLPRuleContentInspection = false
	DLPMatchOnly             = false
	DLPOCREnabled            = true
)

// DLP Dictionaries resource/datasource
const (
	DLPDictionaryResourceName = "this is an acceptance test"
	DLPDictionaryDescription  = "this is an acceptance test"
)

// Admin Users
const (
	AdminUserLoginName = "testAcc@securitygeek.io"
	AdminUserName      = "Test Acc"
	AdminUserEmail     = "testAcc@securitygeek.io"
	AdminUserPassword  = "Password@123!"
)

// User Management
const (
	UserName = "testAcc TF User"
)

// Rule Labels
const (
	RuleLabelName        = "testAcc_rule_label"
	RuleLabelDescription = "testAcc_rule_label"
)
