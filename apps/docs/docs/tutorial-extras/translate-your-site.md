---
sidebar_position: 2
---

# API Reference

Peekaping provides a comprehensive REST API for programmatic access to all features. This reference covers authentication, endpoints, and examples for integrating with Peekaping.

## Base URL and Versioning

All API requests should be made to:
```
https://your-peekaping-instance.com/api/v1
```

The API uses versioning in the URL path. The current version is `v1`.

## Authentication

Peekaping uses JWT (JSON Web Tokens) for API authentication. You can authenticate using:

1. **User credentials** (email/password)
2. **API keys** (recommended for programmatic access)
3. **Existing session** (for web applications)

### Getting an API Key

1. Log into Peekaping web interface
2. Go to **Settings** → **API Keys**
3. Click **"Generate New API Key"**
4. Copy and securely store the key

### Authentication Methods

#### Bearer Token Authentication
```bash
curl -H "Authorization: Bearer your-api-key-here" \
  https://your-peekaping-instance.com/api/v1/monitors
```

#### Login with Credentials
```bash
# Get JWT token
curl -X POST https://your-peekaping-instance.com/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "your-email@example.com",
    "password": "your-password"
  }'

# Response
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-123",
    "email": "your-email@example.com",
    "name": "Your Name"
  }
}
```

## API Endpoints

### Authentication Endpoints

#### POST `/auth/login`
Authenticate user and get JWT tokens.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "accessToken": "jwt-access-token",
  "refreshToken": "jwt-refresh-token",
  "user": {
    "id": "user-id",
    "email": "user@example.com",
    "name": "User Name"
  }
}
```

#### POST `/auth/refresh`
Refresh access token using refresh token.

**Request Body:**
```json
{
  "refreshToken": "your-refresh-token"
}
```

#### POST `/auth/logout`
Logout and invalidate tokens.

### Monitor Endpoints

#### GET `/monitors`
Get all monitors for the authenticated user.

**Query Parameters:**
- `page` (integer): Page number (default: 1)
- `limit` (integer): Items per page (default: 20)
- `type` (string): Filter by monitor type (`http`, `push`)
- `status` (string): Filter by status (`up`, `down`, `pending`)

**Response:**
```json
{
  "monitors": [
    {
      "id": "monitor-123",
      "name": "Website Monitor",
      "type": "http",
      "url": "https://example.com",
      "status": "up",
      "interval": 60,
      "timeout": 30,
      "retries": 3,
      "createdAt": "2024-01-20T10:00:00Z",
      "updatedAt": "2024-01-20T10:30:00Z"
    }
  ],
  "pagination": {
    "page": 1,
    "limit": 20,
    "total": 5,
    "pages": 1
  }
}
```

#### GET `/monitors/{id}`
Get specific monitor by ID.

**Response:**
```json
{
  "id": "monitor-123",
  "name": "Website Monitor",
  "type": "http",
  "url": "https://example.com",
  "method": "GET",
  "headers": {
    "User-Agent": "Peekaping/1.0"
  },
  "expectedStatusCodes": [200, 201],
  "timeout": 30000,
  "interval": 60000,
  "retries": 3,
  "notifications": {
    "enabled": true,
    "channels": ["channel-1", "channel-2"]
  },
  "status": "up",
  "lastCheck": "2024-01-20T10:30:00Z",
  "createdAt": "2024-01-20T10:00:00Z",
  "updatedAt": "2024-01-20T10:30:00Z"
}
```

#### POST `/monitors`
Create a new monitor.

**Request Body:**
```json
{
  "name": "New Website Monitor",
  "type": "http",
  "url": "https://newsite.example.com",
  "method": "GET",
  "timeout": 30000,
  "interval": 60000,
  "retries": 3,
  "expectedStatusCodes": [200],
  "notifications": {
    "enabled": true,
    "channels": ["channel-1"]
  }
}
```

#### PUT `/monitors/{id}`
Update existing monitor.

#### DELETE `/monitors/{id}`
Delete monitor.

#### POST `/monitors/{id}/pause`
Pause monitor.

#### POST `/monitors/{id}/resume`
Resume paused monitor.

### Heartbeat Endpoints

#### GET `/monitors/{id}/heartbeats`
Get heartbeat history for a monitor.

**Query Parameters:**
- `from` (ISO date): Start date
- `to` (ISO date): End date
- `limit` (integer): Number of heartbeats to return

**Response:**
```json
{
  "heartbeats": [
    {
      "id": "heartbeat-123",
      "monitorId": "monitor-123",
      "status": "up",
      "responseTime": 234,
      "statusCode": 200,
      "message": "OK",
      "timestamp": "2024-01-20T10:30:00Z"
    }
  ]
}
```

#### POST `/push/{pushId}`
Send heartbeat for push monitor.

**Request Body (optional):**
```json
{
  "status": "up",
  "message": "Service is running normally",
  "ping": 123
}
```

### Notification Channel Endpoints

#### GET `/notification-channels`
Get all notification channels.

**Response:**
```json
{
  "channels": [
    {
      "id": "channel-123",
      "name": "Email Alerts",
      "type": "email",
      "enabled": true,
      "createdAt": "2024-01-20T09:00:00Z"
    }
  ]
}
```

#### POST `/notification-channels`
Create notification channel.

**Request Body (Email):**
```json
{
  "name": "Gmail Notifications",
  "type": "email",
  "config": {
    "host": "smtp.gmail.com",
    "port": 587,
    "username": "alerts@example.com",
    "password": "app-password",
    "from": "alerts@example.com",
    "to": ["admin@example.com"],
    "encryption": "TLS"
  }
}
```

#### POST `/notification-channels/{id}/test`
Test notification channel.

### Status Page Endpoints

#### GET `/status-pages`
Get all status pages.

#### GET `/status-pages/{slug}`
Get public status page data.

**Response:**
```json
{
  "title": "Service Status",
  "description": "Current status of our services",
  "monitors": [
    {
      "name": "Website",
      "status": "up",
      "responseTime": 234,
      "uptime": {
        "24h": 99.95,
        "7d": 99.89,
        "30d": 99.92
      }
    }
  ],
  "incidents": [
    {
      "title": "Database Performance Issues",
      "status": "resolved",
      "startTime": "2024-01-19T14:00:00Z",
      "endTime": "2024-01-19T15:30:00Z"
    }
  ]
}
```

#### POST `/status-pages`
Create status page.

### Statistics Endpoints

#### GET `/stats/overview`
Get overview statistics.

**Response:**
```json
{
  "totalMonitors": 10,
  "monitorsUp": 9,
  "monitorsDown": 1,
  "avgResponseTime": 345,
  "totalChecks": 12540,
  "uptime": {
    "24h": 99.5,
    "7d": 99.8,
    "30d": 99.9
  }
}
```

#### GET `/monitors/{id}/stats`
Get statistics for specific monitor.

**Query Parameters:**
- `period` (string): Time period (`24h`, `7d`, `30d`, `90d`)

## WebSocket API

Peekaping provides real-time updates via WebSocket connections.

### Connecting to WebSocket

```javascript
const ws = new WebSocket('wss://your-peekaping-instance.com/ws');

