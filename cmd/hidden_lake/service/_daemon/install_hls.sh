#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeService

[Service]
ExecStart=/root/hls_amd64_linux -path=/root -key=/root/priv.key
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service.service

cd /root && \
    rm -f hls_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/download/v1.5.21/hls_amd64_linux && \
    chmod +x hls_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_service.service
systemctl restart hidden_lake_service.service
