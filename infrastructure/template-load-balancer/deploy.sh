#!/bin/bash

set -e

DOMAIN="apps-lgtm.zero-one.cloud"
NGINX_UID="65532"

echo "🚀 Deploying Nginx SSL Reverse Proxy"
echo "Domain: $DOMAIN"
echo ""

# Create directories
mkdir -p {conf,ssl,logs,html,webroot}

# Copy nginx configuration
echo "📝 Copying nginx configuration..."
if [ -f "./nginx.conf" ]; then
    cp ./nginx.conf ./conf/nginx.conf
    echo "✅ nginx.conf copied to ./conf/"
else
    echo "❌ nginx.conf not found in current directory"
    exit 1
fi

# Set ownership
sudo chown -R $NGINX_UID:$NGINX_UID logs ssl html webroot conf 2>/dev/null || true

# SSL Setup
echo "🔐 Checking SSL certificate..."
if [ ! -f "ssl/$DOMAIN.crt" ] || [ ! -f "ssl/$DOMAIN.key" ]; then
    echo "📝 Running SSL setup..."
    chmod +x setup-ssl-host.sh
    ./setup-ssl-host.sh
    [ ! -f "ssl/$DOMAIN.crt" ] && { echo "❌ SSL setup failed"; exit 1; }
else
    # Check expiration
    days_until_expiry=$(( ($(date -d "$(openssl x509 -in ssl/$DOMAIN.crt -noout -enddate | cut -d= -f2)" +%s) - $(date +%s)) / 86400 ))
    echo "📅 Certificate expires in $days_until_expiry days"
    if [ $days_until_expiry -lt 30 ]; then
        echo "⚠️  Certificate expires soon! Run ./setup-ssl-host.sh to renew"
    fi
fi

echo "✅ SSL certificate ready"

# Deploy
echo "🐳 Starting nginx..."
docker compose down 2>/dev/null || true
docker compose up -d

sleep 3

# Health check
if ! docker ps | grep -q nginx-unprivileged-optimized; then
    echo "❌ Container failed to start"
    docker logs nginx-unprivileged-optimized --tail 10
    exit 1
fi

# Test with proper Host headers
echo "🔍 Testing services..."
http_code=$(curl -s -o /dev/null -w "%{http_code}" "http://localhost/" 2>/dev/null || echo "000")
https_apps=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: apps-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")
https_portainer=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: portainer-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")
https_grafana=$(curl -s -k -o /dev/null -w "%{http_code}" -H "Host: grafana-lgtm.zero-one.cloud" "https://localhost/" 2>/dev/null || echo "000")

echo "HTTP redirect: $http_code"
echo "HTTPS apps: $https_apps"
echo "HTTPS portainer: $https_portainer"
echo "HTTPS grafana: $https_grafana"

echo ""
echo "🎉 Deployment completed!"
echo "🌐 Main App: https://$DOMAIN"
echo "🐳 Portainer: https://portainer-lgtm.zero-one.cloud"
echo "🌐 Grafana: https://grafana-lgtm.zero-one.cloud"
echo ""
echo "📝 Test commands:"
echo "  curl -k https://apps-lgtm.zero-one.cloud/"
echo "  curl -k https://portainer-lgtm.zero-one.cloud/"
echo "  curl -k https://grafana-lgtm.zero-one.cloud/"
