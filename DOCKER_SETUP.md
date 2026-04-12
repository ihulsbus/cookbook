# Docker Compose Setup for Cookbook Microservices

This docker-compose configuration provides all the infrastructure dependencies needed to run the Cookbook microservices locally.

## Services Included

### Database (PostgreSQL 15)
- **Single PostgreSQL instance** on Port 5432
- **6 separate databases** (one per microservice):
  - `recipe` (user: recipe)
  - `ingredient` (user: ingredient)
  - `instruction` (user: instruction)
  - `metadata` (user: metadata)
  - `image` (user: image)
  - `search` (user: search)

### Storage
- **MinIO** (S3-compatible):
  - API: Port 9000
  - Console: Port 9001 (http://localhost:9001)
  - Bucket: `cbhbe`
  - Credentials: `admin` / `adminpassword`

### Message Queue
- **RabbitMQ**:
  - AMQP: Port 5672
  - Management UI: Port 15672 (http://localhost:15672)
  - Credentials: `admin` / `wENbwHUirBwZQbHpRB46egb9C`

### Reverse Proxy
- **Traefik**:
  - HTTP: Port 80
  - Dashboard: Port 8080 (http://localhost:8080)

## Quick Start

### 1. Ensure init script is executable
```bash
chmod +x init-db.sh
```

### 2. Start all dependencies
```bash
docker-compose up -d
```

### 3. Verify all services are running
```bash
docker-compose ps
```

### 4. Verify databases were created
```bash
docker exec -it cookbook-postgres psql -U admin -d postgres -c '\l'
```

You should see all 6 databases listed: `recipe`, `ingredient`, `instruction`, `metadata`, `image`, and `search`.

### 5. Check service health
```bash
# Check logs
docker-compose logs -f

# Check specific service
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
```

### 6. Access Management Interfaces

- **MinIO Console**: http://localhost:9001
  - Username: `admin`
  - Password: `adminpassword`

- **RabbitMQ Management**: http://localhost:15672
  - Username: `admin`
  - Password: `wENbwHUirBwZQbHpRB46egb9C`

- **Traefik Dashboard**: http://localhost:8080

## Connecting Your Services

### Update config.yaml for each service

All services connect to the **same PostgreSQL instance** (`localhost:5432`) but use **different databases and users**.

#### Recipe Service
```yaml
database:
  host: localhost
  username: recipe
  password: recipe_pass
  port: 5432
  database: recipe
  sslMode: disable
  timezone: Europe/Amsterdam

rabbitmq:
  host: localhost:5672
  username: admin
  password: wENbwHUirBwZQbHpRB46egb9C
```

#### Ingredient Service
```yaml
database:
  host: localhost
  username: ingredient
  password: ingredient_pass
  port: 5432
  database: ingredient
  sslMode: disable
  timezone: Europe/Amsterdam
```

#### Instruction Service
```yaml
database:
  host: localhost
  username: instruction
  password: instruction_pass
  port: 5432
  database: instruction
  sslMode: disable
  timezone: Europe/Amsterdam
```

#### Metadata Service
```yaml
database:
  host: localhost
  username: metadata
  password: metadata_pass
  port: 5432
  database: metadata
  sslMode: disable
  timezone: Europe/Amsterdam
```

#### Image Service
```yaml
database:
  host: localhost
  username: image
  password: image_pass
  port: 5432
  database: image
  sslMode: disable
  timezone: Europe/Amsterdam

s3:
  region: us-east-1
  AWSAccessKey: admin
  AWSAccessSecret: adminpassword
  bucketName: cbhbe
  endpoint: http://localhost:9000

rabbitmq:
  host: localhost:5672
  username: admin
  password: wENbwHUirBwZQbHpRB46egb9C
```

#### Search Service
```yaml
database:
  host: localhost
  username: search
  password: search_pass
  port: 5432
  database: search
  sslMode: disable
  timezone: Europe/Amsterdam
```

## Database Connections

All databases are in the **same PostgreSQL instance** on port **5432**:

| Service | Host | Port | Database | Username | Password |
|---------|------|------|----------|----------|----------|
| Recipe | localhost | 5432 | recipe | recipe | recipe_pass |
| Ingredient | localhost | 5432 | ingredient | ingredient | ingredient_pass |
| Instruction | localhost | 5432 | instruction | instruction | instruction_pass |
| Metadata | localhost | 5432 | metadata | metadata | metadata_pass |
| Image | localhost | 5432 | image | image | image_pass |
| Search | localhost | 5432 | search | search | search_pass |

### Connect with psql
```bash
# Recipe database
psql -h localhost -p 5432 -U recipe -d recipe

# Ingredient database
psql -h localhost -p 5432 -U ingredient -d ingredient

# List all databases as admin
psql -h localhost -p 5432 -U admin -d postgres -c '\l'

# Switch between databases
psql -h localhost -p 5432 -U admin -d postgres
\c recipe
\c ingredient
```

### Direct docker exec
```bash
# Access PostgreSQL container directly
docker exec -it cookbook-postgres psql -U admin -d postgres

# Inside psql, switch databases:
\c recipe
\c ingredient
\c instruction
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

### Stop and remove all containers
```bash
docker-compose down
```

### Stop and remove all containers + volumes (⚠️ deletes all data)
```bash
docker-compose down -v
```

### Restart just the database
```bash
docker-compose restart postgres
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f postgres
docker-compose logs -f rabbitmq
docker-compose logs -f minio
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
CREATE DATABASE recipe OWNER recipe;
GRANT ALL PRIVILEGES ON DATABASE recipe TO recipe;
```

### MinIO bucket not created
```bash
docker-compose restart minio-init
docker-compose logs minio-init
```

### Database connection refused
Check if the database is healthy:
```bash
docker-compose ps
docker-compose logs postgres
```

### RabbitMQ not starting
Clear the volume and restart:
```bash
docker-compose down
docker volume rm cookbook_rabbitmq-data
docker-compose up -d rabbitmq
```

### Port 5432 already in use
If you have another PostgreSQL instance running:
```bash
# Stop local PostgreSQL
sudo systemctl stop postgresql

# Or change the port in docker-compose.yml
ports:
  - "5433:5432"  # Use 5433 on host instead
```

## Network

All services are connected via the `cookbook-network` bridge network. Services can communicate with each other using their service names (e.g., `postgres`, `rabbitmq`, `minio`).

## Data Persistence

Data is persisted in Docker volumes:
- `postgres-data` - All 6 databases
- `minio-data` - S3 bucket data
- `rabbitmq-data` - Message queue data

### Backup database
```bash
# Backup all databases
docker exec cookbook-postgres pg_dumpall -U admin > backup.sql

# Backup specific database
docker exec cookbook-postgres pg_dump -U recipe recipe > recipe_backup.sql
```

### Restore database
```bash
# Restore all databases
docker exec -i cookbook-postgres psql -U admin < backup.sql

# Restore specific database
docker exec -i cookbook-postgres psql -U recipe recipe < recipe_backup.sql
```

## Optional: Keycloak for Local OAuth Testing

Uncomment the Keycloak section in `docker-compose.yml` to run a local Keycloak instance for OAuth testing.

**Note**: You'll need to add a `keycloak` database to `init-db.sh` first:

```bash
# Add to init-db.sh before uncommenting keycloak in docker-compose.yml
CREATE USER keycloak WITH PASSWORD 'keycloak_pass';
CREATE DATABASE keycloak OWNER keycloak;
GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
```

Then start:
```bash
docker-compose down -v
chmod +x init-db.sh
docker-compose up -d
```

Access Keycloak at http://localhost:8081 with credentials `admin` / `admin`.

## Benefits of Single PostgreSQL Instance

✅ **Simpler management** - One container instead of 6
✅ **Lower resource usage** - Shared PostgreSQL processes
✅ **Easier backups** - Single pg_dumpall command
✅ **Faster startup** - One healthcheck instead of 6
✅ **Better performance** - Shared buffer pool and connections
✅ **Logical isolation** - Each service still has its own database and user

While each microservice has its own database, they share the same PostgreSQL instance, which is ideal for local development.
