services:
  execution:
    build:
      context: .
      dockerfile: ${CLIENT:-geth}/Dockerfile
    ports:
      - "8545:8545" # RPC
      - "8546:8546" # websocket
      - "7301:6060" # metrics
      - "30303:30303" # P2P TCP
      - "30303:30303/udp" # P2P UDP
    command: ["bash", "./execution-entrypoint"]
    volumes:
      - ${HOST_DATA_DIR}:/data
    environment:
      - NODE_TYPE=${NODE_TYPE:-vanilla}
    env_file:
      - ${NETWORK_ENV:-.env.mainnet} # Use .env.mainnet by default, override with .env.sepolia for testnet
  node:
    build:
      context: .
      dockerfile: ${CLIENT:-geth}/Dockerfile
    depends_on:
      - execution
    ports:
      - "7545:8545" # RPC
      - "9222:9222" # P2P TCP
      - "9222:9222/udp" # P2P UDP
      - "7300:7300" # metrics
      - "6060:6060" # pprof
    command: ["bash", "./op-node-entrypoint"]
    env_file:
      - ${NETWORK_ENV:-.env.mainnet} # Use .env.mainnet by default, override with .env.sepolia for testnet
