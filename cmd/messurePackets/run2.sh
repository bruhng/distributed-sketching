#!/bin/bash

# Sketch type fixed for this file
SKETCH_TYPE="count-min"

# Values to test
MERGE_RATES=(0.1 0.5 1.0)
CLIENT_AMOUNTS=(10 20 50)
STREAM_RATES=(1000 5000)

# Host list
HOSTS="sketch@10.42.0.2 sketch@10.42.0.3 sketch@10.42.0.4"

for MERGE_RATE in "${MERGE_RATES[@]}"; do
  for CLIENT_AMOUNT in "${CLIENT_AMOUNTS[@]}"; do
    for STREAM_RATE in "${STREAM_RATES[@]}"; do

      echo "Running $SKETCH_TYPE with MERGE_RATE=$MERGE_RATE, CLIENT_AMOUNT=$CLIENT_AMOUNT, STREAM_RATE=$STREAM_RATE"

      pssh -i -t 1000000000 -H "$HOSTS" \
        "source .profile && cd distributed-sketching/cmd/messurePackets && \
        go run . -mergeRate=${MERGE_RATE} -clientAmount=${CLIENT_AMOUNT} -streamRate=${STREAM_RATE} -type=${SKETCH_TYPE}"

      echo "----------------------------------------"
      echo
    done
  done
done

