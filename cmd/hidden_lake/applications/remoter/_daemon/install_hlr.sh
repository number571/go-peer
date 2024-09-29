#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeRemoter

[Service]
ExecStart=/root/hlr_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_remoter.service

cd /root && \
    rm -f hlr_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/latest/download/hlr_amd64_linux && \
    chmod +x hlr_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_remoter.service
systemctl restart hidden_lake_remoter.service
