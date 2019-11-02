package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/everett-protocol/everett-hackathon/x/ibc-transfer/internal/types"

	"github.com/cosmos/cosmos-sdk/x/ibc"
)

type Keeper struct {
	cdc          *codec.Codec
	key          sdk.StoreKey
	supplyKeeper supply.Keeper
	port         ibc.Port
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, port ibc.Port) Keeper {

	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the IBC transfer module account has not been set")
	}

	return Keeper{
		cdc:          cdc,
		key:          key,
		supplyKeeper: supplyKeeper,
		port:         port,
	}
}

func (k Keeper) SendTransfer(ctx sdk.Context, chainId string, amount sdk.Coins, sender, receiver sdk.AccAddress) sdk.Error {
	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, amount)
	if err != nil {
		return err
	}

	packet := types.PacketTransfer{Receiver: receiver, Amount: amount}
	return k.port.Send(ctx, chainId, packet)
}

func (k Keeper) ReceiveTransferPacket(ctx sdk.Context, packet types.PacketTransfer) sdk.Error {
	err := k.supplyKeeper.MintCoins(ctx, types.ModuleName, packet.Amount)
	if err != nil {
		return err
	}
	return k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, packet.Receiver, packet.Amount)
}
