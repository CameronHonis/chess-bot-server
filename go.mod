module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator

replace github.com/CameronHonis/chess => ../chess

replace github.com/CameronHonis/log => ../log

replace github.com/CameronHonis/set => ../set

require (
	github.com/CameronHonis/chess v0.0.0-20231109054928-c290d2362c29
	github.com/CameronHonis/chess-arbitrator v0.0.0-20240110015316-4a97d4b1cfa3
	github.com/CameronHonis/log v0.0.0-20240103213023-c305c5e09be3
	github.com/CameronHonis/marker v0.0.0-20231220043644-4b47686a2d7b
	github.com/CameronHonis/service v0.0.0-20240103223336-b8385576c790
	github.com/gorilla/websocket v1.5.1
)

require (
	github.com/CameronHonis/set v0.0.0-20231212050345-6dcde9af4710 // indirect
	github.com/google/uuid v1.5.0 // indirect
	golang.org/x/net v0.20.0 // indirect
)
