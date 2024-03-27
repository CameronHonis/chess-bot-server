package uci_client_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/CameronHonis/chess-bot-server/uci_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"io"
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
	time.Sleep(m.delay)
	contents := string(p)
	switch contents {
	case "uci":
		return m.out.Write([]byte("id name Stockfish dev-20240314-fb07281f\n" +
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
	default:
		return 0, fmt.Errorf("unknown command")
	}
}

func MockReaderWriter(resDelay time.Duration) (io.Reader, io.Writer) {
	outBuf := bytes.Buffer{}
	return &outBuf, NewMockWriter(&outBuf, resDelay)
}

var _ = Describe("Client", func() {
	var client *uci_client.Client
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
				mockReader, mockWriter := MockReaderWriter(10 * time.Millisecond)
				client = uci_client.NewUciClient(mockReader, mockWriter)
			})
			It("does not error", func() {
				_, err := client.Init(ctx)
				Expect(err).To(Succeed())
			})
		})
		When("the engine does not respond to 'uci'", func() {
			BeforeEach(func() {
				writeBuf := bytes.Buffer{}
				readBuf := bytes.Buffer{}
				client = uci_client.NewUciClient(&readBuf, &writeBuf)
			})
			It("returns an error", func() {
				_, err := client.Init(ctx)
				fmt.Println(err)
				Expect(err).ToNot(Succeed())
			})
		})
	})
})
