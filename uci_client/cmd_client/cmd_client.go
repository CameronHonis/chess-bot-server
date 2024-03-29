package cmd_client

import (
	"context"
	"fmt"
	"github.com/CameronHonis/marker"
	"io"
	"os/exec"
	"strings"
	"sync"
	"time"
)

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

// Client is a friendly wrapper around a running exec.Cmd that allows easy reads on constantly changing Stdout
type Client struct {
	__static__ marker.Marker
	r          io.Reader
	w          io.Writer

	__config__    marker.Marker // these should be safe to change while processing io
	_readBufSize  uint
	_flushOnWrite bool

	__dynamic__ marker.Marker // should always require mutex lock to manipulate internally
	_isReading  bool
	_lines      []string
	mu          sync.Mutex
}

func NewClient(r io.Reader, w io.Writer) *Client {
	return &Client{
		r:             r,
		w:             w,
		_readBufSize:  4096,
		_flushOnWrite: true,
		_isReading:    false,
		_lines:        make([]string, 0),
		mu:            sync.Mutex{},
	}
}

// ClientFromCmd expects a running command
func ClientFromCmd(cmd *exec.Cmd) (*Client, error) {
	if cmd.Process == nil {
		return nil, fmt.Errorf("cmd must be running before creating Client")
	}

	w, openWriterErr := cmd.StdinPipe()
	if openWriterErr != nil {
		return nil, fmt.Errorf("cannot create Client, could not open writer to cmd: %s", openWriterErr)
	}

	r, openReaderErr := cmd.StdoutPipe()
	if openReaderErr != nil {
		return nil, fmt.Errorf("cannot create Client, coud not open reader to cmd: %s", openReaderErr)
	}

	return NewClient(r, w), nil
}

func (cc *Client) ReadLine(ctx context.Context) (string, error) {
	go cc.readLines(ctx)
	for {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("timeout before next line")
		default:
			line, popErr := cc.popLine()
			if popErr == nil {
				return line, nil
			}
			time.Sleep(time.Millisecond)
		}
	}
}

func (cc *Client) WriteString(s string) error {
	if cc.flushOnWrite() {
		cc.flushLines()
	}
	_, err := cc.w.Write([]byte(s))
	return err
}

func (cc *Client) SetBufSize(size uint) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._readBufSize = size
}

func (cc *Client) SetFlushOnWrite(flushOnWrite bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._flushOnWrite = flushOnWrite
}

func (cc *Client) FlushReader() {
	cc.flushLines()
}

func (cc *Client) readLines(ctx context.Context) {
	if cc.isReading() {
		return
	}
	cc.setIsReading(true)

	br := &BlockingReader{cc.r, ctx}
	var carryLine string

	for {
		p := make(ByteDump, cc.readBufSize())
		n, err := br.Read(p)
		if err != nil {
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
		if n == int(cc.readBufSize()) && !endsWithNewline {
			carryLine = lines[len(lines)-1]
			lines = lines[:len(lines)-1]
		} else if endsWithNewline {
			lines = lines[:len(lines)-1]
		}

		cc.pushLines(lines...)
	}
	cc.setIsReading(false)
}

func (cc *Client) pushLines(lines ...string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._lines = append(cc._lines, lines...)
}

func (cc *Client) popLine() (string, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if len(cc._lines) == 0 {
		return "", fmt.Errorf("no lines to pop")
	}
	line := cc._lines[0]
	cc._lines = cc._lines[1:]
	return line, nil
}

func (cc *Client) flushLines() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._lines = make([]string, 0)
}

func (cc *Client) readBufSize() uint {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._readBufSize
}

func (cc *Client) flushOnWrite() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._flushOnWrite
}

func (cc *Client) isReading() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._isReading
}

func (cc *Client) setIsReading(isReading bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._isReading = isReading
}
