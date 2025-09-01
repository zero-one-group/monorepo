#!/bin/bash

exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

set -e

# Update and install packages
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install unzip -y

# Add 2GB Swap (optimized and robust)
if ! swapon --show | grep -q '/swapfile' && ! grep -q '/swapfile' /etc/fstab; then
    # Check disk space (need 2.5GB free)
    if [ $(df / | awk 'NR==2 {print $4}') -lt 2621440 ]; then
        echo "❌ Insufficient disk space for swap"
    else
        echo "Setting up 2GB swap..."

        # Create swapfile (fallback to dd if fallocate fails)
        sudo fallocate -l 2G /swapfile || sudo dd if=/dev/zero of=/swapfile bs=1M count=2048

        sudo chmod 600 /swapfile
        sudo mkswap /swapfile
        sudo swapon /swapfile
        echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

        # Optimize for server performance
        sudo sysctl vm.swappiness=10
        echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-swappiness.conf

        echo "✅ Swap configured successfully"
    fi
else
    echo "Swap already configured"
fi

free -h

# Enable required kernel parameters for Docker
cat <<EOF | sudo tee /etc/sysctl.d/99-docker-bridge.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

# Load the br_netfilter module
sudo modprobe br_netfilter

# Apply sysctl parameters without reboot
sudo sysctl -p /etc/sysctl.d/99-docker-bridge.conf

# Install Docker
sudo apt-get install -y ca-certificates curl gnupg
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg
echo \
  "deb [arch="$(dpkg --print-architecture)" signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  "$(. /etc/os-release && echo "$VERSION_CODENAME")" stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update
sudo apt-get install -y docker-ce=5:28.3.3-1~ubuntu.$(lsb_release -rs)~$(lsb_release -cs) docker-ce-cli=5:28.3.3-1~ubuntu.$(lsb_release -rs)~$(lsb_release -cs) containerd.io docker-buildx-plugin docker-compose-plugin

# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure Docker permissions
sudo usermod -aG docker ubuntu
sudo chmod 666 /var/run/docker.sock

# Set custom hostname
sudo hostnamectl set-hostname ${hostname}
