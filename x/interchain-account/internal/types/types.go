package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type InterchainAccountTx struct {
	Msgs []sdk.Msg
}
