# Deployment Guide

## Production Checklist

- [ ] Review and set all environment variables (especially `JWT_SECRET`)
- [ ] Configure database backups
- [ ] Set up monitoring/logging
- [ ] Configure SSL/TLS certificates
- [ ] Set up firewall rules
- [ ] Enable database authentication
- [ ] Review security settings
- [ ] Test failover procedures

---

## Environment Configuration

### Critical Settings

```env
# Security
JWT_SECRET=<generate-strong-random-key>
GIN_MODE=release
ENV=production

# Database
DB_HOST=<production-db-host>
DB_PORT=5432
DB_NAME=medpredict
DB_USER=<db-user>
DB_PASSWORD=<strong-password>
DB_SSL_MODE=require  # Always use SSL in production

# Server
PORT=8080
FRONTEND_URL=https://yourdomain.com

# Logging
LOG_LEVEL=info
LOG_DIR=/var/log/medpredict

# Optional: AI Services
ANTHROPIC_API_KEY=<your-api-key>
OPENAI_API_KEY=<your-api-key>

# Optional: Notifications
TWILIO_SID=<your-sid>
TWILIO_AUTH_TOKEN=<your-token>
TWILIO_WHATSAPP_NUMBER=<your-number>
```

### Security Notes

1. **JWT_SECRET**: Generate with `openssl rand -base64 32`
2. **DB_PASSWORD**: Use a strong random password
3. **Never commit .env** to version control
4. Use environment variable management service (AWS Secrets Manager, etc.)

---

## Deployment Options

### Option 1: Linux/Unix Server

#### 1. Build Binary

```bash
GOOS=linux GOARCH=amd64 go build -o med-predict-backend cmd/server/main.go
```

#### 2. Setup systemd Service

Create `/etc/systemd/system/medpredict.service`:

```ini
[Unit]
Description=Med Predict Backend
After=network.target postgresql.service

[Service]
Type=simple
User=medpredict
WorkingDirectory=/opt/medpredict
EnvironmentFile=/opt/medpredict/.env
ExecStart=/opt/medpredict/med-predict-backend
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

#### 3. Start Service

```bash
sudo systemctl daemon-reload
sudo systemctl start medpredict
sudo systemctl enable medpredict
sudo systemctl status medpredict
```

#### 4. View Logs

```bash
sudo journalctl -u medpredict -f
```

---

### Option 2: Docker on Linux

#### 1. Build Image

```bash
docker build -t med-predict-backend:1.0 .
docker tag med-predict-backend:1.0 med-predict-backend:latest
```

#### 2. Push to Registry

```bash
docker tag med-predict-backend:1.0 registry.example.com/med-predict-backend:1.0
docker push registry.example.com/med-predict-backend:1.0
```

#### 3. Run Container

```bash
docker run -d \
  --name medpredict \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file /etc/medpredict/.env \
  -v /var/log/medpredict:/root/logs \
  registry.example.com/med-predict-backend:1.0
```

---

### Option 3: Kubernetes

#### 1. Create ConfigMap

```bash
kubectl create configmap medpredict-config \
  --from-file=.env.production \
  -n medpredict
```

#### 2. Deployment YAML

Create `k8s/deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: medpredict-backend
  namespace: medpredict
spec:
  replicas: 3
  selector:
    matchLabels:
      app: medpredict-backend
  template:
    metadata:
      labels:
        app: medpredict-backend
    spec:
      containers:
      - name: backend
        image: registry.example.com/med-predict-backend:1.0
        ports:
        - containerPort: 8080
        envFrom:
        - configMapRef:
            name: medpredict-config
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        resources:
          requests:
            cpu: 200m
            memory: 256Mi
          limits:
            cpu: 500m
            memory: 512Mi
```

#### 3. Deploy

```bash
kubectl apply -f k8s/deployment.yaml
kubectl get pods -n medpredict
```

---

### Option 4: Cloud Platforms

#### AWS (Elastic Beanstalk)

```bash
eb init -p "Go 1.21 running on 64bit Amazon Linux 2"
eb create medpredict-prod
eb deploy
```

#### Google Cloud (Cloud Run)

```bash
gcloud run deploy medpredict-backend \
  --source . \
  --platform managed \
  --region us-central1 \
  --set-env-vars ENV=production,DB_HOST=cloudsql-instance
