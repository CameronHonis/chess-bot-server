package app

import "github.com/CameronHonis/log"

const ENV_ARBITRATOR_CLIENT = "arbitrator_client"
const ENV_BOT_CLIENT = "bot_client"
const ENV_BOT_CLIENT_MANAGER = "bot_client_manager"

func LoggerConfig() *log.LoggerConfig {
	logConfigBuilder := log.NewConfigBuilder()
	logConfigBuilder.WithDecorator(ENV_ARBITRATOR_CLIENT, log.WrapGreen)
	logConfigBuilder.WithDecorator(ENV_BOT_CLIENT, log.WrapBlue)
	logConfigBuilder.WithDecorator(ENV_BOT_CLIENT_MANAGER, log.WrapCyan)
	//logConfigBuilder.WithMutedEnv("arbitrator_client")
	//logConfigBuilder.WithMutedEnv("bot_manager")

	return logConfigBuilder.Build()
}
