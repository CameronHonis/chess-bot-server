package arbitrator_client

import "github.com/CameronHonis/service"

type ArbitratorClientConfig struct {
	service.ConfigI
	authSecret string
	url        string
}

func NewArbitratorClientConfig(authSecret, url string) *ArbitratorClientConfig {
	return &ArbitratorClientConfig{
		authSecret: authSecret,
		url:        url,
	}
}

func (c *ArbitratorClientConfig) AuthSecret() string {
	return c.authSecret
}

func (c *ArbitratorClientConfig) Url() string {
	return c.url
}
