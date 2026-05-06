# Golang Docker

A production-ready Go REST API demonstrating best practices, including proper project structure following [golang-standards/project-layout](https://github.com/golang-standards/project-layout), Docker containerization, and clean architecture patterns.

## Features

- REST API built with Gin framework
- PostgreSQL database integration
- Docker containerization with multi-stage build
- Docker Compose for local development
- Proper project structure following golang-standards
- Clean architecture pattern (handler → repository → database)
- Centralized initialization in `internal/init`
- Configuration management via environment variables
- Secure credential handling with `.env` file
- CRUD operations for products

## Quick Start

### Prerequisites

- Go 1.25.0+
- Docker & Docker Compose
- Make (optional, but recommended)

### 1. Clone / Download Project

```bash
cd golang-docker
```

### 2. Setup Environment Variables

```bash
# Copy .env.example to .env (creates local credentials)
cp .env.example .env
```

The `.env` file contains default credentials for local development. **Do NOT commit this to Git!** It's already in `.gitignore`.

### 3. Start Everything with Docker Compose

```bash
# Start database + application
make docker-compose-up

# Or without make:
docker-compose -f deployments/docker-compose.yml up
```

This will:
- Start PostgreSQL container (port 5432)
- Build & start Go application (port 8080)
- Automatically create database tables
- Load environment variables from `.env`

### 4. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Get all products
curl http://localhost:8080/products

# Create a product
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":15000000}'

# Get product by ID
curl http://localhost:8080/products/1

# Update product
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Gaming Laptop","price":20000000}'

# Delete product
curl -X DELETE http://localhost:8080/products/1
```

### 5. Stop Everything

```bash
make docker-compose-down

# Or without make:
docker-compose -f deployments/docker-compose.yml down
```

---

## Project Structure

```
golang-docker/
├── cmd/
│   └── api/                 # Application entry point
│       └── main.go          # Minimal main - orchestration only
├── internal/                # Private application code (Go compiler enforced)
│   ├── init/                # Initialization logic
│   │   └── init.go          # Database setup & table creation
│   ├── handler/             # HTTP request handlers
│   │   └── product.go
│   ├── model/               # Domain models
│   │   └── product.go
│   └── repository/          # Data access layer
│       └── product.go
├── configs/                 # Configuration
│   └── config.go            # Load config from env variables
├── build/
│   └── package/
│       └── Dockerfile       # Multi-stage Docker build
├── deployments/
│   └── docker-compose.yml   # Docker Compose configuration
├── go.mod
├── go.sum
├── Makefile                 # Development commands
├── .env.example             # Template for environment variables
├── .env                     # Local environment (NEVER COMMIT!)
├── .gitignore
├── .dockerignore
└── README.md
```

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/products` | Get all products |
| GET | `/products/:id` | Get product by ID |
| POST | `/products` | Create new product |
| PUT | `/products/:id` | Update product |
| DELETE | `/products/:id` | Delete product |

### Example Payloads

**Create Product:**
```json
{
  "name": "MacBook Pro",
  "price": 25000000
}
```

**Update Product:**
```json
{
  "name": "MacBook Pro M3",
  "price": 27000000
}
```

---

## Running the Project

### Option 1: Docker Compose (Recommended - Everything in Docker)

```bash
# 1. Setup env
cp .env.example .env

# 2. Start
make docker-compose-up

# 3. Test
curl http://localhost:8080/health

# 4. Stop
make docker-compose-down
```

**Advantages:**
- No local dependencies needed
- Database automatically managed
- All config in `.env`
- One command to start everything

---

### Option 2: Local Development (Without Docker)

#### Prerequisites

- Go 1.25.0+ installed
- PostgreSQL running locally
- `psql` command available

#### Steps

```bash
# 1. Create database locally
createdb go_app

# 2. Setup environment variables
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=your_postgres_password
export DB_NAME=go_app
export PORT=8080

# 3. Download dependencies
make deps

# 4. Run the app
make run

# 5. In another terminal, test it
curl http://localhost:8080/health
```

---

### Option 3: Build & Run Binary

```bash
# 1. Build binary
make build

# 2. Run binary
./main
```

---

## Makefile Commands

**Common commands to use:**

```bash
# View all available commands
make help

# Development
make run                   # Run app locally
make dev                   # Run with hot reload
make build                 # Build binary

# Testing
make test                  # Run tests
make test-coverage         # Tests with coverage report

# Code Quality
make fmt                   # Format code
make vet                   # Check code errors
make lint                  # Lint code

# Docker
make docker-build          # Build Docker image
make docker-compose-up     # Start containers
make docker-compose-down   # Stop containers
make docker-compose-logs   # View logs

# Cleanup
make clean                 # Remove build artifacts
make deps                  # Download dependencies
```

---

## Environment Variables

Configure these in `.env` file:

```bash
# Server
PORT=8080

# Database
DB_HOST=db                    # 'db' for Docker, 'localhost' for local dev
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres          # Change this in production!
DB_NAME=go_app
```

For production, use secure methods:
- Environment variables set by CI/CD
- Kubernetes Secrets
- AWS Secrets Manager
- HashiCorp Vault

See `ENV_SETUP_GUIDE.md` for detailed security information.

---

## Common Workflows

### Morning: Start Development

```bash
cd golang-docker
cp .env.example .env
make docker-compose-up
# Wait for "Database connected successfully"
make run
```

### During Development

```bash
# In one terminal - run app
make run

# In another terminal - run tests
make test

# Format & lint code
make fmt
make vet
```

### Before Committing

```bash
make clean
make test-coverage
make lint
```

### Cleanup

```bash
make docker-compose-down
make clean
```

---

## Docker Details

### Build Docker Image

```bash
make docker-build
```

Creates image: `golang-docker:latest`

### Run Docker Container

```bash
# Run the container
docker run -p 8080:8080 golang-docker:latest

# With environment variables
docker run -p 8080:8080 \
  -e DB_HOST=db \
  -e DB_PASSWORD=mysecret \
  golang-docker:latest
```

### Multi-Stage Build

The Dockerfile uses multi-stage build:
- **Stage 1**: Compiles Go code (uses full Go image)
- **Stage 2**: Runs binary (uses minimal Alpine image)

Result: Final image is ~20MB instead of 300MB+

---

## Architecture

### Initialization Flow

```
cmd/api/main.go
    ↓
configs.New()                    # Load from .env
    ↓
internal/init.InitDB(cfg)        # Connect to database + create tables
    ↓
repository.NewProductRepository  # Create data layer
    ↓
handler.NewProductHandler        # Create HTTP handlers
    ↓
setupRouter()                    # Configure routes
    ↓
router.Run()                     # Start server
```

### Request Flow

```
HTTP Request
    ↓
Handler (validates, parses JSON)
    ↓
Repository (database queries)
    ↓
Database (PostgreSQL)
    ↓
Response (JSON)
```

---

## Security

### Credentials Management

✅ **`.env.example`** - Template (COMMIT THIS)
✅ **`.env`** - Local credentials (IN .gitignore, NEVER COMMIT)
✅ **Production** - Use secrets management (Vault, K8s Secrets, etc.)

See `ENV_SETUP_GUIDE.md` for detailed security practices.

---

## Troubleshooting

### "Connection refused" error

```bash
# Check if containers are running
docker ps

# View logs
make docker-compose-logs

# Restart
make docker-compose-down
make docker-compose-up
```

### "Database doesn't exist"

```bash
# This should be automatic, but if not:
docker exec go_app_db createdb -U postgres go_app
```

### "Port already in use"

```bash
# Find what's using port 8080
lsof -i :8080

# Or change port in .env
PORT=8081
```

### "go: command not found"

Install Go from https://go.dev/dl/

---

## Next Steps

1. ✅ Clone/download this project
2. ✅ Run `cp .env.example .env`
3. ✅ Run `make docker-compose-up`
4. ✅ Test with curl commands above
5. 📝 Modify for your use case
6. 🚀 Push to GitHub

---

## Technologies

- **Language**: Go 1.25.0
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 15
- **Container**: Docker & Docker Compose
- **Runtime**: Alpine Linux

---

## Learning Resources

- [Golang Standards Project Layout](https://github.com/golang-standards/project-layout)
- [Gin Documentation](https://gin-gonic.com/)
- [Docker Documentation](https://docs.docker.com/)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

---

## License

MIT

---

## Support

For detailed information on:
- **Project structure**: See `STRUCTURE_GUIDE.txt`
- **Environment setup**: See `ENV_SETUP_GUIDE.md`
- **Makefile commands**: See `MAKEFILE_EXPLAINED.md`