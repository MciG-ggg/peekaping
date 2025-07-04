# All-in-One Docker Variants

This document explains the three all-in-one Docker variants available for Peekaping, each with a different database backend.

## Available Variants

### 1. All-in-One SQLite (Recommended for Development)
- **Database**: SQLite (embedded)
- **Container**: Single container with application and database
- **Use Case**: Development, testing, simple deployments
- **Pros**: Simplest setup, no external database required
- **Cons**: Not suitable for high-concurrency scenarios

### 2. All-in-One MongoDB
- **Database**: MongoDB (embedded in container)
- **Container**: Single container with application and MongoDB
- **Use Case**: Production deployments requiring document storage
- **Pros**: Full-featured NoSQL database, good for complex data structures
- **Cons**: Larger container size, more memory usage

### 3. All-in-One PostgreSQL
- **Database**: PostgreSQL (embedded in container)
- **Container**: Single container with application and PostgreSQL
- **Use Case**: Production deployments requiring relational database
- **Pros**: Full-featured SQL database, excellent for complex queries
- **Cons**: Larger container size, more memory usage

## Quick Start

### SQLite Variant
```bash
# Using Docker Compose
docker-compose -f docker-compose.all-in-one-sqlite.yml up -d

# Using Docker directly
docker build -f Dockerfile.all-in-one-sqlite -t peekaping:sqlite .
docker run -p 8383:8383 -v ./data:/app/data peekaping:sqlite
```

### MongoDB Variant
```bash
# Using Docker Compose
docker-compose -f docker-compose.all-in-one-mongo.yml up -d

# Using Docker directly
docker build -f Dockerfile.all-in-one-mongo -t peekaping:mongo .
docker run -p 8383:8383 -v ./data:/data/db peekaping:mongo
```

### PostgreSQL Variant
```bash
# Using Docker Compose
docker-compose -f docker-compose.all-in-one-postgres.yml up -d

# Using Docker directly
docker build -f Dockerfile.all-in-one-postgres -t peekaping:postgres .
docker run -p 8383:8383 -v ./data:/var/lib/postgresql/data peekaping:postgres
```

## Configuration

### Environment Variables

All variants support the following environment variables:

#### Common Variables
- `API_URL`: API URL for the web application (default: `http://localhost:8383`)

#### Database-Specific Variables

**SQLite:**
- `DB_TYPE`: `sqlite` (automatically set)
- `DB_NAME`: Database file path (default: `/app/data/peekaping.db`)

**MongoDB:**
- `DB_TYPE`: `mongodb` (automatically set)
- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `27017`)
- `DB_NAME`: Database name (default: `peekaping`)
- `DB_USER`: Database user (default: `peekaping`)
- `DB_PASS`: Database password (default: `peekaping`)

**PostgreSQL:**
- `DB_TYPE`: `postgres` (automatically set)
- `DB_HOST`: Database host (default: `localhost`)
- `DB_PORT`: Database port (default: `5432`)
- `DB_NAME`: Database name (default: `peekaping`)
- `DB_USER`: Database user (default: `peekaping`)
- `DB_PASS`: Database password (default: `peekaping`)

### Custom Configuration

You can override default values by creating a `.env` file:

```env
# .env file example
DB_NAME=myapp
DB_USER=myuser
DB_PASS=mypassword
API_URL=https://monitoring.example.com
```

## Data Persistence

### SQLite
Data is stored in `/app/data/peekaping.db` inside the container.
Mount `./data:/app/data` to persist data.

### MongoDB
Data is stored in `/data/db` inside the container.
Mount `./data:/data/db` to persist data.

### PostgreSQL
Data is stored in `/var/lib/postgresql/data` inside the container.
Mount `./data:/var/lib/postgresql/data` to persist data.

## Ports

All variants expose port `8383` for the web application.

## Logs

Application logs are stored in `/var/log/supervisor/` inside the container:
- `server.log`: Application server logs
- `caddy.log`: Caddy web server logs
- `mongodb.log`: MongoDB logs (MongoDB variant only)
- `postgres.log`: PostgreSQL logs (PostgreSQL variant only)

Mount `./logs:/var/log/supervisor` to persist logs.

## Performance Considerations

### SQLite
- Best for: < 1000 concurrent users
- Memory usage: ~100-200MB
- Startup time: ~5-10 seconds

### MongoDB
- Best for: 1000-10000 concurrent users
- Memory usage: ~200-500MB
- Startup time: ~10-20 seconds

### PostgreSQL
- Best for: 1000-10000 concurrent users
- Memory usage: ~200-500MB
- Startup time: ~10-20 seconds

## Troubleshooting

### Container won't start
1. Check logs: `docker logs <container_name>`
2. Verify port 8383 is available
3. Check data directory permissions

### Database connection issues
1. Verify database is initialized (check logs)
2. Check environment variables
3. Verify data persistence volume mounts

### Migration issues
1. Check migration logs in container
2. Verify database user has correct permissions
3. Try rebuilding the container

## Security Notes

### Default Credentials
**Important**: Change default database passwords in production!

For MongoDB and PostgreSQL variants, the default password is `peekaping`.
Always set a strong password using the `DB_PASS` environment variable.

### Network Security
- The database ports are not exposed externally
- Only port 8383 is exposed for the web application
- Use reverse proxy with SSL/TLS for production

## Upgrading

To upgrade to a new version:

1. Stop the current container
2. Pull/build the new image
3. Start the new container with the same data volumes

```bash
# Example upgrade process
docker-compose -f docker-compose.all-in-one-sqlite.yml down
docker-compose -f docker-compose.all-in-one-sqlite.yml pull
docker-compose -f docker-compose.all-in-one-sqlite.yml up -d
```

## Backup and Restore

### SQLite
```bash
# Backup
docker exec <container> cp /app/data/peekaping.db /app/data/backup.db

# Restore
docker exec <container> cp /app/data/backup.db /app/data/peekaping.db
```

### MongoDB
```bash
# Backup
docker exec <container> mongodump --db peekaping --out /data/backup

# Restore
docker exec <container> mongorestore --db peekaping /data/backup/peekaping
```

### PostgreSQL
```bash
# Backup
docker exec <container> pg_dump peekaping > backup.sql

# Restore
docker exec <container> psql peekaping < backup.sql
```

## Support

For issues specific to the all-in-one variants, please include:
- Which variant you're using (SQLite/MongoDB/PostgreSQL)
- Your docker-compose.yml or docker run command
- Environment variables (without sensitive data)
- Container logs
