sh #!/bin/bash

# Script to generate Redis TLS certificates for testing
# This creates self-signed certificates for local testing only

set -e

CERT_DIR="./redis-certs"
mkdir -p "$CERT_DIR"

echo "Generating Redis TLS certificates for testing..."

# Generate CA private key
openssl genrsa -out "$CERT_DIR/ca.key" 2048

# Generate CA certificate
openssl req -new -x509 -days 365 -key "$CERT_DIR/ca.key" -sha256 -out "$CERT_DIR/ca.crt" -subj "/C=US/ST=Test/L=Test/O=Test/CN=Redis-Test-CA"

# Generate server private key
openssl genrsa -out "$CERT_DIR/redis.key" 2048

# Generate server certificate signing request
openssl req -new -key "$CERT_DIR/redis.key" -out "$CERT_DIR/redis.csr" -subj "/C=US/ST=Test/L=Test/O=Test/CN=localhost"

# Sign server certificate with CA
openssl x509 -req -in "$CERT_DIR/redis.csr" -CA "$CERT_DIR/ca.crt" -CAkey "$CERT_DIR/ca.key" -CAcreateserial -out "$CERT_DIR/redis.crt" -days 365 -sha256

# Set proper permissions
chmod 600 "$CERT_DIR"/*.key
chmod 644 "$CERT_DIR"/*.crt

echo "Redis TLS certificates generated successfully!"
echo "Certificate files created in: $CERT_DIR"
echo ""
echo "You can now start the Redis containers with:"
echo "docker-compose -f docker-compose.redis-test.yml up -d"
echo ""
echo "Test connection strings:"
echo "- Simple Redis: redis://localhost:6381"
echo "- Redis with auth: redis://:testpassword@localhost:6379"
echo "- Redis with TLS (ignore cert): rediss://:testpassword@localhost:6380 (set ignoreTls=true)"
