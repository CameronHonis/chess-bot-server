package cmd_client_test

import (
	"bytes"
	"context"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("BlockingReader", func() {
	var buf bytes.Buffer
	var br *cmd_client.BlockingReader
	var cancelCtx context.CancelFunc
	BeforeEach(func() {
		var ctx context.Context
		ctx, cancelCtx = context.WithTimeout(context.Background(), 10*time.Millisecond)
		buf = bytes.Buffer{}
		br = cmd_client.NewBlockingReader(&buf, ctx)
	})
	AfterEach(func() {
		cancelCtx()
	})
	When("a message is already available for the reader", func() {
		BeforeEach(func() {
			buf.WriteByte('a')
		})
		It("returns the message", func() {
			p := make([]byte, 1)
			_, readErr := br.Read(p)
			Expect(readErr).To(Succeed())
			Expect(p).To(Equal([]byte{'a'}))
		})
	})
	When("a message is not already available for the reader", func() {
		When("a message becomes available before the context lifetime", func() {
			BeforeEach(func() {
				go func() {
					time.Sleep(5 * time.Millisecond)
					buf.WriteByte('b')
				}()
			})
			It("returns the message", func() {
				p := make([]byte, 1)
				_, readErr := br.Read(p)
				Expect(readErr).To(Succeed())
				Expect(p).To(Equal([]byte{'b'}))
			})
		})
		When("a message does not become available for the reader during the context lifetime", func() {
			It("returns an error", func() {
				p := make([]byte, 1)
				_, readErr := br.Read(p)
				Expect(readErr).ToNot(Succeed())
				expP := make([]byte, 1)
				Expect(p).To(Equal(expP))
			})
		})
	})
})
