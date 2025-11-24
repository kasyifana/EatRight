# EatRight Backend - VPS Deployment Guide

This guide will walk you through deploying the EatRight backend on an Ubuntu VPS.

## Prerequisites

- Ubuntu 20.04+ VPS with root access
- Domain name pointed to your VPS IP
- Supabase project with PostgreSQL database
- Go 1.22+ (for building, if not building locally)

## Step 1: Prepare Your VPS

### 1.1 Update system packages
```bash
sudo apt update
sudo apt upgrade -y
```

### 1.2 Install required packages
```bash
sudo apt install -y nginx git ufw
```

### 1.3 Setup firewall
```bash
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw enable
```

## Step 2: Setup Application Directory

### 2.1 Create application directory
```bash
sudo mkdir -p /var/www/eatright
sudo chown -R $USER:$USER /var/www/eatright
```

### 2.2 Upload your files
You can use `scp`, `rsync`, or `git clone`:

**Using SCP (from your local machine):**
```bash
# Build locally first
./scripts/build.sh

# Upload to VPS
scp -r bin .env user@your-vps-ip:/var/www/eatright/
```

**Using Git:**
```bash
cd /var/www/eatright
git clone <your-repository-url> .
```

## Step 3: Configure Environment Variables

### 3.1 Create .env file
```bash
cd /var/www/eatright
nano .env
```

### 3.2 Add your configuration
```env
PORT=8080
ENV=production

DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres

SUPABASE_URL=https://[PROJECT-REF].supabase.co
SUPABASE_KEY=[YOUR-ANON-KEY]
SUPABASE_SERVICE_KEY=[YOUR-SERVICE-KEY]

JWT_SECRET=[GENERATE-RANDOM-SECRET]
JWT_EXPIRY=24h

ALLOWED_ORIGINS=https://your-frontend-domain.com,https://www.your-frontend-domain.com
```

**Generate a random JWT secret:**
```bash
openssl rand -base64 32
```

### 3.3 Secure the .env file
```bash
chmod 600 .env
```

## Step 4: Run Database Migrations

1. Go to your Supabase dashboard
2. Navigate to SQL Editor
3. Copy the contents of `migrations/001_create_tables.sql`
4. Execute the script

## Step 5: Setup Systemd Service

### 5.1 Copy service file
```bash
sudo cp /var/www/eatright/deploy/eatright.service /etc/systemd/system/
```

### 5.2 Update service file paths (if needed)
```bash
sudo nano /etc/systemd/system/eatright.service
```

Ensure these paths are correct:
- `WorkingDirectory=/var/www/eatright`
- `EnvironmentFile=/var/www/eatright/.env`
- `ExecStart=/var/www/eatright/bin/eatright-server`

### 5.3 Set proper ownership
```bash
sudo chown -R www-data:www-data /var/www/eatright
```

### 5.4 Make binary executable
```bash
chmod +x /var/www/eatright/bin/eatright-server
```

### 5.5 Enable and start the service
```bash
sudo systemctl daemon-reload
sudo systemctl enable eatright
sudo systemctl start eatright
```

### 5.6 Check service status
```bash
sudo systemctl status eatright
```

### 5.7 View logs
```bash
sudo journalctl -u eatright -f
```

## Step 6: Configure Nginx

### 6.1 Copy Nginx configuration
```bash
sudo cp /var/www/eatright/deploy/nginx.conf /etc/nginx/sites-available/eatright
```

### 6.2 Update domain name
```bash
sudo nano /etc/nginx/sites-available/eatright
```

Replace `your-domain.com` with your actual domain.

### 6.3 Enable the site
```bash
sudo ln -s /etc/nginx/sites-available/eatright /etc/nginx/sites-enabled/
```

### 6.4 Test Nginx configuration
```bash
sudo nginx -t
```

### 6.5 Restart Nginx
```bash
sudo systemctl restart nginx
```

