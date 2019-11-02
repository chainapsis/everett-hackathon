package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

func registerInterface(cdc *codec.Codec) {
	cdc.RegisterInterface((*ibc.Packet)(nil), nil)
}

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgTransfer{}, "ibc/Transfer", nil)
	cdc.RegisterConcrete(PacketTransfer{}, "ibc/PacketTransfer", nil)
}

var ModuleCdc = codec.New()

func init() {
	registerInterface(ModuleCdc)
	RegisterCodec(ModuleCdc)
}
