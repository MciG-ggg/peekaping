---
sidebar_position: 1
---

# Getting Started: Your First Monitor

This guide will walk you through setting up your first monitor in Peekaping. You'll learn how to monitor a website or API endpoint and configure basic alerting.

## Prerequisites

Before starting, ensure you have:
- Peekaping running (see [Docker Installation](/self-hosting/docker-installation))
- Access to the web interface
- Admin account created during initial setup

## Creating Your First Monitor

### 1. Access the Dashboard

1. Open your browser and navigate to your Peekaping instance (e.g., `http://localhost:8383`)
2. Log in with your admin credentials
3. You'll see the main dashboard with an empty monitors list

### 2. Add a New Monitor

1. Click the **"Add Monitor"** button or navigate to **Monitors** → **New**
2. You'll see the monitor creation form with several tabs

### 3. Configure Basic Settings

#### General Tab

Fill in the basic information:

- **Name**: Give your monitor a descriptive name (e.g., "Company Website")
- **URL**: Enter the URL to monitor (e.g., `https://example.com`)
- **Monitor Type**: Select the type of monitoring:
  - **HTTP/HTTPS**: For websites and APIs
  - **Push**: For services that send heartbeats to Peekaping

#### HTTP/HTTPS Monitor Settings

For HTTP monitors, configure these essential settings:

**Request Settings:**
- **Method**: Choose HTTP method (GET, POST, PUT, DELETE, etc.)
- **Headers**: Add custom headers if needed
- **Body**: Add request body for POST/PUT requests
- **User Agent**: Custom user agent string (optional)

**Response Validation:**
- **Expected Status Codes**: Which HTTP status codes indicate success (default: 200-399)
- **Expected Response Time**: Maximum acceptable response time in milliseconds
- **Keyword**: Text that should be present in the response body
- **Expected Response**: Exact response text to match

**Authentication:**
- **Basic Auth**: Username and password for basic authentication
- **Bearer Token**: API token for bearer authentication
- **Custom Headers**: For custom authentication schemes

### 4. Configure Monitoring Intervals

In the **Intervals** tab:

- **Check Interval**: How often to check the service (e.g., every 60 seconds)
- **Retry Interval**: How long to wait before retrying a failed check
- **Max Redirects**: Maximum number of redirects to follow
- **Timeout**: Request timeout in seconds

:::tip Recommended Settings
- **Check Interval**: 60-300 seconds for most services
- **Retry Interval**: 60 seconds
- **Timeout**: 30 seconds
- **Max Redirects**: 3
:::

### 5. Set Up Notifications

In the **Notifications** tab:

