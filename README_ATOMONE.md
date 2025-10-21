# Oasis (Osmosis Fork)

This is a fork of Osmosis adapted for Oasis with the following key changes:

## Key Changes from Osmosis

1. **Address Prefix**: Changed from `osmo` to `oasis`
   - Account addresses: `oasis1...`
   - Validator operator addresses: `oasisvaloper...`
   - Consensus addresses: `oasisvalcons...`

2. **Chain ID**: `localosmosis-oasis` for local development

3. **Denomination**: Using standard cosmos `stake` token for development

4. **Status**: ✅ **Fully Functional** - All transactions working correctly with `oasis` prefix

## Quick Start - Local Development

### Prerequisites
- Go 1.25+
- Make
- jq (for genesis manipulation)

### Build
```bash
make build
# or
go install ./cmd/osmosisd
```

### Initialize & Start LocalNet
```bash
# Initialize the chain (if not already done)
bash init-localnet.sh

# Start the node
./build/osmosisd start --home $HOME/.atomone-local
```

The local network:
- Chain ID: `localosmosis-oasis`
- Address prefix: `oasis`
- Test validator mnemonic: `guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host`
- Fast block times (1s) for development
- API server enabled on port 1317
- RPC on port 26657

### Useful Commands

Query validator status:
```bash
./build/osmosisd query staking validators --home $HOME/.atomone-local
```

Check your balance:
```bash
./build/osmosisd query bank balances $(./build/osmosisd keys show validator -a --keyring-backend test --home $HOME/.atomone-local) --home $HOME/.atomone-local
```

Send tokens:
```bash
./build/osmosisd tx bank send validator <recipient_address> 1000stake --keyring-backend test --home $HOME/.atomone-local --chain-id localosmosis-oasis
```

## Code Changes

### Modified Files (Backend)

1. **app/params/config.go**
   - Changed `Bech32PrefixAccAddr` from `"osmo"` to `"oasis"`
   - Changed `HumanCoinUnit` from `"osmo"` to `"oasis"`
   - Changed `BaseCoinUnit` from `"uosmo"` to `"uoasis"`

2. **app/params/proto.go** ⭐ **Critical Fix**
   - Changed hardcoded address prefixes in `MakeEncodingConfig()` from `"osmo"`/`"osmovaloper"` to use `Bech32PrefixAccAddr`/`Bech32PrefixValAddr` constants
   - **This was the root cause** of "hrp does not match bech32 prefix" errors during transaction signing

3. **app/params/amino.go** ⭐ **Critical Fix**
   - Changed hardcoded address prefixes in `MakeEncodingConfig()` to use constants (same as proto.go)

4. **app/app.go**
   - Fixed hardcoded `"osmovaloper"` string to use `appparams.Bech32PrefixValAddr` (line 810)
   - Wrapped hardcoded `osmo` address funding logic in `if appparams.Bech32PrefixAccAddr == "osmo"` block (lines 996-1049)

5. **cmd/osmosisd/cmd/root.go**
   - Changed temp app initialization chain ID from `"osmosis-1"` to `"oasis-1"` (line 361)

6. **x/protorev/types/params.go**
   - Modified validation to allow empty admin address for local testing

7. **go.mod**
   - Updated `github.com/bytedance/sonic` to v1.14.1 for Go 1.25 compatibility

## Development Notes

- The protorev module is disabled by default in localnet as it requires additional setup
- Block time is set to 1s for faster development iterations
- API server is enabled on the default port (1317)
- Validator starts with bonded status and 50M delegated tokens

## Troubleshooting

### "hrp does not match bech32 prefix: expected 'osmo' got 'oasis'" Error

If you encounter this error when signing or broadcasting transactions, it means the encoding config still has hardcoded "osmo" prefixes. Ensure you've:

1. ✅ Modified `app/params/proto.go` and `app/params/amino.go` to use `Bech32PrefixAccAddr` and `Bech32PrefixValAddr` constants instead of hardcoded strings
2. ✅ Rebuilt the binary: `make build`
3. ✅ Restarted the node if it was running

This was the root cause of transaction signing failures and has been fixed in this fork.

## Original Osmosis README

See [README.md](./README.md) for the original Osmosis documentation.

## Frontend

The Osmosis frontend is configured to work with the local Oasis network. **Note**: The frontend will show Osmosis branding but connects to your local Oasis chain with correct `oasis` prefix.

**Setup & Run**:
```bash
cd frontend
yarn install  # First time only
yarn build    # First time only (~5 min)
yarn dev      # Start dev server at http://localhost:3000
```

**Configuration**: `frontend/packages/web/.env.local`
```bash
NEXT_PUBLIC_IS_TESTNET=true
NEXT_PUBLIC_OSMOSIS_RPC_OVERWRITE=http://localhost:26657
NEXT_PUBLIC_OSMOSIS_REST_OVERWRITE=http://localhost:1317
NEXT_PUBLIC_OSMOSIS_CHAIN_ID_OVERWRITE=localosmosis-oasis
NEXT_PUBLIC_OSMOSIS_CHAIN_NAME_OVERWRITE=Oasis Local
```

**Using the Frontend**:

**IMPORTANT - Remove Old Chain from Keplr**:

The frontend now has the correct configuration, but Keplr caches chains. You MUST remove the old chain:

Method 1 - Reset Keplr (Easiest):
1. Keplr extension → Settings → Advanced → "Clear All"  
2. **WARNING**: This removes all connected sites. You'll need to reconnect everywhere.
3. After clearing, refresh browser and connect to frontend

Method 2 - Manual Removal:
1. Keplr → Select "Oasis Local" from chain dropdown
2. Click the three dots (...) → Settings
3. Look for "Delete this chain" or similar option
4. Confirm deletion
5. Refresh browser and reconnect

