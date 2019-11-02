package lsp

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
)

func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgOpenLiquidStakingPosition:
			return handleMsgOpenLiquidStakingPosition(ctx, k, msg)
		default:
			return sdk.ErrUnknownRequest("unknown message").Result()
		}
	}
}

func handleMsgOpenLiquidStakingPosition(ctx sdk.Context, k keeper.Keeper, msg types.MsgOpenLiquidStakingPosition) sdk.Result {
	err := k.OpenLiquidStakingPosition(ctx, msg.TransferChanId, msg.InterchainAccountChanId, msg.Amount, msg.Sender, msg.Validator)
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