## Step 7: Setup SSL with Let's Encrypt (Recommended)

### 7.1 Install Certbot
```bash
sudo apt install -y certbot python3-certbot-nginx
```

### 7.2 Obtain SSL certificate
```bash
sudo certbot --nginx -d your-domain.com -d www.your-domain.com
```

### 7.3 Auto-renewal test
```bash
sudo certbot renew --dry-run
```

The certificate will auto-renew. Certbot adds a systemd timer for this.

### 7.4 Uncomment HTTPS block in Nginx config
After SSL is set up, edit `/etc/nginx/sites-available/eatright` and uncomment the HTTPS server block.

```bash
sudo nano /etc/nginx/sites-available/eatright
sudo nginx -t
sudo systemctl reload nginx
```

## Step 8: Verify Deployment

### 8.1 Test health endpoint
```bash
curl https://your-domain.com/health
```

Expected response:
```json
{
  "status": "healthy",
  "env": "production"
}
```

### 8.2 Test API endpoints
```bash
# Should return 401 (unauthorized, but endpoint works)
curl https://your-domain.com/api/users/me
```

## Maintenance Commands

### View logs
```bash
# Follow logs in real-time
sudo journalctl -u eatright -f

# View last 100 lines
sudo journalctl -u eatright -n 100
```

### Restart service
```bash
sudo systemctl restart eatright
```

### Stop service
```bash
sudo systemctl stop eatright
```

### Update application
```bash
# 1. Stop service
sudo systemctl stop eatright

# 2. Upload new binary or pull code
cd /var/www/eatright
git pull  # if using git

# 3. Rebuild if needed
./scripts/build.sh

# 4. Start service
sudo systemctl start eatright

# 5. Check status
sudo systemctl status eatright
```

### Check service status
```bash
sudo systemctl status eatright
```

## Troubleshooting

### Service won't start
```bash
# Check logs
sudo journalctl -u eatright -n 50

# Check file permissions
ls -la /var/www/eatright/bin/eatright-server
ls -la /var/www/eatright/.env

# Verify environment file
sudo -u www-data cat /var/www/eatright/.env
```

### Database connection errors
- Verify `DATABASE_URL` in `.env` is correct
- Check Supabase dashboard for database status
- Ensure VPS IP is allowed in Supabase network settings

### Nginx errors
```bash
# Test configuration
sudo nginx -t

# View error logs
sudo tail -f /var/log/nginx/eatright_error.log
```

### Port already in use
```bash
# Check what's using port 8080
sudo lsof -i :8080

# Change PORT in .env and restart
sudo systemctl restart eatright
```

## Security Recommendations

1. **Firewall**: Keep UFW enabled with only necessary ports open
2. **SSH**: Disable password authentication, use SSH keys only
3. **Updates**: Regularly update system packages
4. **Backups**: Setup automated database backups in Supabase
5. **Monitoring**: Consider setting up monitoring (e.g., Prometheus, Grafana)
6. **Rate Limiting**: Add rate limiting in Nginx for API protection

## Performance Optimization

### Enable Nginx caching (optional)
Add to Nginx config:
```nginx
proxy_cache_path /var/cache/nginx levels=1:2 keys_zone=api_cache:10m max_size=100m inactive=60m;

location /api/ {
    proxy_cache api_cache;
    proxy_cache_valid 200 5m;
    # ... rest of proxy settings
}
```

### Database connection pooling
The application already configures connection pooling. Adjust in `internal/app/config/database.go` if needed:
```go
sqlDB.SetMaxIdleConns(10)
sqlDB.SetMaxOpenConns(100)
```

## Support

For issues or questions:
- Check application logs: `sudo journalctl -u eatright -f`
- Check Nginx logs: `sudo tail -f /var/log/nginx/eatright_error.log`
- Verify Supabase connection in dashboard

---

**Deployment Complete! 🚀**

Your EatRight backend should now be running at `https://your-domain.com`
