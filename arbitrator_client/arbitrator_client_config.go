package arbitrator_client

import "github.com/CameronHonis/service"

type ArbitratorClientConfig struct {
	service.ConfigI
}

func NewArbitratorClientConfig() *ArbitratorClientConfig {
	return &ArbitratorClientConfig{}
}
