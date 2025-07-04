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
export DB_TYPE=${DB_TYPE:-mongodb}
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-27017}
export DB_NAME=${DB_NAME:-peekaping}
export DB_USER=${DB_USER:-peekaping}
export DB_PASS=${DB_PASS:-peekaping}

# Create data directory if it doesn't exist
mkdir -p /data/db

# Initialize MongoDB if needed
if [ ! -f /data/db/.mongodb_initialized ]; then
    echo "Initializing MongoDB..."

    # Start MongoDB without auth temporarily
    mongod --dbpath /data/db --fork --logpath /var/log/supervisor/mongodb-init.log

    # Wait for MongoDB to be ready
    echo "Waiting for MongoDB to be ready..."
    sleep 5

    # Create admin user and database
    mongo admin --eval "
        db.createUser({
            user: 'admin',
            pwd: '$DB_PASS',
            roles: ['root']
        });
        db.auth('admin', '$DB_PASS');
        use $DB_NAME;
        db.createUser({
            user: '$DB_USER',
            pwd: '$DB_PASS',
            roles: ['readWrite', 'dbAdmin']
        });
    "

    # Stop MongoDB
    mongod --dbpath /data/db --shutdown

    # Mark as initialized
    touch /data/db/.mongodb_initialized
    echo "MongoDB initialization completed!"
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

# Start supervisor to manage MongoDB, server, and Caddy
echo "Starting supervisor to manage MongoDB, server, and Caddy..."
exec /usr/bin/supervisord -c /etc/supervisor/conf.d/supervisord.conf
