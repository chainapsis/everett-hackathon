package swap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/everett-protocol/everett-hackathon/x/swap/internal/types"
)

// GenesisState is the bank state that must be provided at genesis.
type GenesisState struct {
	CoinDenom string       `json:"coin_denom"`
	Pools     []types.Pool `json:"pools"`
}

// NewGenesisState creates a new genesis state.
func NewGenesisState(coinDenom string, pools []types.Pool) GenesisState {
	return GenesisState{
		CoinDenom: coinDenom,
		Pools:     pools,
	}
}

// DefaultGenesisState returns a default genesis state
func DefaultGenesisState() GenesisState {
	return NewGenesisState(
		"uevrt", []types.Pool{
			types.Pool{
				BalanceCoin:  sdk.NewCoin("uevrt", sdk.NewInt(10000000)),
				BalanceToken: sdk.NewCoin("buatom", sdk.NewInt(10000000)),
			},
			types.Pool{
				BalanceCoin:  sdk.NewCoin("uevrt", sdk.NewInt(10000000)),
				BalanceToken: sdk.NewCoin("uatom", sdk.NewInt(10000000)),
			},
		},
	)
}

// InitGenesis sets distribution information for genesis.
func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) {
	config := types.PoolConfig{
		CoinDenom: data.CoinDenom,
	}

	keeper.SetPoolConfig(ctx, config)

	for _, pool := range data.Pools {
		err := keeper.SetPool(ctx, pool)
		if err != nil {
			panic(err.Error())
		}
	}
}
