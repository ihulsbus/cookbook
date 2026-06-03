# Docker Compose Setup for Cookbook Microservices

This docker-compose configuration provides all the infrastructure dependencies needed to run the Cookbook microservices locally.

## Services Included

### Database (PostgreSQL 15)
- **Single PostgreSQL instance** on Port 5432
- **7 separate databases** (one per microservice + Keycloak):
  - `recipes` (user: recipe)
  - `ingredients` (user: ingredient)
  - `instructions` (user: instruction)
  - `metadata` (user: metadata)
  - `images` (user: image)
  - `search` (user: search)
  - `keycloak` (user: keycloak) - for authentication

### Storage
- **MinIO** (S3-compatible):
  - API: Port 9000
  - Console: Port 9001 (http://localhost:9001)
  - Bucket: `cbhbe`
  - Credentials: `admin` / `adminpassword`

### Message Queue
- **RabbitMQ 3**:
  - AMQP: Port 5672
  - Management UI: Port 15672 (http://localhost:15672)
  - Credentials: `admin` / `wENbwHUirBwZQbHpRB46egb9C`

### Authentication
- **Keycloak**:
  - Port: 8081 (http://localhost:8081)
  - Admin Credentials: `administrator` / `changeme` (configurable via env)
  - Uses PostgreSQL for persistence
  - Pre-configured for development

### Reverse Proxy
- **Traefik 2.10**:
  - HTTP: Port 80
  - Dashboard: Port 8080 (http://localhost:8080)
  - Auto-discovers services with Docker labels

## Quick Start

### 1. Ensure init script is executable
```bash
chmod +x init-db.sh
```

### 2. Set environment variables (optional)
```bash
# Create a .env file in the project root
cat > .env <<EOF
KEYCLOAK_ADMIN_PASSWORD=mysecretpassword
KEYCLOAK_POSTGRESQL_PASSWORD=keycloak_pass
EOF
```

### 3. Start all dependencies
```bash
docker-compose up -d
```

This will start:
- PostgreSQL with automatic database initialization
- MinIO with automatic bucket creation
- RabbitMQ
- Keycloak (with health checks)
- Traefik reverse proxy

### 4. Verify all services are running
```bash
docker-compose ps
```

All services should show status as "Up" or "Up (healthy)".

### 5. Verify databases were created
```bash
docker exec -it cookbook-postgres psql -U admin -d postgres -c '\l'
```

You should see all 7 databases listed: `recipes`, `ingredients`, `instructions`, `metadata`, `images`, `search`, and `keycloak`.

### 6. Check service health
```bash
# Check all logs
docker-compose logs -f

# Check specific service
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
docker-compose logs -f keycloak
```

### 7. Access Management Interfaces

- **MinIO Console**: http://localhost:9001
  - Username: `admin`
  - Password: `adminpassword`
  - View and manage S3 buckets and objects

- **RabbitMQ Management**: http://localhost:15672
  - Username: `admin`
  - Password: `wENbwHUirBwZQbHpRB46egb9C`
  - Monitor queues, exchanges, and messages

- **Traefik Dashboard**: http://localhost:8080
  - View routing rules and service health

- **Keycloak Admin Console**: http://localhost:8081
  - Username: `administrator`
  - Password: `changeme` (or value from .env)
  - Configure realms, clients, and users

## Connecting Your Microservices

All services connect to the **same PostgreSQL instance** (`localhost:5432`) but use **different databases and users**.

### Environment Variables for Services

#### Recipe Service
```bash
# Database
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=recipes
export cbb_database_username=recipe
export cbb_database_password=recipe_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam

# RabbitMQ
export cbb_rabbitmq_host=localhost
export cbb_rabbitmq_port=5672
export cbb_rabbitmq_username=admin
export cbb_rabbitmq_password=wENbwHUirBwZQbHpRB46egb9C

# Keycloak (if using)
export cbb_auth0_domain=http://localhost:8081
export cbb_auth0_clientid=cookbook-client
export cbb_auth0_audience=cookbook-api

# Debug
export cbb_debug=true
```

#### Ingredient Service
```bash
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=ingredients
export cbb_database_username=ingredient
export cbb_database_password=ingredient_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam
```

#### Instruction Service
```bash
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=instructions
export cbb_database_username=instruction
export cbb_database_password=instruction_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam
```

#### Metadata Service
```bash
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=metadata
export cbb_database_username=metadata
export cbb_database_password=metadata_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam
```

#### Image Service
```bash
# Database
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=images
export cbb_database_username=image
export cbb_database_password=image_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam

# S3/MinIO
export cbb_s3_endpoint=http://localhost:9000
export cbb_s3_key=admin
export cbb_s3_secret=adminpassword
export cbb_s3_bucket=cbhbe

# RabbitMQ
export cbb_rabbitmq_host=localhost
export cbb_rabbitmq_port=5672
export cbb_rabbitmq_username=admin
export cbb_rabbitmq_password=wENbwHUirBwZQbHpRB46egb9C
```

#### Search Service
```bash
export cbb_database_host=localhost
export cbb_database_port=5432
export cbb_database_database=search
export cbb_database_username=search
export cbb_database_password=search_pass
export cbb_database_sslmode=disable
export cbb_database_timezone=Europe/Amsterdam
```

## Database Connections

All databases are in the **same PostgreSQL instance** on port **5432**:

| Service | Host | Port | Database | Username | Password |
|---------|------|------|----------|----------|----------|
| Recipe | localhost | 5432 | recipes | recipe | recipe_pass |
| Ingredient | localhost | 5432 | ingredients | ingredient | ingredient_pass |
| Instruction | localhost | 5432 | instructions | instruction | instruction_pass |
| Metadata | localhost | 5432 | metadata | metadata | metadata_pass |
| Image | localhost | 5432 | images | image | image_pass |
| Search | localhost | 5432 | search | search | search_pass |
| Keycloak | localhost | 5432 | keycloak | keycloak | keycloak_pass |

### Connect with psql
```bash
# Recipe database
psql -h localhost -p 5432 -U recipe -d recipes

# Ingredient database
psql -h localhost -p 5432 -U ingredient -d ingredients

# List all databases as admin
psql -h localhost -p 5432 -U admin -d postgres -c '\l'

# Connect as admin and switch databases
psql -h localhost -p 5432 -U admin -d postgres
\c recipes
\c ingredients
```

### Direct docker exec
```bash
# Access PostgreSQL container directly
docker exec -it cookbook-postgres psql -U admin -d postgres

# Inside psql, switch databases:
\c recipes
\c ingredients
\c instructions
\l  # list all databases
```

## Managing the Stack

### Stop all services
```bash
docker-compose stop
```

### Start all services
```bash
docker-compose start
```

### Restart all services
```bash
docker-compose restart
```

### Restart a specific service
```bash
docker-compose restart postgres
docker-compose restart rabbitmq
docker-compose restart keycloak
```

### Stop and remove all containers
```bash
docker-compose down
```

### Stop and remove all containers + volumes (⚠️ deletes all data)
```bash
docker-compose down -v
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
docker-compose logs -f minio
docker-compose logs -f keycloak

# Last 100 lines
docker-compose logs --tail=100 postgres
```

### Check resource usage
```bash
docker stats
```

## Keycloak Configuration

Keycloak is included for OAuth2/OIDC authentication.

### Initial Setup

1. Access Keycloak at http://localhost:8081
2. Login with `administrator` / `changeme`
3. Create a new realm called `cookbook`
4. Create a client:
   - Client ID: `cookbook-client`
   - Client Protocol: `openid-connect`
   - Access Type: `public` (for SPAs) or `confidential` (for backend)
   - Valid Redirect URIs: `http://localhost:*`
5. Create test users in the realm
6. Note the realm endpoint: `http://localhost:8081/realms/cookbook`

### Environment Variables for Services

```bash
export cbb_auth0_domain=http://localhost:8081/realms/cookbook
export cbb_auth0_clientid=cookbook-client
export cbb_auth0_audience=cookbook-api
```

## Troubleshooting

### Databases not created

Check if the init script ran:
```bash
docker-compose logs postgres | grep "All databases"
```

If not found, the script may not have executed. Recreate:
```bash
docker-compose down -v
chmod +x init-db.sh
docker-compose up -d
```

### Check which databases exist
```bash
docker exec -it cookbook-postgres psql -U admin -d postgres -c '\l'
```

### Manually create missing database
```bash
docker exec -it cookbook-postgres psql -U admin -d postgres
```
Then run:
```sql
CREATE USER recipe WITH PASSWORD 'recipe_pass';
CREATE DATABASE recipes OWNER recipe;
GRANT ALL PRIVILEGES ON DATABASE recipes TO recipe;
```

### MinIO bucket not created

Check minio-init logs:
```bash
docker-compose logs minio-init
```

Restart bucket creation:
```bash
docker-compose restart minio-init
```

Verify bucket exists:
```bash
docker exec minio-init mc ls myminio
```

### Database connection refused

Check if PostgreSQL is healthy:
```bash
docker-compose ps postgres
docker-compose logs postgres
```

Wait for health check to pass:
```bash
# PostgreSQL may take 30-60 seconds to be fully ready
docker-compose logs -f postgres | grep "database system is ready"
```

### RabbitMQ not starting

Clear the volume and restart:
```bash
docker-compose down
docker volume rm cookbook_rabbitmq-data
docker-compose up -d rabbitmq
```

### Keycloak taking too long to start

Keycloak has a long startup time (2-5 minutes). Check logs:
```bash
docker-compose logs -f keycloak
```

Look for: `Keycloak ... started`

### Port conflicts

If ports are already in use:

```bash
# Check what's using the port
sudo lsof -i :5432  # PostgreSQL
sudo lsof -i :9000  # MinIO
sudo lsof -i :5672  # RabbitMQ
sudo lsof -i :8081  # Keycloak

# Stop conflicting service or change port in docker-compose.yml
```

### Services can't connect to each other

Ensure you're using the correct network. From inside containers, use service names:
- Database: `postgres:5432`
- RabbitMQ: `rabbitmq:5672`
- MinIO: `minio:9000`

From your host machine (running services locally), use `localhost`.

### Clear all data and start fresh

```bash
# Stop everything and remove volumes
docker-compose down -v

# Remove any orphaned volumes
docker volume prune

# Start fresh
docker-compose up -d
```

## Network

All services are connected via the `cookbook-network` bridge network. Services can communicate with each other using their service names:
- `postgres` - PostgreSQL database
- `rabbitmq` - Message queue
- `minio` - S3-compatible storage
- `keycloak` - Authentication server
- `traefik` - Reverse proxy

## Data Persistence

Data is persisted in Docker volumes:
- `postgres-data` - All 7 databases (recipes, ingredients, instructions, metadata, images, search, keycloak)
- `minio-data` - S3 bucket data and objects
- `rabbitmq-data` - Message queue data and configurations

### Backup database

```bash
# Backup all databases
docker exec cookbook-postgres pg_dumpall -U admin > backup_$(date +%Y%m%d).sql

# Backup specific database
docker exec cookbook-postgres pg_dump -U recipe recipes > recipes_backup_$(date +%Y%m%d).sql

# Backup with compression
docker exec cookbook-postgres pg_dump -U recipe recipes | gzip > recipes_backup_$(date +%Y%m%d).sql.gz
```

### Restore database

```bash
# Restore all databases
docker exec -i cookbook-postgres psql -U admin < backup.sql

# Restore specific database
docker exec -i cookbook-postgres psql -U recipe recipes < recipes_backup.sql

# Restore from compressed backup
gunzip -c recipes_backup.sql.gz | docker exec -i cookbook-postgres psql -U recipe recipes
```

### Backup MinIO data

```bash
# List all buckets
docker exec minio-init mc ls myminio

# Mirror bucket to local directory
docker exec minio-init mc mirror myminio/cbhbe /backup/minio

# Or use docker cp
docker cp minio:/data ./minio-backup
```

## Performance Tips

### Increase PostgreSQL shared buffers
Edit docker-compose.yml:
```yaml
postgres:
  command: postgres -c shared_buffers=256MB -c max_connections=200
```

### Increase RabbitMQ memory limit
Edit docker-compose.yml:
```yaml
rabbitmq:
  environment:
    RABBITMQ_VM_MEMORY_HIGH_WATERMARK: 1GB
```

### Enable MinIO caching
Edit docker-compose.yml:
```yaml
minio:
  environment:
    MINIO_CACHE: "on"
    MINIO_CACHE_DRIVES: "/cache"
  volumes:
    - minio-cache:/cache
```

## Benefits of Single PostgreSQL Instance

✅ **Simpler management** - One container instead of 7
✅ **Lower resource usage** - Shared PostgreSQL processes
✅ **Easier backups** - Single pg_dumpall command
✅ **Faster startup** - One healthcheck instead of 7
✅ **Better performance** - Shared buffer pool and connection pooling
✅ **Logical isolation** - Each service still has its own database and user
✅ **Transaction support** - Can use distributed transactions if needed

While each microservice has its own database, they share the same PostgreSQL instance, which is ideal for local development and testing.

## Production Considerations

For production deployments:

1. **Separate PostgreSQL instances** - Consider one database per service for true isolation
2. **Managed services** - Use AWS RDS, Azure Database, or GCP Cloud SQL
3. **Secrets management** - Use Vault, AWS Secrets Manager, or Kubernetes secrets
4. **High availability** - Set up PostgreSQL replication and failover
5. **Monitoring** - Add Prometheus, Grafana, and alerting
6. **Backups** - Automated daily backups with retention policies
7. **SSL/TLS** - Enable SSL for all database connections
8. **Resource limits** - Set CPU and memory limits for all containers
9. **Security** - Change all default passwords and use strong credentials
10. **Network isolation** - Use separate networks for different service tiers
