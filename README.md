# Shawty-UR - URL Shortener Service

A modern URL shortening service built with Go, PostgreSQL, and Redis.

## Features

- 🔗 URL shortening with custom short codes
- 👥 User management and authentication
- 🔐 Google OAuth 2.0 authentication
- 🔒 Bcrypt password hashing for local accounts
- 📊 Click analytics and tracking
- ⚡ Redis caching for rate limiting
- 🎯 Tiered API quotas (authenticated vs free users)
- 📈 Prometheus metrics integration
- 📊 Grafana dashboards for monitoring
- 🐘 PostgreSQL for persistent storage
- 🔄 Database migrations with Goose
- 🐳 Docker support for easy deployment

## Architecture

- **PostgreSQL**: Stores user data, URLs, and analytics
- **Redis**: Handles caching and rate limiting
- **Chi Router**: HTTP routing with middleware support
- **Goose**: Database migration management
- **Prometheus**: Metrics collection and monitoring
- **Grafana**: Visualization and dashboards

## Prerequisites

- Go 1.23.0 or higher
- Docker & Docker Compose
- Make (optional, for convenience)

## Quick Start

### 1. Clone and Setup

```bash
cd /home/onizuka/gophercises/shawty-ur

# Install dependencies
go mod download
```

### 2. Start Docker Containers

```bash
# Using Make (recommended)
make docker-up

# OR manually
docker-compose up -d
```

This starts:
- PostgreSQL on `localhost:5432`
- Redis on `localhost:6379`

### 2b. Start with Monitoring (Optional)

To start the application with Prometheus and Grafana monitoring:

```bash
# Start all services including monitoring stack
docker-compose -f docker-compose.yml -f monitoring/docker-compose.monitoring.yml up -d
```

This additionally starts:
- **Prometheus** on `localhost:9090` - Metrics collection
- **Grafana** on `localhost:3000` - Dashboards (admin/admin)
- **PostgreSQL Exporter** on `localhost:9187` - Database metrics
- **Redis Exporter** on `localhost:9121` - Cache metrics

### 3. Install Goose and Run Migrations

```bash
# Install goose migration tool
make migrate-install

# Run migrations
make migrate-up

# Check migration status
make migrate-status
```

### 4. Build and Run

```bash
# Build the application
make build

# Run the application
make run

# OR run in development mode (no build)
make dev

# OR do everything in one command
make setup  # Docker + Goose + Migrations
make build && make run
```

## Configuration

All configuration is stored in `.env` (create from `.env.example`):

```bash
# Server
ADDR=:8080

# PostgreSQL (User data & persistent storage)
DB_DSN=postgres://howl:turnip_man1234@localhost:5432/social?sslmode=disable
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_MAX_IDLE_TIME=10m

# Redis (Caching & Rate Limiting)
REDIS_ADDR=localhost:6379
REDIS_PASSWORD=waifu_waguri
REDIS_DB=0

# Security
JWT_SECRET=waifu_waguri
SESSION_KEY=waifu-waguri

# Rate Limiting & API Quotas
API_QUOTA=100              # Quota for authenticated users (per 30 minutes)
API_QUOTE_FREE=10          # Quota for unauthenticated users (per 30 minutes)
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=1m

# OAuth 2.0 (Google)
GOOGLE_CLIENT_ID=your-client-id.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=your-client-secret
GOOGLE_REDIRECT_URL=http://localhost:8080/api/v1/auth/google/callback
```

### Setting up Google OAuth

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Create a new project or select existing one
3. Enable "Google+ API" or "Google People API"
4. Go to "Credentials" → "Create Credentials" → "OAuth 2.0 Client ID"
5. Configure OAuth consent screen
6. Add authorized redirect URI: `http://localhost:8080/api/v1/auth/google/callback`
7. Copy the Client ID and Client Secret to your `.env` file

**Important**: Keep your `.env` file secure and never commit it to version control.

## API Endpoints

### Health Check
```bash
GET /health
# Response: OK
```

### Metrics (Prometheus)
```bash
GET /metrics
# Response: Prometheus-formatted metrics
```

### Authentication

