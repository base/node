#!/usr/bin/env bash
# Пример запуска ноды Base с настройкой env-переменных

export BASE_P2P_PORT=${BASE_P2P_PORT:-30303}
export BASE_RPC_PORT=${BASE_RPC_PORT:-8545}

echo "Starting Base node on P2P:$BASE_P2P_PORT RPC:$BASE_RPC_PORT..."

./base-node \
  --p2p-port $BASE_P2P_PORT \
  --rpc-port $BASE_RPC_PORT \
  "$@"
