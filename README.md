# Internal Developer Portal

A production-ready Internal Developer Portal built with Clean Architecture, implementing 15 proven design patterns. Features Go backend with PostgreSQL and Redis, Vue.js 3 frontend with TypeScript, complete Docker setup, and comprehensive CI/CD workflows.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/shennawardana23/example-mcp-pub.git
cd example-mcp-pub

# Copy environment file
cp .env.example .env

# Start the entire stack with Docker Compose
make start

# Or start individual components:
# Backend with hot reload
make dev

# Frontend development server (in another terminal)
make frontend-dev
```

The application will be available at:
- **Frontend**: http://localhost:5173
- **API**: http://localhost:8080
- **API Documentation**: http://localhost:8080/swagger/index.html
- **Health Check**: http://localhost:8080/health

**Demo Credentials:**
- Username: `admin`
- Password: `admin123`

## 📋 Features

### Core Functionality
1. **Service Catalog** - Browse, search, and filter services with advanced filtering
2. **API Documentation Browser** - Interactive Swagger/OpenAPI documentation
3. **Developer Onboarding Hub** - Centralized dashboard for new developers
4. **Team Dashboard** - Real-time metrics and service health monitoring
5. **Metrics Viewer** - Service performance and health metrics

### Architecture & Patterns
Implements 15 proven patterns (see ARCHITECTURE.md for details):
- ✅ Dependency Injection (DI)
- ✅ Repository Pattern
- ✅ Service Layer
- ✅ Controller Pattern
- ✅ Middleware Chain
- ✅ Entity/Model Separation
- ✅ Converter Pattern
- ✅ Query Builder
- ✅ Caching Strategy (Redis)
- ✅ Error Translation
- ✅ Multi-layer Validation
- ✅ Context Propagation
- ✅ Config Management
- ✅ Authentication (JWT + Basic)
- ✅ Interface Segregation

## 🏗️ Project Structure

```
.
├── cmd/
│   └── api/                    # Application entry point
├── internal/
│   ├── app/
│   │   ├── controller/         # HTTP request handlers
│   │   ├── service/            # Business logic layer
│   │   ├── repository/         # Data access layer
│   │   ├── middleware/         # Auth, logging, CORS
│   │   ├── entity/             # Database entities
│   │   ├── model/              # API request/response models
│   │   ├── converter/          # Entity/Model conversion
│   │   └── util/               # Utility functions
│   ├── config/                 # Configuration management
│   └── database/               # Database & Redis clients
├── api/                        # Vue.js 3 frontend
│   ├── src/
│   │   ├── components/         # Reusable Vue components
│   │   ├── views/              # Page components
│   │   ├── composables/        # Composition API logic
│   │   ├── stores/             # Pinia state management
│   │   ├── router/             # Vue Router config
│   │   └── utils/              # API client & utilities
│   ├── vite.config.ts          # Vite configuration
│   └── tailwind.config.js      # Tailwind CSS config
├── migrations/                 # Database migration files
├── docs/                       # Additional documentation
├── docker-compose.yml          # Multi-service Docker setup
├── Dockerfile                  # Backend container
├── Makefile                    # Development commands
└── .github/workflows/          # CI/CD workflows
```

## 🛠️ Technology Stack

### Backend
- **Language**: Go 1.25+
- **Framework**: Gin (HTTP router)
- **Database**: PostgreSQL 16
- **Cache**: Redis 7
- **ORM**: database/sql (standard library)
- **Authentication**: JWT (golang-jwt/jwt)
- **Logging**: Logrus (structured logging)
- **Documentation**: Swagger/OpenAPI
- **Monitoring**: Sentry/NewRelic ready

### Frontend
- **Framework**: Vue.js 3 (Composition API)
- **Language**: TypeScript (strict mode)
- **Build Tool**: Vite
- **State**: Pinia
- **Routing**: Vue Router 4
- **Styling**: Tailwind CSS 3
- **HTTP Client**: Axios (with interceptors)
- **Testing**: Vitest

### Infrastructure
- **Containerization**: Docker & Docker Compose
- **CI/CD**: GitHub Actions
- **Hot Reload**: Air (Go) & Vite (Vue)

## 📦 Installation & Setup

### Prerequisites
- Go 1.25 or higher
- Node.js 18+ and npm
- PostgreSQL 16+
- Redis 7+
- Docker & Docker Compose (optional)
- Make

### Local Development (without Docker)

1. **Install Dependencies**
   ```bash
   # Backend dependencies
   make deps
   
   # Frontend dependencies
   make frontend-install
   ```

2. **Setup Database**
   ```bash
   # Create database
   createdb devportal
   
   # Run migrations
   make migrate
   ```

3. **Configure Environment**
   ```bash
   cp .env.example .env
   # Edit .env with your local settings
   ```

4. **Start Services**
   ```bash
   # Terminal 1: Backend with hot reload
   make dev
   
   # Terminal 2: Frontend dev server
   make frontend-dev
   ```

### Docker Development

```bash
# Start all services
make start