// Send authentication
ws.onopen = function() {
  ws.send(JSON.stringify({
    type: 'auth',
    token: 'your-jwt-token'
  }));
};

// Receive updates
ws.onmessage = function(event) {
  const data = JSON.parse(event.data);
  console.log('Received:', data);
};
```

### WebSocket Events

#### Heartbeat Updates
```json
{
  "type": "heartbeat",
  "data": {
    "monitorId": "monitor-123",
    "status": "up",
    "responseTime": 234,
    "timestamp": "2024-01-20T10:30:00Z"
  }
}
```

#### Monitor Status Changes
```json
{
  "type": "monitor_status",
  "data": {
    "monitorId": "monitor-123",
    "oldStatus": "down",
    "newStatus": "up",
    "timestamp": "2024-01-20T10:30:00Z"
  }
}
```

## SDK and Libraries

### JavaScript/TypeScript SDK

```bash
npm install @peekaping/sdk
```

```javascript
import { PeekapingClient } from '@peekaping/sdk';

const client = new PeekapingClient({
  baseUrl: 'https://your-peekaping-instance.com',
  apiKey: 'your-api-key'
});

// Get monitors
const monitors = await client.monitors.list();

// Create monitor
const monitor = await client.monitors.create({
  name: 'My Website',
  type: 'http',
  url: 'https://example.com'
});

// Subscribe to real-time updates
client.on('heartbeat', (data) => {
  console.log('Heartbeat received:', data);
});
```

### Python SDK

```bash
pip install peekaping-python
```

```python
from peekaping import PeekapingClient

client = PeekapingClient(
    base_url='https://your-peekaping-instance.com',
    api_key='your-api-key'
)

# Get monitors
monitors = client.monitors.list()

# Create monitor
monitor = client.monitors.create({
    'name': 'My Website',
    'type': 'http',
    'url': 'https://example.com'
})
```

### Go SDK

```go
import "github.com/0xfurai/peekaping-go"

client := peekaping.NewClient(&peekaping.Config{
    BaseURL: "https://your-peekaping-instance.com",
    APIKey:  "your-api-key",
})

// Get monitors
monitors, err := client.Monitors.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

// Create monitor
monitor, err := client.Monitors.Create(ctx, &peekaping.Monitor{
    Name: "My Website",
    Type: "http",
    URL:  "https://example.com",
})
```

## Rate Limiting

The API implements rate limiting to ensure fair usage:

- **Standard users**: 1000 requests per hour
- **Burst limit**: 100 requests per minute
- **WebSocket connections**: 5 concurrent connections per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1642694400
```

## Error Handling

