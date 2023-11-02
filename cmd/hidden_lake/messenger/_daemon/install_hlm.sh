#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeMessenger

[Service]
ExecStart=/root/hlm_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_messenger.service

cd /root && \
    rm -f hlm_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/download/v1.5.19/hlm_amd64_linux && \
    chmod +x hlm_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_messenger.service
systemctl restart hidden_lake_messenger.service
