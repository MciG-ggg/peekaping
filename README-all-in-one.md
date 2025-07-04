# Peekaping All-in-One Container

This setup combines the entire Peekaping stack (server, web, and migrations) into a single Docker container for simplified deployment.

## Features

- ✅ Go backend server
- ✅ React frontend (built and served by Caddy)
- ✅ Database migrations (SQLite)
- ✅ Caddy reverse proxy (API calls to backend, static files served directly)
- ✅ Process management via supervisor

## Quick Start

### Using Docker Compose (Recommended)

```bash
# Build and run the all-in-one container
docker-compose -f docker-compose.all-in-one.yml up --build

# Or run in detached mode
docker-compose -f docker-compose.all-in-one.yml up -d --build
```

### Using Docker Commands

```bash
# Build the image
docker build -f Dockerfile.all-in-one -t peekaping-all-in-one .

# Run the container
docker run -d \
  --name peekaping-all-in-one \
  -p 8383:8383 \
  -v $(pwd)/.data/sqlite:/app/data \
  -e DB_TYPE=sqlite \
  -e DB_NAME=/app/data/peekaping.db \
  -e API_URL=http://localhost:8383 \
  --env-file .env \
  peekaping-all-in-one
```

## Access

Once running, access your Peekaping instance at: http://localhost:8383

## Configuration

The container uses the same environment variables as the original multi-container setup. Make sure your `.env` file contains the necessary configuration.

Key environment variables:
- `DB_TYPE=sqlite` (automatically set)
- `DB_NAME=/app/data/peekaping.db` (automatically set)
- `API_URL=http://localhost:8383` (automatically set)

## Data Persistence

SQLite database is stored in `/app/data/peekaping.db` inside the container and mounted to `./.data/sqlite/` on the host.

## Logs

To view logs:
```bash
# All logs
docker-compose -f docker-compose.all-in-one.yml logs -f

# Or for docker run
docker logs -f peekaping-all-in-one
```

## Stopping

```bash
# For docker-compose
docker-compose -f docker-compose.all-in-one.yml down

# For docker run
docker stop peekaping-all-in-one
docker rm peekaping-all-in-one
```

## Architecture

The all-in-one container uses:
- **Supervisor**: Process manager that starts and monitors both the Go server and Caddy
- **Caddy**: Serves static files and proxies API calls to the local Go server
- **Go Server**: Runs on port 8034 internally, accessed via Caddy proxy
- **SQLite Database**: File-based database stored in persistent volume

## Process Flow

1. Container starts and runs `startup.sh`
2. Database migrations are executed
3. Supervisor starts both the Go server and Caddy
4. Caddy listens on port 8383 and serves the React app
5. API calls to `/api/*` are proxied to the Go server on localhost:8034
6. WebSocket connections to `/socket.io/*` are proxied to the Go server
