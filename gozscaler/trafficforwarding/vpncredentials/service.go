package vpncredentials

import (
	"github.com/zscaler/terraform-provider-zia/gozscaler/client"
)

type Service struct {
	Client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{Client: c}
}
