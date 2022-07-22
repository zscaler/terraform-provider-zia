package security_policy_settings

import (
	client "github.com/zscaler/zscaler-sdk-go/zia"
)

type Service struct {
	Client *client.Client
}

func New(c *client.Client) *Service {
	return &Service{Client: c}
}
