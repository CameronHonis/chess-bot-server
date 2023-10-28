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
		return HandleFindBotMatchMessage(content)
	case server.CONTENT_TYPE_MATCH_UPDATE:
		content, ok := msg.Content.(*server.MatchUpdateMessageContent)
		if !ok {
			return fmt.Errorf("could not cast message to MatchUpdateMessageContent")
		}
		return HandleMatchUpdateMessage(content)
	case server.CONTENT_TYPE_UPGRADE_AUTH_DENIED:
		HandleUpgradeAuthDeniedMessage()
		return nil
	case server.CONTENT_TYPE_MOVE:
		return nil
	case server.CONTENT_TYPE_UPGRADE_AUTH_GRANTED:
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
	msg := server.Message{
		Topic:       "",
		ContentType: server.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
		Content: &server.UpgradeAuthRequestMessageContent{
			Secret: botSecret,
		},
	}
	sendErr := GetArbitratorClient().SendMessage(&msg)
	if sendErr != nil {
		return sendErr
	}
	return nil
}

func HandleFindBotMatchMessage(content *server.FindBotMatchMessageContent) error {
	// TODO: handle remote engine lookups
	_, engineErr := engines.GetLocalEngine(content.BotName)
	if engineErr != nil {
		return engineErr
	}
	msg := server.Message{
		Topic:       "findMatch",
		ContentType: server.CONTENT_TYPE_FIND_BOT_MATCH,
		Content: &server.FindBotMatchMessageContent{
			BotName: content.BotName,
		},
	}
	msgErr := GetArbitratorClient().SendMessage(&msg)
	if msgErr != nil {
		return msgErr
	}
	return nil
}

func HandleUpgradeAuthDeniedMessage() {
	panic("arbitrator denied bot client auth")
}

func HandleMatchUpdateMessage(content *server.MatchUpdateMessageContent) error {

	return nil
}
