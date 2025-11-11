# Quick Start Guide - Shawty-UR

## 🚀 Get Running in 3 Minutes

### Step 1: Start Docker (30 seconds)
```bash
cd /home/onizuka/gophercises/shawty-ur
make docker-up
```
✅ PostgreSQL running on `localhost:5432`
✅ Redis running on `localhost:6379`

---

### Step 2: Setup Database (1 minute)
```bash
# Install Goose migration tool
make migrate-install

# Run all migrations (creates tables)
make migrate-up
```
✅ Users table created
✅ URLs table created
✅ Analytics table created

---

### Step 3: Build & Run (30 seconds)
```bash
# Build the app
make build

# Run the app
make run
```
✅ Server running at `http://localhost:8080`

---

## 🧪 Test It Works

Open a new terminal and run:

```bash
# Health check
curl http://localhost:8080/api/v1/health
# Expected: OK

# List users endpoint
curl http://localhost:8080/api/v1/users
# Expected: List users
```

---

## 🎯 One-Command Setup

If you prefer, do everything at once:

```bash
# Complete setup: Docker + Goose + Migrations
make setup

# Then build and run
make build && make run
```

---

## 📊 Available Endpoints

- `GET  /api/v1/health` - Health check
- `GET  /api/v1/users` - List users
- `POST /api/v1/users` - Create user
- `GET  /api/v1/users/{id}` - Get user
- `DELETE /api/v1/users/{id}` - Delete user
- `POST /api/v1/shorten` - Shorten URL
- `GET  /api/v1/resolve` - Resolve short URL

---

## 🛑 Stop Everything

```bash
# Stop the app: Ctrl+C

# Stop Docker containers
make docker-down
```

---

## 🔧 Common Commands

| Command | Description |
|---------|-------------|
| `make help` | Show all commands |
| `make docker-up` | Start PostgreSQL + Redis |
| `make docker-logs` | View database logs |
| `make migrate-status` | Check migration status |
| `make dev` | Run without building |
| `make clean` | Remove build artifacts |

---

## ❓ Troubleshooting

**Port 8080 already in use?**
```bash
# Change port in .env
ADDR=:3000
```

**Database connection error?**
```bash
# Check Docker is running
docker ps

# Restart containers
make restart
```

**Can't connect to Redis?**
```bash
# Test Redis
docker exec -it shawty-redis redis-cli -a waifu_waguri ping
# Should respond: PONG
```

---

## 📚 Need More Info?

See [README.md](README.md) for complete documentation.

---

**Happy coding! 🎉**
