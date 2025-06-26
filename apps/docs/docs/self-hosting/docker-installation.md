---
sidebar_position: 1
---

# Docker Installation

The easiest way to run Peekaping is using Docker. This guide will help you get Peekaping up and running in minutes.

## Prerequisites

- Docker Engine 20.0+
- Docker Compose 2.0+
- At least 512MB RAM
- 1GB of free disk space

## Quick Start

### 1. Download Configuration Files

```bash
# Download the example environment file
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/.env.example -o .env

# Download the production docker-compose file
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/docker-compose.prod.yml -o docker-compose.yml

# Download nginx configuration (optional, for reverse proxy)
curl -L https://raw.githubusercontent.com/0xfurai/peekaping/main/nginx.conf -o nginx.conf
```

### 2. Configure Environment

Edit the `.env` file with your preferred settings:

```env
# Database Configuration
DB_USER=root
DB_PASSWORD=your-secure-password-here
DB_NAME=peekaping
DB_HOST=mongodb
DB_PORT=27017

# Server Configuration
PORT=8034
CLIENT_URL="http://localhost:8383"

# JWT Configuration
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_SECRET_KEY=your-access-token-secret-here
REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_SECRET_KEY=your-refresh-token-secret-here

# Application Settings
MODE=prod
TZ="America/New_York"
```

:::warning Important Security Notes
- **Change all default passwords and secret keys**
- Use strong, unique passwords for the database
- Generate secure JWT secret keys (use a password generator)
- Consider using environment-specific secrets management
:::

### 3. Start Peekaping

```bash
# Start all services
docker compose up -d

# Check status
docker compose ps

# View logs
docker compose logs -f
```

### 4. Access Peekaping

Once all containers are running:

1. Open your browser and go to `http://localhost:8383`
2. Complete the initial setup wizard
3. Create your admin account
4. Start monitoring your services!

## Docker Images

Peekaping provides official Docker images:

