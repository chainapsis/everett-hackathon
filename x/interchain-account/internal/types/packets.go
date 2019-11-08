package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

const RouterKey = ModuleName

type RegisterIBCAccountPacketData struct {
	Salt string `json:"salt"`
}

var _ ibc.Packet = RegisterIBCAccountPacketData{}

func (packet RegisterIBCAccountPacketData) SenderPort() string {
	return RouterKey
}

func (packet RegisterIBCAccountPacketData) ReceiverPort() string {
	return RouterKey
}

func (packet RegisterIBCAccountPacketData) Type() string {
	return "register-interchain-account"
}

func (packet RegisterIBCAccountPacketData) ValidateBasic() sdk.Error {
	if len(packet.Salt) == 0 {
		return sdk.ErrInternal("salt is empty")
	}
	return nil
}

func (packet RegisterIBCAccountPacketData) Timeout() uint64 {
	return 0
}

func (packet RegisterIBCAccountPacketData) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(packet)
}

type RunTxPacketData struct {
	TxBytes []byte `json:"tx_bytes"`
}

var _ ibc.Packet = RunTxPacketData{}

func (packet RunTxPacketData) SenderPort() string {
	return RouterKey
}

func (packet RunTxPacketData) ReceiverPort() string {
	return RouterKey
}

func (packet RunTxPacketData) Type() string {
	return "run-interchainaccount-tx"
}

func (packet RunTxPacketData) ValidateBasic() sdk.Error {
	if len(packet.TxBytes) == 0 {
		return sdk.ErrInternal("tx bytes is empty")
	}
	return nil
}

func (packet RunTxPacketData) Timeout() uint64 {
	return 0
}

func (packet RunTxPacketData) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(packet)
}
