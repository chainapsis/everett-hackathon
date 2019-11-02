package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type MsgOpenLiquidStakingPosition struct {
	TransferChanId          string         `json:"transfer_chan_id" yaml:"transfer_chan_id"`
	InterchainAccountChanId string         `json:"interchain_account_chan_id" yaml:"interchain_account_chan_id"`
	Amount                  sdk.Coin       `json:"amount" yaml:"amount"`
	Validator               sdk.ValAddress `json:"validator" yaml:"validator"`
	Sender                  sdk.AccAddress `json:"sender" yaml:"sender"`
}

var _ sdk.Msg = MsgOpenLiquidStakingPosition{}

func NewMsgOpenLiquidStakingPosition(transferChanId string, interchainAccountChanId string, amount sdk.Coin, validator sdk.ValAddress, sender sdk.AccAddress) MsgOpenLiquidStakingPosition {
	return MsgOpenLiquidStakingPosition{
		TransferChanId:          transferChanId,
		InterchainAccountChanId: interchainAccountChanId,
		Amount:                  amount,
		Validator:               validator,
		Sender:                  sender,
	}
}

func (msg MsgOpenLiquidStakingPosition) Route() string { return RouterKey }

func (msg MsgOpenLiquidStakingPosition) Type() string { return "open-liquid-staking-position" }

func (msg MsgOpenLiquidStakingPosition) ValidateBasic() sdk.Error {
	if len(msg.TransferChanId) == 0 {
		return sdk.ErrInternal("empty chan id")
	}
	if len(msg.InterchainAccountChanId) == 0 {
		return sdk.ErrInternal("empty chan id")
	}
	if !msg.Amount.IsValid() {
		return sdk.ErrInvalidCoins("send amount is invalid: " + msg.Amount.String())
	}
	if msg.Validator.Empty() {
		return sdk.ErrInvalidAddress("missing validator address")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	return nil
}

func (msg MsgOpenLiquidStakingPosition) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgOpenLiquidStakingPosition) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

type MsgCloseLiquidStakingPosition struct {
	InterchainAccountChanId string         `json:"interchain_account_chan_id" yaml:"interchain_account_chan_id"`
	NftId                   string         `json:"nft_id" yaml:"nft_id"`
	Sender                  sdk.AccAddress `json:"sender" yaml:"sender"`
	Recipient               sdk.AccAddress `json:"recipient" yaml:"nft_id"`
}

var _ sdk.Msg = MsgCloseLiquidStakingPosition{}

func NewMsgCloseLiquidStakingPosition(interchainAccountChanId string, nftId string, sender sdk.AccAddress, recipient sdk.AccAddress) MsgCloseLiquidStakingPosition {
	return MsgCloseLiquidStakingPosition{
		InterchainAccountChanId: interchainAccountChanId,
		NftId:                   nftId,
		Sender:                  sender,
		Recipient:               recipient,
	}
}

func (msg MsgCloseLiquidStakingPosition) Route() string { return RouterKey }

func (msg MsgCloseLiquidStakingPosition) Type() string { return "close-liquid-staking-position" }

func (msg MsgCloseLiquidStakingPosition) ValidateBasic() sdk.Error {
	if len(msg.InterchainAccountChanId) == 0 {
		return sdk.ErrInternal("empty chan id")
	}
	if len(msg.NftId) == 0 {
		return sdk.ErrInternal("empty nft id")
	}
	if msg.Sender.Empty() {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if msg.Recipient.Empty() {
		return sdk.ErrInvalidAddress("missing recipient address")
	}
	return nil
}

func (msg MsgCloseLiquidStakingPosition) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCloseLiquidStakingPosition) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
