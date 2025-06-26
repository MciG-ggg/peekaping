---
sidebar_position: 2
---

# Setting Up Notification Channels

Notification channels are essential for staying informed about your service status. Peekaping supports multiple notification methods including email, Slack, Telegram, webhooks, and more.

## Overview

Notification channels allow you to:
- Receive alerts when services go down or recover
- Get notified about maintenance windows
- Stay informed about important system events
- Route different alerts to different teams or channels

## Supported Notification Types

Peekaping supports these notification channels:

- 📧 **Email (SMTP)** - Email notifications via SMTP server
- 💬 **Slack** - Send alerts to Slack channels
- 📱 **Telegram** - Push notifications via Telegram bot
- 🔗 **Webhooks** - Custom HTTP webhooks for integrations
- 📢 **Ntfy** - Push notifications via Ntfy service

## Creating Notification Channels

### 1. Access Notification Channels

1. Navigate to **Notification Channels** in the sidebar
2. Click **"Add Notification Channel"** or **"New"**
3. Choose your notification type from the list

### 2. Email (SMTP) Configuration

Email is the most common notification method:

#### Required Settings:
- **Name**: Descriptive name for this channel (e.g., "Admin Email")
- **SMTP Host**: Your email server hostname (e.g., `smtp.gmail.com`)
- **SMTP Port**: Server port (usually 587 for TLS, 465 for SSL, 25 for plain)
- **Username**: Your email username/address
- **Password**: Your email password or app password
- **From Email**: Email address to send from
- **To Email**: Email address(es) to send to (comma-separated for multiple)

#### Security Settings:
- **Encryption**: Choose TLS, SSL, or None
- **Skip TLS Verify**: Only enable if using self-signed certificates

#### Example Gmail Configuration:
```
Name: Gmail Notifications
SMTP Host: smtp.gmail.com
SMTP Port: 587
Username: your-email@gmail.com
Password: your-app-password
From Email: your-email@gmail.com
To Email: admin@yourcompany.com
Encryption: TLS
```

:::tip Gmail App Passwords
For Gmail, you'll need to:
1. Enable 2-factor authentication
2. Generate an app password
3. Use the app password instead of your regular password
:::

### 3. Slack Configuration

Send alerts directly to Slack channels:

