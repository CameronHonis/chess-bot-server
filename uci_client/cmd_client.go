package uci_client

import (
	"context"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"time"
)

const DEFAULT_STDOUT_BUF_BYTE_SIZE = 4096

type ByteDump []byte

// String trims all null bytes from end of string
func (b ByteDump) String() string {
	var n = 0
	for ; n < len(b); n++ {
		if b[n] == 0 {
			break
		}
	}
	return string(b[:n])
}

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

// CmdClient is a friendly wrapper around a running exec.Cmd that allows easy reads on constantly changing Stdout
type CmdClient struct {
	r           io.Reader
	w           io.Writer
	readBufSize uint
}

func NewCmdClient(r io.Reader, w io.Writer, readBufSize uint) *CmdClient {
	return &CmdClient{
		r:           r,
		w:           w,
		readBufSize: readBufSize,
	}
}

// CmdClientFromCmd takes a running command
func CmdClientFromCmd(cmd exec.Cmd) (*CmdClient, error) {
	wc, openWriterErr := cmd.StdinPipe()

	if openWriterErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, could not open writer to cmd: %s", openWriterErr)
	}

	r, openReaderErr := cmd.StdoutPipe()
	if openReaderErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, coud not open reader to cmd: %s", openReaderErr)
	}

	return &CmdClient{r, wc, DEFAULT_STDOUT_BUF_BYTE_SIZE}, nil
}

func (cc *CmdClient) Readlines(ctx context.Context, ch chan string) {
	br := &BlockingReader{cc.r, ctx}
	var carryLine string

	for {
		p := make(ByteDump, cc.readBufSize)
		n, err := br.Read(p)
		if err != nil {
			close(ch)
			break
		}
		lines := strings.Split(p.String(), "\n")

		// handle last carryLine
		if carryLine != "" {
			if len(lines) == 0 {
				lines = []string{fmt.Sprintf("%s", carryLine), ""}
			} else {
				lines[0] = fmt.Sprintf("%s%s", carryLine, lines[0])
			}
			carryLine = ""
		}
		// set carryLine or trim empty string
		endsWithNewline := lines[len(lines)-1] == ""
		if n == int(cc.readBufSize) && !endsWithNewline {
			carryLine = lines[len(lines)-1]
			lines = lines[:len(lines)-1]
		} else if endsWithNewline {
			lines = lines[:len(lines)-1]
		}
		for _, line := range lines {
			ch <- line
		}
	}
}

func (cc *CmdClient) WriteString(s string) error {
	_, err := cc.w.Write([]byte(s))
	return err
}

func (cc *CmdClient) SetBufSize(size uint) {
	cc.readBufSize = size
}
