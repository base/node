#!/bin/bash
set -eu

get_public_ip() {
  # Define a list of HTTP-based providers
  local PROVIDERS=(
    "http://ifconfig.me"
    "http://api.ipify.org"
    "http://ipecho.net/plain"
    "http://v4.ident.me"
  )
  # Iterate through the providers until an IP is found or the list is exhausted
  for provider in "${PROVIDERS[@]}"; do
    local IP
    IP=$(curl -s --max-time 10 --connect-timeout 5 "$provider")
    # Check if IP contains a valid format (simple regex for an IPv4 address)
    if [[ $IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      echo "$IP"
      return 0
    fi
  done
  return 1
}

if [[ -z "$OP_NODE_NETWORK" && -z "$OP_NODE_ROLLUP_CONFIG" ]]; then
  echo "expected OP_NODE_NETWORK to be set" 1>&2
  exit 1
fi

# wait until local execution client comes up (authed so will return 401 without token)
until [ "$(curl -s --max-time 10 --connect-timeout 5 -w '%{http_code}' -o /dev/null "${OP_NODE_L2_ENGINE_RPC/ws/http}")" -eq 401 ]; do
  echo "waiting for execution client to be ready"
  sleep 5
done

# public-facing P2P node, advertise public IP address
if PUBLIC_IP=$(get_public_ip); then
  echo "fetched public IP is: $PUBLIC_IP"
else
  echo "Could not retrieve public IP."
  exit 8
fi
export OP_NODE_P2P_ADVERTISE_IP=$PUBLIC_IP

echo "$OP_NODE_L2_ENGINE_AUTH_RAW" > "$OP_NODE_L2_ENGINE_AUTH"

exec ./op-node
