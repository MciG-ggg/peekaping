---
sidebar_position: 1
---

# Configuration Reference

This comprehensive reference covers all configuration options available in Peekaping, from environment variables to advanced settings.

## Environment Variables

Peekaping's behavior is controlled through environment variables defined in your `.env` file.

### Required Configuration

These variables must be set for Peekaping to function:

#### Database Configuration
```env
# MongoDB connection settings
DB_USER=root
DB_PASSWORD=your-secure-password
DB_NAME=peekaping
DB_HOST=mongodb
DB_PORT=27017
```

#### Server Configuration
```env
# Server port and client URL
PORT=8034
CLIENT_URL="http://localhost:8383"
```

#### JWT Authentication
```env
# JWT token configuration
ACCESS_TOKEN_EXPIRED_IN=15m
ACCESS_TOKEN_SECRET_KEY=your-very-long-and-secure-secret-key
REFRESH_TOKEN_EXPIRED_IN=60m
REFRESH_TOKEN_SECRET_KEY=your-other-very-long-and-secure-secret-key
```

### Optional Configuration

#### Application Settings
```env
# Logging and mode settings
MODE=prod                    # Options: prod, dev, debug
LOG_LEVEL=info              # Options: debug, info, warn, error
TZ="UTC"                    # Timezone for the application

# Data retention settings
HEARTBEAT_RETENTION_DAYS=90  # How long to keep heartbeat data
STATS_RETENTION_DAYS=365     # How long to keep statistics

# Rate limiting
API_RATE_LIMIT=1000         # Requests per hour per IP
BURST_RATE_LIMIT=100        # Burst requests per minute

# File upload limits
MAX_UPLOAD_SIZE=10MB        # Maximum file upload size
ALLOWED_FILE_TYPES="png,jpg,jpeg,svg"  # Allowed file extensions
```

#### Security Settings
```env
# CORS configuration
CORS_ORIGINS="http://localhost:5173,https://yourdomain.com"
CORS_CREDENTIALS=true

# Session configuration
SESSION_SECRET=another-secure-secret-key
SESSION_TIMEOUT=24h         # Session timeout duration

# Password policy
MIN_PASSWORD_LENGTH=8       # Minimum password length
REQUIRE_PASSWORD_COMPLEXITY=true  # Require special characters

# Two-factor authentication
TOTP_ISSUER="Peekaping"     # Name shown in authenticator apps
TOTP_WINDOW=1               # TOTP validation window (30s intervals)
```

#### Email Configuration
```env
# Default SMTP settings (can be overridden per notification channel)
DEFAULT_SMTP_HOST=smtp.gmail.com
DEFAULT_SMTP_PORT=587
DEFAULT_SMTP_USER=your-email@gmail.com
DEFAULT_SMTP_PASS=your-app-password
DEFAULT_SMTP_FROM=noreply@yourdomain.com
DEFAULT_SMTP_ENCRYPTION=TLS  # Options: TLS, SSL, NONE
```

#### Performance Settings
```env
# Worker configuration
MONITOR_WORKERS=10          # Number of concurrent monitoring workers
MAX_CONCURRENT_CHECKS=50    # Maximum concurrent monitor checks
WORKER_TIMEOUT=300s         # Worker timeout duration

# Cache configuration
CACHE_TTL=300              # Cache TTL in seconds
REDIS_URL=redis://redis:6379  # Redis URL for caching (optional)

# Database connection pool
DB_MAX_CONNECTIONS=100      # Maximum database connections
DB_CONNECTION_TIMEOUT=30s   # Database connection timeout
```

#### Monitoring Defaults
```env
# Default monitor settings
DEFAULT_CHECK_INTERVAL=60s     # Default check interval
DEFAULT_TIMEOUT=30s            # Default request timeout
DEFAULT_RETRY_COUNT=3          # Default retry attempts
DEFAULT_USER_AGENT="Peekaping/1.0 (Monitoring)"

# Push monitor settings
PUSH_TOKEN_LENGTH=32          # Length of push monitor tokens
PUSH_TOKEN_EXPIRY=never       # Push token expiry (never, 1d, 30d, etc.)
```

## Database Configuration

### MongoDB Settings

#### Connection String Format
```javascript
mongodb://[username:password@]host[:port][/database][?options]
```

#### Production Recommendations
```env
# Replica set configuration
DB_REPLICA_SET=rs0
DB_READ_PREFERENCE=primaryPreferred
DB_WRITE_CONCERN=majority

# SSL/TLS settings
DB_SSL=true
DB_SSL_CERT_PATH=/certs/mongodb.pem
DB_SSL_CA_PATH=/certs/ca.pem

# Connection pooling
DB_MIN_POOL_SIZE=5
DB_MAX_POOL_SIZE=50
DB_MAX_IDLE_TIME=30000
```

