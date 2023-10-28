package main

import (
	"fmt"
	"github.com/CameronHonis/chess-arbitrator/server"
	"github.com/gorilla/websocket"
	"time"
)

var arbitratorClient *ArbitratorClient

type ArbitratorClient struct {
	conn *websocket.Conn
}

func GetArbitratorClient() *ArbitratorClient {
	if arbitratorClient == nil {
		arbitratorClient = &ArbitratorClient{}
		go arbitratorClient.listenOnWebsocket()
	}
	return arbitratorClient
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
	msgBytes, marshalErr := msg.Marshal()
	if marshalErr != nil {
		return marshalErr
	}
	fmt.Println("[CLIENT]", ">>", string(msgBytes))
	return ac.conn.WriteMessage(websocket.TextMessage, msgBytes)
}
