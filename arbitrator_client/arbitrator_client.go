package arbitrator_client

import (
	"fmt"
	"github.com/CameronHonis/chess"
	"github.com/CameronHonis/chess-arbitrator/models"
	"github.com/CameronHonis/chess-bot-server/app"
	"github.com/CameronHonis/chess-bot-server/bot_manager"
	"github.com/CameronHonis/log"
	. "github.com/CameronHonis/marker"
	"github.com/CameronHonis/service"
	"github.com/gorilla/websocket"
	"time"
)

type ArbitratorClientI interface {
	service.ServiceI
	PublicKey() models.Key
	SendMessage(msg *models.Message) error
	RequestAuthUpgrade(secret string) error
	RequestSubscribe(topic models.MessageTopic) error
	SendMove(matchId string, move *chess.Move) error
}

type ArbitratorClient struct {
	service.Service

	__dependencies__ Marker
	LogService       log.LoggerServiceI
	BotMngr          bot_manager.BotManagerI

	__state__ Marker
	conn      *websocket.Conn
	pubKey    models.Key
	priKey    models.Key
}

func (ac *ArbitratorClient) PublicKey() models.Key {
	return ac.pubKey
}

func (ac *ArbitratorClient) SendMessage(msg *models.Message) error {
	if ac.conn == nil {
		return fmt.Errorf("cannot send message, connection is nil")
	}
	ac.SignMessage(msg)
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		return marshalErr
	}
	ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("sending message %s", msg))
	return ac.conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func (ac *ArbitratorClient) RequestAuthUpgrade(secret string) error {
	msg := models.Message{
		Topic:       "",
		ContentType: models.CONTENT_TYPE_UPGRADE_AUTH_REQUEST,
		Content: &models.UpgradeAuthRequestMessageContent{
			Secret: secret,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) RequestSubscribe(topic models.MessageTopic) error {
	msg := models.Message{
		Topic:       "subscribe",
		ContentType: models.CONTENT_TYPE_SUBSCRIBE_REQUEST,
		Content: &models.SubscribeRequestMessageContent{
			Topic: topic,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) SendMove(matchId string, move *chess.Move) error {
	msg := models.Message{
		Topic:       models.MessageTopic(fmt.Sprintf("match-%s", matchId)),
		ContentType: models.CONTENT_TYPE_MOVE,
		Content: &models.MoveMessageContent{
			MatchId: matchId,
			Move:    move,
		},
	}
	return ac.SendMessage(&msg)
}

func (ac *ArbitratorClient) Connect() {
	for ac.conn == nil {
		var err error
		ac.conn, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err != nil {
			ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not Connect to arbitrator, retrying in 1 second: %s", err))
			ac.conn = nil
			time.Sleep(time.Second)
		}
	}
}

func (ac *ArbitratorClient) ListenOnWebsocket() {
	for {
		ac.Connect()
		for {
			_, rawMsg, readErr := ac.conn.ReadMessage()
			if readErr != nil {
				ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("error reading message from websocket: %s", readErr))
				// assume all readErrs are disconnects
				ac.conn = nil
				break
			}
			ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("received message from arbitrator: %s", string(rawMsg)))

			msg, unmarshalErr := models.UnmarshalToMessage(rawMsg)
			if unmarshalErr != nil {
				ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not unmarshal message: %s", unmarshalErr))
				continue
			}
			handleMsgErr := ac.HandleMsg(msg)
			if handleMsgErr != nil {
				ac.LogService.Log(app.ENV_ARBITRATOR_CLIENT, fmt.Sprintf("could not handle message: %s", handleMsgErr))
				continue
			}
		}
	}
}

func (ac *ArbitratorClient) SetPublicPrivateKey(publicKey models.Key, privateKey models.Key) {
	ac.pubKey = publicKey
	ac.priKey = privateKey
}

func (ac *ArbitratorClient) SignMessage(msg *models.Message) {
	msg.SenderKey = ac.pubKey
	msg.PrivateKey = ac.priKey
}
