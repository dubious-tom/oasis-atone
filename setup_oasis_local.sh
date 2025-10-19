#!/bin/bash
#
# Oasis Osmosis LocalNet Initialization Script
# This creates a single-node local testnet with the Oasis token
#

set -e

BINARY="./build/osmosisd"
CHAIN_ID="localosmosis-oasis"
HOME_DIR="$HOME/.oasis-local"
KEYRING="test"
MONIKER="oasis-validator"

# Known test mnemonic for reproducibility
MNEMONIC="race draft rival universe maid cheese steel logic crowd fork comic easy truth drift tomorrow eye buddy head time cash swing swift midnight borrow"

echo "🚀 Oasis Osmosis LocalNet Initialization"
echo "=========================================="
echo ""

# Clean previous state
echo "[1/8] Cleaning previous state..."
rm -rf $HOME_DIR
mkdir -p $HOME_DIR

# Initialize node
echo "[2/8] Initializing node..."
$BINARY init $MONIKER --chain-id $CHAIN_ID --home $HOME_DIR > /dev/null 2>&1

# Create validator key
echo "[3/8] Creating validator key..."
echo $MNEMONIC | $BINARY keys add validator --recover --keyring-backend $KEYRING --home $HOME_DIR > /dev/null 2>&1

# Get validator pubkey (must be before adding accounts)
VAL_PUBKEY=$($BINARY comet show-validator --home $HOME_DIR 2>&1)

# Validate pubkey
if ! echo "$VAL_PUBKEY" | jq -e . > /dev/null 2>&1; then
    echo "Error: Failed to get valid validator pubkey: $VAL_PUBKEY"
    exit 1
fi

# Get addresses
VAL_ADDR=$($BINARY keys show validator -a --keyring-backend $KEYRING --home $HOME_DIR)
VAL_VALOPER=$($BINARY keys show validator --bech val -a --keyring-backend $KEYRING --home $HOME_DIR)

echo "    Validator Address: $VAL_ADDR"
echo "    Validator Operator: $VAL_VALOPER"

echo "[4/8] Adding genesis account..."
$BINARY add-genesis-account $VAL_ADDR 100000000000000000stake,100000000000000000uoasis --home $HOME_DIR

echo "[5/8] Configuring genesis..."

# Disable protorev
jq '.app_state.protorev.params.enabled = false' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set bond denom to stake
jq '.app_state.staking.params.bond_denom = "stake"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Fast unbonding for testing
jq '.app_state.staking.params.unbonding_time = "240s"' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add validator to genesis (VAL_PUBKEY is already JSON)
cat $HOME_DIR/config/genesis.json | jq \
  --arg valop "$VAL_VALOPER" \
  --argjson pubkey "$VAL_PUBKEY" \
  '.app_state.staking.validators = [{
    "operator_address": $valop,
    "consensus_pubkey": $pubkey,
    "jailed": false,
    "status": "BOND_STATUS_BONDED",
    "tokens": "50000000000000000",
    "delegator_shares": "50000000000000000.000000000000000000",
    "description": {
      "moniker": "oasis-validator",
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
    "min_self_delegation": "1",
    "unbonding_on_hold_ref_count": "0",
    "unbonding_ids": []
  }]' > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add delegation
cat $HOME_DIR/config/genesis.json | jq \
  --arg val "$VAL_ADDR" \
  --arg valop "$VAL_VALOPER" \
  '.app_state.staking.delegations = [{
    "delegator_address": $val,
    "validator_address": $valop,
    "shares": "50000000000000000.000000000000000000"
  }]' > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set last validator powers
cat $HOME_DIR/config/genesis.json | jq \
  --arg valop "$VAL_VALOPER" \
  '.app_state.staking.last_validator_powers = [{
    "address": $valop,
    "power": "50000000000000000"
  }]' > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Set last total power
jq '.app_state.staking.last_total_power = "50000000000000000"' \
  $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Add bonded tokens to bonded_tokens_pool module account
# Module address for bonded_tokens_pool with oasis prefix (generated using SDK)
BONDED_POOL_ADDR="oasis1fl48vsnmsdzcv85q5d2q4z5ajdha8yu3j5ytq0"
jq --arg addr "$BONDED_POOL_ADDR" \
  '.app_state.bank.balances += [{
    "address": $addr,
    "coins": [{"denom": "stake", "amount": "50000000000000000"}]
  }] |
  .app_state.bank.supply = [
    {"denom": "stake", "amount": "150000000000000000"},
    {"denom": "uoasis", "amount": "100000000000000000"}
  ]' $HOME_DIR/config/genesis.json > $HOME_DIR/config/genesis.tmp
mv $HOME_DIR/config/genesis.tmp $HOME_DIR/config/genesis.json

# Get consensus address from pubkey
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

echo "[6/8] Configuring node..."

# Fast blocks for testing
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $HOME_DIR/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $HOME_DIR/config/config.toml

# Enable API
sed -i '' 's/enable = false/enable = true/g' $HOME_DIR/config/app.toml
sed -i '' 's/minimum-gas-prices = ""/minimum-gas-prices = "0.0025uoasis"/g' $HOME_DIR/config/app.toml

# Enable unsafe CORS for local development
sed -i '' 's/enabled-unsafe-cors = false/enabled-unsafe-cors = true/g' $HOME_DIR/config/app.toml

echo "[7/8] Validating genesis..."
$BINARY validate-genesis --home $HOME_DIR

echo "[8/8] Setup complete!"
echo ""
echo "=========================================="
echo "✅ Oasis Osmosis LocalNet Ready!"
echo "=========================================="
echo ""
echo "Validator Address:   $VAL_ADDR"
echo "Validator Operator:  $VAL_VALOPER"
echo "Chain ID:            $CHAIN_ID"
echo "Home Directory:      $HOME_DIR"
echo ""
echo "To start the chain:"
echo "  $BINARY start --home $HOME_DIR"
echo ""
echo "Or run in background:"
echo "  $BINARY start --home $HOME_DIR > $HOME_DIR/node.log 2>&1 &"
echo ""
echo "To check status:"
echo "  $BINARY status --node tcp://localhost:26657"
echo ""
echo "To view logs (if running in background):"
echo "  tail -f $HOME_DIR/node.log"
echo ""

