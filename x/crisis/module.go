package crisis

import (
	"encoding/json"
	"math/rand"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/cosmos/cosmos-sdk/x/crisis/client/cli"
	"github.com/cosmos/cosmos-sdk/x/crisis/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/crisis/internal/types"
	sim "github.com/cosmos/cosmos-sdk/x/simulation"
)

var (
	_ module.AppModule           = AppModule{}
	_ module.AppModuleBasic      = AppModuleBasic{}
	_ module.AppModuleSimulation = AppModuleSimulation{}
)

// AppModuleBasic defines the basic application module used by the crisis module.
type AppModuleBasic struct{}

// Name returns the crisis module's name.
func (AppModuleBasic) Name() string {
	return ModuleName
}

// RegisterCodec registers the crisis module's types for the given codec.
func (AppModuleBasic) RegisterCodec(cdc *codec.Codec) {
	RegisterCodec(cdc)
}

// DefaultGenesis returns default genesis state as raw bytes for the crisis
// module.
func (AppModuleBasic) DefaultGenesis() json.RawMessage {
	return types.ModuleCdc.MustMarshalJSON(DefaultGenesisState())
}

// ValidateGenesis performs genesis state validation for the crisis module.
func (AppModuleBasic) ValidateGenesis(bz json.RawMessage) error {
	var data types.GenesisState
	if err := types.ModuleCdc.UnmarshalJSON(bz, &data); err != nil {
		return err
	}
	return types.ValidateGenesis(data)
}

// RegisterRESTRoutes registers no REST routes for the crisis module.
func (AppModuleBasic) RegisterRESTRoutes(_ context.CLIContext, _ *mux.Router) {}

// GetTxCmd returns the root tx command for the crisis module.
func (AppModuleBasic) GetTxCmd(cdc *codec.Codec) *cobra.Command {
	return cli.GetTxCmd(cdc)
}

// GetQueryCmd returns no root query command for the crisis module.
func (AppModuleBasic) GetQueryCmd(_ *codec.Codec) *cobra.Command { return nil }

//____________________________________________________________________________

// AppModuleSimulation defines the module simulation functions used by the crisis module.
type AppModuleSimulation struct{}

// RegisterStoreDecoder performs a no-op.
func (AppModuleSimulation) RegisterStoreDecoder(_ sdk.StoreDecoderRegistry) {}

// GenerateGenesisState performs a no-op.
func (AppModuleSimulation) GenerateGenesisState(_ *module.GeneratorInput) {
}

// RandomizedParams doesn't create any randomized crisis param changes for the simulator.
func (AppModuleSimulation) RandomizedParams( _ *rand.Rand) []sim.ParamChange {
	return nil
}

//____________________________________________________________________________

// AppModule implements an application module for the crisis module.
type AppModule struct {
	AppModuleBasic
	AppModuleSimulation

	// NOTE: We store a reference to the keeper here so that after a module
	// manager is created, the invariants can be properly registered and
	// executed.
	keeper *keeper.Keeper
}

// NewAppModule creates a new AppModule object
func NewAppModule(keeper *keeper.Keeper) AppModule {
	return AppModule{
		AppModuleBasic:      AppModuleBasic{},
		AppModuleSimulation: AppModuleSimulation{},
		keeper:              keeper,
	}
}

// Name returns the crisis module's name.
func (AppModule) Name() string {
	return ModuleName
}

// RegisterInvariants performs a no-op.
func (AppModule) RegisterInvariants(_ sdk.InvariantRegistry) {}

// Route returns the message routing key for the crisis module.
func (AppModule) Route() string {
	return RouterKey
}

// NewHandler returns an sdk.Handler for the crisis module.
func (am AppModule) NewHandler() sdk.Handler {
	return NewHandler(*am.keeper)
}

// QuerierRoute returns no querier route.
func (AppModule) QuerierRoute() string { return "" }

// NewQuerierHandler returns no sdk.Querier.
func (AppModule) NewQuerierHandler() sdk.Querier { return nil }

// InitGenesis performs genesis initialization for the crisis module. It returns
// no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, data json.RawMessage) []abci.ValidatorUpdate {
	var genesisState GenesisState
	types.ModuleCdc.MustUnmarshalJSON(data, &genesisState)
	InitGenesis(ctx, *am.keeper, genesisState)

	am.keeper.AssertInvariants(ctx)
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the crisis
// module.
func (am AppModule) ExportGenesis(ctx sdk.Context) json.RawMessage {
	gs := ExportGenesis(ctx, *am.keeper)
	return types.ModuleCdc.MustMarshalJSON(gs)
}

// BeginBlock performs a no-op.
func (AppModule) BeginBlock(_ sdk.Context, _ abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the crisis module. It returns no validator
// updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, *am.keeper)
	return []abci.ValidatorUpdate{}
}
