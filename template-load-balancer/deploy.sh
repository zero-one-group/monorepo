#!/bin/bash

set -e

DOMAIN="{{ apps_domain }}"
NGINX_MONITORING="{{ nginx_monitoring_domain }}"
NGINX_UID="65532"
CONTAINER_NAME="nginx-unprivileged-optimized"

echo "🚀 Deploying Nginx SSL Reverse Proxy"
echo "Domain: $DOMAIN"
echo ""

# Function to calculate file hash
calculate_hash() {
    if [ -f "$1" ]; then
        sha256sum "$1" | cut -d' ' -f1
    else
        echo "missing"
    fi
}

# Function to check if container is running
is_container_running() {
    docker ps --filter "name=$CONTAINER_NAME" --filter "status=running" -q | grep -q .
}

# Create directories
mkdir -p {conf,ssl,logs,html,auth,webroot,naxsi}

# Track if any changes were made
CHANGES_MADE=false

# Copy and check nginx configurations
echo "📝 Checking nginx configuration changes..."

# Check default.conf
if [ -f "./default.conf" ]; then
    current_hash=$(calculate_hash "./conf/default.conf")
    new_hash=$(calculate_hash "./default.conf")

    if [ "$current_hash" != "$new_hash" ]; then
        echo "🔄 default.conf has changes, updating..."
        sudo cp ./default.conf ./conf/default.conf
        CHANGES_MADE=true
    else
        echo "✅ default.conf unchanged"
    fi
else
    echo "❌ default.conf not found in current directory"
    exit 1
fi

# Check nginx.conf
if [ -f "./nginx.conf" ]; then
    current_hash=$(calculate_hash "./conf/nginx.conf")
    new_hash=$(calculate_hash "./nginx.conf")

    if [ "$current_hash" != "$new_hash" ]; then
        echo "🔄 nginx.conf has changes, updating..."
        sudo cp ./nginx.conf ./conf/nginx.conf
        CHANGES_MADE=true
    else
        echo "✅ nginx.conf unchanged"
    fi
else
    echo "❌ nginx.conf not found in current directory"
    exit 1
fi

# Check naxsi_core.rules
if [ -f "./naxsi_core.rules" ]; then
    current_hash=$(calculate_hash "./naxsi/naxsi_core.rules")
    new_hash=$(calculate_hash "./naxsi_core.rules")

    if [ "$current_hash" != "$new_hash" ]; then
        echo "🔄 naxsi_core.rules has changes, updating..."
        sudo cp ./naxsi_core.rules ./naxsi/naxsi_core.rules
        CHANGES_MADE=true
    else
        echo "✅ naxsi_core.rules unchanged"
    fi
else
    echo "❌ naxsi_core.rules not found in current directory"
    exit 1
fi

# Check if .htpasswd already exists
if [ -f "./auth/.htpasswd" ]; then
    echo "✅ Basic auth file already exists, skipping creation"
else
    echo "🔒 Setting up Basic Authentication for $DOMAIN"

    # Check if htpasswd is installed
    if ! command -v htpasswd &> /dev/null; then
        echo "❌ htpasswd not found. Installing apache2-utils..."
        sudo apt-get update && sudo apt-get install -y apache2-utils
    fi

    # Prompt for username and password
    read -p "Enter username for basic auth: " auth_user
    read -s -p "Enter password for basic auth: " auth_pass
    echo ""

    # Create htpasswd file
    htpasswd -b -c ./auth/.htpasswd "$auth_user" "$auth_pass"
    echo "✅ Basic auth credentials created"
    CHANGES_MADE=true
fi

# Set ownership
sudo chown -R $NGINX_UID:$NGINX_UID logs ssl webroot conf html auth naxsi 2>/dev/null || true

# SSL Setup
echo "🔐 Checking SSL certificate..."

# Define possible certificate files to check
DOMAIN_CERT="ssl/$DOMAIN.crt"
DOMAIN_KEY="ssl/$DOMAIN.key"
WILDCARD_CERT="ssl/wildcard.crt"
WILDCARD_KEY="ssl/wildcard.key"
SSL_UPDATED=false

# Check if either domain-specific or wildcard certificates exist
if [[ -f "$DOMAIN_CERT" && -f "$DOMAIN_KEY" ]]; then
    echo "✅ Domain-specific certificates found"
    CERT_FILE="$DOMAIN_CERT"