#### Local Authentication (Email/Password)
```bash
POST /api/v1/users/register
# Body: { "username": "john", "email": "john@example.com", "password": "secure123" }

POST /api/v1/users/login
# Body: { "username": "john", "password": "secure123" }
# OR:   { "email": "john@example.com", "password": "secure123" }
```

#### OAuth Authentication (Google)
```bash
GET /api/v1/auth/google
# Redirects to Google OAuth consent screen

GET /api/v1/auth/google/callback
# OAuth callback endpoint (handled automatically)
```

### User Management
```bash
GET    /api/v1/users       # List all users
POST   /api/v1/users       # Create user (deprecated - use /register)
GET    /api/v1/users/{id}  # Get user by ID
DELETE /api/v1/users/{id}  # Delete user
```

### URL Shortening

**Note**: Authentication is optional. Unauthenticated users get limited quota (API_QUOTE_FREE=10 requests per 30 min), while authenticated users get higher quota (API_QUOTA=100 requests per 30 min).

```bash
POST /api/v1/shorten   # Shorten a URL
GET  /{shortCode}      # Resolve and redirect to original URL
```

## Database Schema

### Users Table
```sql
- id (bigserial, primary key)
- username (varchar(255), unique, not null)
- email (varchar(255), unique, not null)
- password_hash (varchar(255), nullable) -- NULL for OAuth users
- provider (varchar(50), default 'local') -- 'local' or 'google'
- provider_id (varchar(255), nullable) -- OAuth provider's user ID
- avatar_url (varchar(512), nullable) -- Profile picture from OAuth
- email_verified (boolean, default false)
- created_at (timestamp with time zone)
- updated_at (timestamp with time zone)

INDEX: idx_users_provider_id ON (provider, provider_id) WHERE provider != 'local'
```

### URLs Table
```sql
- id (bigserial)
- user_id (bigint, FK to users)
- original_url (text)
- short_code (varchar, unique)
- custom_short (boolean)
- clicks (bigint)
- expires_at (timestamp)
- created_at (timestamp)
- updated_at (timestamp)
```

### URL Analytics Table
```sql
- id (bigserial)
- url_id (bigint, FK to urls)
- ip_address (inet)
- user_agent (text)
- referrer (text)
- country (varchar)
- city (varchar)
- clicked_at (timestamp)
```

## Makefile Commands

```bash
make help              # Show all available commands
make docker-up         # Start PostgreSQL and Redis
make docker-down       # Stop containers
make docker-logs       # View container logs
make docker-clean      # Remove containers and volumes

make migrate-install   # Install Goose
make migrate-up        # Run all migrations
make migrate-down      # Rollback last migration
make migrate-reset     # Rollback all migrations
make migrate-status    # Show migration status
make migrate-create NAME=create_table  # Create new migration

make deps              # Download Go dependencies
make build             # Build application
make run               # Run application
make dev               # Run without building
make clean             # Clean build artifacts

make setup             # Complete setup (Docker + Goose + Migrations)
make start             # Start everything (Docker + Build + Run)
make test              # Run tests
```

## Development Workflow

### 1. Create a New Migration

```bash
make migrate-create NAME=add_user_avatar
```

This creates two files in `migrations/`:
- `XXXXXX_add_user_avatar.sql`

Edit the file:
```sql
-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ADD COLUMN avatar_url VARCHAR(255);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users DROP COLUMN avatar_url;
-- +goose StatementEnd
```

### 2. Run the Migration

```bash
make migrate-up
```

### 3. Rollback if Needed

```bash
make migrate-down  # Rollback last migration
```

## Monitoring and Observability

The application includes comprehensive monitoring using Prometheus and Grafana.

### Available Metrics

The application exposes the following Prometheus metrics at `/metrics`:

#### HTTP Metrics
- `http_requests_total` - Total number of HTTP requests (labeled by method, endpoint, status)
- `http_request_duration_seconds` - Duration of HTTP requests in seconds (histogram)
- `active_connections` - Number of currently active connections (gauge)

#### Application Metrics
- `urls_shortened_total` - Total number of URLs shortened
- `urls_resolved_total` - Total number of URLs resolved

