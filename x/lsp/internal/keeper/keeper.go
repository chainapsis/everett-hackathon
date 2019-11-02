package keeper

import (
	"errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	ibctransfer "github.com/everett-protocol/everett-hackathon/x/ibc-transfer"
	interchainaccount "github.com/everett-protocol/everett-hackathon/x/interchain-account"
	"strconv"

	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
)

type Keeper struct {
	cdc            *codec.Codec
	key            sdk.StoreKey
	accountKeeper  auth.AccountKeeper
	supplyKeeper   supply.Keeper
	nftKeeper      nft.Keeper
	transferKeeper ibctransfer.Keeper
	iaKeeper       interchainaccount.Keeper
}

func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, accountKeeper auth.AccountKeeper, supplyKeeper supply.Keeper, nftKeeper nft.Keeper, transferKeeper ibctransfer.Keeper, iaKeeper interchainaccount.Keeper) Keeper {

	if addr := supplyKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic("the Liquid staking position module account has not been set")
	}

	return Keeper{
		cdc:            cdc,
		key:            key,
		accountKeeper:  accountKeeper,
		supplyKeeper:   supplyKeeper,
		nftKeeper:      nftKeeper,
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

	collateral := sdk.NewCoin(amount.Denom, amount.Amount.Quo(sdk.NewInt(10)))
	minted := sdk.NewCoin("b"+amount.Denom, amount.Amount.Sub(collateral.Amount))

	err = k.supplyKeeper.MintCoins(ctx, types.ModuleName, sdk.Coins{minted})
	if err != nil {
		return err
	}
	err = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{minted})
	if err != nil {
		return err
	}

	lspInfo := types.LspInfo{
		Validator:       validator,
		ChainAddress:    address,
		BondedShare:     amount,
		Collateral:      collateral,
		MintedLiquidity: minted,
	}
	goErr = k.setLspInfo(ctx, salt, lspInfo)
	if goErr != nil {
		return sdk.ErrInternal(goErr.Error())
	}

	// Use token uri to save the information directly.
	lspNft := nft.NewBaseNFT(salt, sender, string(types.ModuleCdc.MustMarshalJSON(lspInfo)))
	err = k.nftKeeper.MintNFT(ctx, "lsp", &lspNft)
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) CloseLiquidStakingPosition(ctx sdk.Context, iaChanId string, nftId string, sender sdk.AccAddress, recipient sdk.AccAddress) sdk.Error {
	lspInfo, goErr := k.getLspInfo(ctx, nftId)
	if goErr != nil {
		return sdk.ErrInternal(goErr.Error())
	}

	account := k.accountKeeper.GetAccount(ctx, sender)
	if account == nil {
		return sdk.ErrUnknownAddress("unknown sender")
	}

	liquidToken := account.GetCoins().AmountOf(lspInfo.MintedLiquidity.Denom)
	if liquidToken.LT(lspInfo.MintedLiquidity.Amount) {
		return sdk.ErrUnauthorized("you don't have enough minted bonded token")
	}

	err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{lspInfo.MintedLiquidity})
	if err != nil {
		return err
	}

	err = k.supplyKeeper.BurnCoins(ctx, types.ModuleName, sdk.Coins{lspInfo.MintedLiquidity})
	if err != nil {
		return err
	}

	err = k.nftKeeper.DeleteNFT(ctx, "lsp", nftId)
	if err != nil {
		return err
	}

	setRecipientMsg := distribution.NewMsgSetWithdrawAddress(lspInfo.ChainAddress, recipient)
	withdrawMsg := distribution.NewMsgWithdrawDelegatorReward(lspInfo.ChainAddress, lspInfo.Validator)
	unbondingMsg := staking.NewMsgUndelegate(lspInfo.ChainAddress, lspInfo.Validator, lspInfo.BondedShare)

	err = k.iaKeeper.SendMsgs(ctx, iaChanId, []sdk.Msg{unbondingMsg, setRecipientMsg, withdrawMsg})
	if err != nil {
		return err
	}

	return nil
}

func (k Keeper) getSalt(ctx sdk.Context) string {
	store := ctx.KVStore(k.key)
	key := []byte("salt")
	lastSalt := uint64(0)
	saltBytes := store.Get(key)
	if len(saltBytes) > 0 {
		types.ModuleCdc.MustUnmarshalBinaryBare(saltBytes, &lastSalt)
	}
	lastSalt++
	store.Set([]byte("salt"), types.ModuleCdc.MustMarshalBinaryBare(lastSalt))
	return strconv.FormatUint(lastSalt, 10)
}

func (k Keeper) setLspInfo(ctx sdk.Context, key string, lspInfo types.LspInfo) error {
	store := ctx.KVStore(k.key)
	bz, err := types.ModuleCdc.MarshalBinaryBare(lspInfo)
	if err != nil {
		return err
	}
	store.Set([]byte(key), bz)

	return nil
}

func (k Keeper) getLspInfo(ctx sdk.Context, key string) (types.LspInfo, error) {
	store := ctx.KVStore(k.key)
	bz := store.Get([]byte(key))
	if len(bz) == 0 {
		return types.LspInfo{}, errors.New("lsp info not exist")
	}

	lspInfo := types.LspInfo{}
	err := types.ModuleCdc.UnmarshalBinaryBare(bz, &lspInfo)

	return lspInfo, err
}
