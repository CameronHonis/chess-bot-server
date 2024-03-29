package cmd_client

import (
	"context"
	"fmt"
	"io"
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
