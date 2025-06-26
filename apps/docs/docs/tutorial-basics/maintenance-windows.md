---
sidebar_position: 4
---

# Maintenance Windows & Advanced Features

Maintenance windows are essential for preventing false alerts during planned downtime. This guide covers setting up maintenance windows and exploring Peekaping's advanced monitoring features.

## Understanding Maintenance Windows

Maintenance windows tell Peekaping when planned downtime is expected, preventing false alerts and keeping your uptime statistics accurate during legitimate maintenance periods.

### Benefits of Maintenance Windows

- **Prevent False Alerts**: No notifications during planned downtime
- **Accurate Uptime Stats**: Maintenance time doesn't count against uptime
- **Team Coordination**: Clear communication about planned maintenance
- **Professional Status Pages**: Customers see "maintenance" instead of "down"
- **Compliance**: Meet SLA requirements for planned vs unplanned downtime

## Creating Maintenance Windows

### 1. Access Maintenance

1. Navigate to **Maintenance** in the sidebar
2. Click **"Add Maintenance"** or **"New"**
3. You'll see the maintenance window creation form

### 2. Basic Configuration

#### General Settings:
- **Title**: Descriptive name for the maintenance (e.g., "Database Migration")
- **Description**: Detailed description of what's being maintained
- **Start Time**: When maintenance begins
- **End Time**: When maintenance is expected to end
- **Timezone**: Timezone for the maintenance window

#### Example Configuration:
```
Title: Weekly Server Maintenance
Description: Routine server updates and security patches
Start Time: 2024-01-20 02:00:00
End Time: 2024-01-20 04:00:00
Timezone: UTC
```

### 3. Affected Monitors

Select which monitors should enter maintenance mode:

- **All Monitors**: Apply to all current and future monitors
- **Selected Monitors**: Choose specific monitors
- **Monitor Groups**: Apply to logical groups of monitors

#### Best Practices:
- Only include monitors that will actually be affected
- Consider dependencies (if API is down, dependent services will fail too)
- Group related services for easier management

### 4. Recurring Maintenance

For regular maintenance windows:

#### Recurrence Options:
- **One-time**: Single maintenance window
- **Daily**: Repeats every day
- **Weekly**: Repeats every week on the same day
- **Monthly**: Repeats monthly on the same date
- **Custom**: Advanced cron-based scheduling

#### Cron Expression Examples:
```bash
# Every Sunday at 2 AM UTC
0 2 * * 0

# First Monday of every month at 3 AM
0 3 * * 1#1

# Every weekday at 6 AM
0 6 * * 1-5

# Every 6 hours
0 */6 * * *
```

### 5. Notification Settings

Configure how users are notified about maintenance:

- **Advance Notice**: Send notifications before maintenance starts
- **Start Notification**: Alert when maintenance begins
- **End Notification**: Alert when maintenance completes
- **Channels**: Which notification channels to use
- **Status Page Updates**: Automatically update status pages

## Advanced Maintenance Features

### 1. Maintenance Templates

Create reusable templates for common maintenance types:

#### Database Maintenance Template:
```yaml
title: "Database Maintenance - {date}"
description: "Routine database optimization and backup verification"
duration: 2 hours
affected_services:
  - "Database Cluster"
  - "API Services"
  - "Web Application"
notification_advance: 24 hours
```

#### Security Update Template:
```yaml
title: "Security Updates - {date}"
description: "Installation of critical security patches"
duration: 30 minutes
affected_services:
  - "All Services"
notification_advance: 4 hours
auto_extend: true
```

### 2. Emergency Maintenance

For unplanned maintenance that needs to start immediately:

1. **Quick Create**: Simplified form for urgent maintenance
2. **Immediate Start**: Maintenance begins as soon as it's created
3. **Auto-notification**: Automatically notify all configured channels
4. **Flexible End Time**: Extend or shorten as needed

#### Emergency Maintenance Workflow:
1. Click **"Emergency Maintenance"**
2. Select affected monitors
3. Add brief description
4. Click **"Start Now"**
5. Update duration as situation develops

### 3. Maintenance Extensions

Sometimes maintenance takes longer than expected:

