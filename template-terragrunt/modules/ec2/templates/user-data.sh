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

# Install PostgreSQL client
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update -y
sudo apt install postgresql-client-17 -y

# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Set custom hostname
sudo hostnamectl set-hostname ${hostname}

# Setup SSH directory and key for ubuntu user
sudo -u ubuntu mkdir -p /home/ubuntu/.ssh
sudo -u ubuntu chmod 700 /home/ubuntu/.ssh
sudo -u ubuntu touch /home/ubuntu/.ssh/id_rsa
sudo -u ubuntu chmod 600 /home/ubuntu/.ssh/id_rsa

# Write the private key
sudo -u ubuntu cat > /home/ubuntu/.ssh/id_rsa << 'EOL'
${priv_key}
EOL
