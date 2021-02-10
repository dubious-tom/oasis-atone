package rest

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"github.com/c-osmosis/osmosis/x/lockup/types"
	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

func registerQueryRoutes(clientCtx client.Context, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/%s", types.QueryModuleBalance), queryModuleBalanceFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s", types.QueryModuleLockedAmount), queryModuleLockedAmountFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}", types.QueryAccountUnlockableCoins, RestOwnerAddress), queryAccountUnlockableCoinsFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}", types.QueryAccountLockedCoins, RestOwnerAddress), queryAccountLockedCoinsFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/{%s}", types.QueryAccountLockedPastTime, RestOwnerAddress, RestTimestamp), queryAccountLockedPastTimeFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/{%s}", types.QueryAccountUnlockedBeforeTime, RestOwnerAddress, RestTimestamp), queryAccountUnlockedBeforeTimeFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/{%s}/{%s}", types.QueryAccountLockedPastTimeDenom, RestOwnerAddress, RestDenom, RestTimestamp), queryAccountLockedPastTimeDenomFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}", types.QueryLockedByID, LockID), queryLockedByIDFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/{%s}", types.QueryAccountLockedLongerThanDuration, RestOwnerAddress, RestDuration), queryAccountLockedLongerThanDurationFn(clientCtx)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/{%s}/{%s}/{%s}", types.QueryAccountLockedLongerThanDurationDenom, RestOwnerAddress, RestDenom, RestDuration), queryAccountLockedLongerThanDurationDenomFn(clientCtx)).Methods("GET")
}

func queryModuleBalanceFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryModuleBalance), nil)
		if rest.CheckNotFoundError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryModuleLockedAmountFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, height, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryModuleLockedAmount), []byte{})
		if rest.CheckInternalServerError(w, err) {
			return
		}

		clientCtx = clientCtx.WithHeight(height)
		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountUnlockableCoinsFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.AccountUnlockableCoinsRequest{Owner: owner}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountUnlockableCoins), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountLockedCoinsFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.AccountLockedCoinsRequest{Owner: owner}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountLockedCoins), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountLockedPastTimeFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		strTimestamp := vars[RestTimestamp]
		timestamp, err := strconv.ParseInt(strTimestamp, 10, 64)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.AccountLockedPastTimeRequest{Owner: owner, Timestamp: time.Unix(timestamp, 0)}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountLockedPastTime), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountUnlockedBeforeTimeFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		strTimestamp := vars[RestTimestamp]
		timestamp, err := strconv.ParseInt(strTimestamp, 10, 64)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.AccountUnlockedBeforeTimeRequest{Owner: owner, Timestamp: time.Unix(timestamp, 0)}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountUnlockedBeforeTime), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountLockedPastTimeDenomFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		strTimestamp := vars[RestTimestamp]
		timestamp, err := strconv.ParseInt(strTimestamp, 10, 64)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		denom := vars[RestDenom]

		params := types.AccountLockedPastTimeDenomRequest{Owner: owner, Timestamp: time.Unix(timestamp, 0), Denom: denom}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountLockedPastTimeDenom), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryLockedByIDFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)

		strLockID := vars[LockID]
		lockID, err := strconv.ParseUint(strLockID, 10, 64)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.LockedRequest{LockId: lockID}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryLockedByID), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountLockedLongerThanDurationFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		strDuration := vars[RestDuration]
		duration, err := time.ParseDuration(strDuration)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		params := types.AccountLockedLongerDurationRequest{Owner: owner, Duration: duration}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountLockedLongerThanDuration), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}

func queryAccountLockedLongerThanDurationDenomFn(clientCtx client.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		strOwnerAddress := vars[RestOwnerAddress]
		owner, err := sdk.AccAddressFromBech32(strOwnerAddress)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		strDuration := vars[RestDuration]
		duration, err := time.ParseDuration(strDuration)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		denom := vars[RestDenom]

		params := types.AccountLockedLongerDurationDenomRequest{Owner: owner, Duration: duration, Denom: denom}

		bz, err := clientCtx.LegacyAmino.MarshalJSON(params)
		if rest.CheckBadRequestError(w, err) {
			return
		}

		res, _, err := clientCtx.QueryWithData(fmt.Sprintf("custom/lockup/%s", types.QueryAccountLockedLongerThanDurationDenom), bz)
		if rest.CheckInternalServerError(w, err) {
			return
		}

		if rest.CheckInternalServerError(w, err) {
			return
		}

		rest.PostProcessResponse(w, clientCtx, res)
	}
}
