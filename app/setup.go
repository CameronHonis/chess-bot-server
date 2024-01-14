package app

import (
	arbc "github.com/CameronHonis/chess-bot-server/arbitrator_client"
	"github.com/CameronHonis/chess-bot-server/bot_manager"
	"github.com/CameronHonis/chess-bot-server/models"
	"github.com/CameronHonis/log"
)

func LoggerConfig() *log.LoggerConfig {
	logConfigBuilder := log.NewConfigBuilder()
	logConfigBuilder.WithDecorator(models.ENV_ARBITRATOR_CLIENT, log.WrapGreen)
	logConfigBuilder.WithDecorator(models.ENV_BOT_CLIENT, log.WrapBlue)
	logConfigBuilder.WithDecorator(models.ENV_BOT_CLIENT_MANAGER, log.WrapCyan)
	//logConfigBuilder.WithMutedEnv("arbitrator_client")
	//logConfigBuilder.WithMutedEnv("bot_manager")

	return logConfigBuilder.Build()
}

func Setup() *AppService {
	logService := log.NewLoggerService(LoggerConfig())

	botManager := bot_manager.NewBotManager(bot_manager.NewBotManagerConfig())
	botManager.AddDependency(logService)

	arbClient := arbc.NewArbitratorClient(arbc.NewArbitratorClientConfig())
	arbClient.AddDependency(botManager)
	arbClient.AddDependency(logService)

	appService := NewAppService(NewAppServiceConfig())
	appService.AddDependency(arbClient)

	return appService
}
