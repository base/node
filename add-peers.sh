#!/usr/bin/env bash
set -euo pipefail

BOOTNODE_URL="${BOOTNODE_URL:-https://chains.base.org}"
NETWORK=""
EXECUTION_RPC="${EXECUTION_RPC:-http://localhost:8545}"
CONSENSUS_RPC="${CONSENSUS_RPC:-http://localhost:7545}"
EXECUTION_ONLY=false
CONSENSUS_ONLY=false
LOOP_INTERVAL=0

show_usage() {
    cat <<EOF
Usage: $(basename "$0") --network <network> [options]

Fetches public node records from chains.base.org and adds them as peers.

Required:
  --network <name>         Network name (base-mainnet, base-sepolia, base-zeronet)

Options:
  --execution-rpc <url>    Execution layer RPC endpoint (default: http://localhost:8545)
  --consensus-rpc <url>    Consensus layer RPC endpoint (default: http://localhost:7545)
  --bootnode-url <url>     Bootnode server URL (default: https://chains.base.org)
  --execution-only         Only add execution layer peers
  --consensus-only         Only add consensus layer peers
  --loop <seconds>         Poll interval in seconds (0 = run once, default: 0)
  --help                   Show this help message

Prerequisites:
  - curl and jq must be installed
  - Execution client must have admin namespace enabled (--http.api includes admin)
  - Consensus client must have admin RPC enabled (--rpc.enable-admin or BASE_NODE_RPC_ENABLE_ADMIN=true)

Examples:
  $(basename "$0") --network base-sepolia
  $(basename "$0") --network base-sepolia --execution-rpc http://localhost:8545 --loop 300
  $(basename "$0") --network base-mainnet --consensus-only
EOF
    exit 0
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --network) NETWORK="$2"; shift 2 ;;
        --execution-rpc) EXECUTION_RPC="$2"; shift 2 ;;
        --consensus-rpc) CONSENSUS_RPC="$2"; shift 2 ;;
        --bootnode-url) BOOTNODE_URL="$2"; shift 2 ;;
        --execution-only) EXECUTION_ONLY=true; shift ;;
        --consensus-only) CONSENSUS_ONLY=true; shift ;;
        --loop) LOOP_INTERVAL="$2"; shift 2 ;;
        --help) show_usage ;;
        *) echo "Unknown option: $1"; show_usage ;;
    esac
done

if [[ -z "$NETWORK" ]]; then
    echo "Error: --network is required"
    show_usage
fi

for cmd in curl jq; do
    if ! command -v "$cmd" &>/dev/null; then
        echo "Error: $cmd is required but not installed"
        exit 1
    fi
done

add_execution_peers() {
    local nodes
    nodes=$(echo "$1" | jq -r '.execution[]? // empty' 2>/dev/null)
    if [[ -z "$nodes" ]]; then
        echo "[$(date -Iseconds)] No execution peers found for $NETWORK"
        return
    fi

    local count=0
    while IFS= read -r enode; do
        result=$(curl -s -X POST "$EXECUTION_RPC" \
            -H "Content-Type: application/json" \
            -d "{\"jsonrpc\":\"2.0\",\"method\":\"admin_addPeer\",\"params\":[\"$enode\"],\"id\":1}" 2>/dev/null)
        success=$(echo "$result" | jq -r '.result // false' 2>/dev/null)
        if [[ "$success" == "true" ]]; then
            count=$((count + 1))
        else
            echo "[$(date -Iseconds)] Failed to add execution peer: $enode"
        fi
    done <<< "$nodes"

    echo "[$(date -Iseconds)] Added $count execution peers for $NETWORK"
}

add_consensus_peers() {
    local nodes
    nodes=$(echo "$1" | jq -r '.consensus[]? // empty' 2>/dev/null)
    if [[ -z "$nodes" ]]; then
        echo "[$(date -Iseconds)] No consensus peers found for $NETWORK"
        return
    fi

    local count=0
    while IFS= read -r addr; do
        result=$(curl -s -X POST "$CONSENSUS_RPC" \
            -H "Content-Type: application/json" \
            -d "{\"jsonrpc\":\"2.0\",\"method\":\"opp2p_connectPeer\",\"params\":[\"$addr\"],\"id\":1}" 2>/dev/null)
        error=$(echo "$result" | jq -r '.error // empty' 2>/dev/null)
        if [[ -z "$error" ]]; then
            count=$((count + 1))
        else
            echo "[$(date -Iseconds)] Failed to add consensus peer: $addr ($error)"
        fi
    done <<< "$nodes"

    echo "[$(date -Iseconds)] Added $count consensus peers for $NETWORK"
}

run_once() {
    echo "[$(date -Iseconds)] Fetching peers from $BOOTNODE_URL for $NETWORK..."

    response=$(curl -sf "$BOOTNODE_URL/$NETWORK/peers" 2>/dev/null)
    if [[ -z "$response" ]]; then
        echo "[$(date -Iseconds)] Error: failed to fetch from $BOOTNODE_URL/$NETWORK/peers"
        return 1
    fi

    if ! echo "$response" | jq -e '.execution' &>/dev/null; then
        echo "[$(date -Iseconds)] Error: invalid response for network $NETWORK"
        return 1
    fi

    if [[ "$CONSENSUS_ONLY" != "true" ]]; then
        add_execution_peers "$response"
    fi

    if [[ "$EXECUTION_ONLY" != "true" ]]; then
        add_consensus_peers "$response"
    fi
}

if [[ "$LOOP_INTERVAL" -gt 0 ]]; then
    echo "[$(date -Iseconds)] Running in loop mode (interval: ${LOOP_INTERVAL}s)"
    while true; do
        run_once || true
        sleep "$LOOP_INTERVAL"
    done
else
    run_once
fi
