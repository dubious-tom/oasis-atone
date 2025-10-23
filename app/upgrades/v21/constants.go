package v21

import (
	"github.com/osmosis-labs/osmosis/v31/app/upgrades"

	store "cosmossdk.io/store/types"
	consensustypes "github.com/cosmos/cosmos-sdk/x/consensus/types"
	// TODO ATOMONE: Crisis module removed in SDK v0.50+ - need to refactor
	// crisistypes "cosmossdk.io/x/crisis/types"
)

// UpgradeName defines the on-chain upgrade name for the Osmosis v21 upgrade.
const (
	UpgradeName    = "v21"
	TestingChainId = "testing-chain-id"
)

var Upgrade = upgrades.Upgrade{
	UpgradeName:          UpgradeName,
	CreateUpgradeHandler: CreateUpgradeHandler,
	StoreUpgrades: store.StoreUpgrades{
		Added: []string{
			// v47 modules
			// TODO ATOMONE: Crisis module removed in SDK v0.50+ - need to refactor
			// crisistypes.ModuleName,
			consensustypes.ModuleName,
		},
		Deleted: []string{},
	},
}
