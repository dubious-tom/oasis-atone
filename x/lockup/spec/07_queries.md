<!--
order: 7
-->

# Queries

In this section we describe the queries required on grpc server.

```protobuf
// Query defines the gRPC querier service.
service Query {
    // Return full balance of the module
	rpc ModuleBalance(ModuleBalanceRequest) returns (ModuleBalanceResponse);
	// Return locked balance of the module
	rpc ModuleLockedAmount(ModuleLockedAmountRequest) returns (ModuleLockedAmountResponse);

	// Returns unlockable coins which are not withdrawn yet
	rpc AccountUnlockableCoins(AccountUnlockableCoinsRequest) returns (AccountUnlockableCoinsResponse);
	// Return a locked coins that can't be withdrawn
	rpc AccountLockedCoins(AccountLockedCoinsRequest) returns (AccountLockedCoinsResponse);

	// Returns locked records of an account with unlock time beyond timestamp
	rpc AccountLockedPastTime(AccountLockedPastTimeRequest) returns (AccountLockedPastTimeResponse);
	// Returns unlocked records with unlock time before timestamp
	rpc AccountUnlockedBeforeTime(AccountUnlockedBeforeTimeRequest) returns (AccountUnlockedBeforeTimeResponse);

	// Returns lock records by address, timestamp, denom
	rpc AccountLockedPastTimeDenom(AccountLockedPastTimeDenomRequest) returns (AccountLockedPastTimeDenomResponse);
	// Returns lock record by id
	rpc LockedByID(LockedRequest) returns (LockedResponse);

	// Returns account locked records with longer duration
	rpc AccountLockedLongerThanDuration(AccountLockedLongerDurationRequest) returns (AccountLockedLongerDurationResponse);
	// Returns account's locked records for a denom with longer duration
	rpc AccountLockedLongerThanDurationDenom(AccountLockedLongerDurationDenomRequest) returns (AccountLockedLongerDurationDenomResponse);
}
```