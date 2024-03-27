package uci_client

import (
	"context"
	"fmt"
	"github.com/CameronHonis/set"
	"io"
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
func (u *Client) Init(ctx context.Context) (*set.Set[string], error) {
	_, writeErr := u.w.Write([]byte("uci"))
	if writeErr != nil {
		return nil, fmt.Errorf("could not write to uci client: %s", writeErr)
	}

	var bytes []byte

	for len(bytes) == 0 {
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("ctx finished before engine response")
		default:
			var readErr error
			bytes, readErr = io.ReadAll(u.r)
			if readErr != nil {
				return nil, fmt.Errorf("could not read from uci client: %s", readErr)
			}
		}
	}

	contents := string(bytes)
	fmt.Println(contents)

	return set.EmptySet[string](), nil

}

//func (u *Client) SetOption()
