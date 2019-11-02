package interchain_account

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case ibc.MsgPacket:
			switch packet := msg.Packet.(type) {
			case types.PacketRegisterInterchainAccount:
				return handleRegisterInterchainAccount(ctx, k, msg.ChannelID, packet)
			case types.PacketRunInterchainAccountTx:
				return handleRunInterchainAccountTx(ctx, k, msg.ChannelID, packet)
			default:
				return sdk.ErrUnknownRequest("unknown packet").Result()
			}
		default:
			return sdk.ErrUnknownRequest("unknown message").Result()
		}
	}
}

func handleRegisterInterchainAccount(ctx sdk.Context, k keeper.Keeper, channelId string, packet types.PacketRegisterInterchainAccount) sdk.Result {
	err := k.RegisterInterchainAccount(ctx, channelId, packet.Salt)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}

func handleRunInterchainAccountTx(ctx sdk.Context, k keeper.Keeper, channelId string, packet types.PacketRunInterchainAccountTx) sdk.Result {
	interchainAccountTx, err := k.UnmarshalTx(ctx, packet.TxBytes)
	if err != nil {
		return sdk.ErrInternal(err.Error()).Result()
	}
	return k.RunTx(ctx, channelId, interchainAccountTx)
}