- **Server**: [`0xfurai/peekaping-server`](https://hub.docker.com/r/0xfurai/peekaping-server)
- **Web**: [`0xfurai/peekaping-web`](https://hub.docker.com/r/0xfurai/peekaping-web)

### Image Tags

- `latest` - Latest stable release
- `main` - Latest development build
- `v1.x.x` - Specific version tags

## Advanced Configuration

### Custom Docker Compose

If you want to customize the setup, here's a complete `docker-compose.yml` example:

```yaml
version: '3.8'

services:
  mongodb:
    image: mongo:7
    container_name: peekaping-mongodb
    restart: unless-stopped
    environment:
      MONGO_INITDB_ROOT_USERNAME: ${DB_USER}
      MONGO_INITDB_ROOT_PASSWORD: ${DB_PASSWORD}
      MONGO_INITDB_DATABASE: ${DB_NAME}
    volumes:
      - mongodb_data:/data/db
    networks:
      - peekaping-network
    healthcheck:
      test: ["CMD", "mongosh", "--eval", "db.adminCommand('ping')"]
      interval: 30s
      timeout: 10s
      retries: 3

  server:
    image: 0xfurai/peekaping-server:latest
    container_name: peekaping-server
    restart: unless-stopped
    environment:
      - DB_USER=${DB_USER}
      - DB_PASSWORD=${DB_PASSWORD}
      - DB_NAME=${DB_NAME}
      - DB_HOST=mongodb
      - DB_PORT=27017
      - PORT=${PORT}
      - CLIENT_URL=${CLIENT_URL}
      - ACCESS_TOKEN_EXPIRED_IN=${ACCESS_TOKEN_EXPIRED_IN}
      - ACCESS_TOKEN_SECRET_KEY=${ACCESS_TOKEN_SECRET_KEY}
      - REFRESH_TOKEN_EXPIRED_IN=${REFRESH_TOKEN_EXPIRED_IN}
      - REFRESH_TOKEN_SECRET_KEY=${REFRESH_TOKEN_SECRET_KEY}
      - MODE=${MODE}
      - TZ=${TZ}
    depends_on:
      mongodb:
        condition: service_healthy
    networks:
      - peekaping-network
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:${PORT}/health"]
      interval: 30s
      timeout: 10s
      retries: 3

  web:
    image: 0xfurai/peekaping-web:latest
    container_name: peekaping-web
    restart: unless-stopped
    ports:
      - "8383:80"
    environment:
      - API_URL=http://server:${PORT}
    depends_on:
      server:
        condition: service_healthy
    networks:
      - peekaping-network

volumes:
  mongodb_data:

networks:
  peekaping-network:
    driver: bridge
```

### Environment Variables Reference

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `DB_USER` | MongoDB username | `root` | Yes |
| `DB_PASSWORD` | MongoDB password | - | Yes |
| `DB_NAME` | Database name | `peekaping` | Yes |
| `DB_HOST` | MongoDB hostname | `mongodb` | Yes |
| `DB_PORT` | MongoDB port | `27017` | Yes |
| `PORT` | Server port | `8034` | Yes |
| `CLIENT_URL` | Frontend URL | `http://localhost:8383` | Yes |
| `ACCESS_TOKEN_EXPIRED_IN` | Access token expiry | `15m` | Yes |
| `ACCESS_TOKEN_SECRET_KEY` | Access token secret | - | Yes |
| `REFRESH_TOKEN_EXPIRED_IN` | Refresh token expiry | `60m` | Yes |
| `REFRESH_TOKEN_SECRET_KEY` | Refresh token secret | - | Yes |
| `MODE` | Log level | `prod` | No |
| `TZ` | Timezone | `UTC` | No |

## Persistent Data

Peekaping stores data in MongoDB. The docker-compose setup creates a named volume `mongodb_data` to persist your monitoring data.

### Backup Data

```bash
# Create backup
docker exec peekaping-mongodb mongodump --uri="mongodb://root:your-password@localhost:27017/peekaping?authSource=admin" --out=/backup

# Copy backup to host
docker cp peekaping-mongodb:/backup ./peekaping-backup
```

### Restore Data

```bash
# Copy backup to container
docker cp ./peekaping-backup peekaping-mongodb:/backup

# Restore backup
docker exec peekaping-mongodb mongorestore --uri="mongodb://root:your-password@localhost:27017/peekaping?authSource=admin" /backup/peekaping
```

## Reverse Proxy Setup

### Nginx

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8383;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_cache_bypass $http_upgrade;
    }
}
```

### Traefik

```yaml
labels:
  - "traefik.enable=true"
  - "traefik.http.routers.peekaping.rule=Host(`your-domain.com`)"
  - "traefik.http.routers.peekaping.entrypoints=web"
  - "traefik.http.services.peekaping.loadbalancer.server.port=80"
```

## Troubleshooting

### Check Container Status

```bash
# View all containers
docker compose ps

# View logs for specific service
docker compose logs -f server
docker compose logs -f web
docker compose logs -f mongodb
```

### Common Issues

#### Container Won't Start

1. Check if ports are already in use:
   ```bash
   netstat -tulpn | grep :8383
   ```

2. Verify environment variables:
   ```bash
   docker compose config
   ```

#### Database Connection Issues

1. Check MongoDB health:
   ```bash
   docker exec peekaping-mongodb mongosh --eval "db.adminCommand('ping')"
   ```

2. Verify credentials in `.env` file

#### Web Interface Not Loading

1. Check if all services are healthy:
   ```bash
   docker compose ps
   ```

2. Verify network connectivity:
   ```bash
   docker exec peekaping-web curl -f http://server:8034/health
   ```

### Updating Peekaping

```bash
# Pull latest images
docker compose pull

# Restart with new images
docker compose up -d

# Clean up old images
docker image prune
```

## Next Steps

- [Configure your first monitor](/tutorial-basics/create-a-document)
- [Set up notification channels](/tutorial-basics/create-a-page)
- [Create status pages](/tutorial-basics/create-a-blog-post)
- [Enable security features](/tutorial-basics/deploy-your-site)
