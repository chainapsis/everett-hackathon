package interchain_account

import (
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/types"
)

const (
	ModuleName         = types.ModuleName
	CosmosSdkChainType = types.CosmosSdkChainType
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper                       = keeper.Keeper
	RunTxPacketData              = types.RunTxPacketData
	RegisterIBCAccountPacketData = types.RegisterIBCAccountPacketData
	InterchainAccountTx          = types.InterchainAccountTx
)