#### Setup Steps:
1. **Create Slack App**:
   - Go to [Slack API](https://api.slack.com/apps)
   - Click "Create New App" → "From scratch"
   - Name your app (e.g., "Peekaping") and select workspace

2. **Configure Bot Permissions**:
   - Go to "OAuth & Permissions"
   - Add these scopes under "Bot Token Scopes":
     - `chat:write`
     - `chat:write.public`
   - Install app to workspace

3. **Get Bot Token**:
   - Copy the "Bot User OAuth Token" (starts with `xoxb-`)

4. **Configure in Peekaping**:
   - **Name**: Channel name (e.g., "Engineering Alerts")
   - **Bot Token**: Paste your bot token
   - **Channel**: Channel name (e.g., `#alerts`) or user ID
   - **Username**: Bot display name (optional)

#### Example Configuration:
```
Name: Engineering Alerts
Bot Token: xoxb-your-bot-token-here
Channel: #engineering-alerts
Username: Peekaping Bot
```

### 4. Telegram Configuration

Get instant notifications on your mobile device:

#### Setup Steps:
1. **Create Bot**:
   - Message [@BotFather](https://t.me/botfather) on Telegram
   - Send `/newbot` command
   - Follow instructions to create bot
   - Save the bot token

2. **Get Chat ID**:
   - Start chat with your bot
   - Send a message to the bot
   - Visit: `https://api.telegram.org/bot<YOUR_BOT_TOKEN>/getUpdates`
   - Find your chat ID in the response

3. **Configure in Peekaping**:
   - **Name**: Descriptive name
   - **Bot Token**: Your bot token from BotFather
   - **Chat ID**: Your chat ID or group chat ID

#### Example Configuration:
```
Name: Personal Alerts
Bot Token: 123456789:ABCdefGHIjklMNOpqrsTUVwxyz
Chat ID: 123456789
```

### 5. Webhook Configuration

Integrate with custom systems or services:

#### Configuration:
- **Name**: Webhook name
- **URL**: Endpoint URL to send POST requests
- **Method**: HTTP method (usually POST)
- **Content Type**: Request content type
- **Headers**: Custom headers (optional)
- **Body Template**: Custom JSON payload template

#### Webhook Payload

Peekaping sends webhook payloads in this format:

```json
{
  "monitor": {
    "id": "monitor-id",
    "name": "Website Monitor",
    "url": "https://example.com",
    "type": "http"
  },
  "heartbeat": {
    "status": "down",
    "time": "2024-01-15T10:30:00Z",
    "msg": "HTTP Error: 500 Internal Server Error",
    "ping": 5000,
    "important": true
  },
  "event": "monitor.down"
}
```

#### Example Discord Webhook:
```json
{
  "content": "🚨 **{{monitor.name}}** is {{heartbeat.status}}",
  "embeds": [{
    "title": "{{monitor.name}}",
    "description": "{{heartbeat.msg}}",
    "color": "{{#if (eq heartbeat.status 'up')}}3066993{{else}}15158332{{/if}}",
    "fields": [
      {
        "name": "Status",
        "value": "{{heartbeat.status}}",
        "inline": true
      },
      {
        "name": "Response Time",
        "value": "{{heartbeat.ping}}ms",
        "inline": true
      }
    ],
    "timestamp": "{{heartbeat.time}}"
  }]
}
```

### 6. Ntfy Configuration

Simple push notifications:

#### Configuration:
- **Name**: Channel name
- **Topic**: Ntfy topic name
- **Server URL**: Ntfy server (default: https://ntfy.sh)
- **Username**: Username (if using authentication)
- **Password**: Password (if using authentication)
- **Priority**: Notification priority (1-5)

## Testing Notification Channels

After creating a channel:

1. **Save** the configuration
2. Click the **"Test"** button
3. Check that you receive the test notification
4. If the test fails, review your configuration

### Common Test Issues:

#### Email Not Received:
- Check spam/junk folder
- Verify SMTP settings
- Ensure firewall allows SMTP traffic
- Try different encryption settings

#### Slack Not Working:
- Verify bot token is correct
- Check bot permissions
- Ensure channel name is correct
- Bot must be invited to private channels

#### Telegram Silent:
- Verify bot token
- Check chat ID (positive for users, negative for groups)
- Ensure you've started a chat with the bot

## Assigning Channels to Monitors

After creating notification channels:

1. **Edit your monitor**
2. Go to the **Notifications** tab
3. **Enable notifications**
4. **Select channels** to notify
5. Configure **notification settings**:
   - **Retry Logic**: Failed checks before alerting
   - **Recovery Logic**: Successful checks before recovery alert
   - **Resend Interval**: How often to resend while down
6. **Save** changes

## Notification Rules and Logic

### When Notifications Are Sent:

- **Monitor Down**: After configured number of failed retries
- **Monitor Recovery**: After configured number of successful checks
- **Maintenance Start**: When maintenance window begins
- **Maintenance End**: When maintenance window ends

### Resend Logic:

- Notifications are resent based on the **Resend Interval**
- Useful for ensuring critical alerts aren't missed
- Set to 0 to disable resending

### Important vs Regular Notifications:

Peekaping distinguishes between:
- **Important**: Status changes (up ↔ down)
- **Regular**: Maintenance windows, manual events

This allows you to route critical alerts differently from informational messages.

## Advanced Configuration

### Multiple Channels per Monitor

You can assign multiple notification channels to a single monitor:

- **Primary channels**: Critical team members
- **Secondary channels**: Management, status boards
- **Escalation channels**: If issues persist

### Channel Groups

For complex setups, consider creating multiple channels for different purposes:

- **Critical Alerts**: Immediate attention required
- **Maintenance Notices**: Planned downtime notifications
- **Status Updates**: General status information
- **Recovery Notifications**: Service restoration alerts

### Conditional Notifications

While Peekaping doesn't have built-in conditional logic, you can achieve this with webhooks:

```javascript
// Example webhook that only alerts for critical services
if (monitor.name.includes('CRITICAL') && heartbeat.status === 'down') {
  // Send to emergency channel
} else if (heartbeat.status === 'down') {
  // Send to regular channel
}
```

## Best Practices

### Email Configuration:
- Use dedicated monitoring email addresses
- Set up email filtering rules
- Consider using distribution lists
- Test email delivery regularly

### Slack Integration:
- Create dedicated monitoring channels
- Use thread notifications for follow-ups
- Pin important status messages
- Set up channel descriptions with escalation info

### Webhook Security:
- Use HTTPS endpoints
- Implement webhook signature verification
- Rate limit webhook endpoints
- Log webhook events for debugging

### General Guidelines:
- **Test all channels regularly**
- **Keep contact information current**
- **Document escalation procedures**
- **Review notification settings periodically**
- **Use descriptive channel names**

## Troubleshooting

### Email Issues:
```bash
# Test SMTP connectivity
telnet smtp.gmail.com 587

# Check TLS/SSL settings
openssl s_client -connect smtp.gmail.com:587 -starttls smtp
```

### Webhook Debugging:
- Use tools like [webhook.site](https://webhook.site) for testing
- Check response status codes and error messages
- Verify content-type headers
- Test payload format with curl

### Rate Limiting:
- Most services have rate limits
- Space out notifications appropriately
- Consider using resend intervals wisely
- Monitor for 429 (Too Many Requests) errors

## Next Steps

With notification channels configured:

1. [Create status pages](/tutorial-basics/create-status-page) for public communication
2. [Set up maintenance windows](/tutorial-basics/maintenance-windows) for planned downtime
3. [Configure advanced monitoring](/tutorial-extras/manage-docs-versions) features
4. [Review security settings](/tutorial-basics/troubleshooting-guide) for your setup

Need help? Check the [troubleshooting guide](/tutorial-basics/troubleshooting-guide) or join our community discussions!
