# Go Inventory System with API Gateway

A microservice-based inventory system built entirely in Go that demonstrates best practices for production-ready backend systems. This project showcases idiomatic Go, microservices architecture, API gateway patterns, and modern software engineering principles.

## System Architecture

### Components

1. **API Gateway (port 8000)**
   - Routes external requests to internal microservices
   - Handles JWT authentication, rate limiting, logging, and metrics
   - Provides unified entry point for all client requests

2. **Auth Service (port 8083)**
   - Handles user registration and login
   - Issues and validates JWT tokens
   - Manages user authentication state

3. **User Service (port 8081)**
   - Handles CRUD operations for user accounts
   - Manages user profiles and preferences
   - Provides user data management

4. **Order Service (port 8082)**
   - Handles CRUD operations for product orders
   - Manages order lifecycle and status
   - Provides order tracking and history

5. **Shared Module**
   - Common models and utility functions
   - Shared configuration and database schemas
   - Reusable components across services

## Technology Stack

- **Go 1.21+** - Primary programming language
- **SQLite** - Lightweight database (can be replaced with PostgreSQL)
- **GORM** - Object-Relational Mapping
- **JWT** - JSON Web Tokens for authentication
- **Gorilla Mux** - HTTP router and URL matcher
- **Prometheus** - Metrics collection and monitoring
- **Docker & Docker Compose** - Containerization and orchestration

## Project Structure

```
go-inventory-system/
├── gateway/                 # API Gateway service
│   ├── config/              # Configuration loading
│   ├── middleware/          # Auth, rate-limit, logging, metrics
│   ├── router/              # Route registration
│   └── main.go              # Entry point
│
├── services/               
│   ├── auth/                # Auth microservice
│   │   ├── handler.go       # Handlers for login/register
│   │   ├── db.go            # Database initialization
│   │   └── main.go          # Service entry point
│   │
│   ├── users/               # User microservice
│   │   ├── handler.go       # CRUD + profile endpoints
│   │   ├── db.go            # Database initialization
│   │   └── main.go          # Service entry point
│   │
│   └── orders/              # Order microservice
│       ├── handler.go       # Order endpoints
│       ├── db.go            # Database initialization
│       └── main.go          # Service entry point
│
├── shared/                 # Common utilities
│   ├── models.go            # Shared model types
│   ├── utils.go             # Hashing, validation, etc.
│   └── config.go            # Shared config structs
│
├── docker-compose.yml      # Service orchestration
├── go.mod                  # Module definition
├── routes.yaml             # Gateway routing config
└── README.md               # Documentation
```

## Getting Started

### Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose (for containerized deployment)

### Local Development

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd go-inventory-system
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Start individual services**
   ```bash
   # Start API Gateway
   go run gateway/main.go

   # Start Auth Service (in new terminal)
   go run services/auth/main.go

   # Start Users Service (in new terminal)
   go run services/users/main.go

   # Start Orders Service (in new terminal)
   go run services/orders/main.go
   ```

### Docker Deployment

1. **Build and start all services**
   ```bash
   docker-compose up --build
   ```

2. **Access the API Gateway**
   - Gateway: http://localhost:8000
   - Auth Service: http://localhost:8083
   - Users Service: http://localhost:8081
   - Orders Service: http://localhost:8082

## API Endpoints

### Authentication

- `POST /auth/register` - Register a new user
- `POST /auth/login` - Login user and get JWT token

### Users

- `GET /users` - List all users
- `POST /users` - Create a new user
- `GET /users/{id}` - Get specific user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user
- `GET /users/me` - Get current user profile

### Orders

- `GET /orders` - List all orders
- `POST /orders` - Create a new order
- `GET /orders/{id}` - Get specific order
- `PUT /orders/{id}` - Update order
- `DELETE /orders/{id}` - Delete order
- `GET /orders/user/{user_id}` - Get orders for specific user

### Health Checks

- `GET /health` - Service health check
- `GET /metrics` - Prometheus metrics

## Configuration

### Environment Variables

- `PORT` - Service port (default: 8080)
- `DATABASE_URL` - Database connection string
- `JWT_SECRET` - JWT signing secret
- `ENVIRONMENT` - Environment (development/production)
- `LOG_LEVEL` - Logging level

### Gateway Routes

Routes are configured in `routes.yaml`:

```yaml
routes:
  - path: /auth
    backend: http://localhost:8083
    methods: ["GET", "POST"]
  - path: /users
    backend: http://localhost:8081
    methods: ["GET", "POST", "PUT", "DELETE"]
  - path: /orders
    backend: http://localhost:8082
    methods: ["GET", "POST", "PUT", "DELETE"]
```

## Testing

### Manual Testing with curl

1. **Register a user**
   ```bash
   curl -X POST http://localhost:8000/auth/register \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123",
       "username": "testuser"
     }'
   ```

2. **Login and get token**
   ```bash
   curl -X POST http://localhost:8000/auth/login \
     -H "Content-Type: application/json" \
     -d '{
       "email": "user@example.com",
       "password": "password123"
     }'
   ```

3. **Access protected endpoint**
   ```bash
   curl -X GET http://localhost:8000/users/me \
     -H "Authorization: Bearer <your-jwt-token>"
   ```

4. **Create an order**
   ```bash
   curl -X POST http://localhost:8000/orders \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer <your-jwt-token>" \
     -d '{
       "product_name": "Laptop",
       "quantity": 1,
       "total_price": 999.99
     }'
   ```

## Security Features

- **JWT Authentication** - Stateless token-based authentication
- **Password Hashing** - bcrypt for secure password storage
- **Rate Limiting** - Prevents abuse with configurable limits
- **CORS Support** - Cross-origin resource sharing
- **Input Validation** - Request validation and sanitization

## Monitoring & Observability

- **Request Logging** - Detailed request/response logging
- **Prometheus Metrics** - HTTP request counts and durations
- **Health Checks** - Service health monitoring
- **Graceful Shutdown** - Proper service termination

## Best Practices Implemented

### Code Organization
- **Modular Design** - Single responsibility principle
- **Interface-driven** - Loose coupling between components
- **Error Handling** - Comprehensive error management
- **Configuration Management** - Environment-based configuration

### Security
- **Authentication Middleware** - JWT validation
- **Password Security** - bcrypt hashing
- **Input Validation** - Request sanitization
- **CORS Configuration** - Cross-origin security

### Performance
- **Connection Pooling** - Database connection management
- **Rate Limiting** - Request throttling
- **Graceful Shutdown** - Proper resource cleanup
- **Metrics Collection** - Performance monitoring

### Deployment
- **Containerization** - Docker for consistent environments
- **Service Orchestration** - Docker Compose for multi-service deployment
- **Environment Configuration** - Flexible configuration management
- **Health Monitoring** - Service health checks

## Production Considerations

### Scalability
- **Microservices Architecture** - Independent service scaling
- **Database Optimization** - Indexing and query optimization
- **Load Balancing** - Multiple service instances
- **Caching** - Redis for session and data caching

### Monitoring
- **Application Metrics** - Custom business metrics
- **Infrastructure Monitoring** - System resource monitoring
- **Log Aggregation** - Centralized logging (ELK stack)
- **Alerting** - Proactive issue detection

### Security
- **HTTPS/TLS** - Encrypted communication
- **API Rate Limiting** - DDoS protection
- **Input Sanitization** - XSS and injection prevention
- **Audit Logging** - Security event tracking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

---
