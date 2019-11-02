package lsp

import (
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
