package bot_manager

import (
	"github.com/CameronHonis/chess-arbitrator/auth"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/engines"
)

type BotClient interface {
	Key() models.Key
	Engine() engines.Engine
}

type LocalBotClient struct {
	key    models.Key
	engine engines.Engine
}

func NewLocalBotClient(engineName string) (BotClient, error) {
	engine, err := engines.GetLocalEngine(engineName)
	if err != nil {
		return nil, err
	}
	pubKey, _ := auth.GenerateKeyset()
	botClient := &LocalBotClient{
		key:    pubKey,
		engine: engine,
	}

	return botClient, nil
}

func (c *LocalBotClient) Key() models.Key {
	return c.key
}

func (c *LocalBotClient) Engine() engines.Engine {
	return c.engine
}
