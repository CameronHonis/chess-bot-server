package uci_client

import (
	"context"
	"fmt"
	"github.com/CameronHonis/set"
	"io"
	"strings"
	"time"
)

// Client represents a client for any engine supporting UCI (Universal Chess Interface)
// This interface is outlined [here](https://www.stmintz.com/ccc/index.php?id=141612)
type Client struct {
	r    io.Reader
	w    io.Writer
	opts *set.Set[string]
}

func NewUciClient(r io.Reader, w io.Writer) *Client {

	return &Client{
		r:    r,
		w:    w,
		opts: set.EmptySet[string](),
	}
}

// Init tells the engine to use the uci protocol and stores the configurable options.
// It returns the set of options that are configurable.
func (c *Client) Init(ctx context.Context) (*set.Set[string], error) {
	c.flushReader()
	_, writeErr := c.w.Write([]byte("uci"))
	if writeErr != nil {
		return nil, fmt.Errorf("could not write to uci client: %s", writeErr)
	}

	resp, readErr := waitForEngineRes(ctx, c.r)
	if readErr != nil {
		return nil, readErr
	}

	for _, line := range strings.Split(resp, "\n") {
		if strings.HasPrefix(line, "option name") {
			optionDetails := line[len("option name "):]
			optionName := strings.Split(optionDetails, " ")[0]
			c.opts.Add(optionName)
		}
	}

	return c.opts.Copy(), nil
}

func (c *Client) IsOption(optName string) bool {
	return c.opts.Has(optName)
}

func (c *Client) SetOption(optName string, optVal string) error {
	c.flushReader()

	_, writeErr := c.w.Write([]byte(fmt.Sprintf("setoption name %s value %s", optName, optVal)))
	if writeErr != nil {
		return fmt.Errorf("could not write to uci client: %s", writeErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	resp, readErr := waitForEngineRes(ctx, c.r)
	if readErr != nil {
		return nil
	}

	return fmt.Errorf("error setting option: %s", resp)
}

func (c *Client) SetPosition(fen string) error {
	_, writeErr := c.w.Write([]byte(fmt.Sprintf("position fen %s", fen)))
	if writeErr != nil {
		return fmt.Errorf("could not write to uci client: %s", writeErr)
	}
	return nil
}

func (c *Client) flushReader() {
	_, _ = io.ReadAll(c.r)
}

func waitForEngineRes(ctx context.Context, r io.Reader) (string, error) {
	var bytes []byte
	for len(bytes) == 0 {
		select {
		case <-ctx.Done():
			return "", fmt.Errorf("ctx finished before engine response")
		default:
			var readErr error
			bytes, readErr = io.ReadAll(r)
			if readErr != nil {
				return "", fmt.Errorf("could not read from uci client: %s", readErr)
			}
		}
	}
	return string(bytes), nil
}
