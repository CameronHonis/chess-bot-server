package bot_client

import . "github.com/CameronHonis/log"

const ENV_ARBITRATOR_CLIENT = "arbitrator_client"
const ENV_BOT_CLIENT = "bot_client"
const ENV_BOT_CLIENT_MANAGER = "bot_client_manager"

func ConfigLogger() {
	logConfigBuilder := NewLogManagerConfigBuilder()
	logConfigBuilder.WithDecorator(ENV_ARBITRATOR_CLIENT, WrapGreen)
	logConfigBuilder.WithDecorator(ENV_BOT_CLIENT, WrapBlue)
	logConfigBuilder.WithDecorator(ENV_BOT_CLIENT_MANAGER, WrapCyan)
	//logConfigBuilder.WithMutedEnv("arbitrator_client")
	//logConfigBuilder.WithMutedEnv("bot_client")

	logConfig := logConfigBuilder.Build()
	GetLogManager().InjectConfig(logConfig)
}
