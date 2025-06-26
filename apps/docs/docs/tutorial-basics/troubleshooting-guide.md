---
sidebar_position: 5
---

# Troubleshooting Guide

This guide covers common issues you might encounter while using Peekaping and provides step-by-step solutions to resolve them.

## Quick Diagnostics

Before diving into specific issues, run through these quick checks:

### Health Check Commands

```bash
# Check if all containers are running
docker compose ps

# View recent logs
docker compose logs --tail=50

# Check container health
docker compose exec server curl -f http://localhost:8034/health

# Check database connectivity
docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"
```

### System Resources

```bash
# Check disk space
df -h

# Check memory usage
free -h

# Check CPU usage
top

# Check network connectivity
ping google.com
```

## Common Installation Issues

### Docker Compose Fails to Start

#### Error: "Port already in use"

**Problem**: The default ports (8383, 8034, 27017) are already occupied.

**Solution**:
1. Check what's using the ports:
   ```bash
   netstat -tulpn | grep :8383
   netstat -tulpn | grep :8034
   netstat -tulpn | grep :27017
   ```

2. Either stop the conflicting service or change Peekaping ports in your `.env` file:
   ```env
   # Use different ports
   PORT=9034
   CLIENT_PORT=9383
   DB_PORT=27018
   ```

3. Update your `docker-compose.yml` to use the new ports:
   ```yaml
   web:
     ports:
       - "9383:80"
   ```

#### Error: "Permission denied" or "Cannot connect to Docker daemon"

**Problem**: Docker permission issues.

**Solution**:
```bash
# Add your user to docker group
sudo usermod -aG docker $USER

# Restart session or run:
newgrp docker

# Or run commands with sudo (not recommended for production)
sudo docker compose up -d
```

#### Error: "Network peekaping-network already exists"

**Problem**: Previous Docker network conflicts.

**Solution**:
```bash
# Remove existing network
docker network rm peekaping-network

# Or use different network name in docker-compose.yml
networks:
  peekaping-network-new:
    driver: bridge
```

### Environment Configuration Issues

#### Error: "Database connection failed"

**Problem**: MongoDB connection issues.

**Solutions**:

1. **Check MongoDB logs**:
   ```bash
   docker compose logs mongodb
   ```

2. **Verify environment variables**:
   ```bash
   # Check if .env file is properly loaded
   docker compose config
   ```

3. **Test MongoDB connectivity**:
   ```bash
   # Test from within server container
   docker compose exec server ping mongodb

   # Test MongoDB authentication
   docker compose exec mongodb mongosh \
     --username ${DB_USER} \
     --password ${DB_PASSWORD} \
     --authenticationDatabase admin
   ```

4. **Reset MongoDB data** (⚠️ This deletes all data):
   ```bash
   docker compose down
   docker volume rm peekaping_mongodb_data
   docker compose up -d
   ```

#### Error: "JWT token errors" or "Authentication failed"

**Problem**: Invalid JWT configuration.

**Solution**:
1. Generate new secure JWT secrets:
   ```bash
   # Generate random strings (64+ characters)
   openssl rand -base64 64
   ```

2. Update `.env` file with new secrets:
   ```env
   ACCESS_TOKEN_SECRET_KEY=your-new-secret-here
   REFRESH_TOKEN_SECRET_KEY=your-other-new-secret-here
   ```

3. Restart services:
   ```bash
   docker compose restart server
   ```

## Application Issues

### Monitors Not Working

#### Monitor shows "DOWN" but service is accessible

**Possible Causes & Solutions**:

1. **Timeout too low**:
   - Increase timeout in monitor settings (try 30-60 seconds)
   - Check service response times manually

2. **Network connectivity**:
   ```bash
   # Test from Peekaping container
   docker compose exec server curl -I https://your-website.com

   # Check DNS resolution
   docker compose exec server nslookup your-website.com
   ```

3. **SSL/TLS certificate issues**:
   ```bash
   # Check certificate validity
   openssl s_client -connect your-website.com:443 -servername your-website.com

   # Skip certificate verification in monitor (not recommended for production)
   # Enable "Skip TLS Verify" in monitor settings
   ```

4. **User Agent blocking**:
   - Some services block monitoring tools
   - Try changing User Agent in monitor settings
   - Use a browser-like User Agent string

5. **Rate limiting**:
   - Service might be rate limiting requests
   - Increase check interval
   - Use different IP addresses or proxies

#### Push monitors not receiving heartbeats

**Solutions**:

1. **Verify heartbeat URL**:
   ```bash
   # Test heartbeat endpoint
   curl -X POST https://your-peekaping-instance.com/api/push/your-monitor-id
   ```

2. **Check firewall settings**:
   - Ensure Peekaping is accessible from your application
   - Check if ports are open

3. **Review application logs**:
   - Check if heartbeat requests are being sent
   - Look for HTTP errors in application logs

