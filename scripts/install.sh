#!/bin/bash

set -e

echo "[EP] Installing..."

mkdir -p /opt/ep-cluster

cp -r controller /opt/ep-cluster/
cp -r agent /opt/ep-cluster/
cp -r payload /opt/ep-cluster/
cp services/*.service /etc/systemd/system/

echo "[EP] Binaries already built as controller_bin, agent_bin, payload_bin."

systemctl daemon-reload
systemctl enable controller
systemctl enable agent

echo "[EP] Installation done."