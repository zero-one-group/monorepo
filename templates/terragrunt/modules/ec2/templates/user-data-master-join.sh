#!/bin/bash

exec > >(tee /var/log/user-data.log | logger -t user-data -s 2>/dev/console) 2>&1

set -e

# Set AWS Region from Terraform
AWS_REGION="${aws_region}"
echo "Initial AWS Region from Terraform: $AWS_REGION"

# Get Region from Instance Metadata with retry
for i in {1..5}; do
    METADATA_REGION=$(curl -s --connect-timeout 5 http://169.254.169.254/latest/meta-data/placement/region)
    if [ ! -z "$METADATA_REGION" ]; then
        echo "AWS Region detected from metadata: $METADATA_REGION"
        AWS_REGION="$METADATA_REGION"
        break
    fi
    echo "Attempt $i: Waiting for AWS region metadata..."
    sleep 2
done

echo "Using AWS Region: $AWS_REGION"

# Set custom hostname
sudo hostnamectl set-hostname ${hostname}

# System Updates
sudo apt-get update -y
sudo apt-get upgrade -y
sudo apt-get install unzip -y

# Add 2GB Swap (bulletproof version)
setup_swap() {
    echo "Setting up 2GB swap..."

    # Remove any existing broken swap
    sudo swapoff /swapfile 2>/dev/null || true
    sudo rm -f /swapfile

    # Check disk space (need at least 3GB free for safety)
    if [ $(df / | awk 'NR==2 {print int($4/1024/1024)}') -lt 3 ]; then
        echo "❌ Need at least 3GB free space, have $(df / | awk 'NR==2 {print int($4/1024/1024)}')GB"
        return 1
    fi

    # Create swap file
    echo "Creating 2GB swapfile..."
    sudo dd if=/dev/zero of=/swapfile bs=1M count=2048 status=progress

    # Verify size
    if [ $(stat -c%s /swapfile | awk '{print int($1/1024/1024)}') -ne 2048 ]; then
        echo "❌ Size mismatch: expected 2048MB, got $(stat -c%s /swapfile | awk '{print int($1/1024/1024)}')MB"
        return 1
    fi

    # Setup swap
    sudo chmod 600 /swapfile
    sudo mkswap /swapfile
    sudo swapon /swapfile

    # Remove old entries and add new one
    sudo sed -i '/\/swapfile/d' /etc/fstab
    echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab

    # Configure performance
    sudo sysctl vm.swappiness=10
    echo 'vm.swappiness=10' | sudo tee /etc/sysctl.d/99-swappiness.conf

    echo "✅ 2GB swap created successfully"
    free -h
}

# Call function (won't exit script if it fails)
setup_swap || echo "❌ Swap setup failed, continuing..."

# Enable required kernel parameters for Docker
cat <<EOF | sudo tee /etc/sysctl.d/99-docker-bridge.conf
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
EOF

# Load the br_netfilter module
sudo modprobe br_netfilter

# Apply sysctl parameters without reboot
sudo sysctl -p /etc/sysctl.d/99-docker-bridge.conf

# Install PostgreSQL Client
sudo sh -c 'echo "deb http://apt.postgresql.org/pub/repos/apt $(lsb_release -cs)-pgdg main" > /etc/apt/sources.list.d/pgdg.list'
wget --quiet -O - https://www.postgresql.org/media/keys/ACCC4CF8.asc | sudo apt-key add -
sudo apt update -y
sudo apt install postgresql-client-17 -y

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

# Get AZ information
AVAILABILITY_ZONE=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)

# Debug information before attempting to join swarm
echo "Current AWS Region before swarm join attempts: $AWS_REGION"
echo "Current AWS identity:"
aws sts get-caller-identity --region "$AWS_REGION"

PROJECT="${project_name}"
echo "Using project name: $PROJECT"

# Get the cluster identifier
CLUSTER_IDENTIFIER="${cluster_identifier}"
echo "Using cluster identifier: $CLUSTER_IDENTIFIER"

# Wait for the master node to be ready
max_attempts=30
attempt=1
while [ $attempt -le $max_attempts ]; do
    echo "Attempt $attempt to retrieve swarm tokens from region: $AWS_REGION"

    if MANAGER_TOKEN=$(aws ssm get-parameter \
        --name "/$PROJECT/swarm/$CLUSTER_IDENTIFIER/manager-token" \
        --with-decryption \
        --query "Parameter.Value" \
        --output text \
        --region "$AWS_REGION") && \
       MASTER_IP=$(aws ssm get-parameter \
        --name "/$PROJECT/swarm/$CLUSTER_IDENTIFIER/master-ip" \
        --query "Parameter.Value" \
        --output text \
        --region "$AWS_REGION"); then

        echo "Successfully retrieved swarm information for cluster: $CLUSTER_IDENTIFIER"
        echo "Master IP: $MASTER_IP"
        echo "Manager Token retrieved successfully (token value hidden for security)"

        # Join the swarm as a manager
        if docker swarm join --token $MANAGER_TOKEN $MASTER_IP; then
            echo "Successfully joined the swarm as a manager for cluster: $CLUSTER_IDENTIFIER"

            # Wait a bit for the node to be recognized in the swarm
            sleep 10

            echo "Swarm join completed successfully"
            exit 0
        else
            echo "Failed to join swarm on attempt $attempt"
        fi
    else
        echo "Failed to retrieve swarm information for cluster: $CLUSTER_IDENTIFIER on attempt $attempt"
        echo "SSM Parameter Store status in region $AWS_REGION:"
        aws ssm describe-parameters --region "$AWS_REGION" || true
    fi

    echo "Waiting for 10 seconds before next attempt..."
    sleep 10
    attempt=$((attempt + 1))
done

echo "Failed to join the swarm after $max_attempts attempts"
exit 1
