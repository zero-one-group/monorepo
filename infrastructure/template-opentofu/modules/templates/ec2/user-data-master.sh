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
sudo apt-get install -y docker-ce=5:28.3.0-1~ubuntu.$(lsb_release -rs)~$(lsb_release -cs) docker-ce-cli=5:28.3.0-1~ubuntu.$(lsb_release -rs)~$(lsb_release -cs) containerd.io docker-buildx-plugin docker-compose-plugin

# Install AWS CLI
curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
unzip awscliv2.zip
sudo ./aws/install

# Configure Docker permissions
sudo usermod -aG docker ubuntu
sudo chmod 666 /var/run/docker.sock

# Get AZ information
AVAILABILITY_ZONE=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)

# Initialize Docker Swarm
docker swarm init > /tmp/swarm-init.txt 2>&1

# Add proxy label to the node
docker node update --label-add proxy=true $(docker node ls --format "{{.ID}}")
# Add datacenter label based on AZ
docker node update --label-add datacenter=$AVAILABILITY_ZONE $(docker node ls --format "{{.ID}}")

# Create overlay network for swarm
docker network create --driver overlay overlay-network

docker service create \
  --name node-exporter \
  --mode global \
  --mount type=bind,source=/proc,target=/host/proc,readonly=true \
  --mount type=bind,source=/sys,target=/host/sys,readonly=true \
  --mount type=bind,source=/,target=/rootfs,readonly=true \
  --network host \
  prom/node-exporter:latest \
  --path.procfs=/host/proc \
  --path.sysfs=/host/sys \
  --path.rootfs=/rootfs \
  --collector.filesystem.mount-points-exclude="^/(dev|proc|sys|var/lib/docker/.+|var/lib/kubelet/.+)($|/)"

docker service create \
  --name portainer-agent \
  --network overlay-network \
  -p 9001:9001/tcp \
  --mode global \
  --constraint 'node.platform.os == linux' \
  --mount type=bind,src=//var/run/docker.sock,dst=/var/run/docker.sock \
  --mount type=bind,src=//var/lib/docker/volumes,dst=/var/lib/docker/volumes \
  --mount type=bind,src=//,dst=/host \
  portainer/agent:2.21.5

# Extract and save the worker join token
docker swarm join-token worker -q > /tmp/worker-token.txt
docker swarm join-token manager -q > /tmp/manager-token.txt
echo "$(hostname -I | awk '{print $1}'):2377" > /tmp/swarm-master-ip.txt

echo "Current AWS Region before SSM operations: $AWS_REGION"
echo "Current AWS identity:"
aws sts get-caller-identity --region "$AWS_REGION"

# Store in SSM Parameter Store with error checking and logging
PROJECT="${project_name}"
echo "Using project name: $PROJECT"

echo "Attempting to store worker token in SSM with project prefix..."
if aws ssm put-parameter \
    --name "/$PROJECT/swarm/worker-token" \
    --type "SecureString" \
    --value "$(cat /tmp/worker-token.txt)" \
    --overwrite \
    --region "$AWS_REGION"; then
    echo "Successfully stored worker token in SSM with project prefix"
else
    echo "Failed to store worker token in SSM with project prefix"
fi

echo "Attempting to store manager token in SSM with project prefix..."
if aws ssm put-parameter \
    --name "/$PROJECT/swarm/manager-token" \
    --type "SecureString" \
    --value "$(cat /tmp/manager-token.txt)" \
    --overwrite \
    --region "$AWS_REGION"; then
    echo "Successfully stored manager token in SSM with project prefix"
else
    echo "Failed to store manager token in SSM with project prefix"
fi

echo "Attempting to store master IP in SSM with project prefix..."
if aws ssm put-parameter \
    --name "/$PROJECT/swarm/master-ip" \
    --type "String" \
    --value "$(cat /tmp/swarm-master-ip.txt)" \
    --overwrite \
    --region "$AWS_REGION"; then
    echo "Successfully stored master IP in SSM with project prefix"
else
    echo "Failed to store master IP in SSM with project prefix"
fi

# Create a script to handle node labeling
cat > /usr/local/bin/update-node-labels.sh << EOF
#!/bin/bash

