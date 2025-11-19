#!/bin/sh
set -e

# Get configuration from environment variables
CLUSTER_ID=${CLUSTER_ID:-0}
REPLICA_ID=${REPLICA_ID:-0}
REPLICA_COUNT=${REPLICA_COUNT:-3}
ADDRESSES=${ADDRESSES:-"tigerbeetle-0:3000,tigerbeetle-1:3001,tigerbeetle-2:3002"}
LISTEN_ADDRESS=${LISTEN_ADDRESS:-"0.0.0.0:3000"}

# Data file path
DATA_FILE="/data/${CLUSTER_ID}_${REPLICA_ID}.tigerbeetle"

# Format the data file if it doesn't exist
if [ ! -f "$DATA_FILE" ]; then
    echo "Replica ${REPLICA_ID}: Data file not found. Formatting new TigerBeetle data file..."
    tigerbeetle format \
        --cluster=$CLUSTER_ID \
        --replica=$REPLICA_ID \
        --replica-count=$REPLICA_COUNT \
        "$DATA_FILE"
    echo "Replica ${REPLICA_ID}: Data file formatted successfully."
else
    echo "Replica ${REPLICA_ID}: Data file already exists. Skipping format."
fi

# Start TigerBeetle
echo "Replica ${REPLICA_ID}: Starting TigerBeetle..."
echo "  Cluster ID: $CLUSTER_ID"
echo "  Replica ID: $REPLICA_ID"
echo "  Addresses: $ADDRESSES"
echo "  Listening on: $LISTEN_ADDRESS"

exec tigerbeetle start --addresses=$ADDRESSES "$DATA_FILE"
