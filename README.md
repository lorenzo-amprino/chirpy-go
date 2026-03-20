# Chirpy

A lightweight Twitter-like social media API built with Go, featuring user authentication, chirp (post) management, and webhook integrations.

## Features

- **User Management**: Create users, authenticate with JWT tokens, update profiles
- **Chirp Management**: Create, retrieve, delete chirps with author filtering and sorting
- **Authentication**: Secure JWT-based authentication with refresh tokens
- **Token Management**: Refresh and revoke token support
- **Content Validation**: Automatic content cleaning and length validation
- **Webhooks**: Polka webhook integration for premium features
- **Admin Tools**: Metrics tracking and database reset functionality
- **Database**: PostgreSQL-backed persistence with migration support

## Tech Stack

- **Language**: Go 1.25.3
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: Argon2id
- **ID Generation**: UUID v4
- **Environment Management**: godotenv

## Installation

### Prerequisites

- Go 1.25.3 or higher
- PostgreSQL 12+

### Setup

1. Clone the repository:
```bash
git clone https://github.com/lorenzo-amprino/chirpy-go.git
cd chirpy-go
```

2. Install dependencies:
```bash
go mod download
go mod tidy
```

3. Set up environment variables:
```bash
cp .env.example .env
```

4. Configure your `.env` file:
```env
DB_URL=postgres://user:password@localhost:5432/chirpy
JWT_SECRET=your-secret-key-here
POLKA_API_KEY=your-polka-api-key
```

5. Run database migrations:
```bash
# Migrations are in sql/schema/
```

6. Build and run:
```bash
go build -o chirpy
./chirpy
```

The server will start on `http://localhost:8080`

## API Endpoints

### Health & Admin

- `GET /api/healthz` - Health check
- `GET /admin/metrics` - View server metrics
- `POST /admin/reset` - Reset database (admin only)

### Users

- `POST /api/users` - Create new user
- `PUT /api/users` - Update user (requires authentication)
- `POST /api/login` - Authenticate user and get JWT token

### Chirps

- `POST /api/chirps` - Create a new chirp (requires authentication)
- `GET /api/chirps` - Get all chirps (supports `?author_id=<uuid>` and `?sort=asc|desc`)
- `GET /api/chirps/{id}` - Get a specific chirp by ID
- `DELETE /api/chirps/{id}` - Delete a chirp (requires authentication, owner only)

### Token Management

- `POST /api/refresh` - Refresh JWT token using refresh token
- `POST /api/revoke` - Revoke refresh token

### Validation & Webhooks

- `POST /api/validate_chirp` - Validate and clean chirp content
- `POST /api/polka/webhooks` - Handle Polka webhook events (premium upgrades)

## Usage Examples

### Create a User
```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"securepassword"}'
```

### Create a Chirp
```bash
curl -X POST http://localhost:8080/api/chirps \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{"body":"Hello, world!"}'
```

### Get All Chirps (Sorted Descending)
```bash
curl "http://localhost:8080/api/chirps?sort=desc"
```

### Get Chirps by Author
```bash
curl "http://localhost:8080/api/chirps?author_id=<uuid>"
```

### Delete a Chirp
```bash
curl -X DELETE http://localhost:8080/api/chirps/<id> \
  -H "Authorization: Bearer <jwt-token>"
```

## Project Structure

```
chirpy-go/
├── main.go                          # Server entry point
├── handle_*.go                      # HTTP request handlers
├── json.go                          # JSON response utilities
├── metrics.go                       # Metrics tracking
├── readiness.go                     # Health checks
├── reset.go                         # Admin reset handler
├── internal/
│   ├── auth/                        # Authentication & JWT logic
│   │   ├── auth.go
│   │   ├── jwt.go
│   │   └── auth_test.go
│   └── database/                    # Database models & queries
│       ├── db.go
│       ├── models.go
│       └── *.sql.go                 # Auto-generated SQLC files
├── sql/
│   ├── schema/                      # Database migrations
│   │   ├── 001_user.sql
│   │   ├── 002_chirps.sql
│   │   ├── 003_users.sql
│   │   ├── 004_refresh_token.sql
│   │   └── 005_users.sql
│   └── queries/                     # SQL queries
│       ├── user.sql
│       ├── chirps.sql
│       └── refresh_token.sql
├── go.mod                           # Go module definition
├── go.sum                           # Dependency checksums
└── sqlc.yaml                        # SQLC configuration
```

## Development

### Running Tests
```bash
go test ./...
```

### Building
```bash
go build -o chirpy
```

### Database Migrations

Migrations are SQL files located in `sql/schema/`. Apply new migrations by updating your database connection.

### Generating Database Code

Database query code is generated using [sqlc](https://sqlc.dev/). Update queries in `sql/queries/` and regenerate with:
```bash
sqlc generate
```

## Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `DB_URL` | PostgreSQL connection string | Yes |
| `JWT_SECRET` | Secret key for JWT token signing | Yes |
| `POLKA_API_KEY` | API key for Polka webhook verification | No |

## Security Notes

- Always use HTTPS in production
- Keep `JWT_SECRET` secure and rotate regularly
- Implement rate limiting for production use
- Validate all user inputs
- Use strong password requirements
- Store sensitive data securely in environment variables

## License

This project is part of an educational learning path and serves as a reference implementation.

## Author

Lorenzo Amprino
