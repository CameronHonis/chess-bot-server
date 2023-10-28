package main

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	"github.com/CameronHonis/chess-bot-server/engines"
)

func HandleMessageFromArbitrator(msg *server.Message) error {
	switch msg.ContentType {
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
	case server.CONTENT_TYPE_MOVE:
		return nil
	default:
		return fmt.Errorf("unhandled message with content type %s", msg.ContentType)
	}
}

func HandleFindBotMatchMessage(content *server.FindBotMatchMessageContent) error {
	// TODO: handle remote engine lookups
	_, engineErr := engines.GetLocalEngine(content.BotName)
	if engineErr != nil {
		return engineErr
	}
	msg := server.Message{
		Topic:       "findMatch",
		ContentType: server.CONTENT_TYPE_FIND_MATCH,
		Content: &server.FindBotMatchMessageContent{
			BotName:   content.BotName,
			PlayerKey: content.PlayerKey,
		},
	}
	msgErr := GetArbitratorClient().SendMessage(&msg)
	if msgErr != nil {
		return msgErr
	}
	return nil
}

func HandleMatchUpdateMessage(content *server.MatchUpdateMessageContent) error {

	return nil
}