### Database Indexes

Peekaping automatically creates these indexes for optimal performance:

```javascript
// Monitors collection
db.monitors.createIndex({ "userId": 1 })
db.monitors.createIndex({ "type": 1 })
db.monitors.createIndex({ "active": 1 })

// Heartbeats collection
db.heartbeats.createIndex({ "monitorId": 1, "createdAt": -1 })
db.heartbeats.createIndex({ "createdAt": 1 }, { expireAfterSeconds: 7776000 }) // 90 days

// Users collection
db.users.createIndex({ "email": 1 }, { unique: true })
db.users.createIndex({ "createdAt": 1 })

// Notification channels collection
db.notificationChannels.createIndex({ "userId": 1 })
db.notificationChannels.createIndex({ "type": 1 })
```

## Server Configuration

### HTTP Server Settings

```go
// Server configuration in Go
type ServerConfig struct {
    Port            string        `env:"PORT" default:"8034"`
    ReadTimeout     time.Duration `env:"READ_TIMEOUT" default:"30s"`
    WriteTimeout    time.Duration `env:"WRITE_TIMEOUT" default:"30s"`
    IdleTimeout     time.Duration `env:"IDLE_TIMEOUT" default:"120s"`
    MaxHeaderBytes  int           `env:"MAX_HEADER_BYTES" default:"1048576"`
}
```

### WebSocket Configuration

```env
# WebSocket settings
WS_READ_BUFFER_SIZE=1024      # WebSocket read buffer size
WS_WRITE_BUFFER_SIZE=1024     # WebSocket write buffer size
WS_PING_PERIOD=54s            # Ping period for WebSocket connections
WS_PONG_WAIT=60s              # Pong wait timeout
WS_WRITE_WAIT=10s             # Write message timeout
```

### API Configuration

```env
# API versioning
API_VERSION=v1                # API version prefix
API_BASE_PATH=/api           # API base path

# Request limits
MAX_REQUEST_SIZE=10MB        # Maximum request body size
MAX_FORM_SIZE=32MB          # Maximum form data size
MAX_MULTIPART_MEMORY=10MB   # Maximum multipart form memory

# Response settings
API_RESPONSE_TIMEOUT=30s     # API response timeout
ENABLE_API_DOCS=true        # Enable Swagger documentation
```

## Monitor Configuration

### Default Monitor Settings

```yaml
# Default monitor configuration
defaults:
  http:
    method: GET
    timeout: 30s
    interval: 60s
    retries: 3
    follow_redirects: true
    max_redirects: 5
    user_agent: "Peekaping/1.0"
    expected_status_codes: [200, 201, 202, 203, 204, 205, 206, 207, 208, 226]

  push:
    grace_period: 300s        # Grace period before marking as down
    heartbeat_interval: 60s   # Expected heartbeat interval

  common:
    retry_interval: 60s       # Interval between retries
    notification_delay: 0s    # Delay before sending notifications
    recovery_notification: true  # Send recovery notifications
```

### Monitor Types Configuration

#### HTTP Monitor Options
```json
{
  "type": "http",
  "url": "https://example.com",
  "method": "GET",
  "headers": {
    "Authorization": "Bearer token",
    "User-Agent": "Custom User Agent"
  },
  "body": "",
  "timeout": 30000,
  "interval": 60000,
  "retries": 3,
  "expectedStatusCodes": [200, 201],
  "expectedResponseTime": 5000,
  "keyword": "success",
  "authentication": {
    "type": "bearer",
    "token": "your-token"
  },
  "followRedirects": true,
  "maxRedirects": 5,
  "ignoreTlsErrors": false
}
```

#### Push Monitor Options
```json
{
  "type": "push",
  "name": "My Service Heartbeat",
  "interval": 300,
  "gracePeriod": 600,
  "retries": 1
}
```

## Notification Configuration

### Email (SMTP) Settings

```json
{
  "type": "email",
  "name": "Admin Email",
  "host": "smtp.gmail.com",
  "port": 587,
  "username": "your-email@gmail.com",
  "password": "your-app-password",
  "from": "alerts@yourdomain.com",
  "to": ["admin@yourdomain.com", "team@yourdomain.com"],
  "encryption": "TLS",
  "skipTlsVerify": false,
  "template": {
    "subject": "🚨 {{monitor.name}} is {{heartbeat.status}}",
    "body": "Monitor: {{monitor.name}}\nStatus: {{heartbeat.status}}\nMessage: {{heartbeat.msg}}\nTime: {{heartbeat.time}}"
  }
}
```

