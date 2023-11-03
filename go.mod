module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator

require (
	github.com/CameronHonis/chess v0.0.0-20231103001456-f5db2e734a3b
	github.com/CameronHonis/chess-arbitrator v0.0.0-20231102142115-406856391bd3
	github.com/gorilla/websocket v1.5.0
)

require github.com/google/uuid v1.4.0 // indirect
