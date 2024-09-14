#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeDistributor

[Service]
ExecStart=/root/hld_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_distributor.service

cd /root && \
    rm -f hld_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/latest/download/hld_amd64_linux && \
    chmod +x hld_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_distributor.service
systemctl restart hidden_lake_distributor.service
