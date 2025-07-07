#!/bin/bash

DOMAIN="{{ apps_domain }}"
WILDCARD_DOMAIN="{{ wildcard_domain }}"
SSL_DIR="./ssl"
NGINX_UID="65532"

echo "🔐 Setting up SSL for $DOMAIN"
echo "================================================"

# Create SSL directory
mkdir -p $SSL_DIR

# Function to check if nginx is running
check_nginx_running() {
    docker ps | grep -q nginx-unprivileged-optimized
}

# Option 1: Self-signed certificate
create_self_signed() {
    echo "📝 Creating self-signed certificate..."
    openssl req -x509 -nodes -days 365 -newkey rsa:2048 \
        -keyout $SSL_DIR/$DOMAIN.key \
        -out $SSL_DIR/$DOMAIN.crt \
        -subj "/C=US/ST=State/L=City/O=Organization/CN=$DOMAIN" \
        -addext "subjectAltName=DNS:$DOMAIN,DNS:$WILDCARD_DOMAIN,DNS:zero-one.cloud"
    echo "✅ Self-signed certificate created"
}

# Option 2: Let's Encrypt standalone
create_letsencrypt_standalone() {
    echo "🌐 Creating Let's Encrypt certificate..."

    local nginx_was_running=false
    if check_nginx_running; then
        nginx_was_running=true
        echo "⏸️  Stopping nginx temporarily..."
        docker compose down
        sleep 3
    fi

    # Install certbot if needed
    if ! command -v certbot &> /dev/null; then
        echo "Installing certbot..."
        sudo apt-get update && sudo apt-get install -y certbot
    fi

    # Generate certificate
    sudo certbot certonly --standalone \
        --preferred-challenges http \
        --http-01-port 80 \
        -d $DOMAIN \
        --email {{ email_generate_certbot }} \
        --agree-tos \
        --non-interactive

    if [ $? -eq 0 ]; then
        sudo cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem $SSL_DIR/$DOMAIN.crt
        sudo cp /etc/letsencrypt/live/$DOMAIN/privkey.pem $SSL_DIR/$DOMAIN.key
        echo "✅ Let's Encrypt certificate created"
    else
        echo "❌ Certificate generation failed"
        [ "$nginx_was_running" = true ] && docker compose up -d
        return 1
    fi

    # Restart nginx if it was running
    [ "$nginx_was_running" = true ] && docker compose up -d
}

# Option 3: Let's Encrypt webroot
create_letsencrypt_webroot() {
    echo "🌐 Creating Let's Encrypt certificate (webroot)..."

    if ! check_nginx_running; then
        echo "❌ Nginx not running. Start nginx first or use standalone method."
        return 1
    fi

    # Install certbot if needed
    if ! command -v certbot &> /dev/null; then
        sudo apt-get update && sudo apt-get install -y certbot
    fi

    # Create webroot directory
    mkdir -p ./webroot/.well-known/acme-challenge
    sudo chown -R $NGINX_UID:$NGINX_UID ./webroot

    # Generate certificate
    sudo certbot certonly --webroot \
        --webroot-path ./webroot \
        -d $DOMAIN \
        --email {{ email_generate_certbot }} \
        --agree-tos \
        --non-interactive

    if [ $? -eq 0 ]; then
        sudo cp /etc/letsencrypt/live/$DOMAIN/fullchain.pem $SSL_DIR/$DOMAIN.crt
        sudo cp /etc/letsencrypt/live/$DOMAIN/privkey.pem $SSL_DIR/$DOMAIN.key
        docker compose exec nginx nginx -s reload
        echo "✅ Certificate created and nginx reloaded"
    else
        echo "❌ Certificate generation failed"
        return 1
    fi
}

