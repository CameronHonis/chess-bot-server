package cmd_client_test

import (
	"bytes"
	"context"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("CmdClient", func() {
	var buf bytes.Buffer
	var cmdClient *cmd_client.CmdClient
	var ctx context.Context
	var ctxCancel context.CancelFunc
	BeforeEach(func() {
		buf = bytes.Buffer{}
		ctx, ctxCancel = context.WithTimeout(context.Background(), 100*time.Millisecond)
		cmdClient = cmd_client.NewCmdClient(&buf, &buf)
	})
	AfterEach(func() {
		ctxCancel()
	})
	Describe("readlines", func() {
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
					Expect(cmdClient.ReadLine(ctx)).To(Equal("this is the first line"))
					Expect(cmdClient.ReadLine(ctx)).To(Equal("the second line"))
				})
			})
			It("immediately channels all available output", func() {
				Expect(cmdClient.ReadLine(ctx)).To(Equal("this is the first line"))
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
					Expect(cmdClient.ReadLine(ctx)).To(Equal("this is a line without the 'newline' char at the end"))
				})
			})
			When("no output becomes available within the lifetime of the context", func() {
				It("returns an error", func() {
					Expect(cmdClient.ReadLine(ctx)).Error().To(HaveOccurred())
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
					Expect(cmdClient.ReadLine(ctx)).To(Equal("123456789"))
				})
			})
			When("the output is 1.5x the buffer size", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(6)
					buf.WriteString("123456789")
				})
				It("channels the output regardless", func() {
					Expect(cmdClient.ReadLine(ctx)).To(Equal("123456789"))
				})
			})
			When("multiple lines in output wrap into the next buffer frame", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(4)
					buf.WriteString("a\nbcdef\nghk")
				})
				It("channels the output irrespective of the buffer size", func() {
					Expect(cmdClient.ReadLine(ctx)).To(Equal("a"))
					Expect(cmdClient.ReadLine(ctx)).To(Equal("bcdef"))
					Expect(cmdClient.ReadLine(ctx)).To(Equal("ghk"))
				})
			})
		})
	})
})

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
