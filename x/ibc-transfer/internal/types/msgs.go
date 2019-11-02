package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

/*
 Not consider the case that multiple chains exist.
 Just assume that only one to one communication exists for prototype.
*/
type MsgTransfer struct {
	// SourcePort    string         `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	// SourceChannel string         `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	ChanId   string         `json:"chain_id" yaml:"chain_id"`
	Amount   sdk.Coins      `json:"amount" yaml:"amount"`     // the tokens to be transferred
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`     // the sender address
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"` // the recipient address on the destination chain
	// Source        bool           `json:"source" yaml:"source"`                 // indicates if the sending chain is the source chain of the tokens to be transferred
}

func NewMsgTransfer(chanId string, amount sdk.Coins, sender sdk.AccAddress, receiver sdk.AccAddress) MsgTransfer {
	return MsgTransfer{
		ChanId:   chanId,
		Amount:   amount,
		Sender:   sender,
		Receiver: receiver,
	}
}

func (msg MsgTransfer) Route() string { return RouterKey }

func (msg MsgTransfer) Type() string { return "ibc-transfer" }

func (msg MsgTransfer) ValidateBasic() sdk.Error {
	if len(msg.ChanId) == 0 {
		return sdk.ErrInternal("empty chan id")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.Receiver.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if !msg.Amount.IsAllPositive() {
		return sdk.ErrInsufficientCoins("send amount must be positive")
	}
	return nil
}

func (msg MsgTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
