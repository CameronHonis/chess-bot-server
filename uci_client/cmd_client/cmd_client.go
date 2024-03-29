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

func NewCmdClient(r io.Reader, w io.Writer) *CmdClient {
	return &CmdClient{
		r:             r,
		w:             w,
		_readBufSize:  4096,
		_flushOnWrite: true,
		_isReading:    false,
		_lines:        make([]string, 0),
		mu:            sync.Mutex{},
	}
}

// CmdClientFromCmd expects a running command
func CmdClientFromCmd(cmd *exec.Cmd) (*CmdClient, error) {
	if cmd.Process == nil {
		return nil, fmt.Errorf("cmd must be running before creating CmdClient")
	}

	w, openWriterErr := cmd.StdinPipe()
	if openWriterErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, could not open writer to cmd: %s", openWriterErr)
	}

	r, openReaderErr := cmd.StdoutPipe()
	if openReaderErr != nil {
		return nil, fmt.Errorf("cannot create CmdClient, coud not open reader to cmd: %s", openReaderErr)
	}

	return NewCmdClient(r, w), nil
}

func (cc *CmdClient) ReadLine(ctx context.Context) (string, error) {
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

func (cc *CmdClient) WriteString(s string) error {
	if cc.flushOnWrite() {
		cc.flushLines()
	}
	_, err := cc.w.Write([]byte(s))
	return err
}

func (cc *CmdClient) SetBufSize(size uint) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._readBufSize = size
}

func (cc *CmdClient) SetFlushOnWrite(flushOnWrite bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._flushOnWrite = flushOnWrite
}

func (cc *CmdClient) FlushReader() {
	cc.flushLines()
}

func (cc *CmdClient) readLines(ctx context.Context) {
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

func (cc *CmdClient) pushLines(lines ...string) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._lines = append(cc._lines, lines...)
}

func (cc *CmdClient) popLine() (string, error) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if len(cc._lines) == 0 {
		return "", fmt.Errorf("no lines to pop")
	}
	line := cc._lines[0]
	cc._lines = cc._lines[1:]
	return line, nil
}

func (cc *CmdClient) flushLines() {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._lines = make([]string, 0)
}

func (cc *CmdClient) readBufSize() uint {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._readBufSize
}

func (cc *CmdClient) flushOnWrite() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._flushOnWrite
}

func (cc *CmdClient) isReading() bool {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc._isReading
}

func (cc *CmdClient) setIsReading(isReading bool) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	cc._isReading = isReading
}
