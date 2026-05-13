#!/bin/bash
# Shared helper functions for consensus entrypoints

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
