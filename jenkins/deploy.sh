#!/bin/bash

export WSBROOT="ws_bucket"

screen -S ws_bucket -X quit
echo "starting ws_bucket"
screen -S ws_bucket -d -m bash -c "cd ${WSBROOT} && chmod +x ws_bucket && source ~/.bashrc && ./ws_bucket"
sleep 1
screen -list
