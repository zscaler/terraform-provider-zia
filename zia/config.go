package zia

import (
	"log"

	"github.com/zscaler/terraform-provider-zia/gozscaler/activation"
	"github.com/zscaler/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/zscaler/terraform-provider-zia/gozscaler/client"
	"github.com/zscaler/terraform-provider-zia/gozscaler/devicegroups"
	"github.com/zscaler/terraform-provider-zia/gozscaler/dlp_engines"
	"github.com/zscaler/terraform-provider-zia/gozscaler/dlp_notification_templates"
	"github.com/zscaler/terraform-provider-zia/gozscaler/dlp_web_rules"
	"github.com/zscaler/terraform-provider-zia/gozscaler/dlpdictionaries"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/filteringrules"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/networkapplications"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
	"github.com/zscaler/terraform-provider-zia/gozscaler/firewallpolicies/timewindow"
	"github.com/zscaler/terraform-provider-zia/gozscaler/locationmanagement"
	"github.com/zscaler/terraform-provider-zia/gozscaler/locationmanagement/locationgroups"
	"github.com/zscaler/terraform-provider-zia/gozscaler/rule_labels"
	"github.com/zscaler/terraform-provider-zia/gozscaler/security_policy_settings"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/greinternalipranges"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/gretunnelinfo"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
	"github.com/zscaler/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
	"github.com/zscaler/terraform-provider-zia/gozscaler/urlcategories"
	"github.com/zscaler/terraform-provider-zia/gozscaler/urlfilteringpolicies"
	"github.com/zscaler/terraform-provider-zia/gozscaler/usermanagement"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	adminuserrolemgmt          *adminuserrolemgmt.Service
	filteringrules             *filteringrules.Service
	ipdestinationgroups        *ipdestinationgroups.Service
	ipsourcegroups             *ipsourcegroups.Service
	networkapplications        *networkapplications.Service
	networkservices            *networkservices.Service
	timewindow                 *timewindow.Service
	urlcategories              *urlcategories.Service
	urlfilteringpolicies       *urlfilteringpolicies.Service
	usermanagement             *usermanagement.Service
	gretunnels                 *gretunnels.Service
	gretunnelinfo              *gretunnelinfo.Service
	greinternalipranges        *greinternalipranges.Service
	staticips                  *staticips.Service
	virtualipaddresslist       *virtualipaddresslist.Service
	vpncredentials             *vpncredentials.Service
	locationmanagement         *locationmanagement.Service
	locationgroups             *locationgroups.Service
	activation                 *activation.Service
	devicegroups               *devicegroups.Service
	dlpdictionaries            *dlpdictionaries.Service
	dlp_engines                *dlp_engines.Service
	dlp_notification_templates *dlp_notification_templates.Service
	dlp_web_rules              *dlp_web_rules.Service
	rule_labels                *rule_labels.Service
	security_policy_settings   *security_policy_settings.Service
}

type Config struct {
	Username   string
	Password   string
	APIKey     string
	ZIABaseURL string
}

func (c *Config) Client() (*Client, error) {
	cli, err := client.NewClientZIA(c.Username, c.Password, c.APIKey, c.ZIABaseURL)
	if err != nil {
		return nil, err
	}

	ziaClient := &Client{
		adminuserrolemgmt:          adminuserrolemgmt.New(cli),
		filteringrules:             filteringrules.New(cli),
		ipdestinationgroups:        ipdestinationgroups.New(cli),
		ipsourcegroups:             ipsourcegroups.New(cli),
		networkapplications:        networkapplications.New(cli),
		networkservices:            networkservices.New(cli),
		timewindow:                 timewindow.New(cli),
		urlcategories:              urlcategories.New(cli),
		urlfilteringpolicies:       urlfilteringpolicies.New(cli),
		usermanagement:             usermanagement.New(cli),
		virtualipaddresslist:       virtualipaddresslist.New(cli),
		vpncredentials:             vpncredentials.New(cli),
		gretunnels:                 gretunnels.New(cli),
		gretunnelinfo:              gretunnelinfo.New(cli),
		greinternalipranges:        greinternalipranges.New(cli),
		staticips:                  staticips.New(cli),
		locationmanagement:         locationmanagement.New(cli),
		locationgroups:             locationgroups.New(cli),
		activation:                 activation.New(cli),
		devicegroups:               devicegroups.New(cli),
		dlpdictionaries:            dlpdictionaries.New(cli),
		dlp_engines:                dlp_engines.New(cli),
		dlp_notification_templates: dlp_notification_templates.New(cli),
		dlp_web_rules:              dlp_web_rules.New(cli),
		rule_labels:                rule_labels.New(cli),
		security_policy_settings:   security_policy_settings.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
