package keeper

import (
	"bytes"
	"context"
	"fmt"

	"github.com/c-osmosis/osmosis/x/gamm/utils"
	"github.com/c-osmosis/osmosis/x/lockup/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type msgServer struct {
	keeper Keeper
}

// NewMsgServerImpl returns an instance of MsgServer
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{
		keeper: keeper,
	}
}

var _ types.MsgServer = msgServer{}

func (server msgServer) LockTokens(goCtx context.Context, msg *types.MsgLockTokens) (*types.MsgLockTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	lock, err := server.keeper.LockTokens(ctx, msg.Owner, msg.Coins, msg.Duration)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtLockTokens,
			sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(lock.ID)),
			sdk.NewAttribute(types.AttributePeriodLockOwner, lock.Owner.String()),
			sdk.NewAttribute(types.AttributePeriodLockAmount, lock.Coins.String()),
			sdk.NewAttribute(types.AttributePeriodLockDuration, lock.Duration.String()),
			sdk.NewAttribute(types.AttributePeriodLockUnlockTime, lock.EndTime.String()),
		),
	})

	return &types.MsgLockTokensResponse{}, nil
}

func (server msgServer) UnlockPeriodLock(goCtx context.Context, msg *types.MsgUnlockPeriodLock) (*types.MsgUnlockPeriodLockResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	lock, err := server.keeper.UnlockPeriodLockByID(ctx, msg.ID)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	if !bytes.Equal(msg.Owner, lock.Owner) {
		return nil, sdkerrors.Wrap(types.ErrNotLockOwner, fmt.Sprintf("msg sender(%s) and lock owner(%s) does not match", msg.Owner.String(), lock.Owner.String()))
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.TypeEvtUnlock,
			sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(lock.ID)),
			sdk.NewAttribute(types.AttributePeriodLockOwner, lock.Owner.String()),
			sdk.NewAttribute(types.AttributePeriodLockDuration, lock.Duration.String()),
			sdk.NewAttribute(types.AttributePeriodLockUnlockTime, lock.EndTime.String()),
		),
	})

	return &types.MsgUnlockPeriodLockResponse{}, nil
}

func (server msgServer) UnlockTokens(goCtx context.Context, msg *types.MsgUnlockTokens) (*types.MsgUnlockTokensResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	unlocks, coins, err := server.keeper.UnlockAllUnlockableCoins(ctx, msg.Owner)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, err.Error())
	}

	events := sdk.Events{
		sdk.NewEvent(
			types.TypeEvtUnlockTokens,
			sdk.NewAttribute(types.AttributePeriodLockOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeUnlockedCoins, coins.String()),
		),
	}
	for _, lock := range unlocks {
		event := sdk.NewEvent(
			types.TypeEvtUnlock,
			sdk.NewAttribute(types.AttributePeriodLockID, utils.Uint64ToString(lock.ID)),
			sdk.NewAttribute(types.AttributePeriodLockOwner, lock.Owner.String()),
			sdk.NewAttribute(types.AttributePeriodLockDuration, lock.Duration.String()),
			sdk.NewAttribute(types.AttributePeriodLockUnlockTime, lock.EndTime.String()),
		)
		events = events.AppendEvent(event)
	}
	ctx.EventManager().EmitEvents(events)

	return &types.MsgUnlockTokensResponse{}, nil
}
