package uci_client

import (
	"context"
	"fmt"
	"github.com/CameronHonis/chess-bot-server/uci_client/cmd_client"
	"github.com/CameronHonis/set"
	"os/exec"
	"strings"
)

// Client represents a client for any engine supporting UCI (Universal Chess Interface)
// This interface is outlined [here](https://www.stmintz.com/ccc/index.php?id=141612)
type Client struct {
	CmdClient *cmd_client.Client
	opts      *set.Set[string]
}

func NewUciClient(client *cmd_client.Client) *Client {
	return &Client{
		CmdClient: client,
		opts:      set.EmptySet[string](),
	}
}

func UciClientFromCmd(cmd *exec.Cmd) (*Client, error) {
	cmdClient, cmdClientErr := cmd_client.ClientFromCmd(cmd)
	if cmdClientErr != nil {
		return nil, fmt.Errorf("could not make CmdClient: %s", cmdClientErr)
	}
	return NewUciClient(cmdClient), nil
}

// Init tells the engine to use the uci protocol and stores the configurable options.
// It returns the set of options that are configurable.
func (c *Client) Init(ctx context.Context) (*set.Set[string], error) {
	c.CmdClient.SetFlushOnWrite(true)
	writeErr := c.CmdClient.WriteLine("uci")
	if writeErr != nil {
		return nil, fmt.Errorf("could not write to uci CmdClient: %s", writeErr)
	}

	for {
		resp, readErr := c.CmdClient.ReadLine(ctx)
		if readErr != nil {
			return nil, fmt.Errorf("could not read output after init: %s", readErr)
		}
		if strings.HasPrefix(resp, "option name") {
			optionDetails := resp[len("option name "):]
			optionName := strings.Split(optionDetails, " ")[0]
			c.opts.Add(optionName)
		}
		if resp == "uciok" {
			break
		}
	}
	return c.opts.Copy(), nil
}

func (c *Client) IsOption(optName string) bool {
	return c.opts.Has(optName)
}

func (c *Client) SetOption(ctx context.Context, optName string, optVal string) error {
	writeErr := c.CmdClient.WriteLine(fmt.Sprintf("setoption name %s value %s", optName, optVal))
	if writeErr != nil {
		return fmt.Errorf("could not write to uci CmdClient: %s", writeErr)
	}

	resp, readErr := c.CmdClient.ReadLine(ctx) // Only expect set config errors to be received here
	if readErr != nil {
		if _, ok := readErr.(*cmd_client.ReaderTimeout); ok {
			return nil
		} else {
			return fmt.Errorf("cannot read output while setting option: %s", readErr)
		}
	}
	return fmt.Errorf("cannot set option: %s", resp)
}

func (c *Client) SetPosition(fen string) error {
	writeErr := c.CmdClient.WriteLine(fmt.Sprintf("position fen %s", fen))
	if writeErr != nil {
		return fmt.Errorf("could not write to uci CmdClient: %s", writeErr)
	}
	return nil
}

func (c *Client) IsReady(ctx context.Context) (bool, error) {
	writeErr := c.CmdClient.WriteLine("isready")
	if writeErr != nil {
		return false, fmt.Errorf("could not write to uci CmdClient %s", writeErr)
	}

	resp, readErr := c.CmdClient.ReadLine(ctx)
	if readErr != nil {
		return false, readErr
	}
	return resp == "readyok", nil
}

func (c *Client) Go(ctx context.Context, opts *SearchOptions) (string, error) {
	cmd, cmdErr := searchOptionsToCmdStr(opts)
	if cmdErr != nil {
		return "", fmt.Errorf("cannot generate search command: %s", cmdErr)
	}

	writeErr := c.CmdClient.WriteLine(cmd)
	if writeErr != nil {
		return "", fmt.Errorf("could not write to uci CmdClient %s", writeErr)
	}

	for {
		resp, readErr := c.CmdClient.ReadLine(ctx)
		if readErr != nil {
			return "", fmt.Errorf("read error while listening for best move: %s", readErr)
		}
		if strings.HasPrefix(resp, "bestmove") {
			bestMoveDetails := strings.Split(resp, " ")
			if len(bestMoveDetails) <= 1 {
				return "", fmt.Errorf("malformed bestmove response: %s", resp)
			}
			return bestMoveDetails[1], nil
		}
	}
}

func (c *Client) End() error {
	return c.CmdClient.End()
}

func searchOptionsToCmdStr(opts *SearchOptions) (string, error) {
	optsVetErr := opts.Vet()
	if optsVetErr != nil {
		return "", fmt.Errorf("cannot convert search options to command string, search options vetting failed: %s", optsVetErr)
	}
	sb := &strings.Builder{}
	sb.WriteString("go ")
	if len(opts.SearchMoves) > 0 {
		sb.WriteString("searchmoves ")
		for _, searchMove := range opts.SearchMoves {
			sb.WriteString(fmt.Sprintf("%s ", searchMove))
		}
	}
	if opts.WhiteMs > 0 {
		sb.WriteString(fmt.Sprintf("wtime %d ", opts.WhiteMs))
	}
	if opts.BlackMs > 0 {
		sb.WriteString(fmt.Sprintf("btime %d ", opts.BlackMs))
	}
	if opts.WhiteIncrMs > 0 {
		sb.WriteString(fmt.Sprintf("winc %d ", opts.WhiteIncrMs))
	}
	if opts.BlackIncrMs > 0 {
		sb.WriteString(fmt.Sprintf("binc %d ", opts.BlackIncrMs))
	}
	if opts.MovesTillIncr > 0 {
		sb.WriteString(fmt.Sprintf("movestogo %d ", opts.MovesTillIncr))
	}
	if opts.Depth > 0 {
		sb.WriteString(fmt.Sprintf("depth %d ", opts.Depth))
	}
	if opts.SearchMs > 0 {
		sb.WriteString(fmt.Sprintf("movetime %d ", opts.SearchMs))
	}
	return strings.TrimSpace(sb.String()), nil
}
