package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type LspInfo struct {
	Validator       sdk.ValAddress `json:"validator"`
	ChainAddress    sdk.AccAddress `json:"chain_address"`
	BondedShare     sdk.Coin       `json:"bonded_share"`
	Collateral      sdk.Coin       `json:"collateral"`
	MintedLiquidity sdk.Coin       `json:"minted_liquidity"`
}
