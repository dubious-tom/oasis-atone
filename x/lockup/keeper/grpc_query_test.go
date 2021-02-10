package keeper_test

import (
	"time"

	"github.com/c-osmosis/osmosis/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (suite *KeeperTestSuite) LockTokens(addr sdk.AccAddress, coins sdk.Coins, duration time.Duration) {
	suite.app.BankKeeper.SetBalances(suite.ctx, addr, coins)
	_, err := suite.app.LockupKeeper.LockTokens(suite.ctx, addr, coins, duration)
	suite.Require().NoError(err)
}

func (suite *KeeperTestSuite) TestModuleBalance() {
	// test for module balance check
	suite.SetupTest()

	// initial module balance check
	res, err := suite.app.LockupKeeper.ModuleBalance(sdk.WrapSDKContext(suite.ctx), &types.ModuleBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// final module balance check
	res, err = suite.app.LockupKeeper.ModuleBalance(sdk.WrapSDKContext(suite.ctx), &types.ModuleBalanceRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)
}

func (suite *KeeperTestSuite) TestModuleLockedAmount() {
	// test for module locked balance check
	suite.SetupTest()

	// initial module locked balance check
	res, err := suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// current module locked balance check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// module locked balance after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// module locked balance after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.ModuleLockedAmount(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.ModuleLockedAmountRequest{})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
}

func (suite *KeeperTestSuite) TestAccountUnlockableCoins() {
	// test for module unlockable coins check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address unlockable coins check
	res, err := suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: sdk.AccAddress{}})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
	// initial account unlockable coins check
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account unlockable coins check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockableCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// account unlockable balance after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.AccountUnlockableCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// account unlockable balance after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountUnlockableCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.AccountUnlockableCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)
}

func (suite *KeeperTestSuite) TestAccountLockedCoins() {
	// test for account locked coins check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locked coins check
	res, err := suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{Owner: sdk.AccAddress{}})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
	// initial account locked coins check
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locked coins check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, coins)

	// account locked coins after 1 second = unlockTime
	now := suite.ctx.BlockTime()
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(time.Second))), &types.AccountLockedCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})

	// account locked coins after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedCoins(sdk.WrapSDKContext(suite.ctx.WithBlockTime(now.Add(2*time.Second))), &types.AccountLockedCoinsRequest{Owner: addr1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Coins, sdk.Coins{})
}

func (suite *KeeperTestSuite) TestAccountLockedPastTime() {
	// test for account locks check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address locks check
	res, err := suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: sdk.AccAddress{}, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial account locks check
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locks check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1, Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTime(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeRequest{Owner: addr1, Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountUnlockedBeforeTime() {
	// test for account unlockables check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address unlockables check
	res, err := suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: sdk.AccAddress{}, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial account unlockables check
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account unlockables check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1, Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account unlockables after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1, Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account unlockables after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountUnlockedBeforeTime(sdk.WrapSDKContext(suite.ctx), &types.AccountUnlockedBeforeTimeRequest{Owner: addr1, Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)
}

func (suite *KeeperTestSuite) TestAccountLockedPastTimeDenom() {
	// test for account locks by denom check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))
	now := suite.ctx.BlockTime()

	// empty address locks by denom check
	res, err := suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: sdk.AccAddress{}, Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial account locks by denom check
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locks by denom check = unlockTime - 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stake", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks by not available denom
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stake2", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks by denom after 1 second = unlockTime
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stake", Timestamp: now.Add(time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks by denom after 2 second = unlockTime + 1s
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stake", Timestamp: now.Add(2 * time.Second)})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// try querying with prefix coins like "stak" for potential attack
	res, err = suite.app.LockupKeeper.AccountLockedPastTimeDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedPastTimeDenomRequest{Owner: addr1, Denom: "stak", Timestamp: now})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestLockedByID() {
	// test for all locks check
	suite.SetupTest()
	now := suite.ctx.BlockTime()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// lock by not avaialble id check
	res, err := suite.app.LockupKeeper.LockedByID(sdk.WrapSDKContext(suite.ctx), &types.LockedRequest{LockId: 0})
	suite.Require().Error(err)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// lock by available availble id check
	res, err = suite.app.LockupKeeper.LockedByID(sdk.WrapSDKContext(suite.ctx), &types.LockedRequest{LockId: 1})
	suite.Require().NoError(err)
	suite.Require().Equal(res.Lock.ID, uint64(1))
	suite.Require().Equal(res.Lock.Owner, addr1)
	suite.Require().Equal(res.Lock.Coins, coins)
	suite.Require().Equal(res.Lock.Duration, time.Second)
	suite.Require().Equal(res.Lock.EndTime, now.Add(time.Second))
}

func (suite *KeeperTestSuite) TestAccountLockedLongerThanDuration() {
	// test for account locks longer than duration check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locks longer than duration check
	res, err := suite.app.LockupKeeper.AccountLockedLongerThanDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: sdk.AccAddress{}, Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial account locks longer than duration check
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1, Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locks longer than duration check, duration = 0s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1, Duration: 0})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1, Duration: time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check, duration = 2s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDuration(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationRequest{Owner: addr1, Duration: 2 * time.Second})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}

func (suite *KeeperTestSuite) TestAccountLockedLongerThanDurationDenom() {
	// test for account locks longer than duration by denom check
	suite.SetupTest()
	addr1 := sdk.AccAddress([]byte("addr1---------------"))

	// empty address locks longer than duration by denom check
	res, err := suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: sdk.AccAddress{}, Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
	// initial account locks longer than duration by denom check
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// lock coins
	coins := sdk.Coins{sdk.NewInt64Coin("stake", 10)}
	suite.LockTokens(addr1, coins, time.Second)

	// account locks longer than duration check by denom, duration = 0s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: 0, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check by denom, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 1)

	// account locks longer than duration check by not available denom, duration = 1s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: time.Second, Denom: "stake2"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// account locks longer than duration check by denom, duration = 2s
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: 2 * time.Second, Denom: "stake"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)

	// try querying with prefix coins like "stak" for potential attack
	res, err = suite.app.LockupKeeper.AccountLockedLongerThanDurationDenom(sdk.WrapSDKContext(suite.ctx), &types.AccountLockedLongerDurationDenomRequest{Owner: addr1, Duration: 0, Denom: "sta"})
	suite.Require().NoError(err)
	suite.Require().Len(res.Locks, 0)
}
