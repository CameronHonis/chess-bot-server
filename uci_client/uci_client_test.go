package uci_client_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/CameronHonis/chess-bot-server/uci_client"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
	"strings"
	"time"
)

type MockWriter struct {
	out   io.Writer
	delay time.Duration
}

func NewMockWriter(out io.Writer, delay time.Duration) *MockWriter {
	return &MockWriter{
		out:   out,
		delay: delay,
	}
}

func (m *MockWriter) Write(p []byte) (int, error) {
	contents := string(p)
	var resp func(w io.Writer)
	switch contents {
	case "uci\n":
		resp = func(w io.Writer) {
			_, _ = w.Write([]byte("id name Stockfish dev-20240314-fb07281f\n" +
				"id author the Stockfish developers (see AUTHORS file)\n\n" +
				//"option name Debug Log File type string default \n" +
				"option name Threads type spin default 1 min 1 max 1024\n" +
				//"option name Hash type spin default 16 min 1 max 33554432\n" +
				//"option name Clear Hash type button\n" +
				"option name Ponder type check default false\n" +
				//"option name MultiPV type spin default 1 min 1 max 256\n" +
				//"option name Skill Level type spin default 20 min 0 max 20\n" +
				//"option name Move Overhead type spin default 10 min 0 max 5000\n" +
				//"option name nodestime type spin default 0 min 0 max 10000\n" +
				//"option name UCI_Chess960 type check default false\n" +
				//"option name UCI_LimitStrength type check default false\n" +
				//"option name UCI_Elo type spin default 1320 min 1320 max 3190\n" +
				//"option name UCI_ShowWDL type check default false\n" +
				//"option name SyzygyPath type string default <empty>\n" +
				//"option name SyzygyProbeDepth type spin default 1 min 1 max 100\n" +
				//"option name Syzygy50MoveRule type check default true\n" +
				//"option name SyzygyProbeLimit type spin default 7 min 0 max 7\n" +
				//"option name EvalFile type string default nn-1ceb1ade0001.nnue\n" +
				//"option name EvalFileSmall type string default nn-baff1ede1f90.nnue\n" +
				"uciok\n"))
		}
	case "setoption name Threads value 2\n":
	case "setoption name Threads value asdf\n":
		resp = func(w io.Writer) {
			_, _ = w.Write([]byte("terminate called after throwing an instance of 'std::invalid_argument'\n" +
				"  what():  stof\n" +
				"Aborted (core dumped)\n"))
		}
	case "setoption name NotAnOption value some-value\n":
		resp = func(w io.Writer) {
			_, _ = w.Write([]byte("No such option: NotAnOption"))
		}
	case "position fen rnbqkbnr/pppppppp/8/8/4P3/8/PPPP1PPP/RNBQKBNR b KQkq - 0 1\n":
	case "isready\n":
		resp = func(w io.Writer) {
			_, _ = w.Write([]byte("readyok\n"))
		}
	case "go wtime 100000\n":
		resp = func(w io.Writer) {
			toWrite := "info string NNUE evaluation using nn-1ceb1ade0001.nnue\n" +
				"info string NNUE evaluation using nn-baff1ede1f90.nnue\n" +
				"info depth 1 seldepth 2 multipv 1 score cp -1 nodes 20 nps 20000 hashfull 0 tbhits 0 time 1 pv d2d4\n" +
				"info depth 2 seldepth 2 multipv 1 score cp 7 nodes 44 nps 44000 hashfull 0 tbhits 0 time 1 pv a2a3\n" +
				"info depth 3 seldepth 2 multipv 1 score cp 30 nodes 69 nps 69000 hashfull 0 tbhits 0 time 1 pv d2d4\n" +
				"info depth 4 seldepth 2 multipv 1 score cp 30 nodes 90 nps 90000 hashfull 0 tbhits 0 time 1 pv d2d4\n" +
				"info depth 5 seldepth 3 multipv 1 score cp 30 nodes 121 nps 121000 hashfull 0 tbhits 0 time 1 pv d2d4 a7a6\n" +
				"info depth 6 seldepth 4 multipv 1 score cp 32 nodes 199 nps 199000 hashfull 0 tbhits 0 time 1 pv d2d4 a7a6 b1d2\n" +
				"info depth 7 seldepth 6 multipv 1 score cp 28 nodes 1430 nps 476666 hashfull 1 tbhits 0 time 3 pv e2e4 d7d5 e4d5 d8d5 g1f3\n" +
				"info depth 8 seldepth 10 multipv 1 score cp 26 nodes 3179 nps 635800 hashfull 2 tbhits 0 time 5 pv d2d4 g8f6 c2c4 d7d6 g1f3 c8f5\n" +
				"info depth 9 seldepth 11 multipv 1 score cp 26 nodes 4275 nps 712500 hashfull 2 tbhits 0 time 6 pv d2d4 g8f6 c2c4 e7e6 g1f3 c7c5\n" +
				"info depth 10 seldepth 13 multipv 1 score cp 27 nodes 6790 nps 848750 hashfull 3 tbhits 0 time 8 pv d2d4 d7d5 e2e3 g8f6 c2c4 e7e6 b1c3 c7c5 c4d5 e6d5\n" +
				"info depth 11 seldepth 13 multipv 1 score cp 24 nodes 12626 nps 901857 hashfull 5 tbhits 0 time 14 pv e2e4 c7c5 g1f3 d7d6 b1c3\n" +
				"info depth 12 seldepth 15 multipv 1 score cp 25 nodes 15131 nps 945687 hashfull 5 tbhits 0 time 16 pv e2e4 c7c5 g1f3 b8c6 b1c3 e7e6 d2d4 c5d4 f3d4 g8f6 d4c6\n" +
				"info depth 13 seldepth 16 multipv 1 score cp 25 nodes 27508 nps 982428 hashfull 12 tbhits 0 time 28 pv e2e4 c7c5 g1f3 e7e6 b1c3 d7d6 d2d4 c5d4 f1b5 c8d7 d1d4 d7b5\n" +
				"info depth 14 seldepth 19 multipv 1 score cp 20 nodes 84518 nps 1112078 hashfull 33 tbhits 0 time 76 pv e2e4 e7e6 d2d4 d7d5 b1c3 f8b4 e4e5 c7c5 a2a3\n" +
				"info depth 15 seldepth 17 multipv 1 score cp 29 nodes 105387 nps 1170966 hashfull 41 tbhits 0 time 90 pv e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6 c1e3 e7e5 d4b3 f8e7 d1d2\n" +
				"info depth 16 seldepth 19 multipv 1 score cp 29 nodes 137798 nps 1230339 hashfull 54 tbhits 0 time 112 pv e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6 c1e3 e7e5 d4b3 f8e7 d1d2 e8g8\n" +
				"info depth 17 seldepth 26 multipv 1 score cp 31 nodes 214516 nps 1316049 hashfull 83 tbhits 0 time 163 pv e2e4 c7c5 g1f3 e7e6 c2c3 g8f6 e4e5 f6d5 a2a3 b8c6 d2d4 c5d4 c3d4 d7d6\n" +
				"info depth 18 seldepth 26 multipv 1 score cp 29 nodes 312779 nps 1365847 hashfull 120 tbhits 0 time 229 pv e2e4 c7c5 g1f3 b8c6 f1b5 d7d6 e1g1 c8d7 c2c3 a7a6 b5a4 g8f6 f1e1\n" +
				"info depth 19 seldepth 27 multipv 1 score cp 32 nodes 358501 nps 1384173 hashfull 135 tbhits 0 time 259 pv e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6 c1g5 e7e6 d1d2 h7h6 g5e3 f6g4 f2f4 g4e3 d2e3 f8e7 e1c1 b8c6\n" +
				"info depth 20 seldepth 34 multipv 1 score cp 27 nodes 766733 nps 1412031 hashfull 303 tbhits 0 time 543 pv e2e4 e7e5 g1f3 b8c6 d2d4 e5d4 f3d4 g8f6 d4c6 b7c6 f1d3 d7d5 b1d2 c8g4 f2f3 g4e6 e1g1 f8c5 g1h1\n" +
				"info depth 21 seldepth 31 multipv 1 score cp 25 nodes 937093 nps 1419837 hashfull 370 tbhits 0 time 660 pv e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6 c1g5 e7e6 f2f4 h7h6 g5h4 f8e7 d1f3 b7b5 e1c1 c8b7 f4f5 e6e5 h4f6 e7f6 d4b3 e8g8 c1b1\n" +
				"info depth 22 seldepth 33 multipv 1 score cp 29 nodes 1054937 nps 1425590 hashfull 403 tbhits 0 time 740 pv e2e4 c7c5 g1f3 d7d6 d2d4 c5d4 f3d4 g8f6 b1c3 a7a6 f1e2 e7e5 d4b3 f8e7 c1e3 c8e6 f2f4 e5f4 e3f4 e8g8 d1d2 b7b5\n" +
				"info depth 23 seldepth 37 multipv 1 score cp 25 nodes 2384807 nps 1394623 hashfull 762 tbhits 0 time 1710 pv e2e4 e7e5 g1f3 b8c6 f1c4 f8c5 d2d3 g8f6 c2c3 c5b6 e1g1 d7d6 f1e1 h7h6 a2a4 e8g8 h2h3 a7a5 c4b3 c6e7 d3d4 e7g6 b1a3 c7c6 d4e5 d6e5 d1d8 f8d8\n" +
				"info depth 24 seldepth 29 multipv 1 score cp 25 nodes 2572467 nps 1393535 hashfull 797 tbhits 0 time 1846 pv e2e4 e7e5 g1f3 b8c6 f1c4 f8c5 d2d3 g8f6 c2c3 c5b6 e1g1 d7d6 f1e1 h7h6 a2a4 e8g8 h2h3 a7a5 c4b3 c6e7 d3d4 e7g6 b1a3 c7c6 d4e5 d6e5 d1d8 f8d8\n" +
				"info depth 25 seldepth 33 multipv 1 score cp 22 nodes 3693100 nps 1399431 hashfull 914 tbhits 0 time 2639 pv d2d4 g8f6 c2c4 e7e6 g1f3 d7d5 b1c3 f8b4 c1g5 e8g8 e2e3 c7c5 c4d5 e6d5 f1d3 c5c4 d3c2 b8d7 e1g1 b4c3 b2c3 f8e8\n" +
				"info depth 26 seldepth 31 multipv 1 score cp 22 nodes 3854356 nps 1398532 hashfull 921 tbhits 0 time 2756 pv d2d4 g8f6 c2c4 e7e6 g1f3 d7d5 b1c3 f8b4 c1g5 e8g8 e2e3 c7c5 c4d5 e6d5 f1d3 c5c4 d3c2 b8d7 e1g1 b4c3 b2c3 f8e8 a1b1 d8a5\n" +
				"info depth 27 seldepth 32 multipv 1 score cp 29 nodes 4157709 nps 1406056 hashfull 935 tbhits 0 time 2957 pv d2d4 g8f6 c2c4 e7e6 g1f3 d7d5 b1c3 c7c5 c4d5 c5d4 d1d4 e6d5 c1g5 f8e7 e2e3 a7a6 d4h4 b8c6 f1d3 c8e6 e1g1\n" +
				"info depth 28 seldepth 31 multipv 1 score cp 31 lowerbound nodes 5092814 nps 1402592 hashfull 975 tbhits 0 time 3631 pv d2d4\n" +
				"info depth 25 currmove d2d4 currmovenumber 1\n" +
				"info depth 25 currmove e2e4 currmovenumber 2\n" +
				"info depth 25 currmove g1f3 currmovenumber 3\n" +
				"info depth 25 currmove b1c3 currmovenumber 4\n" +
				"info depth 25 currmove e2e3 currmovenumber 5\n" +
				"info depth 25 currmove c2c4 currmovenumber 6\n" +
				"info depth 25 currmove c2c3 currmovenumber 7\n" +
				"info depth 25 currmove a2a3 currmovenumber 8\n" +
				"info depth 25 currmove b2b4 currmovenumber 9\n" +
				"info depth 25 currmove h2h3 currmovenumber 10\n" +
				"info depth 25 currmove a2a4 currmovenumber 11\n" +
				"info depth 25 currmove g2g4 currmovenumber 12\n" +
				"info depth 25 currmove f2f4 currmovenumber 13\n" +
				"info depth 25 currmove g2g3 currmovenumber 14\n" +
				"info depth 25 currmove f2f3 currmovenumber 15\n" +
				"info depth 25 currmove h2h4 currmovenumber 16\n" +
				"info depth 25 currmove d2d3 currmovenumber 17\n" +
				"info depth 25 currmove b2b3 currmovenumber 18\n" +
				"info depth 25 currmove b1a3 currmovenumber 19\n" +
				"info depth 25 currmove g1h3 currmovenumber 20\n" +
				"info depth 28 seldepth 34 multipv 1 score cp 31 nodes 5481585 nps 1404094 hashfull 981 tbhits 0 time 3904 pv d2d4 g8f6 c2c4 e7e6 g2g3 d7d5 f1g2 c7c5 c4d5 e6d5 g1f3 b8c6 e1g1 h7h6 c1e3 f6g4 b1c3 c8e6 d1b3 a8b8\n" +
				"info depth 26 currmove d2d4 currmovenumber 1\n" +
				"info depth 26 currmove e2e4 currmovenumber 2\n" +
				"info depth 26 currmove g1f3 currmovenumber 3\ninfo depth 26 currmove b1c3 currmovenumber 4\n" +
				"info depth 26 currmove e2e3 currmovenumber 5\n" +
				"info depth 26 currmove c2c3 currmovenumber 6\n" +
				"info depth 26 currmove c2c4 currmovenumber 7\n" +
				"info depth 26 currmove a2a3 currmovenumber 8\n" +
				"info depth 26 currmove b2b4 currmovenumber 9\n" +
				"info depth 26 currmove g2g4 currmovenumber 10\n" +
				"info depth 26 currmove b2b3 currmovenumber 11\n" +
				"info depth 26 currmove a2a4 currmovenumber 12\n" +
				"info depth 26 currmove h2h3 currmovenumber 13\n" +
				"info depth 26 currmove f2f4 currmovenumber 14\n" +
				"info depth 26 currmove f2f3 currmovenumber 15\n" +
				"info depth 26 currmove d2d3 currmovenumber 16\n" +
				"info depth 26 currmove g2g3 currmovenumber 17\n" +
				"info depth 26 currmove b1a3 currmovenumber 18\n" +
				"info depth 26 currmove h2h4 currmovenumber 19\n" +
				"info depth 26 currmove g1h3 currmovenumber 20\n" +
				"info depth 29 seldepth 31 multipv 1 score cp 28 upperbound nodes 5811586 nps 1403086 hashfull 986 tbhits 0 time 4142 pv d2d4 g8f6\n" +
				"info depth 26 currmove d2d4 currmovenumber 1\n" +
				"info depth 26 currmove e2e4 currmovenumber 2\n" +
				"info depth 26 currmove g1f3 currmovenumber 3\n" +
				"info depth 26 currmove b1c3 currmovenumber 4\n" +
				"info depth 26 currmove e2e3 currmovenumber 5\n" +
				"info depth 26 currmove c2c4 currmovenumber 6\n" +
				"info depth 26 currmove c2c3 currmovenumber 7\n" +
				"info depth 26 currmove a2a3 currmovenumber 8\n" +
				"info depth 26 currmove b2b3 currmovenumber 9\n" +
				"info depth 26 currmove h2h3 currmovenumber 10\n" +
				"info depth 26 currmove d2d3 currmovenumber 11\n" +
				"info depth 26 currmove b2b4 currmovenumber 12\n" +
				"info depth 26 currmove a2a4 currmovenumber 13\n" +
				"info depth 26 currmove g2g3 currmovenumber 14\n" +
				"info depth 26 currmove f2f3 currmovenumber 15\n" +
				"info depth 26 currmove g2g4 currmovenumber 16\n" +
				"info depth 26 currmove f2f4 currmovenumber 17\n" +
				"info depth 26 currmove h2h4 currmovenumber 18\n" +
				"info depth 26 currmove b1a3 currmovenumber 19\n" +
				"info depth 26 currmove g1h3 currmovenumber 20\n" +
				"info depth 29 seldepth 35 multipv 1 score cp 30 nodes 5946288 nps 1402095 hashfull 988 tbhits 0 time 4241 pv d2d4 g8f6 c2c4 e7e6 g2g3 d7d5 f1g2 d5c4 d1a4 c7c6 a4c4 b7b5 c4c2 c8b7 g1f3 b8d7 f3e5 d8c8 e5d7 f6d7 b1c3 f8e7 e1g1\n" +
				"info depth 27 currmove d2d4 currmovenumber 1\n" +
				"info depth 27 currmove e2e4 currmovenumber 2\n" +
				"info depth 27 currmove g1f3 currmovenumber 3\n" +
				"info depth 27 currmove e2e3 currmovenumber 4\n" +
				"info depth 27 currmove b1c3 currmovenumber 5\n" +
				"info depth 27 currmove c2c4 currmovenumber 6\n" +
				"info depth 27 currmove c2c3 currmovenumber 7\n" +
				"info depth 27 currmove a2a3 currmovenumber 8\n" +
				"info depth 27 currmove a2a4 currmovenumber 9\n" +
				"info depth 27 currmove f2f3 currmovenumber 10\n" +
				"info depth 27 currmove h2h3 currmovenumber 11\n" +
				"info depth 27 currmove f2f4 currmovenumber 12\n" +
				"info depth 27 currmove b2b3 currmovenumber 13\n" +
				"info depth 27 currmove g2g3 currmovenumber 14\n" +
				"info depth 27 currmove g2g4 currmovenumber 15\n" +
				"info depth 27 currmove d2d3 currmovenumber 16\n" +
				"info depth 27 currmove b2b4 currmovenumber 17\n" +
				"info depth 27 currmove h2h4 currmovenumber 18\n" +
				"info depth 27 currmove b1a3 currmovenumber 19\n" +
				"info depth 27 currmove g1h3 currmovenumber 20\n" +
				"info depth 30 seldepth 35 multipv 1 score cp 27 upperbound nodes 6210486 nps 1402865 hashfull 990 tbhits 0 time 4427 pv d2d4 g8f6\n" +
				"info depth 27 currmove d2d4 currmovenumber 1\n" +
				"info depth 27 currmove e2e4 currmovenumber 2\n" +
				"info depth 27 currmove g1f3 currmovenumber 3\n" +
				"info depth 27 currmove b1c3 currmovenumber 4\n" +
				"info depth 27 currmove c2c4 currmovenumber 5\n" +
				"info depth 27 currmove c2c3 currmovenumber 6\n" +
				"info depth 27 currmove e2e3 currmovenumber 7\n" +
				"info depth 27 currmove b2b4 currmovenumber 8\n" +
				"info depth 27 currmove a2a4 currmovenumber 9\n" +
				"info depth 27 currmove a2a3 currmovenumber 10\n" +
				"info depth 27 currmove d2d3 currmovenumber 11\n" +
				"info depth 27 currmove h2h4 currmovenumber 12\n" +
				"info depth 27 currmove h2h3 currmovenumber 13\n" +
				"info depth 27 currmove b2b3 currmovenumber 14\n" +
				"info depth 27 currmove g2g3 currmovenumber 15\n" +
				"info depth 27 currmove f2f3 currmovenumber 16\n" +
				"info depth 27 currmove f2f4 currmovenumber 17\n" +
				"info depth 27 currmove g2g4 currmovenumber 18\n" +
				"info depth 27 currmove b1a3 currmovenumber 19\n" +
				"info depth 27 currmove g1h3 currmovenumber 20\n" +
				"info depth 30 seldepth 36 multipv 1 score cp 25 nodes 6475849 nps 1403521 hashfull 991 tbhits 0 time 4614 pv d2d4 g8f6 c2c4 e7e6 g2g3 d7d5 f1g2 d5c4 d1a4 c7c6 a4c4 b7b5 c4c2 c8b7 g1f3 b8d7 f3e5 f8e7 b1c3 d7e5 d4e5 f6d5 e1g1 e8g8 f1d1 d8b8 c1f4 d5f4 g3f4 c6c5 c3e4\n" +
				"info depth 27 currmove d2d4 currmovenumber 1\n" +
				"info depth 27 currmove e2e4 currmovenumber 2\n" +
				"info depth 27 currmove g1f3 currmovenumber 3\n" +
				"info depth 27 currmove b1c3 currmovenumber 4\n" +
				"info depth 27 currmove c2c4 currmovenumber 5\n" +
				"info depth 27 currmove e2e3 currmovenumber 6\n" +
				"info depth 27 currmove a2a3 currmovenumber 7\n" +
				"info depth 27 currmove c2c3 currmovenumber 8\n" +
				"info depth 27 currmove a2a4 currmovenumber 9\n" +
				"info depth 27 currmove b1a3 currmovenumber 10\n" +
				"info depth 27 currmove d2d3 currmovenumber 11\n" +
				"info depth 27 currmove h2h3 currmovenumber 12\n" +
				"info depth 27 currmove g2g4 currmovenumber 13\n" +
				"info depth 27 currmove b2b4 currmovenumber 14\n" +
				"info depth 27 currmove b2b3 currmovenumber 15\n" +
				"info depth 27 currmove f2f4 currmovenumber 16\n" +
				"info depth 27 currmove g2g3 currmovenumber 17\n" +
				"info depth 27 currmove f2f3 currmovenumber 18\n" +
				"info depth 27 currmove h2h4 currmovenumber 19\n" +
				"info depth 27 currmove g1h3 currmovenumber 20\n" +
				"info depth 31 seldepth 38 multipv 1 score cp 23 upperbound nodes 6765760 nps 1404558 hashfull 992 tbhits 0 time 4817 pv d2d4 g8f6\n" +
				"info depth 27 currmove d2d4 currmovenumber 1\n" +
				"info depth 31 seldepth 38 multipv 1 score cp 25 lowerbound nodes 7408641 nps 1395224 hashfull 995 tbhits 0 time 5310 pv d2d4\n" +
				"info depth 26 currmove d2d4 currmovenumber 1\n" +
				"info depth 31 seldepth 38 multipv 1 score cp 29 lowerbound nodes 7624578 nps 1394654 hashfull 996 tbhits 0 time 5467 pv d2d4\n" +
				"info depth 25 currmove d2d4 currmovenumber 1\n" +
				"info depth 25 currmove e2e4 currmovenumber 2\n" +
				"info depth 25 currmove e2e3 currmovenumber 3\n" +
				"info depth 25 currmove g1f3 currmovenumber 4\n" +
				"info depth 25 currmove b1c3 currmovenumber 5\n" +
				"info depth 25 currmove a2a4 currmovenumber 6\n" +
				"info depth 25 currmove a2a3 currmovenumber 7\n" +
				"info depth 25 currmove c2c4 currmovenumber 8\n" +
				"info depth 25 currmove b2b3 currmovenumber 9\n" +
				"info depth 25 currmove c2c3 currmovenumber 10\n" +
				"info depth 25 currmove b2b4 currmovenumber 11\n" +
				"info depth 25 currmove g2g4 currmovenumber 12\n" +
				"info depth 25 currmove f2f3 currmovenumber 13\n" +
				"info depth 25 currmove d2d3 currmovenumber 14\n" +
				"info depth 25 currmove g2g3 currmovenumber 15\n" +
				"info depth 25 currmove f2f4 currmovenumber 16\n" +
				"info depth 25 currmove h2h3 currmovenumber 17\n" +
				"info depth 25 currmove b1a3 currmovenumber 18\n" +
				"info depth 25 currmove h2h4 currmovenumber 19\n" +
				"info depth 25 currmove g1h3 currmovenumber 20\n" +
				"info depth 31 seldepth 38 multipv 1 score cp 31 nodes 8445571 nps 1391133 hashfull 997 tbhits 0 time 6071 pv d2d4 d7d5 c2c4 e7e6 b1c3 f8b4 g1f3 g8f6 c4d5 e6d5 c1g5 e8g8 e2e3 c7c5 d4c5 b8d7 a1c1 d7c5 d1d4 b4c3 d4c3 c5e4 g5f6 e4f6 c3a3 c8g4\n" +
				"bestmove d2d4 ponder d7d5\n"
			toWriteLines := strings.Split(toWrite, "\n")
			for _, toWriteLine := range toWriteLines {
				_, _ = w.Write([]byte(fmt.Sprint(toWriteLine, "\n")))
				//time.Sleep(1 * time.Millisecond)
			}
		}
	default:
		resp = func(w io.Writer) {
			comm := strings.Split(contents, " ")[0]
			resp := fmt.Sprintf("Unknown command: '%s'. Type help for more information.", comm)
			_, _ = w.Write([]byte(resp))
		}
	}
	//go func(w io.Writer, resp func(w io.Writer)) {
	//	if resp == nil {
	//		return
	//	}
	//time.Sleep(m.delay)
	if resp != nil {
		go resp(m.out)
	}
	//}(m.out, resp)
	return 0, nil
}