elif [[ -f "$WILDCARD_CERT" && -f "$WILDCARD_KEY" ]]; then
    echo "✅ Wildcard certificates found"
    CERT_FILE="$WILDCARD_CERT"
else
    echo "📝 No valid certificates found. Running SSL setup..."
    chmod +x setup-ssl-host.sh
    ./setup-ssl-host.sh
    SSL_UPDATED=true
    CHANGES_MADE=true

    # Check again after setup
    if [[ -f "$DOMAIN_CERT" ]]; then
        CERT_FILE="$DOMAIN_CERT"
    elif [[ -f "$WILDCARD_CERT" ]]; then
        CERT_FILE="$WILDCARD_CERT"
    else
        echo "❌ SSL setup failed - no certificate files found"
        exit 1
    fi
fi

# Check expiration
days_until_expiry=$(( ($(date -d "$(openssl x509 -in "$CERT_FILE" -noout -enddate | cut -d= -f2)" +%s) - $(date +%s)) / 86400 ))
echo "📅 Certificate expires in $days_until_expiry days"
if [ $days_until_expiry -lt 30 ]; then
    echo "⚠️  Certificate expires soon! Run ./setup-ssl-host.sh to renew"
fi

echo "✅ SSL certificate ready"

# Check if container is running and if changes were made
if is_container_running; then
    if [ "$CHANGES_MADE" = true ]; then
        echo "🔄 Changes detected, updating container..."
        echo "📦 Reloading nginx configuration..."

        # Try graceful reload first
        if docker exec $CONTAINER_NAME nginx -t 2>/dev/null; then
            echo "✅ Configuration test passed"
            if docker exec $CONTAINER_NAME nginx -s reload 2>/dev/null; then
                echo "✅ Nginx configuration reloaded successfully"
                RESTART_NEEDED=false
            else
                echo "⚠️  Reload failed, container restart required"
                RESTART_NEEDED=true
            fi
        else
            echo "❌ Configuration test failed, container restart required"
            RESTART_NEEDED=true
        fi

        # If reload failed or SSL was updated, restart container
        if [ "$RESTART_NEEDED" = true ] || [ "$SSL_UPDATED" = true ]; then
            echo "🔄 Restarting container..."
            docker compose -f docker-compose-nginx.yml down
            docker compose -f docker-compose-nginx.yml up -d
        fi
    else
        echo "✅ No changes detected, container running normally"
        echo "🔍 Running health checks..."
    fi
else
    echo "🐳 Container not running, starting nginx..."
    docker compose -f docker-compose-nginx.yml down 2>/dev/null || true
    docker compose -f docker-compose-nginx.yml up -d
fi

# Wait for container to be ready
echo "⏳ Waiting for container to be ready..."
sleep 3

# Health check
if ! docker ps | grep -q $CONTAINER_NAME; then
    echo "❌ Container failed to start"
    docker logs $CONTAINER_NAME --tail 10
    exit 1
fi

# Test with proper Host headers
echo "🔍 Testing services..."
http_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/" 2>/dev/null || echo "000")
https_apps=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: apps-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")
https_portainer=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: portainer-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")
https_grafana=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: grafana-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")
https_nginx=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: nginx-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")

echo "HTTP redirect: $http_code"
echo "HTTPS apps: $https_apps"
echo "HTTPS portainer: $https_portainer"
echo "HTTPS grafana: $https_grafana"
echo "HTTPS nginx: $https_nginx"

echo ""
if [ "$CHANGES_MADE" = true ]; then
    echo "🎉 Deployment completed with updates!"
else
    echo "✅ Deployment verified - no changes needed!"
fi
echo "🐳 Main App: https://$DOMAIN"
echo "🐳 Portainer: https://portainer-lgtm.zero-one.cloud"
echo "🌐 Grafana: https://grafana-lgtm.zero-one.cloud"
echo "🌐 Nginx: https://nginx-lgtm.zero-one.cloud"

# Show container status
echo ""
echo "📊 Container Status:"
docker ps --filter "name=$CONTAINER_NAME" --format "table {% raw %}{{.Names}}{% endraw %}\t{% raw %}{{.Status}}{% endraw %}\t{% raw %}{{.Ports}}{% endraw %}"
