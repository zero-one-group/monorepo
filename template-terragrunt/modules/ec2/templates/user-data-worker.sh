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

    if WORKER_TOKEN=$(aws ssm get-parameter \
        --name "/$PROJECT/swarm/$CLUSTER_IDENTIFIER/worker-token" \
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
        echo "Worker Token retrieved successfully (token value hidden for security)"

        # Join the swarm
        if docker swarm join --token $WORKER_TOKEN $MASTER_IP; then
            echo "Successfully joined the swarm for cluster: $CLUSTER_IDENTIFIER"

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
