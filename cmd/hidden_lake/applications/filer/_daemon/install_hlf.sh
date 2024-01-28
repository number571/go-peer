#!/bin/bash

# root mode
echo "
[Unit]
Description=HiddenLakeFiler

[Service]
ExecStart=/root/hlf_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_filer.service

cd /root && \
    rm -f hlf_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/latest/download/hlf_amd64_linux && \
    chmod +x hlf_amd64_linux

systemctl daemon-reload
systemctl enable hidden_lake_filer.service
systemctl restart hidden_lake_filer.service
