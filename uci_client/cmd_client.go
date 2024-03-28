package uci_client

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

// BlockingReader simply reads from the provided reader until no content is left, at which point it waits for more.
type BlockingReader struct {
	r   io.Reader
	ctx context.Context
}

func NewBlockingReader(r io.Reader, ctx context.Context) *BlockingReader {
	return &BlockingReader{
		r:   r,
		ctx: ctx,
	}
}

func (br *BlockingReader) Read(p []byte) (n int, err error) {
	for {
		select {
		case <-br.ctx.Done():
			return 0, fmt.Errorf("timeout before contents")
		default:
			n, err = br.r.Read(p)
			if err != io.EOF {
				return
			}
			readContent := n > 0
			if readContent {
				return
			}
			time.Sleep(time.Millisecond) // dont needlessly hog cpu
		}
	}
}

type CmdClient struct {
	r           io.Reader
	wc          io.WriteCloser
	readBufSize uint
}

// NewCmdClient takes a **running** command
func NewCmdClient(cmd exec.Cmd, readBufSize uint) (*CmdClient, error) {
	wc, openWriterErr := cmd.StdinPipe()

	if openWriterErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, could not open writer to cmd: %s", openWriterErr)
	}

	r, openReaderErr := cmd.StdoutPipe()
	if openReaderErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, coud not open reader to cmd: %s", openReaderErr)
	}

	return &CmdClient{r, wc, readBufSize}, nil
}

func (cc *CmdClient) Readlines(ctx context.Context, ch chan string) {
	br := &BlockingReader{cc.r, ctx}
	var carryLine string

	p := make([]byte, cc.readBufSize)
	for {
		_, err := br.Read(p)
		if err != nil {
			close(ch)
			break
		}
		lines := strings.Split(string(p), "\n")
		if len(lines) > 0 {
			if carryLine != "" {
				lines[0] = fmt.Sprintf("%s%s", carryLine, lines[0])
			}
			endsWithNewline := lines[len(lines)-1] == ""
			if !endsWithNewline {
				carryLine = lines[len(lines)-1]
				lines = lines[:len(lines)-1]
			}
		}
	}
}

func (cc *CmdClient) WriteString(s string) error {
	_, err := cc.wc.Write([]byte(s))
	return err
}