```

#### DigitalOcean (App Platform)

1. Push to GitHub
2. Connect repository to App Platform
3. Configure environment variables
4. Deploy

---

## Database Setup

### PostgreSQL Configuration

```bash
# Create production database
createdb -U postgres medpredict_prod

# Create dedicated user (recommended)
psql -U postgres -c "CREATE USER medpredict_user WITH PASSWORD 'strong-password';"
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE medpredict_prod TO medpredict_user;"

# Run migrations
psql -U medpredict_user -d medpredict_prod -f migrations/001_init_schema.sql

# Setup backups
pg_dump -U medpredict_user -d medpredict_prod > backup.sql
```

### RDS (AWS)

1. Create RDS instance (PostgreSQL 12+)
2. Configure security groups
3. Set master password
4. Run migrations against endpoint

---

## Reverse Proxy Setup (Nginx)

```nginx
upstream medpredict_backend {
  server 127.0.0.1:8080;
}

server {
  listen 443 ssl http2;
  server_name api.yourdomain.com;

  ssl_certificate /path/to/cert.pem;
  ssl_certificate_key /path/to/key.pem;

  # Security headers
  add_header Strict-Transport-Security "max-age=31536000" always;
  add_header X-Content-Type-Options "nosniff" always;
  add_header X-Frame-Options "DENY" always;

  location / {
    proxy_pass http://medpredict_backend;
    proxy_set_header Host $host;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header X-Forwarded-Proto $scheme;
    proxy_read_timeout 90;
  }
}

server {
  listen 80;
  server_name api.yourdomain.com;
  return 301 https://$server_name$request_uri;
}
```

---

## Monitoring & Logging

### Health Checks

```bash
curl https://api.yourdomain.com/health
```

### Prometheus Metrics Endpoint (optional, not implemented)

Extend handlers with metrics if needed.

### Log Aggregation

Example with ELK Stack:

```bash
# Forward logs to logstash
tail -f /var/log/medpredict/error.log | nc logstash.example.com 5000
```

---

## Backup & Disaster Recovery

### Automated PostgreSQL Backups

```bash
# Daily backup script
#!/bin/bash
pg_dump -U medpredict_user -d medpredict_prod | \
  gzip > /backups/medpredict_$(date +%Y%m%d).sql.gz

# Upload to S3
aws s3 cp /backups/medpredict_*.sql.gz s3://backup-bucket/
```

### Restore from Backup

```bash
gunzip < backup.sql.gz | psql -U medpredict_user -d medpredict_prod
```

---

## Performance Tuning

### PostgreSQL Optimization

Add to `postgresql.conf`:
```
shared_buffers = 256MB
effective_cache_size = 1GB
maintenance_work_mem = 64MB
checkpoint_completion_target = 0.9
wal_buffers = 16MB
default_statistics_target = 100
```

### Connection Pooling

Use PgBouncer for connection pooling:
```
[databases]
medpredict_prod = host=db.example.com dbname=medpredict_prod user=medpredict_user

[pgbouncer]
pool_mode = transaction
max_client_conn = 1000
default_pool_size = 25
```

---

## Troubleshooting

### Service Won't Start

```bash
# Check systemd logs
sudo journalctl -u medpredict -n 50

# Test binary directly
./med-predict-backend

# Check file permissions
ls -la /opt/medpredict/
```

### Database Connection Issues

```bash
# Test connection
psql -h db.example.com -U medpredict_user -d medpredict_prod

# Check environment variables
echo $DB_HOST
echo $DB_USER
```

### High Memory Usage

- Check connection pool settings
- Monitor database queries
- Increase log level to debug issues

---

## Rollback Procedure

1. Update Docker image tag
2. Restart service: `systemctl restart medpredict`
3. Or redeploy: `docker pull && docker run`
4. Verify: `curl /health`

---

## Security Hardening

- [ ] Enable firewall (ufw/iptables)
- [ ] Disable SSH password auth
- [ ] Use key pairs
- [ ] Setup fail2ban
- [ ] Implement rate limiting (done in code)
- [ ] Enable audit logging
- [ ] Use VPN/private networks for databases
- [ ] Regular security updates
- [ ] Setup intrusion detection

---

## Support

For issues:
1. Check logs: `journalctl -u medpredict -f`
2. Test API: `curl http://localhost:8080/health`
3. Verify database: `psql ... -c "SELECT version();"`
