### 安装 prometheus
https://github.com/prometheus/prometheus/releases/tag/v2.12.0
```
mkdir -p /usr/local/prometheus /var/lib/prometheus/data
cat > /etc/systemd/system/prometheus.service <<EOF
[Unit]
Description=prometheus
After=network.target
[Service]
Type=simple
ExecStart=/usr/local/prometheus/prometheus --config.file=/usr/local/prometheus/prometheus.yml --storage.tsdb.path=/var/lib/prometheus/data 
Restart=on-failure
RestartSec=5
StartLimitInterval=3600s
StartLimitBurst=3
[Install]
WantedBy=multi-user.target
EOF

systemctl reset-failed
systemctl daemon-reload
systemctl start prometheus
systemctl status prometheus
systemctl enable prometheus
```

### 安装 grafana
```
apt update
apt autoremove -y
dpkg -i grafana_6.3.5_armhf.deb
apt --fix-broken install
systemctl daemon-reload
systemctl enable grafana-server
systemctl start grafana-server
systemctl status grafana-server
netstat -ntlp
```