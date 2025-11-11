# Shawty-UR - URL Shortener Service

A modern URL shortening service built with Go, PostgreSQL, and Redis.

## Features

- 🔗 URL shortening with custom short codes
- 👥 User management and authentication
- 📊 Click analytics and tracking
- ⚡ Redis caching for rate limiting
- 🐘 PostgreSQL for persistent storage
- 🔄 Database migrations with Goose
- 🐳 Docker support for easy deployment

## Architecture

- **PostgreSQL**: Stores user data, URLs, and analytics
- **Redis**: Handles caching and rate limiting
- **Chi Router**: HTTP routing with middleware support
- **Goose**: Database migration management

## Prerequisites

- Go 1.25.1 or higher
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

All configuration is stored in [.env](.env):

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

# Rate Limiting
RATE_LIMIT_REQUESTS=10
RATE_LIMIT_WINDOW=1m
```

## API Endpoints

### Health Check
```bash
GET /api/v1/health
```

### User Management
```bash
GET    /api/v1/users       # List all users
POST   /api/v1/users       # Create user
GET    /api/v1/users/{id}  # Get user by ID
DELETE /api/v1/users/{id}  # Delete user
```

### URL Shortening
```bash
POST /api/v1/shorten   # Shorten a URL
GET  /api/v1/resolve   # Resolve a short URL
```

## Database Schema

### Users Table
```sql
- id (bigserial)
- username (varchar, unique)
- email (varchar, unique)
- password_hash (varchar)
- created_at (timestamp)
- updated_at (timestamp)
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

## Testing the API

### Health Check
```bash
curl http://localhost:8080/api/v1/health
# Response: OK
```

### List Users
```bash
curl http://localhost:8080/api/v1/users
# Response: List users
```

### Get User by ID
```bash
curl http://localhost:8080/api/v1/users/123
# Response: Get user: 123
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
│   ├── routes/                 # Route handlers
│   │   ├── health.go          # Health check
│   │   ├── users.go           # User routes
│   │   └── shorten.go         # URL shortening routes
│   └── utils/
│       ├── db/
│       │   └── db.go          # PostgreSQL connection
│       └── redis/
│           └── redis.go       # Redis connection
├── app/
│   └── app.go                 # Application struct & middleware
├── config/
│   └── config.go              # Configuration structs
├── migrations/                # Database migrations
│   ├── 00001_create_users_table.sql
│   ├── 00002_create_urls_table.sql
│   └── 00003_create_url_analytics_table.sql
├── bin/                       # Compiled binaries
├── .env                       # Environment configuration
├── docker-compose.yml         # Docker services
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
