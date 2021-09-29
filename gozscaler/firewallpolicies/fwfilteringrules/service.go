package fwfilteringrules

import (
	"github.com/willguibr/terraform-provider-zia/gozscaler/client"
)

type Service struct {
	Client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{Client: c}
}
