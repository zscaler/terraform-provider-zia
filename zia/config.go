package zia

import (
	"log"

	"github.com/willguibr/terraform-provider-zia/gozscaler/adminauditlogs"
	"github.com/willguibr/terraform-provider-zia/gozscaler/adminrolemgmt"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
	"github.com/willguibr/terraform-provider-zia/gozscaler/usermanagement"
)

func init() {
	// remove timestamp from Zscaler provider logger, use the timestamp from the default terraform logger
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

type Client struct {
	adminauditlogs adminauditlogs.Service
	adminrolemgmt  adminrolemgmt.Service
	usermanagement usermanagement.Service
}

type Config struct {
	Username   string
	Password   string
	APIKey     string
	ZIABaseURL string
}

func (c *Config) Client() (*Client, error) {
	config, err := client.NewClientZIA(c.Username, c.Password, c.APIKey, c.ZIABaseURL)
	if err != nil {
		return nil, err
	}

	client := &Client{
		adminauditlogs: *adminauditlogs.New(config),
		adminrolemgmt:  *adminrolemgmt.New(config),
		usermanagement: *usermanagement.New(config),
	}

	log.Println("[INFO] initialized ZIA client")
	return client, nil
}
