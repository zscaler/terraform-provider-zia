package staticips

import (
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
)

type Service struct {
	Client *client.Client
}

func New(c *client.Config) *Service {
	return &Service{Client: client.NewClientZIA(c)}
}