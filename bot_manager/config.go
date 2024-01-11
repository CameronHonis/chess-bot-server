package bot_manager

import "github.com/CameronHonis/service"

type BotManagerConfig struct {
	service.ConfigI
}

func NewBotManagerConfig() *BotManagerConfig {
	return &BotManagerConfig{}
}
