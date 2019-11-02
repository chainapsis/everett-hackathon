package ibc_transfer

import (
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/types"
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
