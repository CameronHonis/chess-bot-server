package app

import (
	"github.com/CameronHonis/chess-bot-server/arbitrator_client"
	"github.com/CameronHonis/chess-bot-server/bot_manager"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
)

type AppServiceI interface {
	service.ServiceI
}

type AppService struct {
	service.Service
	__dependencies__ Marker
	ArbitratorClient arbitrator_client.ArbitratorClientI
	BotClient        bot_manager.BotClientI

	__state__ Marker
}

func NewAppService(config *AppServiceConfig) *AppService {
	s := &AppService{}
	s.Service = *service.NewService(s, config)
	return s
}
