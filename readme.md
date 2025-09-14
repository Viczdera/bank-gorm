# 🏦 Simple Bank API

A robust and secure RESTful banking API built with Go, featuring JWT authentication, PostgreSQL database, and comprehensive transaction management. Built this project demonstrate modern Go development practices with clean architecture, thorough testing, and production-ready deployment.

## 🚀 Features

### Core Banking Operations
- **User Management**: User registration and authentication with secure password hashing
- **Account Management**: Create and manage bank accounts with multi-currency support
- **Money Transfers**: Secure transfers between accounts with transaction integrity
- **Account Entries**: Track all account activities and balance changes

### Security & Authentication
- **JWT & PASETO Tokens**: Dual token implementation for flexible authentication
- **Password Security**: Bcrypt password hashing with salt
- **Protected Routes**: Middleware-based route protection
- **User Authorization**: Ownership validation for account operations

### Technical Excellence
- **Database Transactions**: ACID-compliant money transfers with deadlock prevention
- **SQL Generation**: Type-safe database operations using SQLC
- **Comprehensive Testing**: Unit tests with mocking for reliable code
- **API Validation**: Request validation with custom currency validators
- **Docker Support**: Containerized deployment with multi-stage builds
- **Database Migrations**: Version-controlled schema management

## 🛠️ Technology Stack

- **Language**: Go 1.24
- **Web Framework**: Gin HTTP framework
- **Database**: PostgreSQL with SQLC for type-safe queries
- **Authentication**: JWT & PASETO tokens
- **Testing**: Go testing with GoMock for mocking
- **Containerization**: Docker & Docker Compose
- **Configuration**: Viper for environment management
- **Validation**: Go Playground Validator

## 📊 Database Schema

The application uses a normalized PostgreSQL schema with the following entities:

- **Users**: Authentication and profile information
- **Accounts**: Bank accounts with currency and balance tracking
- **Transfers**: Money transfer records between accounts
- **Entries**: Account activity logs for audit trails

## 🚦 API Endpoints

### Public Endpoints
```
POST /users              # Register new user
POST /users/auth/login   # User authentication
```

### Protected Endpoints (Requires Authentication)
```
GET  /users/:username    # Get user profile
POST /accounts           # Create new account
GET  /accounts/:id       # Get account details
GET  /accounts/all       # List user accounts
POST /transfer           # Transfer money between accounts
```

## 🏃‍♂️ Quick Start

### Prerequisites
- Go 1.24+
- Docker & Docker Compose
- PostgreSQL (if running locally)
- Make (for using Makefile commands)

### 1. Clone the Repository
```bash
git clone https://github.com/Viczdera/bank.git
cd bank
```

### 2. Environment Setup
```bash
# Copy and configure environment variables
cp app.env.example app.env
# Edit app.env with your configuration
```

### 3. Run with Docker Compose (Recommended)
```bash
# Start all services
docker-compose up -d

# The API will be available at http://localhost:9090
```

### 4. Run Locally (Alternative)
```bash
# Start PostgreSQL
make postgres

# Create database
make createdb

# Run migrations
make migrateup

# Start the server
make server
```

## 🧪 Testing

```bash
# Run all tests with coverage
make test

# Generate mocks
make mock

# Run specific package tests
go test -v ./api/...
```

## 📝 Configuration

The application uses environment variables for configuration:

```env
DB_DRIVER=postgres
DB_SOURCE=postgres://user:password@localhost:5432/s_bank?sslmode=disable
SERVER_ADDRESS=0.0.0.0:9090
ACCESS_TOKEN_SYMMETRIC_KEY=your-32-character-secret-key
ACCESS_TOKEN_DURATION=15m
```

## 🔧 Development Commands

The project includes a comprehensive Makefile for development tasks:

```bash
make postgres       # Start PostgreSQL container
make createdb       # Create database
make dropdb         # Drop database
make migrateup      # Run all migrations
make migratedown    # Rollback all migrations
make migrateup1     # Run one migration up
make migratedown1   # Rollback one migration
make sqlc           # Generate SQLC code
make test           # Run tests with coverage
make server         # Start development server
make mock           # Generate mocks for testing
```

## 🏗️ Project Structure

```
bank/
├── api/                    # HTTP handlers and routes
│   ├── server.go          # Server setup and routing
│   ├── user.go            # User endpoints
│   ├── account.go         # Account endpoints
│   ├── transfer.go        # Transfer endpoints
│   ├── middleware.go      # Authentication middleware
│   └── validator.go       # Custom validators
├── db/
│   ├── migrations/        # Database migration files
│   ├── query/             # SQL queries for SQLC
│   ├── sqlc/              # Generated Go code from SQLC
│   └── mock/              # Generated mocks for testing
├── token/                 # JWT and PASETO token implementations
├── util/                  # Utility functions and configuration
├── docker-compose.yaml    # Multi-service deployment
├── dockerfile             # Container build instructions
└── Makefile              # Development automation
```

## 🔒 Security Features

- **Password Hashing**: Bcrypt with configurable cost
- **Token-based Auth**: JWT and PASETO token support
- **Request Validation**: Comprehensive input validation
- **SQL Injection Prevention**: Parameterized queries via SQLC
- **Authorization Checks**: User ownership validation
- **Secure Headers**: Production-ready security middleware

## 🚀 Deployment

### Docker Deployment
```bash
# Build and run with Docker Compose
docker-compose up --build -d
```

### Production Considerations
- Set strong `ACCESS_TOKEN_SYMMETRIC_KEY`
- Use environment-specific configuration
- Set up proper logging and monitoring
- Configure database connection pooling
- Implement rate limiting for production

## 🧪 API Testing

### Example API Calls

#### Register User
```bash
curl -X POST http://localhost:9090/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secret123",
    "full_name": "John Doe",
    "email": "john@example.com"
  }'
```

#### Login
```bash
curl -X POST http://localhost:9090/users/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "secret123"
  }'
```

#### Create Account (requires auth token)
```bash
curl -X POST http://localhost:9090/accounts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{
    "currency": "USD"
  }'
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 🙏 Acknowledgments

- Built with [Gin](https://gin-gonic.com/) web framework
- Database operations powered by [SQLC](https://sqlc.dev/)
- Authentication using [golang-jwt](https://github.com/golang-jwt/jwt)
- Configuration management with [Viper](https://github.com/spf13/viper)
- Testing enhanced with [Testify](https://github.com/stretchr/testify)

## 📞 Me

Victor Chidera - [GitHub](https://github.com/Viczdera)