#### Cache Metrics
- `cache_hits_total` - Total number of cache hits
- `cache_misses_total` - Total number of cache misses
- `redis_operation_duration_seconds` - Duration of Redis operations (histogram)

#### Database Metrics
- `database_query_duration_seconds` - Duration of database queries (histogram)

### Accessing Monitoring Dashboards

After starting with monitoring enabled:

#### Prometheus
- URL: http://localhost:9090
- View raw metrics, create queries, check targets

#### Grafana
- URL: http://localhost:3000
- Default credentials: `admin` / `admin`
- Pre-configured dashboards for application, PostgreSQL, and Redis metrics

### Monitoring Stack Services

The monitoring setup includes:

1. **Prometheus** - Scrapes metrics from:
   - Application (`/metrics` endpoint)
   - PostgreSQL Exporter (database metrics)
   - Redis Exporter (cache metrics)

2. **Grafana** - Visualizes metrics with dashboards

3. **PostgreSQL Exporter** - Exports PostgreSQL database metrics

4. **Redis Exporter** - Exports Redis cache metrics

### Example Prometheus Queries

```promql
# Request rate per second
rate(http_requests_total[5m])

# 95th percentile response time
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Cache hit ratio
cache_hits_total / (cache_hits_total + cache_misses_total)

# Active connections over time
active_connections
```

## Testing the API

### Health Check
```bash
curl http://localhost:8080/api/v1/health
# Response: OK
```

### User Registration (Local)
```bash
curl -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "secure123"
  }'

# Response:
# {
#   "id": 1,
#   "username": "testuser",
#   "email": "test@example.com",
#   "provider": "local",
#   "email_verified": false,
#   "created_at": "2025-11-15T10:30:00Z"
# }
```

### User Login
```bash
curl -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "password": "secure123"
  }'

# Response:
# {
#   "status": "success",
#   "ID": 1,
#   "username": "testuser"
# }
```

### Google OAuth Login
```bash
# Open in browser:
open http://localhost:8080/api/v1/auth/google

# This will:
# 1. Redirect to Google login
# 2. Ask for consent
# 3. Redirect back to your callback
# 4. Create/update user in database
# 5. Set session cookie
```

### List Users
```bash
curl http://localhost:8080/api/v1/users
# Response: List of all users with their details
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/1
# Response: User details for ID 1
```

## Docker Commands (Manual)

```bash
# Start containers
docker-compose up -d

# View logs
docker-compose logs -f postgres
docker-compose logs -f redis

# Stop containers
docker-compose down

# Remove containers and volumes (clean slate)
docker-compose down -v

# Connect to PostgreSQL
docker exec -it shawty-postgres psql -U howl -d social

# Connect to Redis
docker exec -it shawty-redis redis-cli -a waifu_waguri
```

## Project Structure

```
shawty-ur/
├── api/
│   ├── main.go                 # Application entry point
│   ├── auth/                   # Authentication logic
│   │   ├── oauth.go           # OAuth configuration
│   │   └── session.go         # Session management
│   ├── middleware/             # HTTP middleware
│   │   ├── auth.go            # Optional authentication
│   │   └── prometheus.go      # Metrics collection
│   ├── metrics/                # Prometheus metrics definitions
│   │   └── metrics.go
│   ├── routes/                 # Route handlers
│   │   ├── auth.go            # OAuth routes
│   │   ├── health.go          # Health check
│   │   ├── users.go           # User routes
│   │   ├── shorten.go         # URL shortening
│   │   └── resolve.go         # URL resolution
│   └── utils/
│       ├── db/
│       │   └── db.go          # PostgreSQL connection
│       └── redisUtil/
│           └── redis.go       # Redis connection
├── app/
│   └── app.go                 # Application struct & middleware
├── config/
│   └── config.go              # Configuration structs
├── migrations/                # Database migrations
│   ├── 00001_create_users_table.sql
│   ├── 00002_create_urls_table.sql
│   └── 00003_create_url_analytics_table.sql
├── monitoring/                # Monitoring stack configuration
│   ├── docker-compose.monitoring.yml
│   ├── prometheus/
│   │   └── prometheus.yml     # Prometheus config
│   └── grafana/
│       ├── provisioning/
│       │   ├── datasources/   # Auto-configure Prometheus
│       │   └── dashboards/    # Dashboard provisioning
│       └── dashboards/        # Pre-built dashboards (JSON)
├── bin/                       # Compiled binaries
├── .env                       # Environment configuration
├── docker-compose.yml         # Docker services (app, db, redis)
├── Makefile                   # Build commands
└── README.md                  # This file
```

