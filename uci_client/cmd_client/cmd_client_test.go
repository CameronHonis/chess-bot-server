package cmd_client_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sync"
	"time"
)

type ReaderWriter struct {
	buf bytes.Buffer
	mu  sync.Mutex
}

func NewReaderWriter() *ReaderWriter {
	return &ReaderWriter{
		buf: bytes.Buffer{},
		mu:  sync.Mutex{},
	}
}

func (rw *ReaderWriter) Read(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.buf.Read(p)
}

func (rw *ReaderWriter) Write(p []byte) (n int, err error) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	return rw.buf.Write(p)
}

func (rw *ReaderWriter) WriteLine(s string) {
	rw.mu.Lock()
	defer rw.mu.Unlock()
	rw.buf.WriteString(fmt.Sprintf("%s\n", s))
}

var _ = Describe("Client", func() {
	var readerWriter *ReaderWriter
	var cmdClient *cmd_client.Client
	var ctx context.Context
	var ctxCancel context.CancelFunc
	BeforeEach(func() {
		readerWriter = NewReaderWriter()
		ctx, ctxCancel = context.WithTimeout(context.Background(), 55*time.Millisecond)
		cmdClient = cmd_client.NewClient(readerWriter, readerWriter)
	})
	AfterEach(func() {
		ctxCancel()
	})
	Describe("ReadLine", func() {
		When("output is available at read time", func() {
			BeforeEach(func() {
				readerWriter.WriteLine("this is the first line")
			})
			When("additional output becomes available within the lifetime of the context", func() {
				BeforeEach(func() {
					go func() {
						time.Sleep(50 * time.Millisecond)
						readerWriter.WriteLine("the second line")
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
						time.Sleep(10 * time.Millisecond)
						Expect(readerWriter.Write([]byte("this is a line without the 'newline' char at the end"))).Error().To(Succeed())
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
					Expect(readerWriter.Write([]byte("123456789"))).Error().To(Succeed())
				})
				It("channels the output regardless", func() {
					Expect(cmdClient.ReadLine(ctx)).To(Equal("123456789"))
				})
			})
			When("the output is 1.5x the buffer size", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(6)
					readerWriter.WriteLine("123456789")
				})
				It("channels the output regardless", func() {
					Expect(cmdClient.ReadLine(ctx)).To(Equal("123456789"))
				})
			})
			When("multiple lines in output wrap into the next buffer frame", func() {
				BeforeEach(func() {
					cmdClient.SetBufSize(4)
					readerWriter.WriteLine("a\nbcdef\nghk")
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
