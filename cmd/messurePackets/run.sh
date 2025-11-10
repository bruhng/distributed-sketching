#!/usr/bin/env bash

MERGE_RATE=${1:-1000}
CLIENT_AMOUNT=${2:-10}
STREAM_RATE=${3:-2000}
SKETCH_TYPE=${4:-kll}

pssh -i  -t 1000000000 -H "sketch@10.42.0.2 sketch@10.42.0.3 sketch@10.42.0.4" \
  "source .profile && cd distributed-sketching/cmd/messurePackets  && \
  go run . -mergeRate=${MERGE_RATE} -clientAmount=${CLIENT_AMOUNT} -streamRate=${STREAM_RATE} -type=${SKETCH_TYPE}"
