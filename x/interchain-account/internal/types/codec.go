package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

func registerInterface(cdc *codec.Codec) {
	cdc.RegisterInterface((*ibc.Packet)(nil), nil)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(RegisterIBCAccountPacketData{}, "ibc/RegisterInterchainAccount", nil)
	cdc.RegisterConcrete(RunTxPacketData{}, "ibc/RunInterchainAccountTx", nil)
}

var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	registerInterface(ModuleCdc)
	RegisterCodec(ModuleCdc)
}
