# Architecture Documentation

## Overview

This Internal Developer Portal is built using **Clean Architecture** principles with **Dependency Injection** throughout. The system is organized into layers with clear separation of concerns and well-defined boundaries.

## 15 Proven Design Patterns

### 1. Dependency Injection (DI)

**Purpose**: Decouple components and enable testability by injecting dependencies rather than creating them internally.

**Implementation**:
```go
// in cmd/api/main.go
serviceRepo := repository.NewPostgresServiceRepository(db.DB)
serviceService := service.NewServiceService(serviceRepo, redisClient)
serviceController := controller.NewServiceController(serviceService)
```

**Benefits**:
- Easy to mock dependencies for testing
- Flexible component replacement
- Clear dependency graph

### 2. Repository Pattern

**Purpose**: Abstract data access logic and provide a collection-like interface for domain entities.

**Implementation**:
```go
// Interface definition
type ServiceRepository interface {
    Create(ctx context.Context, service *entity.Service) error
    GetByID(ctx context.Context, id int64) (*entity.Service, error)
    List(ctx context.Context, filters ServiceFilters) ([]entity.Service, int64, error)
    Update(ctx context.Context, service *entity.Service) error
    Delete(ctx context.Context, id int64) error
}

// PostgreSQL implementation
type PostgresServiceRepository struct {
    db *sql.DB
}
```

**Benefits**:
- Database-agnostic business logic
- Easy to swap data sources
- Centralized data access logic

### 3. Service Layer

**Purpose**: Encapsulate business logic separate from HTTP handlers and data access.

**Implementation**:
```go
type ServiceService interface {
    Create(ctx context.Context, req *model.ServiceRequest, username string) (*model.ServiceResponse, error)
    GetByID(ctx context.Context, id int64) (*model.ServiceResponse, error)
    List(ctx context.Context, req *model.ServiceListRequest) (*model.PaginatedResponse, error)
}

type serviceService struct {
    repo      repository.ServiceRepository
    converter *converter.ServiceConverter
    redis     *redis.Client
}
```

**Benefits**:
- Business logic reusability
- Transaction management
- Cache integration point

### 4. Controller Pattern

**Purpose**: Handle HTTP requests, perform validation, and delegate to service layer.

**Implementation**:
```go
type ServiceController struct {
    service service.ServiceService
}

func (ctrl *ServiceController) Create(c *gin.Context) {
    var req model.ServiceRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Handle validation error
    }
    response, err := ctrl.service.Create(c.Request.Context(), &req, username)
    c.JSON(http.StatusCreated, response)
}
```

**Benefits**:
- Thin controllers focused on HTTP concerns
- Standard error handling
- Request/response transformation

### 5. Middleware Pattern

**Purpose**: Process requests in a chain before reaching controllers.

**Implementation**:
- **AuthMiddleware**: JWT token validation
- **LoggerMiddleware**: Structured logging with correlation IDs
- **RecoveryMiddleware**: Panic recovery
- **CORSMiddleware**: Cross-origin request handling
- **SecurityHeadersMiddleware**: Security headers

**Benefits**:
- Cross-cutting concerns separation
- Reusable components
- Clean request pipeline

### 6. Entity/Model Separation

**Purpose**: Separate database entities from API models to prevent tight coupling.

**Entities** (internal/app/entity/):
- Represent database structure
- Include all database fields
- Used within repository layer

**Models** (internal/app/model/):
- API request/response structures
- Include validation tags
- Used in controller layer

**Benefits**:
- API evolution without database changes
- Security (hide sensitive fields)
- Different validation rules

### 7. Converter Pattern

**Purpose**: Transform between entities and models with centralized conversion logic.

**Implementation**:
```go
type ServiceConverter struct{}

func (c *ServiceConverter) ToResponse(service *entity.Service) *model.ServiceResponse {
    return &model.ServiceResponse{
        ID:          service.ID,
        Name:        service.Name,
        // ... other fields
    }
}

func (c *ServiceConverter) ToEntity(req *model.ServiceRequest) *entity.Service {
    return &entity.Service{
        Name:        req.Name,
        Description: req.Description,
        // ... other fields
    }
}
```

**Benefits**:
- Centralized transformation logic
- Consistent conversions
- Easy to maintain

### 8. Query Builder

**Purpose**: Build dynamic SQL queries safely based on filters.

**Implementation**:
```go
func (r *PostgresServiceRepository) List(ctx context.Context, filters ServiceFilters) ([]entity.Service, int64, error) {
    whereConditions := []string{}
    args := []interface{}{}
    argPos := 1

    if filters.Search != "" {
        whereConditions = append(whereConditions, fmt.Sprintf("(name ILIKE $%d OR description ILIKE $%d)", argPos, argPos))
        args = append(args, "%"+filters.Search+"%")
        argPos++
    }
    // ... build query dynamically
}
```

**Benefits**:
- SQL injection prevention
- Dynamic filtering
- Maintainable queries

### 9. Caching Strategy (Redis)

**Purpose**: Reduce database load and improve response times with intelligent caching.

**Implementation**:
```go
func (s *serviceService) GetByID(ctx context.Context, id int64) (*model.ServiceResponse, error) {
    cacheKey := fmt.Sprintf("service:%d", id)
    
    // Try cache first
    var cached model.ServiceResponse
    if err := s.redis.Get(ctx, cacheKey).Scan(&cached); err == nil {
        return &cached, nil
    }
    
    // Get from database and cache
    service, err := s.repo.GetByID(ctx, id)
    if err == nil {
        s.redis.Set(ctx, cacheKey, service, time.Hour)
    }
    return service, err
}
```