#### Extension Options:
- **Manual Extension**: Manually extend the maintenance window
- **Auto-extension**: Automatically extend if monitors are still down
- **Progressive Extension**: Extend in 30-minute increments
- **Maximum Duration**: Set limits on how long maintenance can run

#### Best Practices:
- Communicate extensions to users promptly
- Set reasonable maximum durations
- Review why extensions were needed for future planning

## Proxy Configuration

Peekaping supports routing monitoring requests through HTTP proxies, useful for network security, geographic distribution, or accessing internal services.

### 1. Creating Proxy Configurations

#### Access Proxies:
1. Navigate to **Proxies** in the sidebar
2. Click **"Add Proxy"** or **"New"**
3. Configure proxy settings

#### Proxy Settings:
- **Name**: Descriptive name (e.g., "Corporate Proxy")
- **Host**: Proxy server hostname or IP
- **Port**: Proxy server port
- **Protocol**: HTTP or HTTPS
- **Authentication**: Username/password if required

#### Example Configuration:
```
Name: Corporate Proxy
Host: proxy.company.com
Port: 8080
Protocol: HTTP
Username: monitoring-user
Password: secure-password
```

### 2. Proxy Types

#### Forward Proxy:
- Routes requests through a central proxy server
- Useful for corporate networks with firewall restrictions
- Can provide additional security and logging

#### Regional Proxies:
- Monitor services from different geographic locations
- Detect regional connectivity issues
- Provide location-specific response times

#### Load-balanced Proxies:
- Distribute monitoring load across multiple proxy servers
- Provide redundancy for critical monitoring paths
- Improve overall monitoring reliability

### 3. Assigning Proxies to Monitors

1. **Edit Monitor**: Go to monitor configuration
2. **Proxy Settings**: Find the proxy dropdown
3. **Select Proxy**: Choose from configured proxies
4. **Test Configuration**: Verify proxy works correctly

#### Use Cases:
- **Internal Services**: Access services behind corporate firewalls
- **Geographic Monitoring**: Monitor from different regions
- **Security Compliance**: Route through approved proxy servers
- **Network Segmentation**: Separate monitoring traffic

## Advanced Monitor Configuration

### 1. Custom Headers and Authentication

#### Custom Headers:
Add custom HTTP headers for:
- API keys and tokens
- Custom authentication schemes
- Request tracking and debugging
- Content negotiation

```http
Authorization: Bearer your-api-token
X-API-Key: your-api-key
X-Request-ID: unique-request-id
Accept: application/json
User-Agent: Peekaping/1.0 (Monitoring)
```

#### Authentication Methods:
- **Basic Auth**: Username and password
- **Bearer Token**: API tokens and JWTs
- **Custom Headers**: Proprietary authentication schemes
- **Client Certificates**: Mutual TLS authentication

### 2. Advanced Response Validation

#### Response Content Validation:
- **Keyword Matching**: Check for specific text in response
- **JSON Path Validation**: Validate JSON response structure
- **Regex Patterns**: Complex pattern matching
- **Response Size**: Validate response body size

#### Example Validations:
```javascript
// Check API response structure
$.status === "ok" && $.data.length > 0

// Validate specific content
response.includes("Server Status: OK")

// Check for absence of error indicators
!response.includes("error") && !response.includes("exception")
```

#### Status Code Ranges:
- **Success Codes**: 200-299 (default)
- **Redirect Codes**: 300-399 (often acceptable)
- **Client Error**: 400-499 (usually failures)
- **Server Error**: 500-599 (definitely failures)
- **Custom Ranges**: Define your own acceptable ranges

### 3. Certificate Monitoring

Monitor SSL/TLS certificate expiration:

#### Certificate Checks:
- **Expiration Date**: Alert before certificates expire
- **Certificate Chain**: Validate full certificate chain
- **Certificate Authority**: Check issuing CA
- **Certificate Transparency**: Verify CT log inclusion

#### Expiration Alerts:
- **30 Days Warning**: Early warning for renewal planning
- **7 Days Critical**: Urgent renewal required
- **1 Day Emergency**: Certificate expires very soon
- **Post-expiration**: Certificate has already expired

### 4. Performance Thresholds

Set custom performance expectations:

