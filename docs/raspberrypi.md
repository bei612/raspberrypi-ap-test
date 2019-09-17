### [时间设置]
```
dpkg-reconfigure tzdata
date --set="10 September 2019 11:47:00"
date
```

### [网络配置]
```
echo "interface wlan0" >> /etc/dhcpcd.conf
echo "static ip_address=172.16.100.X/27" >> /etc/dhcpcd.conf
```

### [手工连接wifi]

#### 配置 /etc/wpa_supplicant/wpa_supplicant.conf
```
ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev
update_config=1
```

#### 新增配置：
```
wpa_passphrase [SSID] [PASSWORD] >> /etc/wpa_supplicant/wpa_supplicant.conf
```

#### wifi切换：
```
rm -rf /etc/wpa_supplicant/wpa_supplicant.conf
echo "ctrl_interface=DIR=/var/run/wpa_supplicant GROUP=netdev" > /etc/wpa_supplicant/wpa_supplicant.conf
echo "update_config=1" >> /etc/wpa_supplicant/wpa_supplicant.conf
wpa_passphrase xxxxx yyyyy >> /etc/wpa_supplicant/wpa_supplicant.conf
more /etc/wpa_supplicant/wpa_supplicant.conf
killall wpa_supplicant
rm -rf /var/run/wpa_supplicant/wlan0
systemctl restart wpa_supplicant.service
wpa_supplicant -i wlan0 -c /etc/wpa_supplicant/wpa_supplicant.conf -B
wpa_cli reconfigure -i wlan0
wpa_cli -i wlan0 list_network
```