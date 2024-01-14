package arbitrator_client

import (
	"github.com/CameronHonis/service"
)

func (ac *ArbitratorClient) OnConnect(ev service.EventI) (willPropagate bool) {
	config := ac.Config().(*ArbitratorClientConfig)
	if err := RequestAuthUpgrade(ac.SendMessage, config.AuthSecret()); err != nil {
		ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, err.Error())
	}
	return false
}
