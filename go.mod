module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator

replace github.com/CameronHonis/chess => ../chess

replace github.com/CameronHonis/log => ../log

replace github.com/CameronHonis/set => ../set

require (
	github.com/CameronHonis/chess v0.0.0-20231109054928-c290d2362c29
	github.com/CameronHonis/chess-arbitrator v0.0.0-20231111032807-be83c86bd5f5
	github.com/gorilla/websocket v1.5.1
)

require (
	github.com/CameronHonis/log v0.0.0-20231111002532-5d3c065f77b8 // indirect
	github.com/CameronHonis/set v0.0.0-20231110043107-dace21619137 // indirect
	github.com/google/uuid v1.4.0 // indirect
	golang.org/x/net v0.18.0 // indirect
)
