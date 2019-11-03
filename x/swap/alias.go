package swap

import (
	"github.com/everett-protocol/everett-hackathon/x/swap/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/swap/internal/types"
)

const (
	ModuleName        = types.ModuleName
	DefaultParamSpace = types.DefaultParamspace
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper = keeper.Keeper
)
