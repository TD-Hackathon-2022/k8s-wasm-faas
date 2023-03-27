#!/bin/bash

FUNCTION_NAME=${FUNCTION_NAME}
TARGET_HOST=${TARGET_HOST}
TARGET_PORT=${TARGET_PORT}
TARGET_USER=${TARGET_USER}
TARGET_PATH=${TARGET_PATH}

source "$HOME/.cargo/env"

mkdir ~/.ssh/
ssh-keyscan -H "$TARGET_HOST" >> ~/.ssh/known_hosts
cp "/opt/k8s-faas-builder/lambda/${FUNCTION_NAME}" /opt/k8s-faas-builder/src/lambda.rs

cargo build --target wasm32-wasi

scp -r -i /opt/k8s-faas-builder/config/ssh-privatekey \
  -P "$TARGET_PORT" \
  /opt/k8s-faas-builder/target/wasm32-wasi/debug/k8s-faas-builder.wasm \
  "${TARGET_USER}@${TARGET_HOST}:${TARGET_PATH}/${FUNCTION_NAME}.wasm"