func MockCmdClient(resDelay time.Duration) *cmd_client.Client {
	outBuf := bytes.Buffer{}
	mockWriter := NewMockWriter(&outBuf, resDelay)
	return cmd_client.DefaultClient(nil, &outBuf, mockWriter)
}

var _ = Describe("UciClient", func() {
	var uciClient *uci_client.Client
	Describe("::Init", func() {
		var ctx context.Context
		var cancelCtx context.CancelFunc
		BeforeEach(func() {
			ctx, cancelCtx = context.WithTimeout(context.Background(), 100*time.Millisecond)
		})
		AfterEach(func() {
			cancelCtx()
		})
		When("the engine responds to 'uci'", func() {
			BeforeEach(func() {
				cmdClient := MockCmdClient(10 * time.Millisecond)
				uciClient = uci_client.NewUciClient(cmdClient)
			})
			It("does not error", func() {
				Expect(uciClient.Init(ctx)).Error().To(Succeed())
			})
			It("saves the configurable options", func() {
				_, _ = uciClient.Init(ctx)
				Expect(uciClient.IsOption("Threads")).To(BeTrue())
				Expect(uciClient.IsOption("Ponder")).To(BeTrue())
				Expect(uciClient.IsOption("NotAnOption")).ToNot(BeTrue())
			})
			It("returns the configurable options", func() {
				opts, _ := uciClient.Init(ctx)
				Expect(opts).ToNot(BeNil())
			})
		})
		When("the engine does not respond to 'uci'", func() {
			BeforeEach(func() {
				writeBuf := bytes.Buffer{}
				readBuf := bytes.Buffer{}
				cmdClient := cmd_client.DefaultClient(nil, &readBuf, &writeBuf)
				uciClient = uci_client.NewUciClient(cmdClient)
			})
			It("returns an error", func() {
				_, err := uciClient.Init(ctx)
				fmt.Println(err)
				Expect(err).ToNot(Succeed())
			})
		})
	})
	Describe("::SetOption", func() {
		var ctx context.Context
		var cancelCtx context.CancelFunc
		BeforeEach(func() {
			cmdClient := MockCmdClient(10 * time.Millisecond)
			uciClient = uci_client.NewUciClient(cmdClient)
			ctx, cancelCtx = context.WithTimeout(context.Background(), 100*time.Millisecond)
			_, err := uciClient.Init(ctx)
			Expect(err).ToNot(HaveOccurred())
		})
		AfterEach(func() {
			cancelCtx()
		})
		When("the option name is valid", func() {
			When("the value is valid", func() {
				It("does not return an error", func() {
					Expect(uciClient.SetOption(ctx, "Threads", "2")).To(Succeed())
				})
			})
			When("the value is invalid", func() {
				It("returns an error", func() {
					Expect(uciClient.SetOption(ctx, "Threads", "asdf")).ToNot(Succeed())
				})
			})
		})
		When("the option name is invalid", func() {
			It("returns an error", func() {
				Expect(uciClient.SetOption(ctx, "NotAnOption", "some-value")).ToNot(Succeed())
			})
		})
	})
	Describe("::IsReady", func() {
		var ctx context.Context
		var cancelCtx context.CancelFunc
		BeforeEach(func() {
			cmdClient := MockCmdClient(10 * time.Millisecond)
			uciClient = uci_client.NewUciClient(cmdClient)
			ctx, cancelCtx = context.WithTimeout(context.Background(), 100*time.Millisecond)
			Expect(uciClient.Init(ctx)).Error().To(Succeed())
		})
		AfterEach(func() {
			cancelCtx()
		})
		When("the engine is ready", func() {
			It("returns true", func() {
				Expect(uciClient.IsReady(ctx)).To(BeTrue())
			})
		})
	})

	Describe("::Go", func() {
		var ctx context.Context
		var cancelCtx context.CancelFunc
		var opts *uci_client.SearchOptions
		BeforeEach(func() {
			cmdClient := MockCmdClient(10 * time.Millisecond)
			uciClient = uci_client.NewUciClient(cmdClient)
			initCtx, cancelInitCtx := context.WithTimeout(context.Background(), 100*time.Millisecond)
			_, err := uciClient.Init(initCtx)
			Expect(err).ToNot(HaveOccurred())
			cancelInitCtx()

			ctx, cancelCtx = context.WithTimeout(context.Background(), 1*time.Second)

			opts = uci_client.NewSearchOptionsBuilder().WithWhiteMs(100000).Build()
		})
		AfterEach(func() {
			cancelCtx()
		})
		When("the engine is ready", func() {
			It("returns the best move in long algebraic notation", func() {
				Expect(uciClient.Go(ctx, opts)).To(Equal("d2d4"))
			})
		})
	})
})
