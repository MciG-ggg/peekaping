---
sidebar_position: 3
---

# Creating Status Pages

Status pages are essential for communicating service status with your users, customers, and stakeholders. Peekaping allows you to create beautiful, branded status pages that show real-time service health information.

## What are Status Pages?

Status pages are public-facing websites that display:
- Real-time status of your services
- Historical uptime data
- Current incidents and maintenance
- Service performance metrics
- Updates and announcements

Popular examples include [GitHub Status](https://githubstatus.com), [Atlassian Status](https://status.atlassian.com), and [Slack Status](https://status.slack.com).

## Benefits of Status Pages

- **Proactive Communication**: Keep users informed before they contact support
- **Transparency**: Build trust by being open about service issues
- **Reduced Support Load**: Users can check status instead of opening tickets
- **Professional Image**: Show that you monitor and care about service quality
- **Customer Confidence**: Demonstrate reliability and professionalism

## Creating Your First Status Page

### 1. Access Status Pages

1. Navigate to **Status Pages** in the sidebar
2. Click **"Add Status Page"** or **"New"**
3. You'll see the status page creation form

### 2. Basic Configuration

#### General Settings:
- **Name**: Internal name for this status page (e.g., "Public API Status")
- **Title**: Public title displayed on the status page (e.g., "Our Service Status")
- **Description**: Brief description of what this page covers
- **Domain/Slug**: URL path for your status page (e.g., `api-status`)

#### Example Configuration:
```
Name: Main Services Status
Title: Our Service Status
Description: Real-time status and performance of our core services
Slug: status
```

This creates a status page accessible at: `https://your-domain.com/status/status`

### 3. Select Monitors

Choose which monitors to display on this status page:

1. **Available Monitors**: List of all your configured monitors
2. **Selected Monitors**: Monitors that will appear on the status page
3. **Display Order**: Drag to reorder how monitors appear

#### Best Practices:
- Only include customer-facing services
- Group related services together
- Use clear, user-friendly monitor names
- Avoid internal or development monitors

### 4. Customization Options

#### Branding:
- **Logo**: Upload your company logo
- **Primary Color**: Main brand color for the status page
- **Header Color**: Background color for the header section
- **Custom CSS**: Advanced styling options

#### Layout Options:
- **Show Incident History**: Display past incidents
- **Show Performance Charts**: Include response time graphs
- **Show Uptime Percentages**: Display uptime statistics
- **Days of History**: How many days of history to show (default: 90)

### 5. Contact Information

Add contact details for your status page:

- **Support Email**: Contact email for questions
- **Support URL**: Link to your help center or support portal
- **Twitter Handle**: Your support Twitter account
- **Phone Number**: Support phone number (optional)

### 6. Save and Publish

1. Review all settings
2. Click **"Create Status Page"**
3. Your status page is now live and accessible via the generated URL

## Status Page Features

### Real-time Status Display

Your status page automatically shows:

- **Current Status**: Green (operational), Yellow (degraded), Red (major outage)
- **Service Names**: User-friendly names for each monitored service
- **Last Check Time**: When each service was last verified
- **Response Times**: Current performance metrics

### Incident Timeline

The status page displays:
- **Active Incidents**: Current ongoing issues
- **Recent Incidents**: Past 30 days of incidents
- **Incident Details**: Start/end times, affected services, description
- **Resolution Updates**: Timeline of incident resolution

### Performance Charts

Interactive charts showing:
- **Response Time Trends**: Performance over time
- **Uptime Percentages**: Reliability statistics
- **Historical Data**: Configurable time ranges (24h, 7d, 30d, 90d)

### Uptime Statistics

Displays uptime percentages for:
- **Today**: Current day's uptime
- **7 Days**: Weekly average uptime
- **30 Days**: Monthly average uptime
- **90 Days**: Quarterly average uptime

## Advanced Configuration

### Custom Domains

To use your own domain (e.g., `status.yourcompany.com`):

1. **DNS Configuration**:
   ```
   CNAME status.yourcompany.com your-peekaping-domain.com
   ```

2. **SSL Certificate**: Ensure your domain has a valid SSL certificate

3. **Update Base URL**: Configure Peekaping with your custom domain

### Multiple Status Pages

You can create multiple status pages for different audiences:

- **Public Status**: Customer-facing services only
- **Internal Status**: All services for team visibility
- **Partner Status**: Services relevant to business partners
- **Regional Status**: Location-specific services

### White-labeling

For a fully branded experience:

#### Custom CSS Examples:

```css
/* Custom header styling */
.status-header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

/* Custom monitor cards */
.monitor-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
}

/* Custom status indicators */
.status-up { background-color: #28a745; }
.status-down { background-color: #dc3545; }
.status-maintenance { background-color: #007bff; }
```

#### Logo Specifications:
- **Format**: PNG, SVG, or JPG
- **Size**: Recommended 200x50px or similar ratio
- **Background**: Transparent or white background works best

## Managing Incidents

### Automatic Incident Creation

When monitors fail, Peekaping automatically:
1. Creates incident entries on status pages
2. Updates incident status based on monitor recovery
3. Calculates incident duration and affected services

### Manual Incident Management

You can also manually manage incidents:

#### Creating Manual Incidents:
1. Go to **Incidents** section
2. Click **"Create Incident"**
3. Fill in incident details:
   - **Title**: Brief description (e.g., "API Performance Issues")
   - **Description**: Detailed explanation of the issue
   - **Affected Services**: Select impacted monitors
   - **Severity**: Minor, Major, or Critical
   - **Status**: Investigating, Identified, Monitoring, Resolved

#### Incident Updates:
- **Post Updates**: Keep users informed with regular updates
- **Change Status**: Update incident status as resolution progresses
- **Resolution**: Mark incidents as resolved with resolution details

### Maintenance Announcements

For planned maintenance:

1. **Schedule Maintenance**: Create maintenance windows in advance
2. **Automatic Updates**: Status page shows maintenance notices
3. **Affected Services**: Clearly indicate which services will be impacted
4. **Duration**: Display expected maintenance duration

## Status Page Analytics

### Built-in Metrics

Peekaping tracks:
- **Page Views**: How many people visit your status page
- **Incident Views**: Which incidents get the most attention
- **Popular Services**: Which monitors users care about most
- **Geographic Data**: Where your users are located

### Integration with Analytics

You can add external analytics:

```html
<!-- Google Analytics -->
<script async src="https://www.googletagmanager.com/gtag/js?id=GA_MEASUREMENT_ID"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'GA_MEASUREMENT_ID');
</script>
```

## Best Practices

### Communication Guidelines

#### During Incidents:
- **Be Transparent**: Acknowledge issues quickly
- **Provide Updates**: Regular updates even if no progress
- **Set Expectations**: Give realistic timeframes for resolution
- **Post-mortem**: Share learnings after major incidents

#### Maintenance Communication:
- **Advance Notice**: Announce maintenance well in advance
- **Clear Impact**: Explain exactly what will be affected
- **Duration Estimates**: Provide expected maintenance windows
- **Updates**: Post updates if maintenance extends beyond planned time

### Status Page Design

#### User Experience:
- **Clean Layout**: Easy to scan and understand at a glance
- **Mobile Friendly**: Ensure it works well on all devices
- **Fast Loading**: Status pages should load quickly
- **Clear Status**: Use obvious colors and language

#### Information Architecture:
- **Priority Order**: Most critical services first
- **Logical Grouping**: Group related services together
- **Clear Naming**: Use customer-facing service names
- **Consistent Language**: Use the same terminology as your product

### Technical Considerations

#### Performance:
- **Caching**: Status pages are cached for fast loading
- **CDN**: Consider using a CDN for global users
- **Fallback**: Ensure status page works even if main service is down

#### SEO and Discoverability:
- **Meta Tags**: Include relevant meta descriptions
- **Structured Data**: Add schema.org markup for search engines
- **Sitemap**: Include status page in your main website sitemap

## Troubleshooting

### Common Issues

#### Status Page Not Loading:
1. Check DNS configuration
2. Verify SSL certificate
3. Ensure Peekaping is running
4. Check firewall rules

#### Monitors Not Appearing:
1. Verify monitors are assigned to status page
2. Check monitor status and configuration
3. Ensure monitors are not paused
4. Review status page filter settings

#### Styling Issues:
1. Validate custom CSS syntax
2. Check for CSS conflicts
3. Test on different browsers
4. Verify image URLs and formats

### Performance Optimization

```nginx
# Nginx configuration for status page caching
location /status/ {
    proxy_pass http://peekaping-backend;
    proxy_cache status_cache;
    proxy_cache_valid 200 1m;
    proxy_cache_key "$scheme$request_method$host$request_uri";
    add_header X-Cache-Status $upstream_cache_status;
}
```

## Integration Examples

### Embed Status Widget

You can embed status information in your main website:

```html
<!-- Status widget embed -->
<iframe
  src="https://status.yourcompany.com/embed"
  width="100%"
  height="300"
  frameborder="0">
</iframe>
```

### API Integration

Fetch status data programmatically:

```javascript
// Fetch current status
fetch('https://your-peekaping-instance.com/api/status-page/your-slug')
  .then(response => response.json())
  .then(data => {
    console.log('Current status:', data);
    // Update your UI based on status
  });
```

### Slack Integration

Get status updates in Slack:

```bash
# Webhook to post status updates
curl -X POST https://hooks.slack.com/services/YOUR/SLACK/WEBHOOK \
  -H "Content-Type: application/json" \
  -d '{
    "text": "🔴 Service Alert: API is experiencing issues",
    "attachments": [{
      "color": "danger",
      "fields": [{
        "title": "Status Page",
        "value": "https://status.yourcompany.com",
        "short": true
      }]
    }]
  }'
```

## Next Steps

Now that you have a status page set up:

1. [Configure maintenance windows](/tutorial-basics/maintenance-windows) for planned downtime
2. [Set up advanced monitoring](/tutorial-extras/manage-docs-versions) features
3. [Review security settings](/tutorial-basics/troubleshooting-guide) for your setup
4. [Explore API integrations](/tutorial-extras/api-reference) for custom workflows

## Examples and Inspiration

### Well-designed Status Pages:
- **GitHub Status**: Clean, technical information
- **Stripe Status**: Excellent incident communication
- **Atlassian Status**: Great use of component grouping
- **Discord Status**: Good mobile experience
- **Cloudflare Status**: Comprehensive performance data

### Status Page Templates:
Peekaping includes several built-in themes:
- **Classic**: Traditional, professional design
- **Modern**: Clean, minimalist layout
- **Dark**: Dark theme for tech-focused brands
- **Colorful**: Vibrant, friendly design

Ready to create your status page? Start with the basics and gradually add more customization as your needs grow!
