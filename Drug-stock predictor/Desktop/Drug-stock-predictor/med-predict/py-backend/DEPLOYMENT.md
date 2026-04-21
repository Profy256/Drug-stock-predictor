# Deployment Guide

## Local Development

### Prerequisites
- Python 3.11+
- PostgreSQL 12+
- Virtual environment (venv or conda)

### Quick Start

```bash
# 1. Navigate to backend
cd py-backend

# 2. Create and activate virtual environment
python -m venv venv
# Windows:
venv\Scripts\activate
# Linux/Mac:
source venv/bin/activate

# 3. Install dependencies
pip install -r requirements.txt

# 4. Setup environment
cp .env.example .env
# Edit .env with your database credentials

# 5. Create database
createdb med_predict

# 6. Run server
make dev
# or: uvicorn app.main:app --reload
```

Visit http://localhost:8000/docs for API documentation.

## Docker Deployment

### Using Docker Compose (Recommended)

```bash
cd py-backend

# Start services
docker-compose up -d

# View logs
docker-compose logs -f backend

# Stop services
docker-compose down
```

The backend will be available at `http://localhost:8000`.

### Manual Docker Build

```bash
cd py-backend

# Build image
docker build -t med-predict-backend:latest .

# Run container
docker run -p 8000:8000 \
  -e DATABASE_URL=postgresql://user:pass@host:5432/db \
  -e JWT_SECRET=your-secret-key \
  med-predict-backend:latest
```

## Production Deployment

### Environment Variables

Create a `.env` file with production values:

```env
DEBUG=False
ENV=production
DATABASE_URL=postgresql://user:password@prod-db-host:5432/med_predict
JWT_SECRET=your-strong-secret-key-here-change-in-production
HOST=0.0.0.0
PORT=8000
CORS_ORIGINS=["https://yourdomain.com"]
```

### Using Gunicorn

```bash
# Install gunicorn
pip install gunicorn

# Run with 4 workers
gunicorn app.main:app \
  -w 4 \
  -k uvicorn.workers.UvicornWorker \
  --bind 0.0.0.0:8000 \
  --access-logfile - \
  --error-logfile - \
  --log-level info
```

### Using systemd (Linux)

Create `/etc/systemd/system/med-predict.service`:

```ini
[Unit]
Description=Med Predict Backend
After=network.target

[Service]
Type=notify
User=www-data
WorkingDirectory=/var/www/med-predict/py-backend
Environment="PATH=/var/www/med-predict/venv/bin"
ExecStart=/var/www/med-predict/venv/bin/gunicorn app.main:app \
  -w 4 \
  -k uvicorn.workers.UvicornWorker \
  --bind 127.0.0.1:8000 \
  --access-logfile /var/log/med-predict/access.log \
  --error-logfile /var/log/med-predict/error.log
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl daemon-reload
sudo systemctl enable med-predict
sudo systemctl start med-predict
```

### Using nginx (Reverse Proxy)

Create `/etc/nginx/sites-available/med-predict`:

```nginx
upstream med_predict_backend {
    server 127.0.0.1:8000;
}

server {
    listen 80;
    server_name api.yourdomain.com;

    location / {
        proxy_pass http://med_predict_backend;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable:
```bash
sudo ln -s /etc/nginx/sites-available/med-predict /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl restart nginx
```

### Using Docker Swarm

```bash
# Build image
docker build -t med-predict-backend:latest .

# Create stack
docker service create \
  --name med-predict \
  --replicas 3 \
  -p 8000:8000 \
  -e DATABASE_URL=postgresql://user:pass@db:5432/med_predict \
  -e JWT_SECRET=your-secret \
  med-predict-backend:latest
```

### Using Kubernetes

Create `deployment.yaml`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: med-predict-backend
spec:
  replicas: 3
  selector:
    matchLabels:
      app: med-predict-backend
  template:
    metadata:
      labels:
        app: med-predict-backend
    spec:
      containers:
      - name: backend
        image: med-predict-backend:latest
        ports:
        - containerPort: 8000
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: med-predict-secrets
              key: database-url
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: med-predict-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8000
          initialDelaySeconds: 10
          periodSeconds: 5
```

Deploy:
```bash
kubectl apply -f deployment.yaml
```

## Database Setup

### Automatic (Docker Compose)
The schema is automatically created when using docker-compose.

### Manual PostgreSQL

```bash
# Create database
createdb -U postgres med_predict

# Run migrations
psql -U postgres -d med_predict -f migrations/001_init_schema.sql

# Verify
psql -U postgres -d med_predict -c "\dt"
```

### Backup Database

```bash
# Full backup
pg_dump -U postgres med_predict > backup.sql

# Restore
psql -U postgres -d med_predict < backup.sql
```

## Monitoring

### Application Health

```bash
# Check health
curl http://localhost:8000/health

# Response: {"status": "ok", "service": "Med Predict Backend"}
```

### Logs

```bash
# Docker Compose
docker-compose logs -f backend

# Kubernetes
kubectl logs -f deployment/med-predict-backend

# Systemd
sudo journalctl -u med-predict -f
```

### Performance Metrics

```bash
# Response time
curl -w "Time: %{time_total}s\n" http://localhost:8000/health

# Database connections
psql -c "SELECT datname, count(*) FROM pg_stat_activity GROUP BY datname;"
```

## Security Checklist

- [ ] Change JWT_SECRET to a strong, random value
- [ ] Use HTTPS/TLS in production
- [ ] Set DEBUG=False in production
- [ ] Use strong database passwords
- [ ] Restrict database access by IP
- [ ] Enable SSL for database connections
- [ ] Use environment variables for secrets
- [ ] Implement rate limiting
- [ ] Add CORS restrictions
- [ ] Enable request logging and monitoring
- [ ] Regular security updates
- [ ] Automated backups
- [ ] Disaster recovery plan

## Scaling

### Horizontal Scaling

1. **Load Balancing**: Use Nginx, HAProxy, or cloud load balancer
2. **Multiple Instances**: Run multiple backend instances
3. **Database Connection Pooling**: Configure in SQLAlchemy
4. **Caching**: Implement Redis caching (optional)

### Vertical Scaling

Increase server resources (CPU, RAM) as needed.

## Troubleshooting

### Backend won't start

```bash
# Check logs
docker-compose logs backend

# Check database connection
psql postgresql://user:password@localhost/med_predict

# Check port
netstat -tuln | grep 8000
```

### Database connection error

```bash
# Verify PostgreSQL is running
psql -U postgres -c "SELECT 1"

# Check DATABASE_URL
echo $DATABASE_URL

# Test connection
python -c "import psycopg2; psycopg2.connect(os.environ['DATABASE_URL'])"
```

### High memory usage

- Reduce worker count in Gunicorn
- Implement connection pooling
- Add caching layer
- Optimize database queries

### Slow responses

- Profile with `py-spy`
- Check database query performance
- Implement caching
- Scale horizontally

## Updates and Maintenance

### Zero-downtime deployment

1. Build new image: `docker build -t med-predict-backend:v2 .`
2. Start new container with new image
3. Update load balancer to point to new container
4. Gracefully shutdown old container

### Database migrations

```bash
# Before updating backend, backup database
pg_dump -U postgres med_predict > backup.sql

# Run migration script if needed
psql -U postgres -d med_predict -f migrations/002_schema_updates.sql

# Start new backend version
```

## Support

For deployment issues:
1. Check logs: `docker-compose logs`
2. Verify environment variables: `.env`
3. Test database connection
4. Check system resources (disk, memory, CPU)
5. Review security group/firewall rules
