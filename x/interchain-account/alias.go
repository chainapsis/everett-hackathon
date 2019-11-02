package interchain_account

import (
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/types"
)

const (
	ModuleName = types.ModuleName
)

var (
	NewKeeper = keeper.NewKeeper
)

type (
	Keeper                          = keeper.Keeper
	PacketRunInterchainAccountTx    = types.PacketRunInterchainAccountTx
	PacketRegisterInterchainAccount = types.PacketRegisterInterchainAccount
	InterchainAccountTx             = types.InterchainAccountTx
)