# Set AWS Region from parent script
AWS_REGION="${aws_region}"
echo "Initial AWS Region from Terraform in update-node-labels: $AWS_REGION"

# Get Region from Instance Metadata with retry
for i in {1..5}; do
    METADATA_REGION=\$(curl -s --connect-timeout 5 http://169.254.169.254/latest/meta-data/placement/region)
    if [ ! -z "\$METADATA_REGION" ]; then
        echo "AWS Region detected from metadata: \$METADATA_REGION"
        AWS_REGION="\$METADATA_REGION"
        break
    fi
    sleep 2
done

echo "Using AWS Region in update-node-labels: \$AWS_REGION"

# Function to get AZ from instance IP
get_az_from_ip() {
    local private_ip=\$1
    local instance_id=\$(curl -s "http://169.254.169.254/latest/meta-data/instance-id")
    local my_private_ip=\$(curl -s "http://169.254.169.254/latest/meta-data/local-ipv4")
    local my_az=\$(curl -s "http://169.254.169.254/latest/meta-data/placement/availability-zone")

    if [ "\$private_ip" == "\$my_private_ip" ]; then
        echo "\$my_az"
        return
    fi

    # For other nodes, try to get AZ from EC2 API
    az=\$(aws ec2 describe-instances --filters "Name=private-ip-address,Values=\$private_ip" --query 'Reservations[0].Instances[0].Placement.AvailabilityZone' --output text --region "\$AWS_REGION" 2>/dev/null)
    if [ "\$az" == "None" ] || [ -z "\$az" ]; then
        echo "\$my_az"
    else
        echo "\$az"
    fi
}

# Function to check if a node is the monitoring node
is_monitoring_node() {
    local node_id=\$1
    local hostname=\$(docker node inspect \$node_id --format '{{.Description.Hostname}}')
    if [[ "\$hostname" == monitoring* ]]; then
        return 0  # true
    else
        return 1  # false
    fi
}

is_database_node() {
    local node_id=\$1
    local hostname=\$(docker node inspect \$node_id --format '{{.Description.Hostname}}')
    if [[ "\$hostname" == *db* ]]; then
        return 0  # true
    else
        return 1  # false
    fi
}

# Function to label nodes
label_nodes() {
    for NODE_ID in \$(docker node ls --format '{{.ID}}'); do
        if is_monitoring_node "\$NODE_ID"; then
            # Apply monitoring label for monitoring node
            docker node update --label-add node_type=monitoring \$NODE_ID
            echo "Applied monitoring label to node \$NODE_ID"
        elif is_database_node "\$NODE_ID"; then
            # Apply database label for database node
            docker node update --label-add node_type=database \$NODE_ID
            echo "Applied database label to node \$NODE_ID"
        else
            # Apply AZ-based label for other nodes
            NODE_IP=\$(docker node inspect \$NODE_ID --format '{{.Status.Addr}}')
            NODE_AZ=\$(get_az_from_ip "\$NODE_IP")

            if [ ! -z "\$NODE_AZ" ] && [ "\$NODE_AZ" != "None" ]; then
                docker node update --label-add datacenter=\$NODE_AZ \$NODE_ID
                echo "Applied datacenter label \$NODE_AZ to node \$NODE_ID"
            fi
        fi
    done
}

# Run initially
label_nodes
EOF

# Set proper ownership and permissions
sudo chown root:root /usr/local/bin/update-node-labels.sh
sudo chmod 700 /usr/local/bin/update-node-labels.sh

# Create the cron job file
sudo tee /etc/cron.d/docker-labels << 'EOF'
SHELL=/bin/bash
PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin

# Run every minute
* * * * * root /usr/local/bin/update-node-labels.sh >> /var/log/node-labels.log 2>&1
EOF

# Set proper permissions for cron file and add newline
sudo chmod 644 /etc/cron.d/docker-labels
echo "" | sudo tee -a /etc/cron.d/docker-labels

# Create log file with proper permissions
sudo touch /var/log/node-labels.log
sudo chmod 644 /var/log/node-labels.log

# Run the script immediately
sudo /usr/local/bin/update-node-labels.sh
