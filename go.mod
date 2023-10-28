module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator
require (
	github.com/CameronHonis/chess v0.0.0-20231026001700-1bce8cf416e4
	github.com/CameronHonis/chess-arbitrator v0.0.0-20231028040251-e60d7a4a8fb6
	github.com/gorilla/websocket v1.5.0
)

require github.com/google/uuid v1.4.0 // indirect
