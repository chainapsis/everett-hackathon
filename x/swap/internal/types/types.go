package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	ModuleName = "swap"
	RouterKey  = ModuleName
)

type Pool struct {
	BalanceCoin  sdk.Coin `json:"balance_coin"` // intermediation for exchange tokens
	BalanceToken sdk.Coin `json:"balance_token"`
}

func NewPool(balanceCoin sdk.Coin, balanceToken sdk.Coin) Pool {
	return Pool{
		BalanceCoin:  balanceCoin,
		BalanceToken: balanceToken,
	}
}
