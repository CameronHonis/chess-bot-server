module github.com/CameronHonis/chess-bot-server

go 1.18

replace github.com/CameronHonis/chess-arbitrator => ../arbitrator

replace github.com/CameronHonis/chess => ../chess

replace github.com/CameronHonis/log => ../log

replace github.com/CameronHonis/set => ../set

replace github.com/CameronHonis/service => ../service

require (
	github.com/CameronHonis/chess v0.0.0-20240209135107-c9b4c60ee9bb
	github.com/CameronHonis/chess-arbitrator v0.0.0-20240209193524-394a425b1dab
	github.com/CameronHonis/log v0.0.0-20240217020729-dd2bc7ede4b6
	github.com/CameronHonis/marker v0.0.0-20231220043644-4b47686a2d7b
	github.com/CameronHonis/service v0.0.0-20240129185253-97bdfd0882f6
	github.com/gorilla/websocket v1.5.1
)

require (
	github.com/CameronHonis/set v0.0.0-20240120001402-957bb72dae93 // indirect
	github.com/google/uuid v1.6.0 // indirect
	golang.org/x/net v0.21.0 // indirect
)
