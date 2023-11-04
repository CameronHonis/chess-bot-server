module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator

replace github.com/CameronHonis/chess => ../chess

require (
	github.com/CameronHonis/chess v0.0.0-20231104040721-1fa63f099091
	github.com/CameronHonis/chess-arbitrator v0.0.0-20231104050243-d08a654a3855
	github.com/gorilla/websocket v1.5.0
)

require github.com/google/uuid v1.4.0 // indirect