**Cache TTLs**:
- User data: 30 minutes
- Services: 1 hour
- Menus: 24 hours

**Benefits**:
- Improved performance
- Reduced database load
- Scalability

### 10. Error Translation

**Purpose**: Convert internal errors to user-friendly API responses.

**Implementation**:
```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Message string                 `json:"message"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// In controller
if err != nil {
    c.JSON(http.StatusInternalServerError, model.ErrorResponse{
        Error:   "creation_failed",
        Message: err.Error(),
    })
    return
}
```

**Benefits**:
- Consistent error format
- Hide internal details
- User-friendly messages

### 11. Multi-layer Validation

**Purpose**: Validate input at multiple layers for defense in depth.

**Layers**:
1. **Controller**: Struct tags (`binding:"required"`)
2. **Service**: Business rules validation
3. **Repository**: Database constraints

**Implementation**:
```go
type ServiceRequest struct {
    Name        string   `json:"name" binding:"required,min=3,max=100"`
    Description string   `json:"description" binding:"required,min=10,max=500"`
    Type        string   `json:"type" binding:"required,oneof=api library microservice frontend"`
}
```

**Benefits**:
- Early error detection
- Security (prevent invalid data)
- Data integrity

### 12. Context Propagation

**Purpose**: Pass request-scoped data through call chain.

**Implementation**:
```go
// Set correlation ID in middleware
c.Set("correlation_id", correlationID)

// Use in logging
logger.WithField("correlation_id", correlationID).Info("Request completed")

// Pass context through layers
func (s *serviceService) Create(ctx context.Context, req *model.ServiceRequest) {
    // ctx carries correlation ID, user info, timeouts
    return s.repo.Create(ctx, entity)
}
```

**Benefits**:
- Request tracing
- Timeout propagation
- User context access

### 13. Config Management

**Purpose**: Centralize configuration with environment variable support.

**Implementation**:
```go
type Config struct {
    Server   ServerConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}

func Load() (*Config, error) {
    cfg := &Config{
        Server: ServerConfig{
            Port: getEnv("SERVER_PORT", "8080"),
        },
    }
    return cfg.Validate()
}
```

**Benefits**:
- 12-factor app compliance
- Environment-specific config
- Type-safe configuration

### 14. Authentication (JWT + Basic)

**Purpose**: Secure API endpoints with industry-standard authentication.

**JWT Implementation**:
```go
type Claims struct {
    UserID   int64  `json:"user_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func (s *authService) Login(ctx context.Context, req *model.LoginRequest) (*model.LoginResponse, error) {
    // Verify password with bcrypt
    bcrypt.CompareHashAndPassword(user.PasswordHash, req.Password)
    
    // Generate JWT token
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return tokenString, nil
}
```

**Basic Auth** available for service-to-service communication.

**Benefits**:
- Stateless authentication
- Secure password storage
- Token expiration

### 15. Interface Segregation

**Purpose**: Define focused interfaces with minimal methods.

**Implementation**:
```go
// Specific interfaces instead of one large interface
type ServiceRepository interface {
    Create(ctx context.Context, service *entity.Service) error
    GetByID(ctx context.Context, id int64) (*entity.Service, error)
}

type UserRepository interface {
    GetByUsername(ctx context.Context, username string) (*entity.User, error)
}
```

**Benefits**:
- Easy to mock for testing
- Clear contracts
- SOLID principles compliance

## Layer Responsibilities

### Controller Layer
- HTTP request/response handling
- Input validation (struct tags)
- Authentication/authorization checks
- Error response formatting

### Service Layer
- Business logic execution
- Transaction management
- Cache interaction
- Business rule validation
- Error handling

### Repository Layer
- Database queries
- Data mapping
- Connection management
- Query optimization

## Data Flow

```
HTTP Request
    ↓
Middleware Chain (Auth, Logging, CORS)
    ↓
Controller (Validation, Request Binding)
    ↓
Service Layer (Business Logic, Caching)
    ↓
Repository (Data Access)
    ↓
Database (PostgreSQL)
```

## Testing Strategy

1. **Unit Tests**: Service and repository layers
2. **Integration Tests**: API endpoints with test database
3. **Component Tests**: Vue.js components with Vitest
4. **Mock Generation**: Interface-based mocking

## Security Considerations

- **Parameterized Queries**: Prevent SQL injection
- **Password Hashing**: bcrypt with cost 10
- **JWT Tokens**: HS256 signing, short-lived access tokens
- **Input Validation**: All layers
- **CORS**: Whitelist-based
- **Security Headers**: CSP, X-Frame-Options, etc.

## Performance Optimizations

- **Redis Caching**: Reduce database hits
- **Connection Pooling**: Reuse database connections
- **Indexes**: On frequently queried columns
- **Pagination**: Limit result sets
- **Graceful Shutdown**: Complete in-flight requests

## Monitoring & Observability

- **Structured Logging**: JSON format with correlation IDs
- **Health Checks**: Database and Redis status
- **Metrics Ready**: Sentry/NewRelic integration points
- **Request Tracing**: Correlation IDs throughout

## Scalability

- **Stateless API**: Horizontal scaling ready
- **Cache Layer**: Redis for shared state
- **Database**: Connection pooling and optimization
- **Containerized**: Easy deployment and scaling

This architecture provides a solid foundation for a production-ready application with excellent maintainability, testability, and scalability characteristics.
