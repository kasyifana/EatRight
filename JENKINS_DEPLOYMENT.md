# EatRight Backend - Jenkins CI/CD Deployment Guide

Complete guide untuk deploy EatRight Go backend ke VPS menggunakan Jenkins dari awal.

## 📋 Prerequisites

- VPS dengan Ubuntu 20.04+ (minimal 2GB RAM)
- Domain (opsional tapi recommended)
- Akses SSH ke VPS
- Repository GitHub untuk source code

---

## 🚀 Part 1: Setup VPS Awal

### 1.1 Connect ke VPS

```bash
ssh root@your-vps-ip
```

### 1.2 Update System

```bash
sudo apt update && sudo apt upgrade -y
```

### 1.3 Install Dependencies

```bash
# Install essential tools
sudo apt install -y git curl wget ufw nginx

# Install Go 1.22+
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify Go installation
go version
```

### 1.4 Setup Firewall

```bash
sudo ufw allow OpenSSH
sudo ufw allow 'Nginx Full'
sudo ufw allow 8080/tcp  # Jenkins port
sudo ufw enable
sudo ufw status
```

### 1.5 Create Application User

```bash
# Create user for running the app
sudo useradd -m -s /bin/bash eatright
sudo usermod -aG sudo eatright

# Create app directory
sudo mkdir -p /var/www/eatright
sudo chown -R eatright:eatright /var/www/eatright
```

---

## 🔧 Part 2: Install & Configure Jenkins

### 2.1 Install Java (Jenkins Requirement)

```bash
sudo apt install -y openjdk-17-jdk
java -version
```

### 2.2 Install Jenkins

```bash
# Add Jenkins repository
wget -q -O - https://pkg.jenkins.io/debian-stable/jenkins.io.key | sudo apt-key add -
sudo sh -c 'echo deb https://pkg.jenkins.io/debian-stable binary/ > /etc/apt/sources.list.d/jenkins.list'

# Install Jenkins
sudo apt update
sudo apt install -y jenkins

# Start Jenkins
sudo systemctl start jenkins
sudo systemctl enable jenkins
sudo systemctl status jenkins
```

### 2.3 Initial Jenkins Setup

1. **Get initial admin password:**
   ```bash
   cat /var/lib/jenkins/secrets/initialAdminPassword
   ```

2. **Access Jenkins:**
   - Open browser: `http://your-vps-ip:8080`
   - Paste the initial admin password

3. **Install suggested plugins** (pilih "Install suggested plugins")

4. **Create admin user:**
   - Username: `admin`
   - Password: (your choice)
   - Full name: `Admin`
   - Email: your-email@example.com

5. **Jenkins URL:** `http://your-vps-ip:8080/` (atau gunakan domain)

### 2.4 Install Required Jenkins Plugins

Go to: **Manage Jenkins** → **Manage Plugins** → **Available**

Install:
- ✅ Git Plugin (usually already installed)
- ✅ Pipeline Plugin
- ✅ SSH Agent Plugin
- ✅ Credentials Binding Plugin
- ✅ GitHub Integration Plugin (optional)

Restart Jenkins setelah install:
```bash
sudo systemctl restart jenkins
```

---

## 📝 Part 3: Create Jenkins Pipeline

### 3.1 Add GitHub Credentials

1. Go to: **Manage Jenkins** → **Manage Credentials**
2. Click **(global)** → **Add Credentials**
3. Kind: **Username with password**
   - Username: GitHub username
   - Password: GitHub Personal Access Token
   - ID: `github-credentials`
4. Click **Create**

### 3.2 Create New Pipeline Job

1. **Dashboard** → **New Item**
2. Name: `eatright-backend`
3. Type: **Pipeline**
4. Click **OK**

### 3.3 Configure Pipeline

**General:**
- ✅ GitHub project: `https://github.com/your-username/eatright-backend`

**Build Triggers:**
- ✅ Poll SCM: `H/5 * * * *` (check every 5 minutes)
- *Atau setup GitHub webhook untuk trigger otomatis*

**Pipeline:**
- Definition: **Pipeline script from SCM**
- SCM: **Git**
- Repository URL: `https://github.com/your-username/eatright-backend.git`
- Credentials: `github-credentials`
- Branch: `*/main` (atau `*/master`)
- Script Path: `Jenkinsfile`

Click **Save**

---

## 📄 Part 4: Create Deployment Files

### 4.1 Create Jenkinsfile

File ini sudah dibuat di: `Jenkinsfile`

```groovy
// Lihat file Jenkinsfile di repository
```

### 4.2 Create Deployment Script

File ini sudah dibuat di: `scripts/deploy.sh`

```bash
// Lihat file scripts/deploy.sh di repository
```

### 4.3 Setup Environment Variables di VPS

```bash
# Login sebagai user eatright
su - eatright

# Create .env file
nano /var/www/eatright/.env
```

Isi dengan:
```env
PORT=8080
ENV=production

DATABASE_URL=postgresql://postgres:[PASSWORD]@db.[PROJECT-REF].supabase.co:5432/postgres

SUPABASE_URL=https://[PROJECT-REF].supabase.co
SUPABASE_KEY=[YOUR-ANON-KEY]
SUPABASE_SERVICE_KEY=[YOUR-SERVICE-KEY]

JWT_SECRET=[GENERATE-RANDOM]
JWT_EXPIRY=24h

ALLOWED_ORIGINS=https://your-frontend-domain.com
```

Generate JWT secret:
```bash
openssl rand -base64 32
```

Secure the file:
```bash
chmod 600 /var/www/eatright/.env
```

---

## 🔒 Part 5: Setup Systemd & Nginx

### 5.1 Create Systemd Service

