---
sidebar_position: 2
---

# Manual Installation

This guide covers installing Peekaping manually without Docker. This method gives you more control over the installation but requires more setup steps.

## Prerequisites

### System Requirements
- **Operating System**: Linux (Ubuntu 20.04+, Debian 11+, CentOS 8+) or macOS
- **CPU**: 2 cores minimum, 4 cores recommended
- **Memory**: 2GB RAM minimum, 4GB recommended
- **Storage**: 10GB free disk space minimum
- **Network**: Internet access for downloading packages

### Required Software

#### MongoDB
Peekaping requires MongoDB 4.4 or later.

**Ubuntu/Debian:**
```bash
# Import public key
wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | sudo apt-key add -

# Add repository
echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/6.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-6.0.list

# Install MongoDB
sudo apt-get update
sudo apt-get install -y mongodb-org

# Start MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod
```

**CentOS/RHEL:**
```bash
# Add repository
cat > /etc/yum.repos.d/mongodb-org-6.0.repo << EOF
[mongodb-org-6.0]
name=MongoDB Repository
baseurl=https://repo.mongodb.org/yum/redhat/8/mongodb-org/6.0/x86_64/
gpgcheck=1
enabled=1
gpgkey=https://www.mongodb.org/static/pgp/server-6.0.asc
EOF

# Install MongoDB
sudo yum install -y mongodb-org

# Start MongoDB
sudo systemctl start mongod
sudo systemctl enable mongod
```

#### Node.js (for web interface)
```bash
# Install Node.js 18.x
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

#### Go (for server compilation)
```bash
# Download and install Go 1.21+
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify installation
go version
```

#### Additional Dependencies
```bash
# Ubuntu/Debian
sudo apt-get update
sudo apt-get install -y curl wget git build-essential

# CentOS/RHEL
sudo yum groupinstall -y "Development Tools"
sudo yum install -y curl wget git
```

## Installation Steps

### 1. Create User and Directories

```bash
# Create peekaping user
sudo useradd -r -s /bin/false peekaping

# Create directories
sudo mkdir -p /opt/peekaping/{server,web}
sudo mkdir -p /var/log/peekaping
sudo mkdir -p /etc/peekaping

# Set permissions
sudo chown -R peekaping:peekaping /opt/peekaping
sudo chown -R peekaping:peekaping /var/log/peekaping
```

### 2. Download and Build Server

```bash
# Clone repository
cd /tmp
git clone https://github.com/0xfurai/peekaping.git
cd peekaping

# Build server
cd apps/server
go mod download
go build -o peekaping-server ./src

# Install server binary
sudo cp peekaping-server /opt/peekaping/server/
sudo chown peekaping:peekaping /opt/peekaping/server/peekaping-server
sudo chmod +x /opt/peekaping/server/peekaping-server
```

### 3. Build Web Interface

```bash
# Build web interface
cd ../web
npm install
npm run build

