package ibc_transfer

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case ibc.MsgPacket:
			switch packet := msg.Packet.(type) {
			case types.PacketTransfer:
				return handlePacketTransfer(ctx, k, packet)
			default:
				return sdk.ErrUnknownRequest("unknown packet").Result()
			}
		case types.MsgTransfer:
			return handleMsgTransfer(ctx, k, msg)
		default:
			return sdk.ErrUnknownRequest("unknown message").Result()
		}
	}
}

func handlePacketTransfer(ctx sdk.Context, k keeper.Keeper, packet types.PacketTransfer) sdk.Result {
	err := k.ReceiveTransferPacket(ctx, packet)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleMsgTransfer(ctx sdk.Context, k keeper.Keeper, msg types.MsgTransfer) sdk.Result {
	err := k.SendTransfer(ctx, msg.ChanId, msg.Amount, msg.Sender, msg.Receiver)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