```bash
sudo cp /var/www/eatright/deploy/eatright.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable eatright
```

### 5.2 Setup Nginx

```bash
sudo cp /var/www/eatright/deploy/nginx.conf /etc/nginx/sites-available/eatright
sudo ln -s /etc/nginx/sites-available/eatright /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### 5.3 Setup SSL (Let's Encrypt)

```bash
# Install certbot
sudo apt install -y certbot python3-certbot-nginx

# Get certificate
sudo certbot --nginx -d your-domain.com -d www.your-domain.com

# Test auto-renewal
sudo certbot renew --dry-run
```

---

## 🎯 Part 6: First Deployment

### 6.1 Push Code to GitHub

```bash
# Di local machine
git add .
git commit -m "Add Jenkins CI/CD configuration"
git push origin main
```

### 6.2 Trigger Jenkins Build

**Option 1: Manual**
1. Go to Jenkins dashboard
2. Click `eatright-backend`
3. Click **Build Now**

**Option 2: Automatic**
- Jenkins will poll GitHub every 5 minutes
- Atau setup webhook untuk instant trigger

### 6.3 Monitor Build

1. Click build number (e.g., #1)
2. Click **Console Output**
3. Watch the build process

Build stages:
1. ✅ Checkout code from Git
2. ✅ Install Go dependencies
3. ✅ Run tests
4. ✅ Build binary
5. ✅ Deploy to /var/www/eatright
6. ✅ Restart systemd service

---

## ✅ Part 7: Verification

### 7.1 Check Service Status

```bash
sudo systemctl status eatright
```

### 7.2 Check Logs

```bash
# Application logs
sudo journalctl -u eatright -f

# Nginx logs
sudo tail -f /var/log/nginx/eatright_access.log
sudo tail -f /var/log/nginx/eatright_error.log
```

### 7.3 Test API

```bash
# Health check
curl https://your-domain.com/health

# List restaurants
curl https://your-domain.com/api/restaurants

# Swagger UI
curl https://your-domain.com/swagger/index.html
```

---

## 🔄 Part 8: CI/CD Workflow

### Development Workflow

1. **Developer pushes code** ke GitHub
   ```bash
   git add .
   git commit -m "Feature: Add new endpoint"
   git push origin main
   ```

2. **Jenkins detects change** (poll atau webhook)

3. **Automated pipeline runs:**
   - Clone repository
   - Install dependencies
   - Run tests
   - Build binary
   - Deploy to VPS
   - Restart service

4. **Application updated** automatically!

### Rollback Strategy

Jika deployment gagal:

```bash
# Stop current service
sudo systemctl stop eatright

# Restore previous binary
cd /var/www/eatright
cp bin/eatright-server.backup bin/eatright-server

# Start service
sudo systemctl start eatright
```

---

## 🐛 Troubleshooting

### Jenkins Issues

**Jenkins not accessible:**
```bash
sudo systemctl status jenkins
sudo journalctl -u jenkins -f
```

**Jenkins port conflict:**
Edit `/etc/default/jenkins` and change `HTTP_PORT`

### Build Failures

**Go not found:**
- Add Go to Jenkins PATH in pipeline
- Or install Go plugin for Jenkins

**Permission denied:**
```bash
sudo chown -R eatright:eatright /var/www/eatright
sudo chmod +x /var/www/eatright/scripts/deploy.sh
```

### Application Issues

**Service won't start:**
```bash
sudo journalctl -u eatright -n 50
```

**Database connection failed:**
- Check DATABASE_URL in `.env`
- Verify Supabase allows VPS IP

**Port already in use:**
```bash
sudo lsof -i :8080
sudo kill -9 <PID>
```

---

## 📊 Monitoring & Maintenance

### Setup Log Rotation

```bash
sudo nano /etc/logrotate.d/eatright
```

```
/var/log/nginx/eatright*.log {
    daily
    rotate 14
    compress
    delaycompress
    notifempty
    sharedscripts
    postrotate
        systemctl reload nginx
    endscript
}
```

### Backup Strategy

```bash
# Backup script
cat > /root/backup-eatright.sh << 'EOF'
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
tar -czf /backup/eatright-$DATE.tar.gz /var/www/eatright
find /backup -name "eatright-*.tar.gz" -mtime +7 -delete
EOF

chmod +x /root/backup-eatright.sh

# Add to crontab
crontab -e
# Add: 0 2 * * * /root/backup-eatright.sh
```

---

## 🎉 Summary

Setelah setup lengkap, workflow Anda:

1. ✅ **VPS ready** dengan Jenkins, Nginx, SSL
2. ✅ **CI/CD automated** - push code → auto deploy
3. ✅ **Production ready** dengan monitoring & backup
4. ✅ **Secure** dengan HTTPS, firewall, systemd

**Next deployments:**
```bash
# Just push to GitHub!
git push origin main

# Jenkins handles the rest automatically! 🚀
```

---

## 📚 Quick Reference

**Important Commands:**
```bash
# Restart application
sudo systemctl restart eatright

# View logs
sudo journalctl -u eatright -f

# Rebuild manually
cd /var/www/eatright
go build -o bin/eatright-server cmd/server/main.go
sudo systemctl restart eatright

# Check Jenkins
sudo systemctl status jenkins
```

**Important Paths:**
- Application: `/var/www/eatright`
- Logs: `/var/log/nginx/eatright_*.log`
- Service: `/etc/systemd/system/eatright.service`
- Nginx config: `/etc/nginx/sites-available/eatright`

**Important URLs:**
- Jenkins: `http://your-vps-ip:8080`
- Application: `https://your-domain.com`
- Swagger: `https://your-domain.com/swagger/index.html`

---

**Happy Deploying! 🎊**
