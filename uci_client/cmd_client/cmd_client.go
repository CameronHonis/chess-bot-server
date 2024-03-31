package cmd_client

import (
	"context"
	"fmt"
	"github.com/CameronHonis/marker"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const PRINT_IO = true

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

// ReaderWriterProxy is solely used for debugging purposes, although it is made generalized for any reader/writer
// intercept tasks.
type ReaderWriterProxy struct {
	r       io.Reader
	w       io.Writer
	onRead  func(p ByteDump, n int, err error)
	onWrite func(p ByteDump, n int, err error)
	onClose func(error)
}

func DebuggingReaderWriterProxy(r io.Reader, w io.Writer) *ReaderWriterProxy {
	var onRead = func(p ByteDump, n int, err error) {
		if err != nil {
			return
		}
		if PRINT_IO {
			s := string(p[:n])
			fmt.Println(">", s)
		}
	}
	var onWrite = func(p ByteDump, n int, err error) {
		if err != nil {
			return
		}
		if PRINT_IO {
			s := string(p[:n])
			fmt.Println("<", s)
		}
	}
	return &ReaderWriterProxy{
		r:       r,
		w:       w,
		onRead:  onRead,
		onWrite: onWrite,
	}
}

func (rwp *ReaderWriterProxy) Read(p []byte) (n int, err error) {
	n, err = rwp.r.Read(p)
	if rwp.onRead != nil {
		rwp.onRead(p, n, err)
	}
	return
}

func (rwp *ReaderWriterProxy) Write(p []byte) (n int, err error) {
	n, err = rwp.w.Write(p)
	if rwp.onWrite != nil {
		rwp.onWrite(p, n, err)
	}
	return
}

// Client is a friendly wrapper around a running exec.Cmd that allows easy reads on constantly changing Stdout
type Client struct {
	__static__ marker.Marker
	cmd        *exec.Cmd
	stdout     io.Reader
	stdin      io.Writer

	__config__    marker.Marker // these should be safe to change while processing io
	_readBufSize  uint
	_flushOnWrite bool

	__dynamic__ marker.Marker // should always require mutex lock to manipulate internally
	_isReading  bool
	_lines      []string
	mu          sync.Mutex
}

func DefaultClient(cmd *exec.Cmd, r io.Reader, w io.Writer) *Client {
	return &Client{
		cmd:           cmd,
		stdout:        r,
		stdin:         w,
		_readBufSize:  4096,
		_flushOnWrite: true,
		_isReading:    false,
		_lines:        make([]string, 0),
		mu:            sync.Mutex{},
	}
}

// ClientFromCmd expects a running command
func ClientFromCmd(cmd *exec.Cmd) (*Client, error) {
	if cmd.Process != nil {
		return nil, fmt.Errorf("cmd must not be running before creating Client")
	}

	w, openWriterErr := cmd.StdinPipe()
	if openWriterErr != nil {
		return nil, fmt.Errorf("cannot create Client, could not open writer to cmd: %s", openWriterErr)
	}

	r, openReaderErr := cmd.StdoutPipe()
	if openReaderErr != nil {
		return nil, fmt.Errorf("cannot create Client, coud not open reader to cmd: %s", openReaderErr)
	}

	readerWriterProxy := DebuggingReaderWriterProxy(r, w)

	return DefaultClient(cmd, readerWriterProxy, readerWriterProxy), nil
}

// ReadLine is a blocking read on the next line from Stdout. If the context expires, ReadLine
// will return an error indicating as such.
func (cc *Client) ReadLine(ctx context.Context) (string, error) {
	if !cc.isReading() {
		go cc.readLines(context.Background())
	}
	for {
		select {
		case <-ctx.Done():
			return "", NewReaderTimeout("timed out before next line")
		default:
			line, popErr := cc.popLine()
			if popErr == nil {
				return line, nil
			}
			time.Sleep(time.Millisecond)
		}
	}
}

func (cc *Client) WriteLine(s string) error {
	if cc.flushOnWrite() {
		cc.flushLines()
	}

	line := fmt.Sprintf("%s\n", s)
	_, err := cc.stdin.Write([]byte(line))
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

func (cc *Client) End() error {
	if cc.cmd.Process.Signal(os.Interrupt) != nil {
		fmt.Println("WARN: error sending interrupt signal, killing process instead")
		return cc.cmd.Process.Kill()
	}
	return nil
}

func (cc *Client) IsRunning() bool {
	return cc.cmd.ProcessState == nil
}

func (cc *Client) readLines(ctx context.Context) {
	if cc.isReading() {
		return
	}
	cc.setIsReading(true)

	br := &BlockingReader{cc.stdout, ctx}
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
