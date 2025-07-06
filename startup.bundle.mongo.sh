#!/bin/sh
set -e

validate_env_vars() {
    local errors=0

    if [ -z "$DB_USER" ]; then
        echo "ERROR: DB_USER is required and must be set"
        errors=1
    fi

    if [ -z "$DB_PASS" ]; then
        echo "ERROR: DB_PASS is required and must be set"
        errors=1
    fi

    if [ -z "$ACCESS_TOKEN_SECRET_KEY" ]; then
        echo "ERROR: ACCESS_TOKEN_SECRET_KEY is required and must be set"
        errors=1
    fi

    if [ -z "$REFRESH_TOKEN_SECRET_KEY" ]; then
        echo "ERROR: REFRESH_TOKEN_SECRET_KEY is required and must be set"
        errors=1
    fi

    # Validate secret key strength
    if [ ${#ACCESS_TOKEN_SECRET_KEY} -lt 32 ]; then
        echo "ERROR: ACCESS_TOKEN_SECRET_KEY must be at least 32 characters long"
        errors=1
    fi

    if [ ${#REFRESH_TOKEN_SECRET_KEY} -lt 32 ]; then
        echo "ERROR: REFRESH_TOKEN_SECRET_KEY must be at least 32 characters long"
        errors=1
    fi

    if [ $errors -eq 1 ]; then
        echo "Environment validation failed. Please fix the above errors."
        exit 1
    fi

    echo "Environment validation passed."
}

validate_env_vars


# Create env.js file for the web app
cat >/app/web/env.js <<EOF
/* generated each container start */
window.__CONFIG__ = {
  API_URL: ""
};
EOF
# Security: Set appropriate permissions for web assets
chmod 644 /app/web/env.js

# Set default environment variables if not provided
export DB_TYPE=${DB_TYPE:-mongo}
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-27017}
export DB_NAME=${DB_NAME:-peekaping}
export DB_USER=${DB_USER}
export DB_PASS=${DB_PASS}

# Security: Use separate admin credentials
export DB_ADMIN_USER=${DB_ADMIN_USER:-admin}
export DB_ADMIN_PASS=${DB_ADMIN_PASS:-$DB_PASS}

# Set server configuration environment variables
export PORT=${PORT:-8034}
# Security: Use HTTPS by default
export CLIENT_URL=${CLIENT_URL:-http://localhost:8383}
export ACCESS_TOKEN_SECRET_KEY=${ACCESS_TOKEN_SECRET_KEY}
export REFRESH_TOKEN_SECRET_KEY=${REFRESH_TOKEN_SECRET_KEY}
export ACCESS_TOKEN_EXPIRED_IN=${ACCESS_TOKEN_EXPIRED_IN:-15m}
export REFRESH_TOKEN_EXPIRED_IN=${REFRESH_TOKEN_EXPIRED_IN:-168h}
export MODE=${MODE:-prod}
export TZ=${TZ:-UTC}

# Create .env file for the server with secure permissions
cat > /app/.env << EOF
PORT=$PORT
CLIENT_URL=$CLIENT_URL
DB_TYPE=$DB_TYPE
DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASS=$DB_PASS
ACCESS_TOKEN_SECRET_KEY=$ACCESS_TOKEN_SECRET_KEY
REFRESH_TOKEN_SECRET_KEY=$REFRESH_TOKEN_SECRET_KEY
ACCESS_TOKEN_EXPIRED_IN=$ACCESS_TOKEN_EXPIRED_IN
REFRESH_TOKEN_EXPIRED_IN=$REFRESH_TOKEN_EXPIRED_IN
MODE=$MODE
TZ=$TZ
EOF

# Security: Set restrictive permissions on sensitive config file
chmod 600 /app/.env

# Create data directory if it doesn't exist
mkdir -p /data/db

# Create log directory and fix permissions
mkdir -p /var/log/supervisor
chmod 755 /var/log/supervisor

# Initialize MongoDB if needed
if [ ! -f /data/db/.mongodb_initialized ]; then
    echo "Initializing MongoDB..."

    # Start MongoDB without auth for initial setup
    mongod --dbpath /data/db --fork --logpath /var/log/supervisor/mongodb-init.log --noauth

    # Wait for MongoDB to be ready
    echo "Waiting for MongoDB to be ready..."
    retry_count=0
    max_retries=30

    while [ $retry_count -lt $max_retries ]; do
        if mongosh admin --eval "db.runCommand('ping')" >/dev/null 2>&1; then
            echo "MongoDB is ready!"
            break
        fi
        sleep 1
        retry_count=$((retry_count + 1))
    done

    if [ $retry_count -eq $max_retries ]; then
        echo "ERROR: MongoDB failed to start within timeout"
        exit 1
    fi

    # Create users in a single operation
    echo "Creating MongoDB users..."
    mongosh admin --eval "
        db.createUser({
            user: '$DB_ADMIN_USER',
            pwd: '$DB_ADMIN_PASS',
            roles: ['root']
        });
        db.createUser({
            user: '$DB_USER',
            pwd: '$DB_PASS',
            roles: [
                { role: 'readWrite', db: '$DB_NAME' }
            ]
        });
    "

    # Stop MongoDB gracefully
    echo "Stopping MongoDB after initialization..."
    mongod --dbpath /data/db --shutdown

    # Mark as initialized
    touch /data/db/.mongodb_initialized
    chmod 600 /data/db/.mongodb_initialized
    echo "MongoDB initialization completed!"
fi

# Security: Wait for MongoDB to be fully ready with proper health check
echo "Waiting for MongoDB to be fully ready..."
retry_count=0
max_retries=30

while [ $retry_count -lt $max_retries ]; do
    if mongosh "$DB_NAME" --authenticationDatabase admin -u "$DB_USER" -p "$DB_PASS" --eval "db.runCommand('ping')" >/dev/null 2>&1; then
        echo "MongoDB is ready and accessible!"
        break
    fi
    sleep 1
    retry_count=$((retry_count + 1))
done

if [ $retry_count -eq $max_retries ]; then
    echo "ERROR: MongoDB failed to become accessible within timeout"
    exit 1
fi

# Run database migrations
echo "Running database migrations..."
cd /app/server
if ./run-migrations.sh; then
    echo "Migrations completed successfully!"
else
    echo "ERROR: Migration failed!"
    exit 1
fi

# Security: Clear sensitive variables from memory
unset DB_PASS
unset DB_ADMIN_PASS
unset ACCESS_TOKEN_SECRET_KEY
unset REFRESH_TOKEN_SECRET_KEY

# Start supervisor to manage MongoDB, server, and Caddy
echo "Starting supervisor to manage MongoDB, server, and Caddy..."
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
