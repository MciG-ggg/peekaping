#!/bin/sh
set -e

# Create env.js file for the web app
API_URL=${API_URL:-}
cat >/app/web/env.js <<EOF
/* generated each container start */
window.__CONFIG__ = {
  API_URL: "$API_URL"
};
EOF

# Set default environment variables if not provided
export DB_TYPE=${DB_TYPE:-postgres}
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_NAME=${DB_NAME:-peekaping}
export DB_USER=${DB_USER:-peekaping}
export DB_PASS=${DB_PASS:-peekaping}

# Create data directory if it doesn't exist
mkdir -p /var/lib/postgresql/data

# Initialize PostgreSQL if needed
if [ ! -f /var/lib/postgresql/data/.postgres_initialized ]; then
    echo "Initializing PostgreSQL..."

    # Clear data directory if it exists but is not initialized
    if [ -d /var/lib/postgresql/data ]; then
        rm -rf /var/lib/postgresql/data/*
        rm -rf /var/lib/postgresql/data/.[^.]*
    fi

    # Initialize PostgreSQL cluster
    su-exec postgres initdb -D /var/lib/postgresql/data

    # Configure PostgreSQL
    echo "host all all 0.0.0.0/0 md5" >> /var/lib/postgresql/data/pg_hba.conf
    echo "listen_addresses = '*'" >> /var/lib/postgresql/data/postgresql.conf

    # Start PostgreSQL temporarily
    su-exec postgres pg_ctl -D /var/lib/postgresql/data -l /var/log/supervisor/postgres-init.log start

    # Wait for PostgreSQL to be ready
    echo "Waiting for PostgreSQL to be ready..."
    sleep 5

    # Create database and user
    su-exec postgres psql -c "CREATE USER $DB_USER WITH PASSWORD '$DB_PASS';"
    su-exec postgres psql -c "CREATE DATABASE $DB_NAME OWNER $DB_USER;"
    su-exec postgres psql -c "GRANT ALL PRIVILEGES ON DATABASE $DB_NAME TO $DB_USER;"

    # Stop PostgreSQL
    su-exec postgres pg_ctl -D /var/lib/postgresql/data stop

    # Mark as initialized
    touch /var/lib/postgresql/data/.postgres_initialized
    echo "PostgreSQL initialization completed!"
fi

# Run database migrations
echo "Running database migrations..."
cd /app/server
if ./run-migrations.sh; then
    echo "Migrations completed successfully!"
else
    echo "Migration failed!"
    exit 1
fi

# Start supervisor to manage PostgreSQL, server, and Caddy
echo "Starting supervisor to manage PostgreSQL, server, and Caddy..."
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
