# AtomOne Fork - OasisAtone

This is a fork of Osmosis adapted for AtomOne with the following key changes:

## Key Changes from Osmosis

1. **Address Prefix**: Changed from `osmo` to `atone`
   - Account addresses: `atone...`
   - Validator operator addresses: `atonevaloper...`
   - Consensus addresses: `atonevalcons...`

2. **Denomination**: Changed from `uosmo` to `uatone`
   - Base denom: `uatone`
   - Display denom: `atone`

3. **Chain ID**: Using `atone-local-1` for local development

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
# Initialize the chain
bash init-localnet.sh

# Start the node
./build/osmosisd start --home $HOME/.atomone-local
```

The initialization script:
- Creates a validator with mnemonic: `guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host`
- Funds the validator with 100000000000000000000stake
- Sets up fast block times (1s) for development
- Enables API server

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
./build/osmosisd tx bank send validator <recipient_address> 1000stake --keyring-backend test --home $HOME/.atomone-local --chain-id atone-local-1
```

## Code Changes

### Modified Files

1. **app/params/config.go**
   - Changed `Bech32PrefixAccAddr` from `"osmo"` to `"atone"`

2. **x/protorev/types/params.go**
   - Modified validation to allow empty admin address for local testing
   - This fixes initialization issues with the protorev module

3. **app/app.go** (line 810)
   - Fixed hardcoded `"osmovaloper"` string to use `appparams.Bech32PrefixValAddr`
   - This ensures validator addresses use the correct prefix

4. **go.mod**
   - Updated `github.com/bytedance/sonic` to v1.14.1 for Go 1.25 compatibility

## Development Notes

- The protorev module is disabled by default in localnet as it requires additional setup
- Block time is set to 1s for faster development iterations
- API server is enabled on the default port (1317)
- Validator starts with bonded status and 50M delegated tokens

## Original Osmosis README

See [README.md](./README.md) for the original Osmosis documentation.

## Repository

- **Origin**: https://github.com/dubious-tom/oasis-atone (this fork)
- **Upstream**: https://github.com/osmosis-labs/osmosis (original)

