#!/bin/sh
mkdir -p config
wget -q -O config/ipv4.txt https://raw.githubusercontent.com/surfrpt1/clean-ip-scanner/main/config/ipv4.txt
wget -q -O config/ipv6.txt https://raw.githubusercontent.com/surfrpt1/clean-ip-scanner/main/config/ipv6.txt
wget -q -O clean-ip-scanner https://github.com/surfrpt1/clean-ip-scanner/releases/download/v1.1.0/clean-ip-scanner
chmod +x clean-ip-scanner
./clean-ip-scanner -save