Method 3 - Incognito/Private Window:
1. Open http://localhost:3000 in incognito/private browsing
2. Install Keplr in incognito (or enable existing extension)
3. Connect - should suggest correct chain
4. Once confirmed working, use Method 1 or 2 for your main browser

**Connect Wallet**:
1. Open http://localhost:3000 in your browser
2. Click "Connect Wallet" and choose Keplr
3. Approve the chain addition (it will now auto-suggest with correct `oasis` prefix)
4. Import your test validator account using the mnemonic or use an existing account
5. You're ready to use the frontend!

**Current Status**:
- ✅ Address prefix: Correctly shows `oasis` 
- ✅ Chain ID: Correctly shows `localosmosis-oasis`
- ✅ Chain Name: Correctly shows "Oasis Local"
- ✅ Currencies: Correctly shows STAKE token with `stake` denomination
- ✅ Stake Currency: Correctly shows `stake`
- ✅ Fee Currencies: Correctly shows `stake`
- ✅ **Fully Functional** - All basic wallet operations working
- ⚠️ **Swaps Disabled** - Local chains don't have SQS (Sidecar Query Server) for swap routing

**What was fixed (Frontend)**:
1. Modified `frontend/packages/web/config/utils.ts` to override chain configuration for local Oasis chain
2. Modified `frontend/packages/web/config/generate-lists.ts` to generate custom asset list with STAKE token
3. All chain configurations (both for Keplr and cosmos-kit) now use correct `oasis` prefix and `stake` token
4. **Reduced query polling frequency** - Set default refetch interval to 2 minutes (120s) for local/testnet to minimize log clutter. Production uses 30s. Specific queries that need real-time updates (like swap quotes) override this with shorter intervals.

### Modified Files (Frontend)

1. **frontend/packages/web/.env.local**
   - Created with environment variable overrides for local Oasis network

2. **frontend/packages/web/config/utils.ts**
   - Modified `getKeplrCompatibleChain` to detect and override config for local Oasis chain
   - Modified `getChainList` to override raw chain object for cosmos-kit compatibility

3. **frontend/packages/web/config/generate-lists.ts**
   - Modified `generateAssetListFile` to create custom asset list for local Oasis chain with STAKE token

4. **frontend/packages/stores/src/queries-external/icns/index.ts**
   - Disabled ICNS queries for local chains (ICNS contract uses `osmo` addresses not available locally)

5. **frontend/packages/web/hooks/queries/osmosis/use-icns-name.ts**
   - Disabled ICNS hook for local chains to prevent errors

6. **frontend/packages/web/hooks/use-swap.tsx** ⭐ **Critical Fix**
   - Changed default tokens from ATOM/uosmo to `stake`/`uosmo`
   - This makes the frontend work with local Oasis chain by default
   - Users visiting without URL params will see `stake` and `uosmo` as initial token selection

7. **frontend/packages/trpc/src/parameter-types.ts** ⭐ **Critical Fix**
   - Removed hardcoded `startsWith("osmo")` validation from `OsmoAddressSchema` and `UserOsmoAddressSchema`
   - **This was blocking API calls** with `oasis` addresses
   - Now accepts any bech32 address prefix

8. **frontend/packages/web/utils/trpc.ts** - Query Polling Configuration
   - Set default refetch interval to 2 minutes (120s) for local/testnet environments
   - Set default staleTime to 1 minute (60s) for local/testnet
   - Reduces log clutter from frequent polling of `getUserAssets`, `getAssetPrice`, `routeTokenOutGivenIn`, and `oneClickTrading.getParameters`
   - Production still uses 30s/15s intervals for real-time feel

### Known Limitations

**Swap Functionality**:
The swap feature requires an SQS (Sidecar Query Server) to calculate optimal swap routes and quotes. The mainnet SQS at `https://sqs.osmosis.zone` only recognizes tokens from the Osmosis mainnet, so it will return errors like `TRPCClientError: denom is not a valid chain denom (stake)` when trying to swap local tokens.

To enable swaps on your local chain, you would need to:
1. Run a local SQS instance (see [osmosis/ingest/sqs](https://github.com/osmosis-labs/osmosis/tree/main/ingest/sqs))
2. Set `NEXT_PUBLIC_SIDECAR_BASE_URL=http://localhost:9092` in `.env.local`

For now, the frontend will work for:
- ✅ Viewing balances
- ✅ Sending tokens (use the CLI: `osmosisd tx bank send`)
- ✅ Staking operations
- ✅ Governance
- ❌ Swaps (requires local SQS)

## Git History & LFS

**Important Note**: This repository's git history has been cleaned using `git-filter-repo` to remove Git LFS pointer files that referenced unavailable large genesis files from the upstream Osmosis repository.

The upstream Osmosis repository exceeded its GitHub LFS storage quota, making those LFS objects unavailable. Since these large genesis files are not needed for development or deployment of the Oasis fork, they were removed from the entire git history while preserving all code changes and commit lineage.

**What this means**:
- ✅ All code changes from Osmosis and Oasis fork are preserved
- ✅ Git history and commit lineage remains intact
- ✅ Connection to upstream Osmosis repository is maintained
- ✅ Repository can be cloned and pushed without LFS errors
- ❌ Some large genesis files from upstream history are no longer accessible (but weren't needed anyway)

If you need to access historical genesis files, they may be available from the [Osmosis networks repository](https://github.com/osmosis-labs/networks) or other archive sources.

## Repository

- **Origin**: https://github.com/dubious-tom/oasis-atone (this fork)
- **Frontend**: https://github.com/dubious-tom/osmosis-frontend (frontend fork)
- **Upstream**: https://github.com/osmosis-labs/osmosis (original)

