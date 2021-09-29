package zia

import (
	"log"

	"github.com/willguibr/terraform-provider-zia/gozscaler/activation"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminauditlogs"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminuserrolemgmt"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/dlpdictionaries"
	"github.com/willguibr/terraform-provider-zia/gozscaler/fwfilteringrules"
	"github.com/willguibr/terraform-provider-zia/gozscaler/locationmanagement"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/gretunnels"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/staticips"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/virtualipaddresslist"
	"github.com/willguibr/terraform-provider-zia/gozscaler/trafficforwarding/vpncredentials"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	adminauditlogs       *adminauditlogs.Service
	adminuserrolemgmt    *adminuserrolemgmt.Service
	fwfilteringrules     *fwfilteringrules.Service
	dlpdictionaries      *dlpdictionaries.Service
	usermanagement       *usermanagement.Service
	gretunnels           *gretunnels.Service
	staticips            *staticips.Service
	virtualipaddresslist *virtualipaddresslist.Service
	vpncredentials       *vpncredentials.Service
	locationmanagement   *locationmanagement.Service
	activation           *activation.Service
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
		adminauditlogs:       adminauditlogs.New(cli),
		adminuserrolemgmt:    adminuserrolemgmt.New(cli),
		fwfilteringrules:     fwfilteringrules.New(cli),
		dlpdictionaries:      dlpdictionaries.New(cli),
		usermanagement:       usermanagement.New(cli),
		virtualipaddresslist: virtualipaddresslist.New(cli),
		vpncredentials:       vpncredentials.New(cli),
		gretunnels:           gretunnels.New(cli),
		staticips:            staticips.New(cli),
		locationmanagement:   locationmanagement.New(cli),
		activation:           activation.New(cli),
	}

	log.Println("[INFO] initialized ZIA client")
	return ziaClient, nil
}
