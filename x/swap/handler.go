package swap

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/everett-protocol/everett-hackathon/x/swap/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/swap/internal/types"
)

func NewHandler(k keeper.Keeper, bk bank.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		case types.MsgSwap:
			return handleMsgSwap(ctx, k, bk, msg)
		default:
			return sdk.ErrUnknownRequest("unknown message").Result()
		}
	}
}

func handleMsgSwap(ctx sdk.Context, keeper Keeper, bk bank.Keeper, msg types.MsgSwap) sdk.Result {
	_, err := bk.SubtractCoins(ctx, msg.Sender, sdk.Coins{msg.Asset})
	if err != nil {
		return err.Result()
	}

	result, err := keeper.Swap(ctx, msg.Asset, msg.TargetDenom)
	if err != nil {
		return err.Result()
	}
	_, err = bk.AddCoins(ctx, msg.Sender, sdk.Coins{result})
	if err != nil {
		return err.Result()
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
