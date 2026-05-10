#!/bin/bash
# Shared helper functions used by Base node entrypoints

# get_public_ip attempts to discover the node's public IPv4 address
# by querying a list of HTTP-based IP detection services.
# Prints the IP to stdout and returns 0 on success, 1 on failure.
get_public_ip() {
  local PROVIDERS=(
    "http://ifconfig.me"
    "http://api.ipify.org"
    "http://ipecho.net/plain"
    "http://v4.ident.me"
  )

  for provider in "${PROVIDERS[@]}"; do
    local IP
    IP=$(curl -s --max-time 10 --connect-timeout 5 "$provider")
    if [[ $IP =~ ^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
      echo "$IP"
      return 0
    fi
  done
  return 1
}
