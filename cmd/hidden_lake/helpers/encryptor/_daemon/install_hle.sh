#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeEncryptor

[Service]
ExecStart=/root/hll_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_encryptor.service

cd /root && \
    rm -f hll_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/latest/download/hll_amd64_linux && \
    chmod +x hll_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_encryptor.service
systemctl restart hidden_lake_encryptor.service
