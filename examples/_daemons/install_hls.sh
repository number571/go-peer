# root mode

echo "
[Unit]
Description=HiddenLakeService

[Service]
ExecStart=/root/hls_linux -path=/root
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
" > /etc/systemd/system/hidden_lake.service

chdir /root
wget https://github.com/number571/go-peer/releases/download/v1.5.11/hls_linux

chmod +x /root/hls_linux
systemctl daemon-reload
systemctl enable hidden_lake.service
systemctl restart hidden_lake.service
watch systemctl status hidden_lake.service
