#!/usr/bin/env bash

# Sketch type fixed for this file
SKETCH_TYPE="kll"

# Values to test
MERGE_RATES=(1000 10000)
CLIENT_AMOUNTS=(1 10 50 100)
STREAM_RATES=(1000 3334222 6667444 10000666 13333888 16667110 20000332 23333554 26666776 30000000)

# STREAM_START=1000
# STREAM_MAX=30000000
# STREAM_MULTIPLIER=1.05

# Host list
HOSTS="sketch@10.42.0.2 sketch@10.42.0.3 sketch@10.42.0.4"

# # Generate stream rate values dynamically
# stream_rates=()
# current_rate=$STREAM_START
# while (( $(echo "$current_rate <= $STREAM_MAX" | bc -l) )); do
#   stream_rates+=($(printf "%.0f" "$current_rate"))
#   current_rate=$(echo "$current_rate * $STREAM_MULTIPLIER" | bc -l)
# done


for MERGE_RATE in "${MERGE_RATES[@]}"; do
  for CLIENT_AMOUNT in "${CLIENT_AMOUNTS[@]}"; do
    for STREAM_RATE in "${STREAM_RATES[@]}"; do

      echo "Running $SKETCH_TYPE with MERGE_RATE=$MERGE_RATE, CLIENT_AMOUNT=$CLIENT_AMOUNT, STREAM_RATE=$STREAM_RATE"

      pssh -i -t 1000000000 -H "$HOSTS" \
        "source .profile && cd distributed-sketching/cmd/messurePackets && \
        go run . -mergeRate=${MERGE_RATE} -clientAmount=${CLIENT_AMOUNT} -streamRate=${STREAM_RATE} -type=${SKETCH_TYPE}"

      echo "----------------------------------------"
      echo
      grpcurl -plaintext -d '{"numMsg": 1}' -proto ../../proto/sketch.proto localhost:8080 proto.Sketcher/RestartServer
      sleep 1
    done
  done
done