The API uses standard HTTP status codes and returns consistent error responses:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid monitor configuration",
    "details": {
      "url": ["URL is required"],
      "interval": ["Interval must be at least 30 seconds"]
    }
  }
}
```

### Common Error Codes

- `400 Bad Request`: Invalid request data
- `401 Unauthorized`: Missing or invalid authentication
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Resource not found
- `409 Conflict`: Resource already exists
- `422 Unprocessable Entity`: Validation errors
- `429 Too Many Requests`: Rate limit exceeded
- `500 Internal Server Error`: Server error

## Webhook Integration

Peekaping can send webhooks for various events:

### Webhook Events

- `monitor.up`: Monitor recovered
- `monitor.down`: Monitor went down
- `monitor.created`: New monitor created
- `monitor.updated`: Monitor configuration changed
- `monitor.deleted`: Monitor deleted
- `heartbeat.received`: New heartbeat data (push monitors)

### Webhook Payload

```json
{
  "event": "monitor.down",
  "timestamp": "2024-01-20T10:30:00Z",
  "data": {
    "monitor": {
      "id": "monitor-123",
      "name": "Website Monitor",
      "url": "https://example.com",
      "type": "http"
    },
    "heartbeat": {
      "status": "down",
      "message": "Connection timeout",
      "responseTime": null,
      "timestamp": "2024-01-20T10:30:00Z"
    }
  }
}
```

### Webhook Security

Webhooks include an `X-Peekaping-Signature` header with HMAC-SHA256 signature:

```javascript
const crypto = require('crypto');

function verifyWebhook(payload, signature, secret) {
  const expectedSignature = crypto
    .createHmac('sha256', secret)
    .update(payload)
    .digest('hex');

  return crypto.timingSafeEqual(
    Buffer.from(signature),
    Buffer.from(expectedSignature)
  );
}
```

## Examples and Use Cases

### Automated Monitor Management

```bash
#!/bin/bash
# Script to create monitors for all microservices

SERVICES=("api" "web" "database" "cache")
BASE_URL="https://myapp.com"
API_KEY="your-api-key"

for service in "${SERVICES[@]}"; do
  curl -X POST "https://peekaping.com/api/v1/monitors" \
    -H "Authorization: Bearer $API_KEY" \
    -H "Content-Type: application/json" \
    -d "{
      \"name\": \"$service Service\",
      \"type\": \"http\",
      \"url\": \"$BASE_URL/$service/health\",
      \"interval\": 30000,
      \"notifications\": {
        \"enabled\": true,
        \"channels\": [\"alerts-channel\"]
      }
    }"
done
```

### Custom Dashboard Integration

```javascript
// React component for custom dashboard
import React, { useState, useEffect } from 'react';
import { PeekapingClient } from '@peekaping/sdk';

function MonitorDashboard() {
  const [monitors, setMonitors] = useState([]);
  const client = new PeekapingClient({ apiKey: 'your-api-key' });

  useEffect(() => {
    // Load initial data
    client.monitors.list().then(setMonitors);

    // Subscribe to real-time updates
    client.on('heartbeat', (data) => {
      setMonitors(prev => prev.map(monitor =>
        monitor.id === data.monitorId
          ? { ...monitor, status: data.status, responseTime: data.responseTime }
          : monitor
      ));
    });

    return () => client.disconnect();
  }, []);

  return (
    <div className="monitor-grid">
      {monitors.map(monitor => (
        <div key={monitor.id} className={`monitor-card ${monitor.status}`}>
          <h3>{monitor.name}</h3>
          <div className="status">{monitor.status}</div>
          <div className="response-time">{monitor.responseTime}ms</div>
        </div>
      ))}
    </div>
  );
}
```

### Automated Incident Response

```python
from peekaping import PeekapingClient
import slack_sdk

def handle_monitor_down(event_data):
    """Automated incident response when monitor goes down"""
    monitor = event_data['monitor']

    # Create incident ticket
    create_incident_ticket(
        title=f"{monitor['name']} is down",
        description=f"Monitor {monitor['name']} ({monitor['url']}) is not responding"
    )

    # Scale up resources if it's a performance issue
    if 'timeout' in event_data['heartbeat']['message'].lower():
        scale_service(monitor['name'])

    # Send alert to on-call engineer
    notify_on_call_engineer(monitor)

# Set up webhook handler
@app.route('/webhook', methods=['POST'])
def webhook_handler():
    payload = request.get_json()

    if payload['event'] == 'monitor.down':
        handle_monitor_down(payload['data'])

    return 'OK'
```

## Best Practices

### API Usage
- Use API keys instead of user credentials for production integrations
- Implement proper error handling and retry logic
- Cache responses when appropriate
- Use webhook subscriptions for real-time updates instead of polling

### Security
- Store API keys securely (environment variables, secret managers)
- Use HTTPS for all API requests
- Validate webhook signatures
- Implement rate limiting in your applications

### Performance
- Use pagination for large datasets
- Filter API responses to only include needed data
- Implement connection pooling for high-volume usage
- Monitor your API usage to avoid rate limits

### Monitoring Your Monitors
- Set up monitoring for Peekaping itself
- Monitor API response times and error rates
- Set up alerts for API quota usage
- Regularly review and update monitor configurations

This API reference provides everything you need to integrate Peekaping into your infrastructure and workflows. For additional examples and community-contributed integrations, visit our [GitHub repository](https://github.com/0xfurai/peekaping).
