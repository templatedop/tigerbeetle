#!/bin/sh
set -e

# Data file path
DATA_FILE="/data/0_0.tigerbeetle"
CLUSTER_ID=${CLUSTER_ID:-0}
REPLICA_ID=${REPLICA_ID:-0}
ADDRESSES=${ADDRESSES:-0.0.0.0:3000}

# Format the data file if it doesn't exist
if [ ! -f "$DATA_FILE" ]; then
    echo "Data file not found. Formatting new TigerBeetle data file..."
    tigerbeetle format --cluster=$CLUSTER_ID --replica=$REPLICA_ID --replica-count=1 "$DATA_FILE"
    echo "Data file formatted successfully."
else
    echo "Data file already exists. Skipping format."
fi

# Start TigerBeetle
echo "Starting TigerBeetle..."
exec tigerbeetle start --addresses=$ADDRESSES "$DATA_FILE"
