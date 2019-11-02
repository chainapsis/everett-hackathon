package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

const RouterKey = ModuleName

type PacketRegisterInterchainAccount struct {
	Salt string `json:"salt"`
}

var _ ibc.Packet = PacketRegisterInterchainAccount{}

func (packet PacketRegisterInterchainAccount) SenderPort() string {
	return RouterKey
}

func (packet PacketRegisterInterchainAccount) ReceiverPort() string {
	return RouterKey
}

func (packet PacketRegisterInterchainAccount) Type() string {
	return "register-interchain-account"
}

func (packet PacketRegisterInterchainAccount) ValidateBasic() sdk.Error {
	if len(packet.Salt) == 0 {
		return sdk.ErrInternal("salt is empty")
	}
	return nil
}

func (packet PacketRegisterInterchainAccount) Timeout() uint64 {
	return 0
}

func (packet PacketRegisterInterchainAccount) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(packet)
}

type PacketRunInterchainAccountTx struct {
	TxBytes []byte `json:"tx_bytes"`
}

var _ ibc.Packet = PacketRunInterchainAccountTx{}

func (packet PacketRunInterchainAccountTx) SenderPort() string {
	return RouterKey
}

func (packet PacketRunInterchainAccountTx) ReceiverPort() string {
	return RouterKey
}

func (packet PacketRunInterchainAccountTx) Type() string {
	return "run-interchainaccount-tx"
}

func (packet PacketRunInterchainAccountTx) ValidateBasic() sdk.Error {
	if len(packet.TxBytes) == 0 {
		return sdk.ErrInternal("tx bytes is empty")
	}
	return nil
}

func (packet PacketRunInterchainAccountTx) Timeout() uint64 {
	return 0
}

func (packet PacketRunInterchainAccountTx) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(packet)
}