## Adding New Routes

The project uses a scalable route registration pattern:

### 1. Create Route Handler

Create `api/routes/products.go`:
```go
package routes

import (
    "net/http"
    "shawty-ur/app"
    "github.com/go-chi/chi/v5"
)

func RegisterProductRoutes(r chi.Router, application *app.Application) {
    r.Get("/products", listProducts(application))
    r.Post("/products", createProduct(application))
}

func listProducts(app *app.Application) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Access app.DbConnector, app.RedisClient
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("List products"))
    }
}
```

### 2. Register in main.go

```go
application.RegisterRoutes(
    routes.RegisterHealthRoutes,
    routes.RegisterUserRoutes,
    routes.RegisterServiceRoutes,
    routes.RegisterProductRoutes,  // Add this line
)
```

That's it! No changes to Application struct needed.

## Troubleshooting

### Database Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Check PostgreSQL logs
docker logs shawty-postgres

# Test connection
docker exec -it shawty-postgres pg_isready -U howl
```

### Redis Connection Issues

```bash
# Check if Redis is running
docker ps | grep redis

# Test Redis connection
docker exec -it shawty-redis redis-cli -a waifu_waguri ping
# Should respond: PONG
```

### Migration Issues

```bash
# Check migration status
make migrate-status

# Reset database (WARNING: deletes all data)
make migrate-reset
make migrate-up
```

### OAuth Issues

#### "invalid_client" Error
- Check that `GOOGLE_CLIENT_ID` has no leading/trailing spaces
- Verify `GOOGLE_CLIENT_SECRET` is the actual secret (not a placeholder)
- Ensure redirect URI in Google Console matches `GOOGLE_REDIRECT_URL` in `.env`

#### "column password does not exist" Error
- Database schema uses `password_hash`, not `password`
- Run migrations: `make migrate-up`
- Check migration 00004 was applied: `make migrate-status`

#### Password Hash Mismatch
- If you get bcrypt errors, existing users may have corrupted password hashes
- Delete test users and re-register:
  ```bash
  docker exec -it shawty-postgres psql -U howl -d social \
    -c "DELETE FROM users WHERE provider = 'local';"
  ```

### Port Already in Use

```bash
# Find process using port 8080
sudo lsof -i :8080

# Kill the process
kill -9 <PID>

# Or change port in .env
ADDR=:3000
```

## Production Deployment

### Build for Production

```bash
# Build optimized binary
CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/shawty-ur ./api

# Binary size will be ~6-8MB (vs 9.7MB debug build)
```

### Security Checklist for Production

- [ ] Use strong, random values for `JWT_SECRET` and `SESSION_KEY`
- [ ] Enable HTTPS with valid SSL certificates
- [ ] Update `GOOGLE_REDIRECT_URL` to production domain
- [ ] Set `DB_DSN` to use SSL mode: `sslmode=require`
- [ ] Use strong passwords for PostgreSQL and Redis
- [ ] Enable Redis authentication (`requirepass` in redis.conf)
- [ ] Set up rate limiting for authentication endpoints
- [ ] Configure CORS appropriately for your frontend domain
- [ ] Never commit `.env` file to version control
- [ ] Use environment variables or secrets management in production

### Docker Production Setup

Update `docker-compose.yml` for production:
- Use environment-specific `.env` files
- Enable SSL for PostgreSQL
- Use Redis password authentication
- Set up volume backups
- Configure logging

## License

MIT

## Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing`)
3. Commit changes (`git commit -am 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing`)
5. Open Pull Request

---

Built with ❤️ using Go, PostgreSQL, and Redis
