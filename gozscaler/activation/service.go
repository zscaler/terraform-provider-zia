package activation

import (
	"github.com/willguibr/terraform-provider-zia/gozscaler"
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
)

type Service struct {
	Client *client.Client
}

func New(c *gozscaler.Config) *Service {
	return &Service{Client: client.NewClient(c)}
}