1. **Enable Notifications**: Toggle to enable alerting
2. **Retry Logic**: Number of failed checks before marking as down
3. **Recovery Logic**: Number of successful checks before marking as up
4. **Notification Channels**: Select which channels to notify (you'll need to set these up first)
5. **Resend Interval**: How often to resend notifications while down

:::warning Important
You need to configure notification channels before you can receive alerts. See [Setting Up Notifications](/tutorial-basics/create-notification-channel) for details.
:::

### 6. Advanced Settings (Optional)

#### Proxy Configuration
If you need to route requests through a proxy:
1. Create a proxy configuration first
2. Select it in the **Proxy** dropdown

#### Maintenance Windows
To prevent false alerts during maintenance:
1. Create maintenance windows in advance
2. Associate them with your monitors

### 7. Save Your Monitor

1. Review all settings in the tabs
2. Click **"Create Monitor"** to save
3. Your monitor will appear in the dashboard and start checking immediately

## Understanding Monitor Status

Your monitor can have several states:

- 🟢 **UP**: Service is responding normally
- 🔴 **DOWN**: Service is not responding or failing checks
- 🟡 **PENDING**: Newly created monitor, not enough data yet
- 🔵 **MAINTENANCE**: Monitor is in maintenance window
- ⚫ **PAUSED**: Monitor is manually paused

## Monitor Dashboard Features

Once your monitor is created, you'll see:

### Real-time Status
- Current status indicator
- Response time graph
- Uptime percentage (24h, 7d, 30d)
- Recent check history

### Heartbeat Chart
- Visual representation of check results
- Green bars = successful checks
- Red bars = failed checks
- Gray bars = maintenance periods

### Statistics
- **Average Response Time**: Rolling average
- **Uptime Percentage**: Success rate over time periods
- **Total Checks**: Number of checks performed
- **Incidents**: Number of downtime periods

## Push Monitors

For applications that need to send heartbeats to Peekaping:

### 1. Create a Push Monitor
1. Select **"Push"** as the monitor type
2. Configure basic settings (name, etc.)
3. Set the **Expected Interval** (how often heartbeats should arrive)
4. Save the monitor

### 2. Get the Heartbeat URL
After creation, you'll receive a unique URL like:
```
https://your-peekaping-instance.com/api/push/your-unique-id
```

### 3. Send Heartbeats
Your application should send HTTP requests to this URL:

```bash
# Simple heartbeat
curl https://your-peekaping-instance.com/api/push/your-unique-id

# Heartbeat with status
curl -X POST https://your-peekaping-instance.com/api/push/your-unique-id \
  -H "Content-Type: application/json" \
  -d '{"status": "up", "msg": "All systems operational"}'
```

### Heartbeat Parameters
- **status**: "up", "down", or "maintenance"
- **msg**: Optional message describing the status
- **ping**: Response time in milliseconds

## Monitor Management

### Editing Monitors
1. Click on a monitor name or use the **Edit** button
2. Modify any settings as needed
3. Click **"Update Monitor"** to save changes

### Pausing Monitors
1. Click the **Pause** button to temporarily stop monitoring
2. Click **Resume** to restart monitoring
3. Paused monitors don't count toward downtime

### Deleting Monitors
1. Click the **Delete** button (⚠️ **Warning**: This is permanent)
2. Confirm deletion
3. All historical data will be removed

## Best Practices

### Monitor Naming
- Use descriptive names that identify the service
- Include environment if you have multiple (e.g., "API - Production")
- Consider using prefixes for grouping (e.g., "Website - ", "API - ")

### Check Intervals
- **Critical services**: 30-60 seconds
- **Important services**: 60-300 seconds
- **Less critical services**: 300-600 seconds
- **Internal tools**: 600+ seconds

### Response Time Thresholds
- **Websites**: 2000-5000ms
- **APIs**: 1000-3000ms
- **Database queries**: 500-2000ms
- **Internal services**: Based on SLA requirements

### Retry Logic
- **Web services**: 3-5 retries
- **APIs**: 2-3 retries
- **Unreliable networks**: 5+ retries

## Common Issues and Solutions

### Monitor Shows as Down But Service is Working

**Possible causes:**
- Timeout too low for slow responses
- Network connectivity issues
- Firewall blocking requests
- SSL certificate problems

**Solutions:**
- Increase timeout value
- Check network connectivity
- Verify SSL certificates
- Add IP whitelist if needed

### False Positives During Load Spikes

**Solutions:**
- Increase retry count
- Extend timeout values
- Use longer check intervals during peak hours
- Set up maintenance windows for known load periods

### Missing Notifications

**Check:**
- Notification channels are configured
- Monitor has notifications enabled
- Retry logic settings
- Notification channel limits/quotas

## Next Steps

Now that you have your first monitor set up:

1. [Configure notification channels](/tutorial-basics/create-notification-channel) to receive alerts
2. [Create a status page](/tutorial-basics/create-status-page) to share status with users
3. [Set up maintenance windows](/tutorial-basics/maintenance-windows) for planned downtime
4. [Explore advanced monitoring features](/tutorial-extras/manage-docs-versions)

## Getting Help

If you encounter issues:
- Check the [Troubleshooting Guide](/tutorial-basics/troubleshooting-guide)
- Review monitor logs in the dashboard
- Join the community discussions
- Report bugs on [GitHub](https://github.com/0xfurai/peekaping/issues)