#### Response Time Monitoring:
- **Warning Threshold**: Slow but acceptable (e.g., 2000ms)
- **Critical Threshold**: Unacceptably slow (e.g., 5000ms)
- **Baseline Monitoring**: Compare against historical averages
- **Percentile Monitoring**: 95th percentile response times

#### Uptime Requirements:
- **99.9% Uptime**: High availability requirement
- **99.95% Uptime**: Mission-critical services
- **99.99% Uptime**: Ultra-high availability
- **Custom SLAs**: Match your business requirements

## Security and Compliance

### 1. User Management and RBAC

#### User Roles:
- **Admin**: Full system access and configuration
- **Editor**: Can create and modify monitors
- **Viewer**: Read-only access to dashboards
- **Guest**: Limited access to specific status pages

#### Permission Levels:
- **Monitor Management**: Create, edit, delete monitors
- **Notification Configuration**: Set up notification channels
- **Status Page Management**: Create and customize status pages
- **User Management**: Add/remove users and set permissions
- **System Configuration**: Global settings and maintenance

### 2. Two-Factor Authentication (2FA)

Enhanced security for user accounts:

#### Setup Process:
1. **Enable 2FA**: In user security settings
2. **Scan QR Code**: Use authenticator app (Google Authenticator, Authy)
3. **Verify Setup**: Enter verification code
4. **Backup Codes**: Save emergency backup codes

#### Supported Authenticators:
- Google Authenticator
- Authy
- Microsoft Authenticator
- 1Password
- Bitwarden

### 3. API Security

Secure API access:

#### API Key Management:
- **Generate API Keys**: Create keys for programmatic access
- **Scope Limitations**: Limit API key permissions
- **Rotation Policy**: Regular key rotation
- **Usage Monitoring**: Track API key usage

#### Rate Limiting:
- **Request Limits**: Prevent API abuse
- **IP-based Limiting**: Limit requests per IP address
- **User-based Limiting**: Limit requests per user/API key
- **Burst Protection**: Handle traffic spikes gracefully

## Monitoring Best Practices

### 1. Monitor Organization

#### Naming Conventions:
```
[Environment] - [Service] - [Component]
Examples:
- "Production - Website - Homepage"
- "Staging - API - Authentication"
- "Dev - Database - Primary Cluster"
```

#### Service Grouping:
- **Critical Services**: Customer-facing, revenue-impacting
- **Important Services**: Internal tools, dependencies
- **Supporting Services**: Monitoring, logging, backups
- **Development Services**: Staging, testing environments

### 2. Alert Tuning

#### Reducing False Positives:
- **Appropriate Timeouts**: Match service characteristics
- **Retry Logic**: Multiple failures before alerting
- **Maintenance Windows**: Planned downtime handling
- **Environmental Factors**: Consider network conditions

#### Escalation Strategies:
- **Immediate Alerts**: Critical services (< 1 minute)
- **Delayed Alerts**: Less critical services (5-15 minutes)
- **Progressive Escalation**: Escalate based on duration
- **Recovery Notifications**: Confirm when issues are resolved

### 3. Performance Optimization

#### Monitoring Efficiency:
- **Appropriate Intervals**: Balance between detection speed and load
- **Resource Usage**: Monitor Peekaping's own resource consumption
- **Database Optimization**: Regular cleanup of old data
- **Network Efficiency**: Optimize network requests

#### Scaling Considerations:
- **Monitor Limits**: Understand system limits
- **Resource Planning**: Plan for growth
- **Geographic Distribution**: Consider multiple monitoring locations
- **High Availability**: Ensure monitoring system reliability

## Next Steps

With advanced features configured:

1. [Review troubleshooting guide](/tutorial-basics/troubleshooting-guide) for common issues
2. [Explore configuration options](/tutorial-extras/manage-docs-versions) for fine-tuning
3. [Check API documentation](/tutorial-extras/api-reference) for integration possibilities
4. [Set up monitoring for Peekaping itself](/tutorial-extras/manage-docs-versions) (meta-monitoring)

## Getting Help

If you encounter issues with advanced features:
- Check the troubleshooting guide for common solutions
- Review logs for specific error messages
- Join community discussions for tips and tricks
- Report bugs or request features on GitHub

Remember: Start simple and gradually add complexity as your monitoring needs grow!
