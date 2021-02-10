package cli_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"
	tmcli "github.com/tendermint/tendermint/libs/cli"

	"github.com/c-osmosis/osmosis/app"
	"github.com/c-osmosis/osmosis/x/lockup/client/cli"
	lockuptestutil "github.com/c-osmosis/osmosis/x/lockup/client/testutil"
	"github.com/c-osmosis/osmosis/x/lockup/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	"github.com/cosmos/cosmos-sdk/simapp"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	"github.com/cosmos/cosmos-sdk/testutil/network"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktestutil "github.com/cosmos/cosmos-sdk/x/bank/client/testutil"
	dbm "github.com/tendermint/tm-db"
)

type IntegrationTestSuite struct {
	suite.Suite

	cfg     network.Config
	network *network.Network
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.T().Log("setting up integration test suite")

	encCfg := app.MakeEncodingConfig()

	s.cfg = network.Config{
		Codec:             encCfg.Marshaler,
		TxConfig:          encCfg.TxConfig,
		LegacyAmino:       encCfg.Amino,
		InterfaceRegistry: encCfg.InterfaceRegistry,
		AccountRetriever:  authtypes.AccountRetriever{},
		AppConstructor: func(val network.Validator) servertypes.Application {
			return app.NewOsmosisApp(
				val.Ctx.Logger, dbm.NewMemDB(), nil, true, make(map[int64]bool), val.Ctx.Config.RootDir, 0,
				encCfg,
				simapp.EmptyAppOptions{},
				baseapp.SetMinGasPrices(val.AppConfig.MinGasPrices),
			)
		},
		GenesisState:    app.ModuleBasics.DefaultGenesis(encCfg.Marshaler),
		TimeoutCommit:   2 * time.Second,
		ChainID:         "osmosis-1",
		NumValidators:   1,
		BondDenom:       sdk.DefaultBondDenom,
		MinGasPrices:    fmt.Sprintf("0.000006%s", sdk.DefaultBondDenom),
		AccountTokens:   sdk.TokensFromConsensusPower(1000),
		StakingTokens:   sdk.TokensFromConsensusPower(500),
		BondedTokens:    sdk.TokensFromConsensusPower(100),
		PruningStrategy: storetypes.PruningOptionNothing,
		CleanupDir:      true,
		SigningAlgo:     string(hd.Secp256k1Type),
		KeyringOptions:  []keyring.Option{},
	}

	s.network = network.New(s.T(), s.cfg)

	_, err := s.network.WaitForHeight(1)
	s.Require().NoError(err)

	dayLockAmt, err := sdk.ParseCoinNormalized("200stake")
	s.Require().NoError(err)
	secLockAmt, err := sdk.ParseCoinNormalized("11stake")

	val := s.network.Validators[0]

	// lock tokens for a day
	_, err = lockuptestutil.MsgLockTokens(val.ClientCtx, val.Address, dayLockAmt, "24h")
	s.Require().NoError(err)

	// lock tokens for a second
	_, err = lockuptestutil.MsgLockTokens(val.ClientCtx, val.Address, secLockAmt, "1s")
	s.Require().NoError(err)

	// lock tokens for a second
	_, err = lockuptestutil.MsgLockTokens(val.ClientCtx, val.Address, secLockAmt, "1s")
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	_, err = lockuptestutil.MsgUnlockByID(val.ClientCtx, val.Address, "3")
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)
}

func (s *IntegrationTestSuite) TearDownSuite() {
	s.T().Log("tearing down integration test suite")
	s.network.Cleanup()
}

func (s *IntegrationTestSuite) TestNewLockTokensCmd() {
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("NewValidator", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)

	newAddr := sdk.AccAddress(info.GetPubKey().Address())

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(20000))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"lock 201stake tokens for 1 day",
			[]string{
				sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(201))).String(),
				fmt.Sprintf("--%s=%s", cli.FlagDuration, "24h"),
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewLockTokensCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewUnlockTokensCmd() {
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("UnlockTokensAcc", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)

	newAddr := sdk.AccAddress(info.GetPubKey().Address())

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(20000))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)

	// lock tokens for a second
	_, err = lockuptestutil.MsgLockTokens(val.ClientCtx, newAddr, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), "1s")
	s.Require().NoError(err)

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"unlock tokens",
			[]string{
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewUnlockTokensCmd()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestNewUnlockByIDCmd() {
	val := s.network.Validators[0]

	info, _, err := val.ClientCtx.Keyring.NewMnemonic("UnlockByIDCmd", keyring.English, sdk.FullFundraiserPath, hd.Secp256k1)
	s.Require().NoError(err)

	newAddr := sdk.AccAddress(info.GetPubKey().Address())

	_, err = banktestutil.MsgSendExec(
		val.ClientCtx,
		val.Address,
		newAddr,
		sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(20000))), fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
		fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
		fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(10))).String()),
	)
	s.Require().NoError(err)

	// lock tokens for a second
	res, err := lockuptestutil.MsgLockTokens(val.ClientCtx, newAddr, sdk.NewCoins(sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))), "0.0001s")
	s.Require().NoError(err)

	clientCtx := val.ClientCtx
	txResp := sdk.TxResponse{}
	s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(res.Bytes(), &txResp), res.String())
	s.Require().Equal(txResp.Logs[0].Events[0].Attributes[0].Key, "period_lock_id")
	lockID := txResp.Logs[0].Events[0].Attributes[0].Value

	_, err = s.network.WaitForHeight(1)
	s.Require().NoError(err)

	testCases := []struct {
		name         string
		args         []string
		expectErr    bool
		respType     proto.Message
		expectedCode uint32
	}{
		{
			"unlock by id",
			[]string{
				lockID,
				fmt.Sprintf("--%s=%s", flags.FlagFrom, newAddr),
				// common args
				fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
				fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
				fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
			},
			false, &sdk.TxResponse{}, 0,
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.NewUnlockByIDCmd()

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err, out.String())
				s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), tc.respType), out.String())

				txResp := tc.respType.(*sdk.TxResponse)
				s.Require().Equal(tc.expectedCode, txResp.Code, out.String())
			}
		})
	}
}

