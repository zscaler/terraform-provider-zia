package zia

import (
	"log"

	"github.com/willguibr/terraform-provider-zia/gozscaler/activation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/devicegroups"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/filteringrules"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipdestinationgroups"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/ipsourcegroups"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkapplications"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/networkservices"
	"github.com/willguibr/terraform-provider-zia/gozscaler/firewallpolicies/timewindow"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement/locationgroups"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/greinternalipranges"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/gretunnelinfo"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlcategories"
	"github.com/willguibr/terraform-provider-zia/gozscaler/urlfilteringpolicies"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	adminuserrolemgmt    *adminuserrolemgmt.Service
	filteringrules       *filteringrules.Service
	ipdestinationgroups  *ipdestinationgroups.Service
	ipsourcegroups       *ipsourcegroups.Service
	networkapplications  *networkapplications.Service
	networkservices      *networkservices.Service
	dlpdictionaries      *dlpdictionaries.Service
	timewindow           *timewindow.Service
	urlcategories        *urlcategories.Service
	urlfilteringpolicies *urlfilteringpolicies.Service
	usermanagement       *usermanagement.Service
	gretunnels           *gretunnels.Service
	gretunnelinfo        *gretunnelinfo.Service
	greinternalipranges  *greinternalipranges.Service
	staticips            *staticips.Service
	virtualipaddresslist *virtualipaddresslist.Service
	vpncredentials       *vpncredentials.Service
	locationmanagement   *locationmanagement.Service
	locationgroups       *locationgroups.Service
	activation           *activation.Service
	devicegroups         *devicegroups.Service
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
		adminuserrolemgmt:    adminuserrolemgmt.New(cli),
		filteringrules:       filteringrules.New(cli),
		ipdestinationgroups:  ipdestinationgroups.New(cli),
		ipsourcegroups:       ipsourcegroups.New(cli),
		networkapplications:  networkapplications.New(cli),
		networkservices:      networkservices.New(cli),
		dlpdictionaries:      dlpdictionaries.New(cli),
		timewindow:           timewindow.New(cli),
		urlcategories:        urlcategories.New(cli),
		urlfilteringpolicies: urlfilteringpolicies.New(cli),
		usermanagement:       usermanagement.New(cli),
		virtualipaddresslist: virtualipaddresslist.New(cli),
		vpncredentials:       vpncredentials.New(cli),
		gretunnels:           gretunnels.New(cli),
		gretunnelinfo:        gretunnelinfo.New(cli),
		greinternalipranges:  greinternalipranges.New(cli),
		staticips:            staticips.New(cli),
		locationmanagement:   locationmanagement.New(cli),
		locationgroups:       locationgroups.New(cli),
		activation:           activation.New(cli),
		devicegroups:         devicegroups.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
