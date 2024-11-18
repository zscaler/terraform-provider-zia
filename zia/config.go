package zia

import (
	"log"

	gozscaler "github.com/zscaler/zscaler-sdk-go/v2/zia"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/forwarding_control_policy/zpa_gateways"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/departments"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/groups"
	"github.com/zscaler/zscaler-sdk-go/v2/zia/services/usermanagement/users"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	admins                        *services.Service
	roles                         *services.Service
	filteringrules                *services.Service
	ipdestinationgroups           *services.Service
	ipsourcegroups                *services.Service
	networkapplicationgroups      *services.Service
	networkapplications           *services.Service
	networkservices               *services.Service
	networkservicegroups          *services.Service
	applicationservices           *services.Service
	appservicegroups              *services.Service
	timewindow                    *services.Service
	urlcategories                 *services.Service
	urlfilteringpolicies          *services.Service
	users                         *users.Service
	groups                        *groups.Service
	departments                   *departments.Service
	gretunnels                    *services.Service
	gretunnelinfo                 *services.Service
	greinternalipranges           *services.Service
	staticips                     *services.Service
	virtualipaddress              *services.Service
	vpncredentials                *services.Service
	locationmanagement            *services.Service
	locationgroups                *services.Service
	locationlite                  *services.Service
	activation                    *services.Service
	cloudappcontrol               *services.Service
	devicegroups                  *services.Service
	dlpdictionaries               *services.Service
	dlp_engines                   *services.Service
	dlp_idm_profiles              *services.Service
	dlp_idm_profile_lite          *services.Service
	dlp_exact_data_match          *services.Service
	dlp_icap_servers              *services.Service
	dlp_incident_receiver_servers *services.Service
	dlp_notification_templates    *services.Service
	dlp_web_rules                 *services.Service
	pacfiles                      *services.Service
	rule_labels                   *services.Service
	security_policy_settings      *services.Service
	user_authentication_settings  *services.Service
	forwarding_rules              *services.Service
	zpa_gateways                  *zpa_gateways.Service
	sandbox_settings              *services.Service
	sandbox_report                *services.Service
	sandbox_submission            *services.Service
	cloudbrowserisolation         *services.Service
	workloadgroups                *services.Service
}

type Config struct {
	Username   string
	Password   string
	APIKey     string
	ZIABaseURL string
	UserAgent  string
}

func (c *Config) Client() (*Client, error) {
	cli, err := gozscaler.NewClient(c.Username, c.Password, c.APIKey, c.ZIABaseURL, c.UserAgent)
	if err != nil {
		return nil, err
	}

	ziaClient := &Client{
		activation:                    services.New(cli),
		admins:                        services.New(cli),
		cloudappcontrol:               services.New(cli),
		roles:                         services.New(cli),
		filteringrules:                services.New(cli),
		ipdestinationgroups:           services.New(cli),
		ipsourcegroups:                services.New(cli),
		networkapplicationgroups:      services.New(cli),
		networkapplications:           services.New(cli),
		networkservices:               services.New(cli),
		networkservicegroups:          services.New(cli),
		applicationservices:           services.New(cli),
		appservicegroups:              services.New(cli),
		timewindow:                    services.New(cli),
		urlcategories:                 services.New(cli),
		urlfilteringpolicies:          services.New(cli),
		users:                         users.New(cli),
		groups:                        groups.New(cli),
		departments:                   departments.New(cli),
		pacfiles:                      services.New(cli),
		virtualipaddress:              services.New(cli),
		vpncredentials:                services.New(cli),
		gretunnels:                    services.New(cli),
		gretunnelinfo:                 services.New(cli),
		greinternalipranges:           services.New(cli),
		staticips:                     services.New(cli),
		locationmanagement:            services.New(cli),
		locationgroups:                services.New(cli),
		locationlite:                  services.New(cli),
		devicegroups:                  services.New(cli),
		dlpdictionaries:               services.New(cli),
		dlp_engines:                   services.New(cli),
		dlp_idm_profile_lite:          services.New(cli),
		dlp_idm_profiles:              services.New(cli),
		dlp_exact_data_match:          services.New(cli),
		dlp_icap_servers:              services.New(cli),
		dlp_incident_receiver_servers: services.New(cli),
		dlp_notification_templates:    services.New(cli),
		dlp_web_rules:                 services.New(cli),
		rule_labels:                   services.New(cli),
		security_policy_settings:      services.New(cli),
		user_authentication_settings:  services.New(cli),
		forwarding_rules:              services.New(cli),
		zpa_gateways:                  zpa_gateways.New(cli),
		sandbox_settings:              services.New(cli),
		sandbox_report:                services.New(cli),
		sandbox_submission:            services.New(cli),
		cloudbrowserisolation:         services.New(cli),
		workloadgroups:                services.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
