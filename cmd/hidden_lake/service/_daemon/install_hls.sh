# root mode

echo "
[Unit]
Description=HiddenLakeService

[Service]
ExecStart=/root/hls_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_service.service

cd /root && \
    rm -f hls_amd64_linux && \
    wget https://github.com/number571/go-peer/releases/download/v1.5.18/hls_amd64_linux && \
    chmod +x hls_amd64_linux

chmod +x /root/hls_amd64_linux
systemctl daemon-reload
systemctl enable hidden_lake_service.service
systemctl restart hidden_lake_service.service
watch -c SYSTEMD_COLORS=1 systemctl status hidden_lake_service.service
