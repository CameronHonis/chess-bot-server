package uci_client_test

import (
	"bytes"
	"context"
	"github.com/CameronHonis/chess-bot-server/uci_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = FDescribe("CmdClient", func() {
	var buf bytes.Buffer
	var cmdClient *uci_client.CmdClient
	var ctx context.Context
	var ctxCancel context.CancelFunc
	BeforeEach(func() {
		buf = bytes.Buffer{}
		ctx, ctxCancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		cmdClient = uci_client.NewCmdClient(&buf, &buf, 4096)
	})
	AfterEach(func() {
		ctxCancel()
	})
	Describe("Readlines", func() {
		When("output is available at read time", func() {
			BeforeEach(func() {
				buf.WriteString("this is the first line\n")
			})
			When("additional output becomes available within the lifetime of the context", func() {
				BeforeEach(func() {
					go func() {
						time.Sleep(50 * time.Millisecond)
						buf.WriteString("the second line\n")
					}()
				})
				It("channels all available output within the context lifetime", func() {
					var recvStr string
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("this is the first line"))
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("the second line"))
				})
			})
			It("immediately channels all available output", func() {
				var recvStr string
				ch := make(chan string)
				go cmdClient.Readlines(ctx, ch)
				Eventually(ch).Should(Receive(&recvStr))
				Expect(recvStr).To(Equal("this is the first line"))
			})
		})
		When("output is not available at read time", func() {
			When("output becomes available within the lifetime of the context", func() {
				BeforeEach(func() {
					go func() {
						time.Sleep(50 * time.Millisecond)
						buf.WriteString("this is a line without the 'newline' char at the end")
					}()
				})
				It("channels all available output within the context lifetime", func() {
					var recvStr string
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("this is a line without the 'newline' char at the end"))
				})
			})
			When("no output becomes available within the lifetime of the context", func() {
				It("does not channel any output", func() {
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Consistently(ch).ShouldNot(Receive())
				})
				It("closes the channel", func() {
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(BeClosed())
				})
			})
		})
		When("the output exceeds the buffer size", func() {
			When("the output is more than double the buffer size", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(4)
					buf.WriteString("123456789")
				})
				It("channels the output regardless", func() {
					var recvStr string
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("123456789"))
				})
			})
			When("the output is 1.5x the buffer size", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(6)
					buf.WriteString("123456789")
				})
				It("channels the output regardless", func() {
					var recvStr string
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("123456789"))
				})
			})
			When("multiple lines in output wrap into the next buffer frame", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(4)
					buf.WriteString("a\nbcdef\nghk")
				})
				It("channels the output irrespective of the buffer size", func() {
					var recvStr string
					ch := make(chan string)
					go cmdClient.Readlines(ctx, ch)
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("a"))
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("bcdef"))
					Eventually(ch).Should(Receive(&recvStr))
					Expect(recvStr).To(Equal("ghk"))
				})
			})
		})
	})
})

var _ = Describe("BlockingReader", func() {
	var buf bytes.Buffer
	var br *uci_client.BlockingReader
	var cancelCtx context.CancelFunc
	BeforeEach(func() {
		var ctx context.Context
		ctx, cancelCtx = context.WithTimeout(context.Background(), 10*time.Millisecond)
		buf = bytes.Buffer{}
		br = uci_client.NewBlockingReader(&buf, ctx)
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
