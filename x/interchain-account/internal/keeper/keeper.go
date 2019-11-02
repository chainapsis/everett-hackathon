package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/everett-protocol/everett-hackathon/x/interchain-account/internal/types"
	"strings"

	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/ibc"

	"github.com/tendermint/tendermint/crypto/tmhash"
)

type Keeper struct {
	cdc             *codec.Codec
	counterpartyCdc *codec.Codec
	key             sdk.StoreKey
	router          sdk.Router
	ak              auth.AccountKeeper
	port            ibc.Port
}

func NewKeeper(cdc *codec.Codec, counterpartyCdc *codec.Codec, key sdk.StoreKey, router sdk.Router, ak auth.AccountKeeper, port ibc.Port) Keeper {
	return Keeper{
		cdc:             cdc,
		counterpartyCdc: counterpartyCdc,
		key:             key,
		router:          router,
		ak:              ak,
		port:            port,
	}
}

func (k Keeper) RegisterInterchainAccount(ctx sdk.Context, channelId string, salt string) sdk.Error {
	/*
		path = "{packet.sourcePort}/{packet.sourceChannel}"
		address = sha256(path + packet.salt)

		// Should not block even if there is normal account,
		// because attackers can distrupt to create an ibc managed account
		// by sending some assets to estimated address in advance.
		// And IBC managed account has no public key, but its sequence is 1.
		// It can be mark for Interchain account, becuase normal account can't be sequence 1 without publish public key.
		account = accountKeeper.getAccount(account)
		if (account != null) {
		  abortTransactionUnless(account.sequence === 1 && account.pubKey == null)
		} else {
		  accountKeeper.newAccount(address)
		}

		addressesRegisteredChannel[signer] = path

		// set account's sequence to 1
		accountKeeper.setAccount(address, {
		  ...,
		  sequence: 1,
		  ...,
		})
	*/

	// Currently, it seems that there is no way to get the information of counterparty chain.
	// So, just don't use path for hackathon version.
	address, err := k.CalcAddress("", salt)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}

	account := k.ak.GetAccount(ctx, address)
	if account != nil {
		if account.GetSequence() != 1 || account.GetPubKey() != nil {
			return types.ErrAccountAlreadyExist(types.DefaultCodespace)
		}
	} else {
		account := k.ak.NewAccountWithAddress(ctx, address)
		err := account.SetSequence(1)
		if err != nil {
			return sdk.ErrInternal(err.Error())
		}
		k.ak.NewAccount(ctx, account)
	}

	store := ctx.KVStore(k.key)
	// Ignore that which chain makes the interchain account.
	// Assume that only one to one communication exists for prototype version.
	store.Set(address, []byte{1})

	return nil
}

// Determine account's address that will be created.
func (k Keeper) CalcAddress(path string, salt string) ([]byte, error) {
	hash := tmhash.NewTruncated()
	hashsum := hash.Sum([]byte(path + salt))
	return hashsum, nil
}

func (k Keeper) SendMsgs(ctx sdk.Context, chanId string, msgs []sdk.Msg) sdk.Error {
	interchainAccountTx := types.InterchainAccountTx{Msgs: msgs}
	txBytes, err := k.counterpartyCdc.MarshalBinaryLengthPrefixed(interchainAccountTx)
	if err != nil {
		return sdk.ErrInternal(err.Error())
	}
	packet := types.PacketRunInterchainAccountTx{TxBytes: txBytes}

	return k.port.Send(ctx, chanId, packet)
}

func (k Keeper) UnmarshalTx(ctx sdk.Context, txBytes []byte) (types.InterchainAccountTx, error) {
	tx := types.InterchainAccountTx{}
	err := k.counterpartyCdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
	return tx, err
}

func (k Keeper) RunTx(ctx sdk.Context, channelId string, tx types.InterchainAccountTx) sdk.Result {
	err := k.AuthenticateTx(ctx, channelId, tx)
	if err != nil {
		return err.Result()
	}

	msgs := tx.Msgs

	logs := make([]string, 0, len(msgs))
	data := make([]byte, 0, len(msgs))
	var (
		code      sdk.CodeType
		codespace sdk.CodespaceType
	)
	events := ctx.EventManager().Events()

	for _, msg := range msgs {
		result := k.RunMsg(ctx, msg)
		if result.IsOK() == false {
			return result
		}

		data = append(data, result.Data...)

		events = events.AppendEvents(result.Events)

		if len(result.Log) > 0 {
			logs = append(logs, result.Log)
		}

		if !result.IsOK() {
			code = result.Code
			codespace = result.Codespace
			break
		}
	}

	return sdk.Result{
		Code:      code,
		Codespace: codespace,
		Data:      data,
		Log:       strings.TrimSpace(strings.Join(logs, ",")),
		GasUsed:   ctx.GasMeter().GasConsumed(),
		Events:    events,
	}
}

func (k Keeper) AuthenticateTx(ctx sdk.Context, channelId string, tx types.InterchainAccountTx) sdk.Error {
	msgs := tx.Msgs

	seen := map[string]bool{}
	var signers []sdk.AccAddress
	for _, msg := range msgs {
		for _, addr := range msg.GetSigners() {
			if !seen[addr.String()] {
				signers = append(signers, addr)
				seen[addr.String()] = true
			}
		}
	}

	store := ctx.KVStore(k.key)

	for _, signer := range signers {
		path := store.Get(signer)
		// Ignore that which chain makes the interchain account.
		// Assume that only one to one communication exists for prototype version.
		if len(path) == 0 {
			return sdk.ErrUnauthorized("unauthorized")
		}
	}

	return nil
}

func (k Keeper) RunMsg(ctx sdk.Context, msg sdk.Msg) sdk.Result {
	hander := k.router.Route(msg.Route())
	if hander == nil {
		return sdk.ErrInternal("invalid route").Result()
	}

	return hander(ctx, msg)
}

func (k Keeper) SendPacket(ctx sdk.Context, chanId string, packet ibc.Packet) sdk.Error {
	return k.port.Send(ctx, chanId, packet)
}
