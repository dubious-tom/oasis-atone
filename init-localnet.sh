#!/bin/bash
set -e

# Configuration
CHAIN_ID="atone-local-1"
BINARY="./build/osmosisd"
HOME_DIR="$HOME/.atomone-local"
KEYRING="test"
MONIKER="validator"

# Validator mnemonic for key recovery
MNEMONIC="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"

echo "===== Initializing AtomOne LocalNet ====="

# Clean up old data
echo "Cleaning up old data..."
rm -rf $HOME_DIR

# Initialize chain
echo "Initializing chain..."
$BINARY init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR

# Recover validator key from mnemonic
echo "Recovering validator key..."
echo "$MNEMONIC" | $BINARY keys add validator --recover --keyring-backend $KEYRING --home $HOME_DIR

# Get validator address
VAL_ADDR=$($BINARY keys show validator -a --keyring-backend $KEYRING --home $HOME_DIR)
echo "    Validator Address: $VAL_ADDR"

# Get validator operator address (valoper)
VAL_VALOPER=$($BINARY keys show validator --bech val -a --keyring-backend $KEYRING --home $HOME_DIR)
echo "    Validator Operator: $VAL_VALOPER"

# Add genesis account with tokens
echo "Adding genesis account..."
$BINARY genesis add-genesis-account validator 100000000000000000000stake --keyring-backend $KEYRING --home $HOME_DIR

# Get validator consensus pubkey as JSON
echo "Getting validator public key..."
VAL_PUBKEY=$($BINARY comet show-validator --home $HOME_DIR)

# Validate pubkey was captured correctly
if [ -z "$VAL_PUBKEY" ] || [ "$VAL_PUBKEY" = "null" ]; then
    echo "Error: Failed to get validator pubkey: '$VAL_PUBKEY'"
    exit 1
fi

echo "    Validator Pubkey: $VAL_PUBKEY"

# Manually add validator to genesis
echo "Adding validator to genesis..."

# Set stake denom and other params
jq '.app_state.staking.params.bond_denom = "stake"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

jq '.app_state.staking.params.unbonding_time = "240s"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set crisis denom
jq '.app_state.crisis.constant_fee.denom = "stake"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set gov deposit denom
jq '.app_state.gov.params.min_deposit[0].denom = "stake"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set mint denom
jq '.app_state.mint.params.mint_denom = "stake"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Disable protorev module (set admin to empty to bypass validation)
jq '.app_state.protorev.params.admin = ""' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

jq '.app_state.protorev.params.enabled = false' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add validator with 50M stake
jq --arg valoper "$VAL_VALOPER" \
  --argjson pubkey "$VAL_PUBKEY" \
  '.app_state.staking.validators += [{
    "operator_address": $valoper,
    "consensus_pubkey": $pubkey,
    "jailed": false,
    "status": "BOND_STATUS_BONDED",
    "tokens": "50000000000000000",
    "delegator_shares": "50000000000000000.000000000000000000",
    "description": {
      "moniker": "validator",
      "identity": "",
      "website": "",
      "security_contact": "",
      "details": ""
    },
    "unbonding_height": "0",
    "unbonding_time": "1970-01-01T00:00:00Z",
    "commission": {
      "commission_rates": {
        "rate": "0.100000000000000000",
        "max_rate": "0.200000000000000000",
        "max_change_rate": "0.010000000000000000"
      },
      "update_time": "1970-01-01T00:00:00Z"
    },
    "min_self_delegation": "1"
  }]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add validator power
jq --arg valoper "$VAL_VALOPER" \
  '.app_state.staking.last_validator_powers += [{
    "address": $valoper,
    "power": "50000000000000000"
  }]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set last total power
jq '.app_state.staking.last_total_power = "50000000000000000"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add delegation
jq --arg valoper "$VAL_VALOPER" \
  --arg deladdr "$VAL_ADDR" \
  '.app_state.staking.delegations += [{
    "delegator_address": $deladdr,
    "validator_address": $valoper,
    "shares": "50000000000000000.000000000000000000"
  }]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add bonded tokens to bonded_tokens_pool module account
# Module address for bonded_tokens_pool is derived from the module name
BONDED_POOL_ADDR="atone1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3mm4g00"
jq --arg addr "$BONDED_POOL_ADDR" \
  '.app_state.bank.balances += [{
    "address": $addr,
    "coins": [{"denom": "stake", "amount": "50000000000000000"}]
  }]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Get consensus address from pubkey
# Note: The command outputs to stdout, capture it
VAL_CONSENSUS_ADDR=$($BINARY comet show-address --home $HOME_DIR 2>&1 | tail -1)

# Validate consensus address
if [ -z "$VAL_CONSENSUS_ADDR" ] || [ "$VAL_CONSENSUS_ADDR" = "null" ]; then
    echo "Error: Failed to get consensus address: '$VAL_CONSENSUS_ADDR'"
    exit 1
fi

# Strip any whitespace
VAL_CONSENSUS_ADDR=$(echo "$VAL_CONSENSUS_ADDR" | tr -d '[:space:]')

echo "    Consensus Address: $VAL_CONSENSUS_ADDR"

# Add validator signing info to slashing module
jq --arg consaddr "$VAL_CONSENSUS_ADDR" \
  --argjson pubkey "$VAL_PUBKEY" \
  '.app_state.slashing.signing_infos += [{
    "address": $consaddr,
    "validator_signing_info": {
      "address": $consaddr,
      "start_height": "0",
      "index_offset": "0",
      "jailed_until": "1970-01-01T00:00:00Z",
      "tombstoned": false,
      "missed_blocks_counter": "0"
    }
  }]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Configure fast blocks for testing
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/' $HOME_DIR/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/' $HOME_DIR/config/config.toml

# Enable API
sed -i '' 's/enable = false/enable = true/' $HOME_DIR/config/app.toml

# Set gas prices
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025stake"/' $HOME_DIR/config/app.toml

echo "===== Initialization Complete ====="
echo ""
echo "To start the chain, run:"
echo "  $BINARY start --home $HOME_DIR"
echo ""
echo "Validator info:"
echo "  Address:    $VAL_ADDR"
echo "  Valoper:    $VAL_VALOPER"
echo "  Consensus:  $VAL_CONSENSUS_ADDR"
echo "  Mnemonic:   $MNEMONIC"

