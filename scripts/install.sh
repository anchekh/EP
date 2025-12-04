#!/bin/bash

set -e

echo "[EP] Installing..."

mkdir -p /opt/ep-cluster

cp -r controller /opt/ep-cluster/
cp -r agent /opt/ep-cluster/
cp -r payload /opt/ep-cluster/
cp services/*.service /etc/systemd/system/

echo "[EP] Building binaries..."

go build -o /usr/bin/controller ./controller
go build -o /usr/bin/agent ./agent
go build -o /usr/bin/payload ./payload

systemctl daemon-reload

echo "[EP] Done."