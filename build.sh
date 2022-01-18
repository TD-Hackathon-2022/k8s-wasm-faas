#!/bin/bash

FUNCTION_NAME=${FUNCTION_NAME}
TARGET_HOST=${TARGET_HOST}
TARGET_PORT=${TARGET_PORT}
TARGET_USER=${TARGET_USER}
TARGET_PATH=${TARGET_PATH}

chmod 600 /opt/k8s-faas-builder/private-key
cat /opt/k8s-faas-builder/private-key

GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o "$FUNCTION_NAME" "/opt/k8s-faas-builder/lambda/$FUNCTION_NAME"

mkdir ~/.ssh/
ssh-keyscan -H "$TARGET_HOST" >> ~/.ssh/known_hosts

scp -r -i /opt/k8s-faas-builder/private-key \
  -P "$TARGET_PORT" \
  "$FUNCTION_NAME" \
  "${TARGET_USER}@${TARGET_HOST}:${TARGET_PATH}"