func (s *IntegrationTestSuite) TestCmdAccountUnlockableCoins() {
	val := s.network.Validators[0]

	testCases := []struct {
		name  string
		args  []string
		coins sdk.Coins
	}{
		{
			"query validator account balance",
			[]string{
				val.Address.String(),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			sdk.Coins{sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(11))},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountUnlockableCoins()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountUnlockableCoinsResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Equal(tc.coins.String(), result.Coins.String())
		})
	}
}

func (s IntegrationTestSuite) TestCmdModuleBalance() {
	val := s.network.Validators[0]

	testCases := []struct {
		name  string
		args  []string
		coins sdk.Coins
	}{
		{
			"query module balance",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			sdk.Coins{sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(211))},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdModuleBalance()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.ModuleBalanceResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Equal(tc.coins.String(), result.Coins.String())
		})
	}
}

func (s IntegrationTestSuite) TestCmdModuleLockedAmount() {
	val := s.network.Validators[0]

	testCases := []struct {
		name  string
		args  []string
		coins sdk.Coins
	}{
		{
			"query module locked balance",
			[]string{
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			sdk.Coins{sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdModuleLockedAmount()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.ModuleLockedAmountResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Equal(tc.coins.String(), result.Coins.String())
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountLockedCoins() {
	val := s.network.Validators[0]

	testCases := []struct {
		name  string
		args  []string
		coins sdk.Coins
	}{
		{
			"query account locked coins",
			[]string{
				val.Address.String(),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
			sdk.Coins{sdk.NewCoin(s.cfg.BondDenom, sdk.NewInt(200))},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountLockedCoins()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.ModuleLockedAmountResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Equal(tc.coins.String(), result.Coins.String())
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountLockedPastTime() {
	val := s.network.Validators[0]

	timestamp := time.Now().Unix()
	testCases := []struct {
		name string
		args []string
	}{
		{
			"query account locked coins past time",
			[]string{
				val.Address.String(),
				fmt.Sprintf("%d", timestamp),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountLockedPastTime()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountLockedPastTimeResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Len(result.Locks, 1)
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountUnlockedBeforeTime() {
	val := s.network.Validators[0]

	timestamp := time.Now().Unix()
	testCases := []struct {
		name string
		args []string
	}{
		{
			"query account locked coins before time",
			[]string{
				val.Address.String(),
				fmt.Sprintf("%d", timestamp),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountUnlockedBeforeTime()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountUnlockedBeforeTimeResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Len(result.Locks, 1)
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountLockedPastTimeDenom() {
	val := s.network.Validators[0]

	timestamp := time.Now().Unix()
	testCases := []struct {
		name string
		args []string
	}{
		{
			"query account locked coins past time denom",
			[]string{
				val.Address.String(),
				fmt.Sprintf("%d", timestamp),
				s.cfg.BondDenom,
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountLockedPastTimeDenom()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountLockedPastTimeDenomResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Len(result.Locks, 1)
		})
	}
}

func (s IntegrationTestSuite) TestCmdLockedByID() {
	val := s.network.Validators[0]

	testCases := []struct {
		name string
		args []string
	}{
		{
			"get lock by id",
			[]string{
				fmt.Sprintf("%d", 1),
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdLockedByID()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.LockedResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Equal(result.Lock.ID, uint64(1))
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountLockedLongerThanDuration() {
	val := s.network.Validators[0]

	testCases := []struct {
		name string
		args []string
	}{
		{
			"get account locked longer than duration",
			[]string{
				val.Address.String(),
				"1s",
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountLockedLongerThanDuration()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountLockedLongerDurationResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Len(result.Locks, 2)
		})
	}
}

func (s IntegrationTestSuite) TestCmdAccountLockedLongerThanDurationDenom() {
	val := s.network.Validators[0]

	testCases := []struct {
		name string
		args []string
	}{
		{
			"get account locked longer than duration denom",
			[]string{
				val.Address.String(),
				"1s",
				s.cfg.BondDenom,
				fmt.Sprintf("--%s=json", tmcli.OutputFlag),
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		s.Run(tc.name, func() {
			cmd := cli.GetCmdAccountLockedLongerThanDurationDenom()
			clientCtx := val.ClientCtx

			out, err := clitestutil.ExecTestCLICmd(clientCtx, cmd, tc.args)
			s.Require().NoError(err)

			var result types.AccountLockedLongerDurationDenomResponse
			s.Require().NoError(clientCtx.JSONMarshaler.UnmarshalJSON(out.Bytes(), &result))
			s.Require().Len(result.Locks, 2)
		})
	}
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
