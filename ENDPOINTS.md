# API Endpoints Reference

## Base URL
```
http://localhost:8080
```

---

## Quick Reference

### Root & Health
```bash
# API Documentation (JSON)
curl http://localhost:8080/

# Quick Health Check
curl http://localhost:8080/health

# Detailed Health Check
curl http://localhost:8080/api/v1/health
```

---

## Available Endpoints

### 1. Health & Info

| Method | Endpoint | Description | Response |
|--------|----------|-------------|----------|
| GET | `/` | API documentation | JSON with all endpoints |
| GET | `/health` | Quick health check | `OK` |
| GET | `/api/v1/health` | Detailed health check | `OK` |

**Example:**
```bash
# Root - shows all available endpoints
curl http://localhost:8080/

# Quick health check
curl http://localhost:8080/health
```

---

### 2. User Management

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| GET | `/api/v1/users` | List all users | - |
| POST | `/api/v1/users` | Create user | JSON with user data |
| GET | `/api/v1/users/{id}` | Get user by ID | - |
| DELETE | `/api/v1/users/{id}` | Delete user | - |

**Examples:**
```bash
# List all users
curl http://localhost:8080/api/v1/users

# Get specific user
curl http://localhost:8080/api/v1/users/123

# Create user (example)
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/api/v1/users/123
```

---

### 3. URL Shortening

| Method | Endpoint | Description | Request Body |
|--------|----------|-------------|--------------|
| POST | `/api/v1/shorten` | Shorten a URL | JSON with URL |
| GET | `/api/v1/resolve` | Resolve short URL | JSON with short code |

**Examples:**
```bash
# Shorten a URL
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/very/long/url",
    "custom_short": "mylink",
    "expiry": "24h"
  }'

# Resolve a short URL
curl -X GET http://localhost:8080/api/v1/resolve \
  -H "Content-Type: application/json" \
  -d '{"short_code": "mylink"}'
```

---

## Common Issues & Solutions

### ❌ 404 Not Found

**Problem:**
```bash
curl http://localhost:8080/users
# 404 Not Found
```

**Solution:**
All API endpoints are under `/api/v1/`:
```bash
curl http://localhost:8080/api/v1/users
# ✅ Works!
```

**Exception:** Health check works at both:
- `/health` (quick shortcut)
- `/api/v1/health` (versioned endpoint)

---

### ❌ Wrong HTTP Method

**Problem:**
```bash
curl http://localhost:8080/api/v1/users/123
# Trying to delete but using GET
```

**Solution:**
```bash
curl -X DELETE http://localhost:8080/api/v1/users/123
# ✅ Works!
```

---

### ❌ Missing Content-Type

**Problem:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -d '{"username":"alice"}'
# May fail or return error
```

**Solution:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"username":"alice"}'
# ✅ Works!
```

---

## Testing with HTTPie (Alternative)

If you have [HTTPie](https://httpie.io/) installed:

```bash
# Root documentation
http GET localhost:8080/

# Health check
http GET localhost:8080/health

# List users
http GET localhost:8080/api/v1/users

# Create user
http POST localhost:8080/api/v1/users \
  username=alice \
  email=alice@example.com

# Delete user
http DELETE localhost:8080/api/v1/users/123

# Shorten URL
http POST localhost:8080/api/v1/shorten \
  url=https://example.com/long/url \
  custom_short=mylink
```

---

## Quick Test Script

Save this as `test-api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing Shawty-UR API..."
echo ""

echo "1. Root Documentation:"
curl -s $BASE_URL/ | jq .
echo ""

echo "2. Health Check (root):"
curl -s $BASE_URL/health
echo ""

echo "3. Health Check (versioned):"
curl -s $BASE_URL/api/v1/health
echo ""

echo "4. List Users:"
curl -s $BASE_URL/api/v1/users
echo ""

echo "5. Get User by ID:"
curl -s $BASE_URL/api/v1/users/123
echo ""

echo "Done!"
```

Run it:
```bash
chmod +x test-api.sh
./test-api.sh
```

---

## Summary Table

| Endpoint | Method | Purpose | Works? |
|----------|--------|---------|--------|
| `/` | GET | API docs | ✅ |
| `/health` | GET | Quick health | ✅ |
| `/api/v1/health` | GET | Health check | ✅ |
| `/api/v1/users` | GET | List users | ✅ |
| `/api/v1/users` | POST | Create user | ✅ |
| `/api/v1/users/{id}` | GET | Get user | ✅ |
| `/api/v1/users/{id}` | DELETE | Delete user | ✅ |
| `/api/v1/shorten` | POST | Shorten URL | ✅ |
| `/api/v1/resolve` | GET | Resolve URL | ✅ |

---

## Need More Help?

- Full documentation: [README.md](README.md)
- Quick start: [QUICKSTART.md](QUICKSTART.md)
- Setup guide: [SETUP_COMPLETE.md](SETUP_COMPLETE.md)

---

**Remember:** Always use `/api/v1/` prefix for API endpoints! 🚀
