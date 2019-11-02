package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	ibctransfer "github.com/everett-protocol/everett-hackathon/x/ibc-transfer"
	interchainaccount "github.com/everett-protocol/everett-hackathon/x/interchain-account"

	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
)

type Keeper struct {
	cdc            *codec.Codec
	key            sdk.StoreKey
	supplyKeeper   supply.Keeper
	transferKeeper ibctransfer.Keeper
	iaKeeper       interchainaccount.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, supplyKeeper supply.Keeper, transferKeeper ibctransfer.Keeper, iaKeeper interchainaccount.Keeper) Keeper {
	return Keeper{
		cdc:            cdc,
		key:            key,
		supplyKeeper:   supplyKeeper,
		transferKeeper: transferKeeper,
		iaKeeper:       iaKeeper,
	}
}

func (k Keeper) OpenLiquidStakingPosition(ctx sdk.Context, transferChanId string, iaChanId string, amount sdk.Coin, sender sdk.AccAddress, validator sdk.ValAddress) sdk.Error {
	salt := k.getSalt(ctx)
	registerPacket := interchainaccount.PacketRegisterInterchainAccount{
		Salt: salt,
	}

	// XXX: Currently, there is no way to get packet's result (acknowledgement packet),
	// So, it is impossible to react according to the prior packet's result.
	// So, just send packets at once. Therefore it can't handle the case that packet failed.
	// And, order of relay is very important.

	err := k.iaKeeper.SendPacket(ctx, iaChanId, registerPacket)
	if err != nil {
		return err
	}

	// Predict the address will be made
	address, goErr := k.iaKeeper.CalcAddress("", salt)
	if goErr != nil {
		return sdk.ErrInternal(goErr.Error())
	}

	// TODO: Check that coin will be sent is a staking coin.
	err = k.transferKeeper.SendTransfer(ctx, transferChanId, sdk.Coins{amount}, sender, address)
	if err != nil {
		return err
	}

	delegateMsg := staking.NewMsgDelegate(address, validator, amount)
	err = k.iaKeeper.SendMsgs(ctx, iaChanId, []sdk.Msg{delegateMsg})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) getSalt(ctx sdk.Context) string {
	store := ctx.KVStore(k.key)
	key := []byte("salt")
	lastSalt := 0
	saltBytes := store.Get(key)
	if len(saltBytes) > 0 {
		types.ModuleCdc.MustUnmarshalBinaryBare(saltBytes, &lastSalt)
	}
	lastSalt++
	store.Set([]byte("salt"), types.ModuleCdc.MustMarshalBinaryBare(lastSalt))
	return string(types.ModuleCdc.MustMarshalJSON(lastSalt))
}
