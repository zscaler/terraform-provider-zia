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
	CategoryDescription = "this is an acceptance test"
	ConfiguredName      = "this is an acceptance test"
	CustomCategory      = true
)

// Firewall Filtering Rule resource/datasource
const (
	FWRuleResourceName        = "this is an acceptance test"
	FWRuleResourceDescription = "this is an acceptance test"
	FWRuleResourceAction      = "ALLOW"
	FWRuleResourceState       = "ENABLED"
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
)

// Traffic Forwarding Static IP resource/datasource
const (
	StaticIPComment  = "this is an acceptance test"
	StaticIPAddress  = "118.189.211.221"
	StaticRoutableIP = true
)

// Traffic Forwarding VPN Credentials resource/datasource
const (
	VPNCredentialComments     = "this is an acceptance test"
	VPNCredentialTypeUFQDN    = "UFQDN"
	VPNCredentialFQDN         = "test@securitygeek.io"
	VPNCredentialPreSharedKey = "Password@123!"
	VPNCredentialTypeIP       = "IP"
	VPNCredentialIPAddress    = "118.189.211.221"
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

// Rule Labels
const (
	RuleLabelName        = "testAcc_rule_label"
	RuleLabelDescription = "testAcc_rule_label"
)
