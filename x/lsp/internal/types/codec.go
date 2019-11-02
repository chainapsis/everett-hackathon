package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgOpenLiquidStakingPosition{}, "lsp/open", nil)
	cdc.RegisterConcrete(MsgCloseLiquidStakingPosition{}, "lsp/close", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
}