### Notification Issues

#### Email notifications not being sent

**Debugging Steps**:

1. **Test SMTP configuration**:
   ```bash
   # Test SMTP connectivity
   telnet smtp.gmail.com 587

   # Test with openssl for TLS
   openssl s_client -connect smtp.gmail.com:587 -starttls smtp
   ```

2. **Check email logs**:
   ```bash
   # Look for email-related errors
   docker compose logs server | grep -i mail
   ```

3. **Common Gmail issues**:
   - Enable 2FA and use App Password
   - Check "Less secure app access" (deprecated)
   - Verify account isn't locked

4. **Test with different SMTP settings**:
   ```env
   # Try different ports/encryption
   SMTP_PORT=465  # For SSL
   SMTP_ENCRYPTION=SSL
   ```

#### Slack notifications not working

**Debugging Steps**:

1. **Verify bot token**:
   ```bash
   # Test Slack API with your token
   curl -X POST https://slack.com/api/auth.test \
     -H "Authorization: Bearer xoxb-your-token"
   ```

2. **Check bot permissions**:
   - Ensure bot has `chat:write` permission
   - Bot must be invited to private channels
   - Verify channel name is correct (include #)

3. **Test bot in Slack**:
   - Try sending a direct message to the bot
   - Check if bot appears online in Slack

#### Webhook notifications failing

**Debugging Steps**:

1. **Test webhook endpoint**:
   ```bash
   # Test your webhook URL
   curl -X POST https://your-webhook-url.com \
     -H "Content-Type: application/json" \
     -d '{"test": "message"}'
   ```

2. **Check webhook logs**:
   ```bash
   # Look for webhook errors
   docker compose logs server | grep -i webhook
   ```

3. **Verify webhook format**:
   - Check if your endpoint expects specific format
   - Review webhook payload in Peekaping logs

### Status Page Issues

#### Status page not loading

**Solutions**:

1. **Check web container**:
   ```bash
   # Verify web container is running
   docker compose ps web

   # Check web container logs
   docker compose logs web
   ```

2. **Test internal connectivity**:
   ```bash
   # Test from web container to server
   docker compose exec web curl -f http://server:8034/health
   ```

3. **Verify reverse proxy configuration**:
   ```nginx
   # Example working Nginx config
   location / {
       proxy_pass http://localhost:8383;
       proxy_set_header Host $host;
       proxy_set_header X-Real-IP $remote_addr;
   }
   ```

#### Status page shows incorrect information

**Solutions**:

1. **Clear browser cache**:
   - Hard refresh (Ctrl+F5 or Cmd+Shift+R)
   - Clear browser cache and cookies

2. **Check monitor assignments**:
   - Verify correct monitors are assigned to status page
   - Check monitor configurations

3. **Review status page settings**:
   - Verify slug and URL configuration
   - Check custom CSS for conflicts

## Performance Issues

### Slow Response Times

**Optimization Steps**:

1. **Database optimization**:
   ```bash
   # Check MongoDB performance
   docker compose exec mongodb mongosh --eval "db.stats()"

   # Check database size
   docker compose exec mongodb mongosh --eval "db.runCommand({dbStats: 1})"
   ```

2. **Clean up old data**:
   ```bash
   # Manual cleanup (be careful!)
   docker compose exec mongodb mongosh peekaping --eval "
     db.heartbeats.deleteMany({
       createdAt: { \$lt: new Date(Date.now() - 90*24*60*60*1000) }
     })
   "
   ```

3. **Resource monitoring**:
   ```bash
   # Monitor container resources
   docker stats

   # Check system resources
   htop
   ```

### High Memory Usage

**Solutions**:

1. **Restart containers**:
   ```bash
   docker compose restart
   ```

2. **Limit container memory**:
   ```yaml
   # In docker-compose.yml
   services:
     server:
       deploy:
         resources:
           limits:
             memory: 512M
   ```

3. **Configure MongoDB memory**:
   ```yaml
   mongodb:
     command: mongod --wiredTigerCacheSizeGB 0.5
   ```

## Data Issues

### Data Loss Prevention

**Backup Strategies**:

1. **Regular database backups**:
   ```bash
   # Create backup
   docker exec peekaping-mongodb mongodump \
     --uri="mongodb://root:password@localhost:27017/peekaping?authSource=admin" \
     --out=/backup/$(date +%Y%m%d_%H%M%S)

   # Copy backup to host
   docker cp peekaping-mongodb:/backup ./backups/
   ```

2. **Automated backup script**:
   ```bash
   #!/bin/bash
   BACKUP_DIR="/backups/peekaping"
   DATE=$(date +%Y%m%d_%H%M%S)

   # Create backup
   docker exec peekaping-mongodb mongodump \
     --uri="mongodb://root:${DB_PASSWORD}@localhost:27017/peekaping?authSource=admin" \
     --out=/backup/$DATE

   # Copy to host
   docker cp peekaping-mongodb:/backup/$DATE $BACKUP_DIR/

   # Cleanup old backups (keep last 7 days)
   find $BACKUP_DIR -type d -mtime +7 -exec rm -rf {} \;
   ```

### Data Recovery

**If you lose data**:

1. **Check if data still exists**:
   ```bash
   # List collections
   docker compose exec mongodb mongosh peekaping --eval "show collections"

   # Count documents
   docker compose exec mongodb mongosh peekaping --eval "
     db.monitors.countDocuments();
     db.heartbeats.countDocuments();
   "
   ```

2. **Restore from backup**:
   ```bash
   # Copy backup to container
   docker cp ./backups/20240120_143000 peekaping-mongodb:/restore/

   # Restore data
   docker exec peekaping-mongodb mongorestore \
     --uri="mongodb://root:password@localhost:27017/peekaping?authSource=admin" \
     /restore/20240120_143000/peekaping
   ```

## Network and Security Issues

### Firewall Configuration

**Common ports to open**:

```bash
# For external access
sudo ufw allow 8383/tcp  # Web interface
sudo ufw allow 443/tcp   # HTTPS (if using reverse proxy)
sudo ufw allow 80/tcp    # HTTP (if using reverse proxy)

# For API access (if needed)
sudo ufw allow 8034/tcp  # API server
```

### SSL/TLS Issues

**Certificate problems**:

1. **Check certificate validity**:
   ```bash
   # Check certificate expiration
   echo | openssl s_client -connect your-domain.com:443 2>/dev/null | \
     openssl x509 -noout -dates
   ```

2. **Verify certificate chain**:
   ```bash
   # Check full certificate chain
   openssl s_client -connect your-domain.com:443 -showcerts
   ```

3. **Test TLS configuration**:
   ```bash
   # Test TLS versions
   nmap --script ssl-enum-ciphers -p 443 your-domain.com
   ```

## Advanced Troubleshooting

### Debugging with Logs

**Enable debug logging**:

1. **Change log level in .env**:
   ```env
   MODE=debug
   LOG_LEVEL=debug
   ```

2. **View detailed logs**:
   ```bash
   # Follow logs in real-time
   docker compose logs -f server

   # Search for specific errors
   docker compose logs server | grep -i error

   # View logs for specific time period
   docker compose logs --since="2024-01-20T10:00:00" server
   ```

### Performance Profiling

**Monitor application performance**:

1. **Database performance**:
   ```bash
   # Enable MongoDB profiling
   docker compose exec mongodb mongosh peekaping --eval "
     db.setProfilingLevel(2);
     db.system.profile.find().limit(5).sort({ts:-1}).pretty();
   "
   ```

2. **API performance**:
   ```bash
   # Monitor API response times
   curl -w "@curl-format.txt" -s -o /dev/null https://your-domain.com/api/health
   ```

   Create `curl-format.txt`:
   ```
   time_namelookup:  %{time_namelookup}\n
   time_connect:     %{time_connect}\n
   time_appconnect:  %{time_appconnect}\n
   time_pretransfer: %{time_pretransfer}\n
   time_redirect:    %{time_redirect}\n
   time_starttransfer: %{time_starttransfer}\n
   time_total:       %{time_total}\n
   ```

### Container Debugging

**Access container shells**:

```bash
# Access server container
docker compose exec server /bin/sh

# Access MongoDB container
docker compose exec mongodb /bin/bash

# Access web container
docker compose exec web /bin/sh
```

**Inspect container configurations**:

```bash
# View container details
docker inspect peekaping-server

# Check container networks
docker network ls
docker network inspect peekaping-network
```

## Getting More Help

### Community Resources

- **GitHub Issues**: [Report bugs and feature requests](https://github.com/0xfurai/peekaping/issues)
- **Discussions**: Join community discussions for tips and help
- **Documentation**: Review the complete documentation
- **Examples**: Check example configurations in the repository

### Professional Support

For production environments or complex setups:
- Consider professional support options
- Review enterprise features
- Consult with infrastructure specialists

### Reporting Bugs

When reporting issues, include:

1. **System information**:
   ```bash
   # Collect system info
   uname -a
   docker --version
   docker compose version
   ```

2. **Configuration details**:
   - Sanitized `.env` file (remove secrets)
   - `docker-compose.yml` file
   - Browser and version (for web issues)

3. **Error logs**:
   ```bash
   # Collect relevant logs
   docker compose logs --tail=100 > peekaping-logs.txt
   ```

4. **Steps to reproduce**:
   - Clear description of the issue
   - Steps to reproduce the problem
   - Expected vs actual behavior

Remember: The community is here to help! Don't hesitate to ask questions and share your experiences.