# Option 4: DNS Challenge for wildcard ONLY
create_letsencrypt_dns() {
    echo "🌐 Creating Let's Encrypt wildcard certificate (DNS challenge)..."
    echo "⚠️  This will create a wildcard certificate for $WILDCARD_DOMAIN ONLY"
    echo "    (This covers all *.zero-one.cloud subdomains including $DOMAIN)"

    # Install certbot if needed
    if ! command -v certbot &> /dev/null; then
        echo "Installing certbot..."
        sudo apt-get update && sudo apt-get install -y certbot
    fi

    echo ""
    echo "DNS Challenge options:"
    echo "1) Manual (you add TXT record manually)"
    echo "2) Cloudflare API (automated)"
    read -p "Choose [1-2]: " dns_choice

    local certbot_exit_code=0

    case $dns_choice in
        1)
            echo "📝 Manual DNS challenge for wildcard domain..."
            echo "🎯 Requesting certificate for: $WILDCARD_DOMAIN"
            sudo certbot certonly --manual \
                --preferred-challenges dns \
                -d $WILDCARD_DOMAIN \
                --email {{ email_generate_certbot }} \
                --agree-tos \
                --manual-public-ip-logging-ok

            certbot_exit_code=$?
            ;;
        2)
            read -p "Enter Cloudflare API Token: " cf_token
            if [ -z "$cf_token" ]; then
                echo "❌ API token required"
                return 1
            fi

            # Install cloudflare plugin
            if ! command -v pip3 &> /dev/null; then
                sudo apt-get install -y python3-pip
            fi
            pip3 install certbot-dns-cloudflare 2>/dev/null || sudo apt-get install -y python3-certbot-dns-cloudflare

            # Create temp credentials
            CF_CREDS="/tmp/cf-creds-$$.ini"
            echo "dns_cloudflare_api_token = $cf_token" > $CF_CREDS
            chmod 600 $CF_CREDS

            echo "📝 Automated DNS challenge for wildcard domain..."
            echo "🎯 Requesting certificate for: $WILDCARD_DOMAIN"
            sudo certbot certonly --dns-cloudflare \
                --dns-cloudflare-credentials $CF_CREDS \
                --dns-cloudflare-propagation-seconds 60 \
                -d $WILDCARD_DOMAIN \
                --email {{ email_generate_certbot }} \
                --agree-tos \
                --non-interactive

            certbot_exit_code=$?
            rm -f $CF_CREDS
            ;;
        *)
            echo "Invalid choice"
            return 1
            ;;
    esac

    # Check if certbot succeeded
    if [ $certbot_exit_code -eq 0 ]; then
        echo "✅ Certbot completed successfully"

        # Find the certificate directory
        CERT_DIR=$(sudo find /etc/letsencrypt/live/ -name "{{ wildcard_domain }}" -type d | head -n1)

        if [ -n "$CERT_DIR" ] && [ -d "$CERT_DIR" ]; then
            echo "📂 Found certificate directory: $CERT_DIR"

            # Copy wildcard certificate to wildcard-named files
            echo "➡️  Copying certificate to $SSL_DIR/wildcard.crt and $SSL_DIR/wildcard.key"
            sudo cp "$CERT_DIR/fullchain.pem" $SSL_DIR/wildcard.crt
            sudo cp "$CERT_DIR/privkey.pem" $SSL_DIR/wildcard.key

           # Fix ownership and permissions immediately
           echo "🔒 Setting proper permissions for nginx..."
           sudo chown $NGINX_UID:$NGINX_UID $SSL_DIR/wildcard.crt $SSL_DIR/wildcard.key
           sudo chmod 644 $SSL_DIR/wildcard.crt
           sudo chmod 600 $SSL_DIR/wildcard.key

            # Reload nginx if running
            if check_nginx_running; then
                docker compose exec nginx nginx -s reload
                echo "🔄 Nginx reloaded"
            fi
        else
            echo "❌ Could not find certificate directory"
            echo "📋 Available directories:"
            sudo ls -la /etc/letsencrypt/live/ || echo "No certificates found"
            return 1
        fi
    else
        echo "❌ Certificate generation failed (exit code: $certbot_exit_code)"
        return 1
    fi
}

# Check existing certificate
if [ -f "$SSL_DIR/$DOMAIN.crt" ] && [ -f "$SSL_DIR/$DOMAIN.key" ]; then
    echo "📋 Certificate exists:"
    openssl x509 -in $SSL_DIR/$DOMAIN.crt -text -noout | grep -E "(Subject:|DNS:|Not After)" | head -3

    read -p "Recreate certificate? [y/N]: " recreate
    [[ ! $recreate =~ ^[Yy]$ ]] && { echo "Using existing certificate"; exit 0; }
fi

# Choose certificate type
echo ""
echo "Choose certificate type:"
echo "1) Self-signed (testing - covers specific domain + wildcard)"
echo "2) Let's Encrypt - Standalone (single domain: $DOMAIN)"
echo "3) Let's Encrypt - Webroot (single domain: $DOMAIN)"
echo "4) Let's Encrypt - DNS Challenge (wildcard: $WILDCARD_DOMAIN)"
read -p "Enter choice [1-4]: " choice

# Set default paths for the final report
CERT_FINAL_PATH="$SSL_DIR/$DOMAIN.crt"
KEY_FINAL_PATH="$SSL_DIR/$DOMAIN.key"

case $choice in
    1) create_self_signed ;;
    2) create_letsencrypt_standalone ;;
    3) create_letsencrypt_webroot ;;
    4)
        create_letsencrypt_dns
        # If DNS challenge was chosen and succeeded, update the paths for the final report
        if [ $? -eq 0 ]; then
             CERT_FINAL_PATH="$SSL_DIR/wildcard.crt"
             KEY_FINAL_PATH="$SSL_DIR/wildcard.key"
        fi
        ;;
    *) echo "Invalid choice"; exit 1 ;;
esac

# Set permissions
sudo chown $NGINX_UID:$NGINX_UID $SSL_DIR/* 2>/dev/null
sudo chmod 600 $SSL_DIR/*.key 2>/dev/null
sudo chmod 644 $SSL_DIR/*.crt 2>/dev/null

# Show certificate info (always check the file nginx will use)
if [ -f "$CERT_FINAL_PATH" ]; then
    echo ""
    echo "📋 Certificate Information for $CERT_FINAL_PATH:"
    openssl x509 -in "$CERT_FINAL_PATH" -text -noout | grep -E "(Subject:|DNS:|Not After)"

    echo ""
    echo "🌐 Domains covered by this certificate:"
    openssl x509 -in "$CERT_FINAL_PATH" -text -noout | grep -A1 "Subject Alternative Name" | grep DNS || echo "  - $DOMAIN"
fi

echo ""
echo "✅ SSL setup completed!"
echo "Certificate: $CERT_FINAL_PATH"
echo "Private key: $KEY_FINAL_PATH"

