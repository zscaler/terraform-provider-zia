package zia

import (
	"log"

	client "github.com/zscaler/zscaler-sdk-go/zia"
	"github.com/zscaler/zscaler-sdk-go/zia/services/activation"
	"github.com/zscaler/zscaler-sdk-go/zia/services/adminuserrolemgmt"
	"github.com/zscaler/zscaler-sdk-go/zia/services/devicegroups"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_engines"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_icap_servers"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_idm_profiles"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_incident_receiver_servers"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_notification_templates"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/zia/services/dlpdictionaries"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/applicationservices"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/applicationservicesgroup"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/networkapplications"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/zia/services/firewallpolicies/timewindow"
	"github.com/zscaler/zscaler-sdk-go/zia/services/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/zia/services/locationmanagement/locationgroups"
	"github.com/zscaler/zscaler-sdk-go/zia/services/locationmanagement/locationlite"
	"github.com/zscaler/zscaler-sdk-go/zia/services/rule_labels"
	"github.com/zscaler/zscaler-sdk-go/zia/services/security_policy_settings"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/greinternalipranges"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/gretunnelinfo"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/virtualipaddresslist"
	"github.com/zscaler/zscaler-sdk-go/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/zia/services/urlcategories"
	"github.com/zscaler/zscaler-sdk-go/zia/services/urlfilteringpolicies"
	"github.com/zscaler/zscaler-sdk-go/zia/services/user_authentication_settings"
	"github.com/zscaler/zscaler-sdk-go/zia/services/usermanagement"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	adminuserrolemgmt             *adminuserrolemgmt.Service
	filteringrules                *filteringrules.Service
	ipdestinationgroups           *ipdestinationgroups.Service
	ipsourcegroups                *ipsourcegroups.Service
	networkapplications           *networkapplications.Service
	networkservices               *networkservices.Service
	applicationservices           *applicationservices.Service
	applicationservicesgroup      *applicationservicesgroup.Service
	timewindow                    *timewindow.Service
	urlcategories                 *urlcategories.Service
	urlfilteringpolicies          *urlfilteringpolicies.Service
	usermanagement                *usermanagement.Service
	gretunnels                    *gretunnels.Service
	gretunnelinfo                 *gretunnelinfo.Service
	greinternalipranges           *greinternalipranges.Service
	staticips                     *staticips.Service
	virtualipaddresslist          *virtualipaddresslist.Service
	vpncredentials                *vpncredentials.Service
	locationmanagement            *locationmanagement.Service
	locationgroups                *locationgroups.Service
	locationlite                  *locationlite.Service
	activation                    *activation.Service
	devicegroups                  *devicegroups.Service
	dlpdictionaries               *dlpdictionaries.Service
	dlp_engines                   *dlp_engines.Service
	dlp_idm_profiles              *dlp_idm_profiles.Service
	dlp_icap_servers              *dlp_icap_servers.Service
	dlp_incident_receiver_servers *dlp_incident_receiver_servers.Service
	dlp_notification_templates    *dlp_notification_templates.Service
	dlp_web_rules                 *dlp_web_rules.Service
	rule_labels                   *rule_labels.Service
	security_policy_settings      *security_policy_settings.Service
	user_authentication_settings  *user_authentication_settings.Service
}

type Config struct {
	Username   string
	Password   string
	APIKey     string
	ZIABaseURL string
	UserAgent  string
}

func (c *Config) Client() (*Client, error) {
	cli, err := client.NewClient(c.Username, c.Password, c.APIKey, c.ZIABaseURL, c.UserAgent)
	if err != nil {
		return nil, err
	}

	ziaClient := &Client{
		adminuserrolemgmt:             adminuserrolemgmt.New(cli),
		filteringrules:                filteringrules.New(cli),
		ipdestinationgroups:           ipdestinationgroups.New(cli),
		ipsourcegroups:                ipsourcegroups.New(cli),
		networkapplications:           networkapplications.New(cli),
		networkservices:               networkservices.New(cli),
		applicationservices:           applicationservices.New(cli),
		applicationservicesgroup:      applicationservicesgroup.New(cli),
		timewindow:                    timewindow.New(cli),
		urlcategories:                 urlcategories.New(cli),
		urlfilteringpolicies:          urlfilteringpolicies.New(cli),
		usermanagement:                usermanagement.New(cli),
		virtualipaddresslist:          virtualipaddresslist.New(cli),
		vpncredentials:                vpncredentials.New(cli),
		gretunnels:                    gretunnels.New(cli),
		gretunnelinfo:                 gretunnelinfo.New(cli),
		greinternalipranges:           greinternalipranges.New(cli),
		staticips:                     staticips.New(cli),
		locationmanagement:            locationmanagement.New(cli),
		locationgroups:                locationgroups.New(cli),
		locationlite:                  locationlite.New(cli),
		activation:                    activation.New(cli),
		devicegroups:                  devicegroups.New(cli),
		dlpdictionaries:               dlpdictionaries.New(cli),
		dlp_engines:                   dlp_engines.New(cli),
		dlp_idm_profiles:              dlp_idm_profiles.New(cli),
		dlp_icap_servers:              dlp_icap_servers.New(cli),
		dlp_incident_receiver_servers: dlp_incident_receiver_servers.New(cli),
		dlp_notification_templates:    dlp_notification_templates.New(cli),
		dlp_web_rules:                 dlp_web_rules.New(cli),
		rule_labels:                   rule_labels.New(cli),
		security_policy_settings:      security_policy_settings.New(cli),
		user_authentication_settings:  user_authentication_settings.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
