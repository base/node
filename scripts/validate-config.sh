#!/bin/bash
# Config validator for base/node
# Validates environment variables, detects deprecated/unknown vars, and provides actionable suggestions

set -eu

# Colors for output
RED='\033[0;31m'
YELLOW='\033[1;33m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

# If running from Docker, adjust paths
if [[ -d "/workspace" ]]; then
    REPO_ROOT="/workspace"
fi

DEPRECATIONS_FILE="${REPO_ROOT}/config/deprecations.json"

# Default env file
ENV_FILE="${1:-${NETWORK_ENV:-.env.mainnet}}"

# Check if env file exists
if [[ ! -f "$ENV_FILE" ]]; then
    echo -e "${RED}ERROR:${NC} Environment file not found: $ENV_FILE" >&2
    echo "Usage: $0 [env-file]" >&2
    exit 1
fi

# Load deprecations if file exists
declare -A DEPRECATIONS
if [[ -f "$DEPRECATIONS_FILE" ]]; then
    # Try to use jq if available, otherwise parse JSON manually
    if command -v jq &> /dev/null; then
        while IFS='=' read -r old new; do
            [[ "$old" =~ ^#.*$ ]] && continue
            [[ -z "$old" ]] && continue
            [[ "$old" =~ ^_comment ]] && continue
            if [[ "$new" != "null" ]]; then
                DEPRECATIONS["$old"]="$new"
            fi
        done < <(jq -r 'to_entries[] | select(.key | startswith("_") | not) | "\(.key)=\(.value)"' "$DEPRECATIONS_FILE" 2>/dev/null || true)
    else
        # Fallback: simple JSON parsing (handles basic key-value pairs)
        while IFS= read -r line; do
            # Skip comments and empty lines
            [[ "$line" =~ ^[[:space:]]*\"_ ]] && continue
            [[ "$line" =~ ^[[:space:]]*// ]] && continue
            [[ -z "$line" ]] && continue
            # Extract key-value pairs (simple regex)
            if [[ "$line" =~ \"([^\"]+)\"[[:space:]]*:[[:space:]]*\"([^\"]+)\" ]]; then
                old="${BASH_REMATCH[1]}"
                new="${BASH_REMATCH[2]}"
                if [[ "$new" != "null" ]]; then
                    DEPRECATIONS["$old"]="$new"
                fi
            fi
        done < "$DEPRECATIONS_FILE"
    fi
fi

# Load environment variables from file
declare -A ENV_VARS
while IFS='=' read -r key value || [[ -n "$key" ]]; do
    # Skip comments and empty lines
    [[ "$key" =~ ^[[:space:]]*# ]] && continue
    [[ -z "$key" ]] && continue
    # Remove leading/trailing whitespace
    key=$(echo "$key" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    # Extract value, handling quotes and inline comments
    value=$(echo "$value" | sed 's/^[[:space:]]*//;s/[[:space:]]*$//')
    # Remove inline comments (everything after #)
    value=$(echo "$value" | sed 's/#.*$//')
    # Remove surrounding quotes - multiple passes to handle all cases
    value=$(echo "$value" | sed -e 's/^"//' -e 's/"$//' -e "s/^'//" -e "s/'$//")
    # Remove any remaining quotes (in case of malformed input)
    value=$(echo "$value" | sed -e 's/"//g' -e "s/'//g")
    # Remove trailing whitespace
    value=$(echo "$value" | sed 's/[[:space:]]*$//')
    ENV_VARS["$key"]="$value"
done < "$ENV_FILE"

# Detect client type (check environment variable first, then env file)
if [[ -z "${CLIENT:-}" ]] && [[ -n "${ENV_VARS[CLIENT]:-}" ]]; then
    CLIENT="${ENV_VARS[CLIENT]}"
fi
CLIENT="${CLIENT:-geth}"

# Define required variables (common to all clients)
REQUIRED_COMMON=(
    "OP_NODE_NETWORK"
    "OP_NODE_L2_ENGINE_AUTH_RAW"
    "OP_NODE_L1_ETH_RPC"
    "OP_NODE_L1_BEACON"
    "OP_NODE_L1_BEACON_ARCHIVER"
)

# Client-specific required variables
declare -A REQUIRED_CLIENT
REQUIRED_CLIENT["reth"]="RETH_CHAIN RETH_SEQUENCER_HTTP"
REQUIRED_CLIENT["geth"]="OP_GETH_SEQUENCER_HTTP"
REQUIRED_CLIENT["nethermind"]="OP_SEQUENCER_HTTP"

# Known optional variables
KNOWN_OPTIONAL=(
    "OP_NODE_L1_RPC_KIND"
    "OP_NODE_ROLLUP_CONFIG"
    "OP_NODE_L2_ENGINE_RPC"
    "OP_NODE_P2P_ADVERTISE_IP"
    "OP_NODE_INTERNAL_IP"
    "OP_NODE_LOG_FORMAT"
    "OP_NODE_LOG_LEVEL"
    "OP_NODE_METRICS_ADDR"
    "OP_NODE_METRICS_ENABLED"
    "OP_NODE_METRICS_PORT"
    "OP_NODE_ROLLUP_LOAD_PROTOCOL_VERSIONS"
    "OP_NODE_RPC_ADDR"
    "OP_NODE_RPC_PORT"
    "OP_NODE_SNAPSHOT_LOG"
    "OP_NODE_SYNCMODE"
    "OP_NODE_VERIFIER_L1_CONFS"
    "OP_NODE_P2P_LISTEN_IP"
    "OP_NODE_P2P_LISTEN_TCP_PORT"
    "OP_NODE_P2P_LISTEN_UDP_PORT"
    "OP_NODE_P2P_BOOTNODES"
    "OP_NODE_P2P_AGENT"
    "OP_NODE_L1_TRUST_RPC"
    "OP_NODE_L1_BEACON_FETCH_ALL_SIDECARS"
    "OP_NODE_L2_ENGINE_KIND"
    "OP_NODE_L2_ENGINE_AUTH"
    "OP_SEQUENCER_HTTP"
    "RETH_CHAIN"
    "RETH_SEQUENCER_HTTP"
    "HOST_DATA_DIR"
    "NETWORK_ENV"
    "CLIENT"
    "RPC_PORT"
    "WS_PORT"
    "AUTHRPC_PORT"
    "METRICS_PORT"
    "P2P_PORT"
    "DISCOVERY_PORT"
    "GETH_VERBOSITY"
    "GETH_DATA_DIR"
    "GETH_CACHE"
    "GETH_CACHE_DATABASE"
    "GETH_CACHE_GC"
    "GETH_CACHE_SNAPSHOT"
    "GETH_CACHE_TRIE"
    "OP_GETH_GCMODE"
    "OP_GETH_SYNCMODE"
    "OP_GETH_ETH_STATS"
    "OP_GETH_ALLOW_UNPROTECTED_TXS"
    "OP_GETH_STATE_SCHEME"
    "OP_GETH_BOOTNODES"
    "OP_GETH_NET_RESTRICT"
    "OP_GETH_OP_NETWORK"
    "RETH_DATA_DIR"
    "RETH_FB_WEBSOCKET_URL"
    "RETH_PRUNING_ARGS"
    "OP_RETH_DISABLE_DISCOVERY"
    "OP_RETH_DISABLE_TX_POOL_GOSSIP"
    "OP_RETH_OP_NETWORK"
    "OP_RETH_SEQUENCER_HTTP"
    "NETHERMIND_DATA_DIR"
    "NETHERMIND_LOG_LEVEL"
    "OP_NETHERMIND_BOOTNODES"
    "OP_NETHERMIND_ETHSTATS_ENABLED"
    "OP_NETHERMIND_ETHSTATS_ENDPOINT"
    "OP_NETHERMIND_ETHSTATS_NODE_NAME"
    "HOST_IP"
    "STATSD_ADDRESS"
)

# Validation results
ERRORS=()
WARNINGS=()
INFO=()

# Helper function to find similar variable names
find_similar() {
    local var="$1"
    local candidates=("${REQUIRED_COMMON[@]}" "${KNOWN_OPTIONAL[@]}")
    local best_match=""
    local min_distance=999
    
    for candidate in "${candidates[@]}"; do
        # Simple Levenshtein-like distance (approximate)
        local distance=0
        if [[ "$var" == *"${candidate}"* ]] || [[ "${candidate}" == *"${var}"* ]]; then
            distance=1
        elif [[ "${var,,}" == "${candidate,,}" ]]; then
            distance=0
        else
            # Count character differences (simplified)
            distance=$(echo "$var" "$candidate" | awk '{
                len1=length($1); len2=length($2)
                if (len1 > len2) diff=len1-len2
                else diff=len2-len1
                print diff
            }')
        fi
        
        if [[ $distance -lt $min_distance ]]; then
            min_distance=$distance
            best_match="$candidate"
        fi
    done
    
    if [[ $min_distance -lt 3 && -n "$best_match" ]]; then
        echo "$best_match"
    fi
}

# Validate required variables
validate_required() {
    local missing=()
    
    # Check common required vars
    for var in "${REQUIRED_COMMON[@]}"; do
        if [[ -z "${ENV_VARS[$var]:-}" ]]; then
            missing+=("$var")
        fi
    done
    
    # Check client-specific required vars
    local client_req="${REQUIRED_CLIENT[$CLIENT]:-}"
    if [[ -n "$client_req" ]]; then
        for var in $client_req; do
            if [[ -z "${ENV_VARS[$var]:-}" ]]; then
                missing+=("$var")
            fi
        done
    fi
    
    if [[ ${#missing[@]} -gt 0 ]]; then
        ERRORS+=("Missing required variables:")
        for var in "${missing[@]}"; do
            ERRORS+=("  - $var")
        done
    fi
}

# Check for deprecated variables
check_deprecated() {
    for var in "${!ENV_VARS[@]}"; do
        if [[ -n "${DEPRECATIONS[$var]:-}" ]]; then
            WARNINGS+=("Deprecated variable '$var' found. Use '${DEPRECATIONS[$var]}' instead.")
        fi
    done
}

# Check for unknown variables
check_unknown() {
    local all_known=("${REQUIRED_COMMON[@]}" "${KNOWN_OPTIONAL[@]}")
    local client_req="${REQUIRED_CLIENT[$CLIENT]:-}"
    for var in $client_req; do
        all_known+=("$var")
    done
    
    for var in "${!ENV_VARS[@]}"; do
        local is_known=false
        for known in "${all_known[@]}"; do
            if [[ "$var" == "$known" ]]; then
                is_known=true
                break
            fi
        done
        
        # Check if it's a deprecated variable
        if [[ -n "${DEPRECATIONS[$var]:-}" ]]; then
            is_known=true
        fi
        
        if [[ "$is_known" == false ]]; then
            local similar=$(find_similar "$var")
            if [[ -n "$similar" ]]; then
                WARNINGS+=("Unknown variable '$var' found. Did you mean '$similar'?")
            else
                WARNINGS+=("Unknown variable '$var' found. This may be a typo or unused variable.")
            fi
        fi
    done
}

# Validate URL format
validate_url() {
    local var="$1"
    local value="${ENV_VARS[$var]:-}"
    
    if [[ -z "$value" ]]; then
        return
    fi
    
    if [[ ! "$value" =~ ^https?:// ]] && [[ ! "$value" =~ ^ws:// ]] && [[ ! "$value" =~ ^wss:// ]]; then
        ERRORS+=("Invalid URL format for '$var': '$value' (must start with http://, https://, ws://, or wss://)")
    fi
}

# Validate port format
validate_port() {
    local var="$1"
    local value="${ENV_VARS[$var]:-}"
    
    if [[ -z "$value" ]]; then
        return
    fi
    
    if [[ ! "$value" =~ ^[0-9]+$ ]] || [[ "$value" -lt 1 ]] || [[ "$value" -gt 65535 ]]; then
        ERRORS+=("Invalid port number for '$var': '$value' (must be between 1 and 65535)")
    fi
}

# Validate numeric format
validate_numeric() {
    local var="$1"
    local value="${ENV_VARS[$var]:-}"
    
    if [[ -z "$value" ]]; then
        return
    fi
    
    if [[ ! "$value" =~ ^[0-9]+$ ]]; then
        WARNINGS+=("Non-numeric value for '$var': '$value' (expected numeric)")
    fi
}

# Validate specific variable formats
validate_formats() {
    # Validate URLs
    validate_url "OP_NODE_L1_ETH_RPC"
    validate_url "OP_NODE_L1_BEACON"
    validate_url "OP_NODE_L1_BEACON_ARCHIVER"
    validate_url "OP_NODE_L2_ENGINE_RPC"
    validate_url "OP_GETH_SEQUENCER_HTTP"
    validate_url "RETH_SEQUENCER_HTTP"
    validate_url "OP_SEQUENCER_HTTP"
    validate_url "RETH_FB_WEBSOCKET_URL"
    
    # Validate ports
    validate_port "RPC_PORT"
    validate_port "WS_PORT"
    validate_port "AUTHRPC_PORT"
    validate_port "METRICS_PORT"
    validate_port "P2P_PORT"
    validate_port "DISCOVERY_PORT"
    
    # Validate numeric values
    validate_numeric "GETH_CACHE"
    validate_numeric "GETH_CACHE_DATABASE"
    validate_numeric "GETH_CACHE_GC"
    validate_numeric "GETH_CACHE_SNAPSHOT"
    validate_numeric "GETH_CACHE_TRIE"
    validate_numeric "GETH_VERBOSITY"
}

# Validate RPC_KIND values
validate_rpc_kind() {
    local value="${ENV_VARS[OP_NODE_L1_RPC_KIND]:-}"
    if [[ -z "$value" ]]; then
        return
    fi
    
    local valid_kinds=("alchemy" "quicknode" "infura" "parity" "nethermind" "debug_geth" "erigon" "basic" "any" "standard")
    local is_valid=false
    for kind in "${valid_kinds[@]}"; do
        if [[ "$value" == "$kind" ]]; then
            is_valid=true
            break
        fi
    done
    
    if [[ "$is_valid" == false ]]; then
        WARNINGS+=("Invalid OP_NODE_L1_RPC_KIND value: '$value'. Valid values: ${valid_kinds[*]}")
    fi
}

# Main validation
main() {
    echo -e "${BLUE}Validating configuration file: $ENV_FILE${NC}"
    echo -e "${BLUE}Detected client: $CLIENT${NC}"
    echo ""
    
    validate_required
    check_deprecated
    check_unknown
    validate_formats
    validate_rpc_kind
    
    # Print results
    local has_errors=false
    local has_warnings=false
    
    if [[ ${#ERRORS[@]} -gt 0 ]]; then
        has_errors=true
        echo -e "${RED}❌ ERRORS:${NC}"
        for error in "${ERRORS[@]}"; do
            echo -e "${RED}  $error${NC}"
        done
        echo ""
    fi
    
    if [[ ${#WARNINGS[@]} -gt 0 ]]; then
        has_warnings=true
        echo -e "${YELLOW}⚠️  WARNINGS:${NC}"
        for warning in "${WARNINGS[@]}"; do
            echo -e "${YELLOW}  $warning${NC}"
        done
        echo ""
    fi
    
    if [[ "$has_errors" == false ]] && [[ "$has_warnings" == false ]]; then
        echo -e "${GREEN}✅ Configuration is valid!${NC}"
        exit 0
    elif [[ "$has_errors" == true ]]; then
        echo -e "${RED}Validation failed. Please fix the errors above.${NC}"
        exit 1
    else
        echo -e "${YELLOW}Validation completed with warnings. Please review the warnings above.${NC}"
        exit 0
    fi
}

# Run main function
main

