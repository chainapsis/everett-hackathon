package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
)

/*
 Not consider the case that multiple chains exist.
 Just assume that only one to one communication exists for prototype.
*/
type PacketTransfer struct {
	Amount   sdk.Coins      `json:"amount" yaml:"amount"`     // the tokens to be transferred
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"` // the recipient address on the destination chain
}

var _ ibc.Packet = PacketTransfer{}

func (packet PacketTransfer) SenderPort() string {
	return RouterKey
}

func (packet PacketTransfer) ReceiverPort() string {
	return RouterKey
}

func (packet PacketTransfer) Type() string {
	return "ibc-transfer"
}

func (packet PacketTransfer) ValidateBasic() sdk.Error {
	if packet.Receiver.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !packet.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + packet.Amount.String())
	}
	if !packet.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

func (packet PacketTransfer) Timeout() uint64 {
	return 0
}

func (packet PacketTransfer) Marshal() []byte {
	return ModuleCdc.MustMarshalBinaryBare(packet)
}
