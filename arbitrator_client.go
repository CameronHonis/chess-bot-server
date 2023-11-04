package main

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/server"
	"github.com/gorilla/websocket"
	"time"
)

var arbitratorClient *ArbitratorClient

type ArbitratorClient struct {
	conn       *websocket.Conn
	publicKey  string
	privateKey string
}

func GetArbitratorClient() *ArbitratorClient {
	if arbitratorClient == nil {
		arbitratorClient = &ArbitratorClient{}
		go arbitratorClient.listenOnWebsocket()
	}
	return arbitratorClient
}

func (ac *ArbitratorClient) GetPublicKey() string {
	return ac.publicKey
}

func (ac *ArbitratorClient) connect() {
	for ac.conn == nil {
		var err error
		ac.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			fmt.Println("[CLIENT]", "could not connect to arbitrator, retrying in 1 second:", err)
			ac.conn = nil
			time.Sleep(time.Second)
		}
	}
}

func (ac *ArbitratorClient) listenOnWebsocket() {
	for {
		ac.connect()
		for {
			_, rawMsg, readErr := ac.conn.ReadMessage()
			if readErr != nil {
				fmt.Println("[CLIENT]", "error reading message from websocket:", readErr)
				// assume all readErrs are disconnects
				ac.conn = nil
				break
			}
			fmt.Println("[CLIENT]", "arbitrator >>", string(rawMsg))

			msg, unmarshalErr := server.UnmarshalToMessage(rawMsg)
			if unmarshalErr != nil {
				fmt.Println("[CLIENT]", "could not unmarshal message:", unmarshalErr)
				continue
			}
			handleMsgErr := HandleMessageFromArbitrator(msg)
			if handleMsgErr != nil {
				fmt.Println("[CLIENT]", "could not handle message:", handleMsgErr)
				continue
			}
		}
	}
}

func (ac *ArbitratorClient) SendMessage(msg *server.Message) error {
	if ac.conn == nil {
		return fmt.Errorf("cannot send message, connection is nil")
	}
	ac.SignMessage(msg)
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		return marshalErr
	}
	fmt.Println("[CLIENT]", ">>", string(msgBytes))
	return ac.conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func (ac *ArbitratorClient) SetPublicPrivateKey(publicKey string, privateKey string) {
	ac.publicKey = publicKey
	ac.privateKey = privateKey
}

func (ac *ArbitratorClient) SignMessage(msg *server.Message) {
	msg.SenderKey = ac.publicKey
	msg.PrivateKey = ac.privateKey
}

func (ac *ArbitratorClient) RequestAuthUpgrade(upgradeSecret string) error {
	msg := server.Message{
		Topic:       "",
		ContentType: server.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
		Content: &server.UpgradeAuthRequestMessageContent{
			Secret: upgradeSecret,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) RequestSubscribe(topic server.MessageTopic) error {
	msg := server.Message{
		Topic:       "subscribe",
		ContentType: server.CONTENT_TYPE_SUBSCRIBE_REQUEST,
		Content: &server.SubscribeRequestMessageContent{
			Topic: topic,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) FindBotMatch(botName string) error {
	msg := server.Message{
		Topic:       "findBotMatch",
		ContentType: server.CONTENT_TYPE_FIND_BOT_MATCH,
		Content: &server.FindBotMatchMessageContent{
			BotName: botName,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) FailInitBotRequest(matchId string, botName string, reason string) error {
	msg := server.Message{
		Topic:       "",
		ContentType: server.CONTENT_TYPE_INIT_BOT_MATCH_FAILURE,
		Content: &server.InitBotMatchFailureMessageContent{
			BotName: botName,
			MatchId: matchId,
			Reason:  reason,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) SucceedInitBotRequest(matchId string, botName string) error {
	msg := server.Message{
		Topic:       "",
		ContentType: server.CONTENT_TYPE_INIT_BOT_MATCH_SUCCESS,
		Content: &server.InitBotMatchSuccessMessageContent{
			MatchId: matchId,
			BotName: botName,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) SendMove(matchId string, move *chess.Move) error {
	msg := server.Message{
		Topic:       server.MessageTopic(fmt.Sprintf("match-%s", matchId)),
		ContentType: server.CONTENT_TYPE_MOVE,
		Content: &server.MoveMessageContent{
			MatchId: matchId,
			Move:    move,
		},
	}
	return ac.SendMessage(&msg)
}
