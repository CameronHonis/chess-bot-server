package app

import (
	arbc "github.com/CameronHonis/chess-bot-server/arbitrator_client"
	botmgr "github.com/CameronHonis/chess-bot-server/bot_manager"
	"github.com/CameronHonis/log"
)

func LoggerConfig() *log.LoggerConfig {
	logConfigBuilder := log.NewConfigBuilder()
	logConfigBuilder.WithDecorator(arbc.ENV_ARBITRATOR_CLIENT, log.WrapGreen)
	logConfigBuilder.WithDecorator(botmgr.ENV_BOT_CLIENT, log.WrapBlue)
	logConfigBuilder.WithDecorator(botmgr.ENV_BOT_MANAGER, log.WrapCyan)
	//logConfigBuilder.WithMutedEnv("arbitrator_client")
	//logConfigBuilder.WithMutedEnv("bot_manager")

	return logConfigBuilder.Build()
}

func ArbitratorClientConfig() *arbc.ArbitratorClientConfig {
	return arbc.NewArbitratorClientConfig("secret", "127.0.0.1:8080")
}

func Setup() *AppService {
	logService := log.NewLoggerService(LoggerConfig())

	botManager := botmgr.NewBotManager(botmgr.NewBotManagerConfig())
	botManager.AddDependency(logService)

	arbClient := arbc.NewArbitratorClient(ArbitratorClientConfig())
	arbClient.AddDependency(botManager)
	arbClient.AddDependency(logService)

	appService := NewAppService(NewAppServiceConfig())
	appService.AddDependency(arbClient)

	return appService
}
