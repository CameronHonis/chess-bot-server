package arbitrator_client

import (
	"github.com/CameronHonis/service"
)

var OnConnSuccess = func(self service.ServiceI, event service.EventI) bool {
	ac := self.(*ArbitratorClient)
	for i := 0; i < 3; i++ {
		sendErr := RefreshAuthCreds(ac.SendMessage)
		if sendErr == nil {
			break
		}
		if i == 2 {
			ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, "could not refresh auth creds: ", sendErr)
		} else {
			ac.LogService.LogRed(ENV_ARBITRATOR_CLIENT, "retrying refresh auth creds...")
		}
	}
	return true
}
