package bot_manager

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"sync"
)

type BotManagerII interface {
	service.ServiceI
	BotClient(key models.Key) (BotClient, error)
	BotClientByMatchId(matchId string) (BotClient, error)
	MatchIdByBotClientKey(key models.Key) (string, error)

	InitBotClient(botName string) (BotClient, error)
	AttachMatchToClient(matchId string) error
	RemoveClient(key models.Key) error
}

type BotManager struct {
	service.Service
	__dependencies__ Marker
	LogService       log.LoggerServiceI

	__state__             Marker
	botClientByMatchId    map[string]BotClient
	matchIdByBotClientKey map[models.Key]string
	botClientByKey        map[models.Key]BotClient
	mu                    sync.Mutex
}

func NewBotManager(config *BotManagerConfig) *BotManager {
	m := &BotManager{
		botClientByMatchId: make(map[string]BotClient),
		botClientByKey:     make(map[models.Key]BotClient),
		mu:                 sync.Mutex{},
	}
	m.Service = *service.NewService(m, config)
	return m
}

func (bm *BotManager) BotClient(key models.Key) (BotClient, error) {
	bm.mu.Lock()
	botClient, ok := bm.botClientByKey[key]
	bm.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no bot client found by key %s", key)
	}
	return botClient, nil
}

func (bm *BotManager) BotClientFromMatchId(matchId string) (BotClient, error) {
	bm.mu.Lock()
	botClient, ok := bm.botClientByMatchId[matchId]
	bm.mu.Unlock()
	if !ok {
		return nil, fmt.Errorf("no bot client found for match %s", matchId)
	}
	return botClient, nil
}

func (bm *BotManager) MatchIdFromBotClientKey(key models.Key) (string, error) {
	bm.mu.Lock()
	matchId, ok := bm.matchIdByBotClientKey[key]
	if !ok {
		return "", fmt.Errorf("no match id exists for key %s", matchId)
	}
	return matchId, nil
}

func (bm *BotManager) InitBotClient(botName string) (BotClient, error) {
	// TODO: add lookups for remote bot clients
	botClient, botClientErr := NewLocalBotClient(botName)
	if botClientErr != nil {
		return nil, botClientErr
	}
	//initErr := botClient.Initialize(match)
	//if initErr != nil {
	//	return initErr
	//}
	bm.mu.Lock()
	bm.botClientByKey[botClient.Key()] = botClient
	bm.mu.Unlock()
	return botClient, nil
}

func (bm *BotManager) AttachMatchToClient(matchId string, key models.Key) error {
	botClient, botClientErr := bm.BotClient(key)
	if botClientErr != nil {
		return botClientErr
	}
	existingMatchId, _ := bm.MatchIdFromBotClientKey(key)
	if existingMatchId == "" {
		return fmt.Errorf("client %s already belongs to match %s", key, existingMatchId)
	}
	bm.mu.Lock()
	bm.botClientByMatchId[matchId] = botClient
	bm.mu.Unlock()
	return nil
}

func (bm *BotManager) RemoveClient(key models.Key) error {
	client, clientErr := bm.BotClient(key)
	if clientErr != nil {
		return fmt.Errorf("bot client with key %s does not exist", key)
	}
	matchId, _ := bm.MatchIdFromBotClientKey(key)

	bm.mu.Lock()
	if matchId != "" {
		delete(bm.matchIdByBotClientKey, key)
		delete(bm.botClientByMatchId, matchId)
	}
	delete(bm.botClientByKey, key)
	bm.mu.Unlock()

	client.Engine().Terminate()
	return nil
}
