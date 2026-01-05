#!/bin/bash

# Variables
TARGET_DIR="workspace/protohackers/smoke-test"
APP_NAME="server"
REMOTE_HOST="homelab"

echo "Step 1: Compiling Go binary for Linux/amd64..."
GOOS=linux GOARCH=amd64 go build -o $APP_NAME .

echo "Step 2: Ensuring remote directory exists..."
# Using ssh to run mkdir -p (creates parents, no error if exists)
ssh $REMOTE_HOST "mkdir -p $TARGET_DIR"

echo "Step 3: Uploading via SFTP..."
# Using a 'here-document' to pass commands to sftp
sftp $REMOTE_HOST <<EOF
put $APP_NAME $TARGET_DIR/$APP_NAME
chmod 755 $TARGET_DIR/$APP_NAME
bye
EOF

echo "Deployment successful!"