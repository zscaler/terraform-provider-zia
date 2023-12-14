package zia

import (
	"log"

	client "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/activation"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/adminuserrolemgmt/admins"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/adminuserrolemgmt/roles"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/cloudbrowserisolation"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/devicegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_engines"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_exact_data_match"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_icap_servers"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_idm_profiles"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_incident_receiver_servers"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_notification_templates"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlp_web_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/dlp/dlpdictionaries"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/applicationservices"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/appservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/filteringrules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipdestinationgroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/ipsourcegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplicationgroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkapplications"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservicegroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/networkservices"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/firewallpolicies/timewindow"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/forwarding_rules"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationgroups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationlite"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/location/locationmanagement"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/rule_labels"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_report"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_settings"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/sandbox/sandbox_submission"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/security_policy_settings"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/greinternalipranges"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnelinfo"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/gretunnels"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/staticips"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/virtualipaddress"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/trafficforwarding/vpncredentials"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlcategories"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/urlfilteringpolicies"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/user_authentication_settings"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/departments"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/groups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/users"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	admins                        *admins.Service
	roles                         *roles.Service
	filteringrules                *filteringrules.Service
	ipdestinationgroups           *ipdestinationgroups.Service
	ipsourcegroups                *ipsourcegroups.Service
	networkapplicationgroups      *networkapplicationgroups.Service
	networkapplications           *networkapplications.Service
	networkservices               *networkservices.Service
	networkservicegroups          *networkservicegroups.Service
	applicationservices           *applicationservices.Service
	appservicegroups              *appservicegroups.Service
	timewindow                    *timewindow.Service
	urlcategories                 *urlcategories.Service
	urlfilteringpolicies          *urlfilteringpolicies.Service
	users                         *users.Service
	groups                        *groups.Service
	departments                   *departments.Service
	gretunnels                    *gretunnels.Service
	gretunnelinfo                 *gretunnelinfo.Service
	greinternalipranges           *greinternalipranges.Service
	staticips                     *staticips.Service
	virtualipaddress              *virtualipaddress.Service
	vpncredentials                *vpncredentials.Service
	locationmanagement            *locationmanagement.Service
	locationgroups                *locationgroups.Service
	locationlite                  *locationlite.Service
	activation                    *activation.Service
	devicegroups                  *devicegroups.Service
	dlpdictionaries               *dlpdictionaries.Service
	dlp_engines                   *dlp_engines.Service
	dlp_idm_profiles              *dlp_idm_profiles.Service
	dlp_exact_data_match          *dlp_exact_data_match.Service
	dlp_icap_servers              *dlp_icap_servers.Service
	dlp_incident_receiver_servers *dlp_incident_receiver_servers.Service
	dlp_notification_templates    *dlp_notification_templates.Service
	dlp_web_rules                 *dlp_web_rules.Service
	rule_labels                   *rule_labels.Service
	security_policy_settings      *security_policy_settings.Service
	user_authentication_settings  *user_authentication_settings.Service
	forwarding_rules              *forwarding_rules.Service
	zpa_gateways                  *zpa_gateways.Service
	sandbox_settings              *sandbox_settings.Service
	sandbox_report                *sandbox_report.Service
	sandbox_submission            *sandbox_submission.Service
	cloudbrowserisolation         *cloudbrowserisolation.Service
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
		admins:                        admins.New(cli),
		roles:                         roles.New(cli),
		filteringrules:                filteringrules.New(cli),
		ipdestinationgroups:           ipdestinationgroups.New(cli),
		ipsourcegroups:                ipsourcegroups.New(cli),
		networkapplicationgroups:      networkapplicationgroups.New(cli),
		networkapplications:           networkapplications.New(cli),
		networkservices:               networkservices.New(cli),
		networkservicegroups:          networkservicegroups.New(cli),
		applicationservices:           applicationservices.New(cli),
		appservicegroups:              appservicegroups.New(cli),
		timewindow:                    timewindow.New(cli),
		urlcategories:                 urlcategories.New(cli),
		urlfilteringpolicies:          urlfilteringpolicies.New(cli),
		users:                         users.New(cli),
		groups:                        groups.New(cli),
		departments:                   departments.New(cli),
		virtualipaddress:              virtualipaddress.New(cli),
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
		dlp_exact_data_match:          dlp_exact_data_match.New(cli),
		dlp_icap_servers:              dlp_icap_servers.New(cli),
		dlp_incident_receiver_servers: dlp_incident_receiver_servers.New(cli),
		dlp_notification_templates:    dlp_notification_templates.New(cli),
		dlp_web_rules:                 dlp_web_rules.New(cli),
		rule_labels:                   rule_labels.New(cli),
		security_policy_settings:      security_policy_settings.New(cli),
		user_authentication_settings:  user_authentication_settings.New(cli),
		forwarding_rules:              forwarding_rules.New(cli),
		zpa_gateways:                  zpa_gateways.New(cli),
		sandbox_settings:              sandbox_settings.New(cli),
		sandbox_report:                sandbox_report.New(cli),
		sandbox_submission:            sandbox_submission.New(cli),
		cloudbrowserisolation:         cloudbrowserisolation.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