# View logs
make docker-logs

# Stop all services
make stop
```

## 🧪 Testing

```bash
# Backend tests
make test

# Backend test coverage
make test-coverage
make test-coverage-html

# Frontend tests
make frontend-test

# Run all checks (format, vet, test)
make check
```

## 📝 Available Make Commands

Run `make help` to see all available commands:

```bash
make help
```

Key commands:
- `make build` - Build the API binary
- `make dev` - Start backend with Air hot reload
- `make test` - Run backend tests
- `make frontend-dev` - Start frontend dev server
- `make start` - Start entire stack with Docker
- `make migrate` - Run database migrations
- `make check` - Run all quality checks

## 🔒 Security Features

- **JWT Bearer Token Authentication** with refresh tokens
- **Password Hashing** using bcrypt
- **Input Validation** at all layers (controller, service, repository)
- **SQL Injection Prevention** via parameterized queries
- **CORS Configuration** with whitelist
- **Security Headers** (CSP, X-Frame-Options, etc.)
- **Rate Limiting** ready (middleware available)
- **Correlation IDs** for request tracing

## 🗄️ Database

### Schema
- **users**: User accounts and authentication
- **services**: Service catalog entries
- **teams**: Team organization
- **metrics**: Service metrics and monitoring data

### Caching Strategy
- **User Data**: 30 minutes TTL
- **Services/Projects**: 1 hour TTL
- **Menu/Navigation**: 24 hours TTL

### Multi-tenancy
PostgreSQL schema-based multi-tenancy with:
- Permissions schema for access control
- Audit schema for change tracking

## 🔄 API Endpoints

### Authentication
- `POST /api/v1/auth/login` - User login (returns JWT)
- `GET /api/v1/auth/me` - Get current user info

### Services
- `GET /api/v1/services` - List services (with filtering)
- `GET /api/v1/services/:id` - Get service by ID
- `POST /api/v1/services` - Create service (auth required)
- `PUT /api/v1/services/:id` - Update service (auth required)
- `DELETE /api/v1/services/:id` - Delete service (auth required)

### Health
- `GET /health` - Health check endpoint

See `/swagger/index.html` for complete API documentation.

## 🎨 Frontend Features

- **Reactive State Management** with Pinia
- **Type-Safe** with TypeScript strict mode
- **Component Testing** with Vitest
- **Responsive Design** with Tailwind CSS
- **Route Guards** for authentication
- **API Interceptors** for token management
- **Error Handling** with user-friendly messages

## 🚢 Deployment

### Docker Deployment
```bash
# Build images
docker-compose build

# Deploy
docker-compose up -d

# Scale services
docker-compose up -d --scale api=3
```

### Environment Variables
See `.env.example` for all configuration options. Key variables:
- `JWT_SECRET` - **Must be changed in production**
- `DB_PASSWORD` - Database password
- `SENTRY_DSN` - Error tracking (optional)
- `LOG_LEVEL` - Logging verbosity

## 📚 Documentation

- **[ARCHITECTURE.md](docs/ARCHITECTURE.md)** - Detailed architecture and design patterns
- **[API Documentation](http://localhost:8080/swagger/index.html)** - Interactive API docs
- **Code Comments** - Inline documentation throughout codebase

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see LICENSE file for details

## 🙏 Acknowledgments

Built with Clean Architecture principles and industry best practices for production-ready applications.
