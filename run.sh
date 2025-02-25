#!/bin/bash

SESSION_NAME="go_project"
NUM_CLIENTS=${1:-5}  # Default to 5 clients if not provided

tmux new-session -d -s $SESSION_NAME -n server

tmux send-keys -t $SESSION_NAME:server "go run ." C-m

for ((i=1; i<=NUM_CLIENTS; i++)); do
    tmux new-window -t $SESSION_NAME -n "client_$i"
    tmux send-keys -t $SESSION_NAME:client_$i "go run . -client" C-m

done

tmux attach-session -t $SESSION_NAME
