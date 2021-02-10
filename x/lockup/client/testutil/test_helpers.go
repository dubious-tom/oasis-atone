package testutil

import (
	"fmt"

	lockupcli "github.com/c-osmosis/osmosis/x/lockup/client/cli"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/testutil"
	clitestutil "github.com/cosmos/cosmos-sdk/testutil/cli"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// commonArgs is args for CLI test commands
var commonArgs = []string{
	fmt.Sprintf("--%s=true", flags.FlagSkipConfirmation),
	fmt.Sprintf("--%s=%s", flags.FlagBroadcastMode, flags.BroadcastBlock),
	fmt.Sprintf("--%s=%s", flags.FlagFees, sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(10))).String()),
}

// MsgLockTokens creates a redelegate message.
func MsgLockTokens(clientCtx client.Context, owner fmt.Stringer, amount fmt.Stringer, duration string, extraArgs ...string) (testutil.BufferWriter, error) {

	args := []string{
		amount.String(),
		fmt.Sprintf("--%s=%s", lockupcli.FlagDuration, duration),
		fmt.Sprintf("--%s=%s", flags.FlagFrom, owner.String()),
	}

	args = append(args, commonArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, lockupcli.NewLockTokensCmd(), args)
}

// MsgUnlockTokens unlock all unlockable tokens from an account
func MsgUnlockTokens(clientCtx client.Context, owner fmt.Stringer, extraArgs ...string) (testutil.BufferWriter, error) {

	args := []string{
		fmt.Sprintf("--%s=%s", flags.FlagFrom, owner.String()),
	}

	args = append(args, commonArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, lockupcli.NewUnlockTokensCmd(), args)
}

// MsgUnlockByID unlock unlockable tokens
func MsgUnlockByID(clientCtx client.Context, owner fmt.Stringer, ID string, extraArgs ...string) (testutil.BufferWriter, error) {

	args := []string{
		ID,
		fmt.Sprintf("--%s=%s", flags.FlagFrom, owner.String()),
	}

	args = append(args, commonArgs...)
	return clitestutil.ExecTestCLICmd(clientCtx, lockupcli.NewUnlockByIDCmd(), args)
}
