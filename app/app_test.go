package app

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	db "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/simapp"

	abci "github.com/tendermint/tendermint/abci/types"
)

func Test3zdExport(t *testing.T) {
	db := db.NewMemDB()
	tczapp := NewThreeChainzApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")
	err := setGenesis(tczapp)
	require.NoError(t, err)

	// Making a new app object with the db, so that initchain hasn't been called
	newTczapp := NewThreeChainzApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")
	_, _, _, err = newTczapp.ExportAppStateAndValidators(false, []string{})
	require.NoError(t, err, "ExportAppStateAndValidators should not have an error")
}

// ensure that black listed addresses are properly set in bank keeper
func TestBlackListedAddrs(t *testing.T) {
	db := db.NewMemDB()
	tczapp := NewThreeChainzApp(log.NewTMLogger(log.NewSyncWriter(os.Stdout)), db, nil, true, 0, map[int64]bool{}, "")

	for acc := range maccPerms {
		require.Equal(t, !allowedReceivingModAcc[acc], tczapp.bankKeeper.BlacklistedAddr(tczapp.accountKeeper.GetModuleAddress(acc)))
	}
}

func TestGetMaccPerms(t *testing.T) {
	dup := GetMaccPerms()
	require.Equal(t, maccPerms, dup, "duplicated module account permissions differed from actual module account permissions")
}

func setGenesis(tczapp *ThreeChainzApp) error {
	genesisState := simapp.NewDefaultGenesisState()
	stateBytes, err := codec.MarshalJSONIndent(tczapp.Codec(), genesisState)
	if err != nil {
		return err
	}

	// Initialize the chain
	tczapp.InitChain(
		abci.RequestInitChain{
			Validators:    []abci.ValidatorUpdate{},
			AppStateBytes: stateBytes,
		},
	)

	tczapp.Commit()
	return nil
}
