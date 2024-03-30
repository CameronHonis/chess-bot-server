package uci_client

import (
	"fmt"
	"strings"
)

type SearchOptions struct {
	SearchMoves   []string // list of moves to search from root formatted in UCI-compliant, long algebraic notation
	WhiteMs       uint
	BlackMs       uint
	WhiteIncrMs   uint
	BlackIncrMs   uint
	MovesTillIncr uint
	Depth         uint // search exactly until this depth
	SearchMs      uint // search exactly this amount of ms
}

func (so *SearchOptions) Vet() error {
	if so.Depth != 0 && so.SearchMs != 0 {
		return fmt.Errorf("cannot set both Depth and SearchMs")
	}
	for _, searchMove := range so.SearchMoves {
		var isValid = true
		if len(searchMove) == 4 {
			isValid = true
		} else if len(searchMove) == 5 {
			isUpgrade := searchMove[3] == '8' || searchMove[3] == '1' && (strings.HasSuffix(searchMove, "q") ||
				strings.HasSuffix(searchMove, "Q") ||
				strings.HasSuffix(searchMove, "r") ||
				strings.HasSuffix(searchMove, "R") ||
				strings.HasSuffix(searchMove, "b") ||
				strings.HasSuffix(searchMove, "B") ||
				strings.HasSuffix(searchMove, "n") ||
				strings.HasSuffix(searchMove, "N"))
			isValid = isUpgrade
		} else {
			isValid = false
		}
		if searchMove[0] < 'a' || searchMove[0] > 'h' {
			isValid = false
		} else if searchMove[1] < '1' || searchMove[1] > '8' {
			isValid = false
		} else if searchMove[2] < 'a' || searchMove[2] > 'h' {
			isValid = false
		} else if searchMove[3] < '1' || searchMove[3] > '8' {
			isValid = false
		}
		if !isValid {
			return fmt.Errorf("invalid search move %s", searchMove)
		}
	}
	return nil
}

type SearchOptionsBuilder struct {
	searchOptions *SearchOptions
}

func (sob *SearchOptionsBuilder) WithSearchMoves(searchMoves []string) *SearchOptionsBuilder {
	sob.searchOptions.SearchMoves = searchMoves
	return sob
}

func (sob *SearchOptionsBuilder) WithWhiteMs(ms uint) *SearchOptionsBuilder {
	sob.searchOptions.WhiteMs = ms
	return sob
}

func (sob *SearchOptionsBuilder) WithBlackMs(ms uint) *SearchOptionsBuilder {
	sob.searchOptions.BlackMs = ms
	return sob
}

func (sob *SearchOptionsBuilder) WithWhiteIncrMs(incrMs uint) *SearchOptionsBuilder {
	sob.searchOptions.WhiteIncrMs = incrMs
	return sob
}

func (sob *SearchOptionsBuilder) WithBlackIncrMs(incrMs uint) *SearchOptionsBuilder {
	sob.searchOptions.BlackIncrMs = incrMs
	return sob
}

func (sob *SearchOptionsBuilder) WithMovesTillIncr(movesTillIncr uint) *SearchOptionsBuilder {
	sob.searchOptions.MovesTillIncr = movesTillIncr
	return sob
}

func (sob *SearchOptionsBuilder) WithDepth(depth uint) *SearchOptionsBuilder {
	if sob.searchOptions.SearchMs != 0 {
		fmt.Println("WARNING: cannot set both Depth and SearchMs, unsetting SearchMs")
		sob.searchOptions.SearchMs = 0
	}
	sob.searchOptions.Depth = depth
	return sob
}

func (sob *SearchOptionsBuilder) WithSearchMs(ms uint) *SearchOptionsBuilder {
	if sob.searchOptions.Depth != 0 {
		fmt.Println("WARNING: cannot set both Depth and SearchMs, unsetting Depth")
		sob.searchOptions.Depth = 0
	}
	sob.searchOptions.SearchMs = ms
	return sob
}

func (sob *SearchOptionsBuilder) Build() *SearchOptions {
	return sob.searchOptions
}

func NewSearchOptionsBuilder() *SearchOptionsBuilder {
	return &SearchOptionsBuilder{
		searchOptions: &SearchOptions{},
	}
}