### Slack Configuration

```json
{
  "type": "slack",
  "name": "Engineering Alerts",
  "botToken": "xoxb-your-bot-token",
  "channel": "#alerts",
  "username": "Peekaping",
  "iconEmoji": ":warning:",
  "template": {
    "text": "🚨 *{{monitor.name}}* is {{heartbeat.status}}",
    "attachments": [
      {
        "color": "{{#if (eq heartbeat.status 'up')}}good{{else}}danger{{/if}}",
        "fields": [
          {
            "title": "Status",
            "value": "{{heartbeat.status}}",
            "short": true
          },
          {
            "title": "Response Time",
            "value": "{{heartbeat.ping}}ms",
            "short": true
          }
        ]
      }
    ]
  }
}
```

### Webhook Configuration

```json
{
  "type": "webhook",
  "name": "Custom Integration",
  "url": "https://your-webhook-endpoint.com",
  "method": "POST",
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token"
  },
  "body": {
    "event": "{{event}}",
    "monitor": {
      "id": "{{monitor.id}}",
      "name": "{{monitor.name}}",
      "url": "{{monitor.url}}"
    },
    "status": "{{heartbeat.status}}",
    "message": "{{heartbeat.msg}}",
    "timestamp": "{{heartbeat.time}}"
  }
}
```

## Status Page Configuration

### Basic Status Page Settings

```json
{
  "name": "Main Status Page",
  "title": "Our Service Status",
  "description": "Real-time status of our core services",
  "slug": "status",
  "customDomain": "status.yourdomain.com",
  "branding": {
    "logoUrl": "https://yourdomain.com/logo.png",
    "primaryColor": "#3b82f6",
    "backgroundColor": "#f8fafc",
    "textColor": "#1e293b"
  },
  "settings": {
    "showIncidentHistory": true,
    "showPerformanceCharts": true,
    "showUptimePercentages": true,
    "daysOfHistory": 90,
    "refreshInterval": 30
  },
  "contact": {
    "email": "support@yourdomain.com",
    "url": "https://yourdomain.com/support",
    "twitter": "@yourdomain",
    "phone": "+1-555-0123"
  }
}
```

### Custom CSS for Status Pages

```css
/* Status page custom styling */
:root {
  --primary-color: #3b82f6;
  --success-color: #10b981;
  --warning-color: #f59e0b;
  --error-color: #ef4444;
  --background-color: #f8fafc;
  --text-color: #1e293b;
  --border-color: #e2e8f0;
}

.status-page {
  font-family: 'Inter', -apple-system, BlinkMacSystemFont, sans-serif;
  background-color: var(--background-color);
  color: var(--text-color);
}

.status-header {
  background: linear-gradient(135deg, var(--primary-color) 0%, #1e40af 100%);
  padding: 2rem 0;
  text-align: center;
  color: white;
}

.monitor-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  margin: 2rem 0;
}

.monitor-card {
  background: white;
  border-radius: 8px;
  padding: 1rem;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  border-left: 4px solid var(--success-color);
}

.monitor-card.down {
  border-left-color: var(--error-color);
}

.monitor-card.maintenance {
  border-left-color: var(--warning-color);
}
```

## Security Configuration

### Authentication Settings

```yaml
authentication:
  # JWT configuration
  jwt:
    access_token:
      secret: "your-access-token-secret"
      expiry: "15m"
    refresh_token:
      secret: "your-refresh-token-secret"
      expiry: "7d"

  # Session configuration
  session:
    secret: "your-session-secret"
    max_age: "24h"
    secure: true        # Set to true in production with HTTPS
    same_site: "strict"

  # Password policy
  password:
    min_length: 8
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_special: true
    max_attempts: 5     # Account lockout threshold
    lockout_duration: "30m"
```

### Two-Factor Authentication

```yaml
totp:
  issuer: "Peekaping"
  account_name: "{{user.email}}"
  period: 30            # TOTP period in seconds
  digits: 6             # Number of digits in TOTP
  algorithm: "SHA1"     # TOTP algorithm
  skew: 1               # Allowed time skew (periods)
  backup_codes:
    count: 10           # Number of backup codes to generate
    length: 8           # Length of each backup code
```

### Rate Limiting Configuration

