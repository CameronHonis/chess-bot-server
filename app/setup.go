package app

import (
	"fmt"
	arbc "github.com/CameronHonis/chess-bot-server/arbitrator_client"
	botmgr "github.com/CameronHonis/chess-bot-server/bot_manager"
	"github.com/CameronHonis/log"
	"os"
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
	domainVal, domainExists := os.LookupEnv("ARBITRATOR_DOMAIN")
	if !domainExists {
		domainVal = "127.0.0.1"
	}
	portVal, portExists := os.LookupEnv("ARBITRATOR_PORT")
	if !portExists {
		portVal = "8080"
	}
	return arbc.NewArbitratorClientConfig("secret", fmt.Sprint("%s:%s", domainVal, portVal))
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