# Install web files
sudo cp -r dist/* /opt/peekaping/web/
sudo chown -R peekaping:peekaping /opt/peekaping/web/
```

### 4. Configuration

Create the main configuration file:

```bash
sudo tee /etc/peekaping/peekaping.env << EOF
# Database Configuration
DB_USER=peekaping
DB_PASSWORD=secure-password-here
DB_NAME=peekaping
DB_HOST=localhost
DB_PORT=27017

# Server Configuration
PORT=8034
CLIENT_URL="http://localhost:8383"

# JWT Configuration
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_SECRET_KEY=$(openssl rand -base64 64)
REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_SECRET_KEY=$(openssl rand -base64 64)

# Application Settings
MODE=prod
TZ="UTC"
LOG_LEVEL=info
EOF

# Set permissions
sudo chown peekaping:peekaping /etc/peekaping/peekaping.env
sudo chmod 600 /etc/peekaping/peekaping.env
```

### 5. Configure MongoDB

```bash
# Create database and user
mongosh << EOF
use admin
db.createUser({
  user: "peekaping",
  pwd: "secure-password-here",
  roles: [
    { role: "readWrite", db: "peekaping" }
  ]
})

use peekaping
db.createCollection("monitors")
EOF
```

### 6. Create Systemd Services

#### Server Service

```bash
sudo tee /etc/systemd/system/peekaping-server.service << EOF
[Unit]
Description=Peekaping Server
After=network.target mongod.service
Requires=mongod.service

[Service]
Type=simple
User=peekaping
Group=peekaping
WorkingDirectory=/opt/peekaping/server
ExecStart=/opt/peekaping/server/peekaping-server
EnvironmentFile=/etc/peekaping/peekaping.env
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal

# Security settings
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/var/log/peekaping

[Install]
WantedBy=multi-user.target
EOF
```

#### Web Service (Nginx)

First, install Nginx:

```bash
# Ubuntu/Debian
sudo apt-get install -y nginx

# CentOS/RHEL
sudo yum install -y nginx
```

Configure Nginx:

```bash
sudo tee /etc/nginx/sites-available/peekaping << EOF
server {
    listen 8383;
    server_name _;

    root /opt/peekaping/web;
    index index.html;

    # API proxy
    location /api/ {
        proxy_pass http://localhost:8034;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
    }

    # WebSocket proxy
    location /ws {
        proxy_pass http://localhost:8034;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
    }

    # Static files
    location / {
        try_files \$uri \$uri/ /index.html;
        add_header Cache-Control "no-cache, no-store, must-revalidate";
        add_header Pragma "no-cache";
        add_header Expires "0";
    }

    # Static assets with caching
    location ~* \.(js|css|png|jpg|jpeg|gif|ico|svg|woff|woff2|ttf|eot)$ {
        expires 1y;
        add_header Cache-Control "public, immutable";
    }
}
EOF

# Enable site
sudo ln -s /etc/nginx/sites-available/peekaping /etc/nginx/sites-enabled/
sudo nginx -t
```

### 7. Start Services

```bash
# Start and enable services
sudo systemctl start peekaping-server
sudo systemctl enable peekaping-server

sudo systemctl restart nginx
sudo systemctl enable nginx

# Check service status
sudo systemctl status peekaping-server
sudo systemctl status nginx
```

## Configuration Management

### Environment Variables

You can override configuration by editing `/etc/peekaping/peekaping.env`:

```bash
sudo nano /etc/peekaping/peekaping.env

# Restart server after changes
sudo systemctl restart peekaping-server
```

### Logging

View logs using journalctl:

```bash
# Server logs
sudo journalctl -u peekaping-server -f

# Nginx logs
sudo tail -f /var/log/nginx/access.log
sudo tail -f /var/log/nginx/error.log
```

### Database Management

#### Backup Database

```bash
#!/bin/bash
# Create backup script
BACKUP_DIR="/var/backups/peekaping"
DATE=$(date +%Y%m%d_%H%M%S)

sudo -u peekaping mkdir -p $BACKUP_DIR

mongodump \
  --host localhost:27017 \
  --username peekaping \
  --password secure-password-here \
  --authenticationDatabase admin \
  --db peekaping \
  --out $BACKUP_DIR/$DATE

# Compress backup
tar -czf $BACKUP_DIR/peekaping_$DATE.tar.gz -C $BACKUP_DIR $DATE
rm -rf $BACKUP_DIR/$DATE

# Keep only last 7 days
find $BACKUP_DIR -name "peekaping_*.tar.gz" -mtime +7 -delete
```

#### Restore Database

```bash
# Extract backup
tar -xzf /var/backups/peekaping/peekaping_20240120_143000.tar.gz -C /tmp

# Restore database
mongorestore \
  --host localhost:27017 \
  --username peekaping \
  --password secure-password-here \
  --authenticationDatabase admin \
  --db peekaping \
  /tmp/20240120_143000/peekaping
```

## Security Hardening

### Firewall Configuration

```bash
# Ubuntu/Debian (UFW)
sudo ufw allow 8383/tcp  # Web interface
sudo ufw allow ssh
sudo ufw enable

# CentOS/RHEL (firewalld)
sudo firewall-cmd --permanent --add-port=8383/tcp
sudo firewall-cmd --permanent --add-service=ssh
sudo firewall-cmd --reload
```

### MongoDB Security

```bash
# Edit MongoDB configuration
sudo nano /etc/mongod.conf

# Add security settings:
security:
  authorization: enabled

net:
  bindIp: 127.0.0.1

# Restart MongoDB
sudo systemctl restart mongod
```

### SSL/TLS Configuration

For production, configure SSL/TLS with Let's Encrypt:

```bash
# Install Certbot
sudo apt-get install -y certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com

# Auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

## Monitoring and Maintenance

### Health Checks

Create health check script:

```bash
sudo tee /usr/local/bin/peekaping-health << EOF
#!/bin/bash

# Check if server is running
if ! systemctl is-active --quiet peekaping-server; then
    echo "ERROR: Peekaping server is not running"
    exit 1
fi

# Check if API is responding
if ! curl -f http://localhost:8034/health > /dev/null 2>&1; then
    echo "ERROR: Peekaping API is not responding"
    exit 1
fi

# Check MongoDB connection
if ! mongosh --quiet --eval "db.adminCommand('ping')" > /dev/null 2>&1; then
    echo "ERROR: MongoDB is not responding"
    exit 1
fi

echo "OK: All services are healthy"
EOF

sudo chmod +x /usr/local/bin/peekaping-health
```

### Log Rotation

```bash
sudo tee /etc/logrotate.d/peekaping << EOF
/var/log/peekaping/*.log {
    daily
    missingok
    rotate 30
    compress
    delaycompress
    notifempty
    create 644 peekaping peekaping
    postrotate
        systemctl reload peekaping-server
    endscript
}
EOF
```

### Performance Tuning

#### MongoDB Optimization

```bash
# Edit MongoDB configuration
sudo nano /etc/mongod.conf

# Optimize for monitoring workload
storage:
  wiredTiger:
    engineConfig:
      cacheSizeGB: 1
    indexConfig:
      prefixCompression: true

operationProfiling:
  slowOpThresholdMs: 100
  mode: slowOp
```

#### System Limits

```bash
# Increase file descriptor limits
sudo tee -a /etc/security/limits.conf << EOF
peekaping soft nofile 65536
peekaping hard nofile 65536
EOF

# Add to systemd service
sudo systemctl edit peekaping-server
# Add:
# [Service]
# LimitNOFILE=65536
```

## Troubleshooting

### Common Issues

#### Server Won't Start

```bash
# Check logs
sudo journalctl -u peekaping-server -n 50

# Check configuration
sudo -u peekaping /opt/peekaping/server/peekaping-server --check-config

# Verify permissions
ls -la /opt/peekaping/server/
ls -la /etc/peekaping/
```

#### Database Connection Issues

```bash
# Test MongoDB connection
mongosh "mongodb://peekaping:password@localhost:27017/peekaping"

# Check MongoDB logs
sudo journalctl -u mongod -f

# Verify MongoDB is running
sudo systemctl status mongod
```

#### Web Interface Not Loading

```bash
# Check Nginx configuration
sudo nginx -t

# Check Nginx logs
sudo tail -f /var/log/nginx/error.log

# Verify static files
ls -la /opt/peekaping/web/
```

### Performance Issues

```bash
# Monitor system resources
htop
iostat -x 1
free -h

# Monitor MongoDB performance
mongosh --eval "db.currentOp()"
mongosh --eval "db.serverStatus()"

# Check disk space
df -h
```

## Updating Peekaping

### Update Process

```bash
# Stop services
sudo systemctl stop peekaping-server

# Backup current installation
sudo cp -r /opt/peekaping /opt/peekaping.backup.$(date +%Y%m%d)

# Download new version
cd /tmp
git clone https://github.com/0xfurai/peekaping.git
cd peekaping

# Build and install server
cd apps/server
go build -o peekaping-server ./src
sudo cp peekaping-server /opt/peekaping/server/
sudo chown peekaping:peekaping /opt/peekaping/server/peekaping-server

# Build and install web interface
cd ../web
npm install
npm run build
sudo cp -r dist/* /opt/peekaping/web/
sudo chown -R peekaping:peekaping /opt/peekaping/web/

# Start services
sudo systemctl start peekaping-server
```

### Database Migrations

```bash
# Run database migrations (if needed)
sudo -u peekaping /opt/peekaping/server/peekaping-server --migrate

# Verify update
curl http://localhost:8034/health
```

This manual installation method provides full control over your Peekaping deployment while requiring more hands-on management. For easier maintenance, consider using the [Docker installation](/self-hosting/docker) method instead.
