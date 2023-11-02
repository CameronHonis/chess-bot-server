package main

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	"github.com/CameronHonis/chess-bot-server/engines"
	"os"
)

func HandleMessageFromArbitrator(msg *server.Message) error {
	switch msg.ContentType {
	case server.CONTENT_TYPE_AUTH:
		content, ok := msg.Content.(*server.AuthMessageContent)
		if !ok {
			return fmt.Errorf("could not cast message to AuthMessageContent")
		}
		return HandleAuthMessage(content)
	case server.CONTENT_TYPE_FIND_BOT_MATCH:
		content, ok := msg.Content.(*server.FindBotMatchMessageContent)
		if !ok {
			return fmt.Errorf("could not cast message to FindBotMatchMessageContent")
		}
		return HandleFindBotMatchMessage(msg.SenderKey, content)
	case server.CONTENT_TYPE_MATCH_UPDATE:
		content, ok := msg.Content.(*server.MatchUpdateMessageContent)
		if !ok {
			return fmt.Errorf("could not cast message to MatchUpdateMessageContent")
		}
		return HandleMatchUpdateMessage(content)
	case server.CONTENT_TYPE_UPGRADE_AUTH_DENIED:
		HandleUpgradeAuthDeniedMessage()
		return nil
	case server.CONTENT_TYPE_UPGRADE_AUTH_GRANTED:
		return HandleAuthUpgradeGrantedMessage()
	case server.CONTENT_TYPE_SUBSCRIBE_REQUEST_DENIED:
		content, ok := msg.Content.(*server.SubscribeRequestDeniedMessageContent)
		if !ok {
			return fmt.Errorf("could not cast message to SubscribeRequestDeniedMessageContent")
		}
		HandleSubscribeDeniedMessage(content.Topic)
		return nil
	case server.CONTENT_TYPE_SUBSCRIBE_REQUEST_GRANTED:
		return nil
	case server.CONTENT_TYPE_MOVE:
		return nil
	default:
		return fmt.Errorf("unhandled message with content type %s", msg.ContentType)
	}
}

func HandleAuthMessage(content *server.AuthMessageContent) error {
	GetArbitratorClient().SetPublicPrivateKey(content.PublicKey, content.PrivateKey)
	botSecret, ok := os.LookupEnv("BOT_CLIENT_SECRET")
	if !ok {
		panic("could not determine bot client secret")
	}
	return GetArbitratorClient().RequestAuthUpgrade(botSecret)
}

func HandleAuthUpgradeGrantedMessage() error {
	return GetArbitratorClient().RequestSubscribe("findBotMatch")
}

func HandleFindBotMatchMessage(senderKey string, content *server.FindBotMatchMessageContent) error {
	// TODO: handle remote engine lookups
	_, engineErr := engines.GetLocalEngine(content.BotName)
	if engineErr != nil {
		_ = GetArbitratorClient().FailInitBotRequest(senderKey, content.BotName, engineErr.Error())
		return engineErr
	}
	return GetArbitratorClient().SucceedInitBotRequest(senderKey, content.BotName)
}

func HandleUpgradeAuthDeniedMessage() {
	panic("arbitrator denied bot client auth")
}

func HandleSubscribeDeniedMessage(topic server.MessageTopic) {
	panic(fmt.Sprintf("arbitrator denied bot client subscription%s", topic))
}

func HandleMatchUpdateMessage(content *server.MatchUpdateMessageContent) error {

	return nil
}
