package lsp

import (
	"encoding/json"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/everett-protocol/everett-hackathon/x/lsp/client/cli"
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/keeper"
	"github.com/everett-protocol/everett-hackathon/x/lsp/internal/types"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"

	abci "github.com/tendermint/tendermint/abci/types"
)

var (
	_ module.AppModuleBasic = AppModuleBasic{}
)

type AppModuleBasic struct{}

func (AppModuleBasic) Name() string { return types.ModuleName }

func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) { types.RegisterCodec(cdc) }

func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return nil
}

func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	return nil
}

func (AppModuleBasic) RegisterRESTRoutes(ctx context.CLIContext, rtr *mux.Router) {
	// noop
}

func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

func (AppModuleBasic) GetQueryCmd(cdc *codec.Codec) *cobra.Command {
	return nil
}

type AppModule struct {
	AppModuleBasic
	k keeper.Keeper
}

func NewAppModule(k keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		k:              k,
	}
}

func (AppModule) RegisterInvariants(ir sdk.InvariantRegistry) {

}

func (AppModule) Name() string {
	return types.ModuleName
}

func (AppModule) Route() string {
	return types.ModuleName
}

func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(am.k)
}

func (AppModule) QuerierRoute() string {
	return types.ModuleName
}

func (AppModule) NewQuerierHandler() sdk.Querier {
	return nil
}

func (AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

func (AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	return nil
}

func (AppModule) BeginBlock(ctx sdk.Context, req abci.RequestBeginBlock) {
	// noop
}

func (AppModule) EndBlock(ctx sdk.Context, req abci.RequestEndBlock) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
