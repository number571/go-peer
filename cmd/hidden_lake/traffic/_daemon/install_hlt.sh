# root mode

echo "
[Unit]
Description=HiddenLakeTraffic

[Service]
ExecStart=/root/hlt_amd64_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake_traffic.service

cd /root && wget https://github.com/number571/go-peer/releases/download/v1.5.15/hlt_amd64_linux

chmod +x /root/hlt_amd64_linux
systemctl daemon-reload
systemctl enable hidden_lake_traffic.service
systemctl restart hidden_lake_traffic.service
watch -c SYSTEMD_COLORS=1 systemctl status hidden_lake_traffic.service