```yaml
rate_limiting:
  api:
    requests_per_hour: 1000
    burst_limit: 100

  auth:
    login_attempts: 5
    lockout_duration: "30m"

  status_page:
    requests_per_minute: 60

  push:
    heartbeats_per_minute: 10
```

## Performance Optimization

### Caching Configuration

```env
# Redis caching (optional)
REDIS_URL=redis://redis:6379
REDIS_PASSWORD=your-redis-password
REDIS_DB=0
REDIS_MAX_CONNECTIONS=100

# Cache TTL settings
CACHE_MONITOR_LIST=300        # Monitor list cache (5 minutes)
CACHE_STATUS_PAGE=60          # Status page cache (1 minute)
CACHE_HEARTBEAT_DATA=30       # Heartbeat data cache (30 seconds)
CACHE_USER_SESSION=3600       # User session cache (1 hour)
```

### Database Optimization

```yaml
database:
  # Connection pooling
  pool:
    min_size: 5
    max_size: 50
    max_idle_time: "30m"
    max_lifetime: "1h"

  # Query optimization
  query:
    timeout: "30s"
    max_time_ms: 30000

  # Indexing strategy
  indexes:
    # Automatically created by Peekaping
    auto_create: true

    # Custom indexes for performance
    custom:
      - collection: "heartbeats"
        keys: { "monitorId": 1, "status": 1, "createdAt": -1 }
      - collection: "monitors"
        keys: { "userId": 1, "active": 1 }
```

### Worker Configuration

```yaml
workers:
  # Monitor checking workers
  monitor:
    count: 10           # Number of monitor workers
    timeout: "5m"       # Worker timeout
    max_concurrent: 50  # Max concurrent checks

  # Notification workers
  notification:
    count: 5            # Number of notification workers
    timeout: "30s"      # Notification timeout
    retry_attempts: 3   # Retry failed notifications

  # Cleanup workers
  cleanup:
    interval: "1h"      # Cleanup interval
    batch_size: 1000    # Records to process per batch
    max_age: "90d"      # Maximum age for heartbeat data
```

## Logging Configuration

### Log Levels and Output

```yaml
logging:
  level: "info"         # debug, info, warn, error
  format: "json"        # json, text
  output: "stdout"      # stdout, stderr, file

  # File logging (if output is file)
  file:
    path: "/var/log/peekaping.log"
    max_size: "100MB"
    max_age: "7d"
    max_backups: 5
    compress: true

  # Log rotation
  rotation:
    enabled: true
    interval: "24h"

  # Component-specific logging
  components:
    database: "warn"
    auth: "info"
    monitors: "debug"
    notifications: "info"
```

### Structured Logging Fields

```json
{
  "timestamp": "2024-01-20T10:30:00Z",
  "level": "info",
  "component": "monitor",
  "message": "Monitor check completed",
  "fields": {
    "monitorId": "monitor-123",
    "status": "up",
    "responseTime": 234,
    "statusCode": 200
  }
}
```

## Advanced Configuration

### Multi-tenancy Settings

```yaml
multi_tenancy:
  enabled: false        # Enable multi-tenant mode
  isolation: "database" # database, schema, or row_level

  # Tenant identification
  identification:
    method: "subdomain"  # subdomain, header, or path
    header_name: "X-Tenant-ID"

  # Resource limits per tenant
  limits:
    monitors: 100
    notification_channels: 20
    status_pages: 5
    users: 10
```

### API Gateway Integration

```yaml
api_gateway:
  # Rate limiting headers
  rate_limit_headers:
    enable: true
    header_prefix: "X-RateLimit-"

  # Request tracking
  request_id:
    header: "X-Request-ID"
    generate: true

  # Load balancing
  load_balancer:
    health_check_path: "/health"
    health_check_interval: "30s"
```

### Monitoring Peekaping Itself

```yaml
# Meta-monitoring configuration
meta_monitoring:
  enabled: true

  # Self-monitoring endpoints
  endpoints:
    - name: "Peekaping API Health"
      url: "http://localhost:8034/health"
      interval: "30s"

    - name: "Peekaping Web Interface"
      url: "http://localhost:8383"
      interval: "60s"

    - name: "Database Connection"
      url: "http://localhost:8034/health/db"
      interval: "60s"

  # Resource monitoring
  resources:
    cpu_threshold: 80     # Alert if CPU > 80%
    memory_threshold: 80  # Alert if memory > 80%
    disk_threshold: 90    # Alert if disk > 90%
```

This configuration reference provides a comprehensive overview of all available settings in Peekaping. Start with the basic required settings and gradually add more advanced configuration as your needs grow.
