# Marko - Country Arrival Notifier Backend

A minimal backend service for a mobile app that lets users create groups, join groups, and trigger country-based notifications when they arrive or leave a country.

## 🚀 Features

- **User Authentication**: JWT-based authentication via Supabase Auth
- **Group Management**: Create, join, and list user groups
- **Location Updates**: Track country arrivals/departures with notifications
- **Push Notifications**: Integration-ready with Expo Push API
- **Health Monitoring**: Built-in health check endpoint
- **Graceful Shutdown**: Proper server lifecycle management

## 📁 Project Structure (Mono‑repo)

```
marko/
├── backend/                     # Backend (Go)
│   ├── cmd/
│   │   └── server/
│   │       └── main.go          # Main server entry point
│   ├── internal/
│   │   ├── auth/                # Authentication middleware
│   │   ├── config/              # Configuration management
│   │   ├── db/                  # Database connection and models
│   │   ├── groups/              # Group CRUD operations
│   │   ├── locations/           # Location update handling
│   │   └── notifications/       # Notification service and handlers
│   ├── migrations/              # Database migrations
│   ├── go.mod                   # Go module dependencies
│   └── go.sum                   # Dependency checksums
├── mobile/                      # Mobile app (Expo/React Native)
├── Dockerfile                   # Builds backend from backend/
├── docker-compose.yml           # Local dev services (e.g., Postgres)
└── .env.example                 # Environment variables template
```

## 🛠️ Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Docker and Docker Compose (for local development)
- Supabase account (for authentication)

## 🔧 Local Development Setup

### 1. Clone and Setup

```bash
# Clone the repository
git clone <your-repo-url>
cd country-arrival-notifier

# Copy environment template
cp .env.example .env

# Edit .env file with your configuration
nano .env
```

### 2. Start PostgreSQL

```bash
# Start PostgreSQL using Docker Compose
docker compose up -d postgres

# Verify PostgreSQL is running
docker compose ps
```

### 3. Install Backend Dependencies

```bash
cd backend
go mod download
```

### 4. Run Database Migrations

```bash
# Install goose (if not already installed)
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir backend/migrations postgres "postgres://user:password@localhost:5432/marko?sslmode=disable" up
```

### 5. Start the Backend Server

```bash
# Run the server (from repo root)
go run ./backend/cmd/server/main.go

# Or build and run (from repo root)
go build -o server ./backend/cmd/server/main.go
./server
```

The server will start on port 8080 (configurable via `PORT` environment variable).

## 🧪 API Endpoints

### Health Check
```
GET /healthz
```

### Authentication
All API endpoints require Bearer token authentication via Supabase Auth.

### Groups
```
POST   /api/v1/groups          # Create group
GET    /api/v1/groups          # List user groups
POST   /api/v1/groups/:id/join # Join group
GET    /api/v1/groups/:id/members # Get group members
```

### Locations
```
POST   /api/v1/location        # Update location (country arrival/departure)
```

Request body:
```json
{
  "countryCode": "US",
  "status": "arrived"  // or "left"
}
```

### Notifications
```
GET    /api/v1/notifications   # Get user notifications
```

Query parameters:
- `limit`: Number of notifications to return (default: 50)

## 🚀 Deployment

### Fly.io Deployment

1. **Install Fly CLI**:
```bash
curl -L https://fly.io/install.sh | sh
```

2. **Login to Fly.io**:
```bash
fly auth login
```

3. **Create Fly.io App**:
```bash
fly launch
```

4. **Set Environment Variables**:
```bash
fly secrets set DATABASE_URL="your-production-db-url"
fly secrets set SUPABASE_JWT_SECRET="your-supabase-jwt-secret"
fly secrets set SUPABASE_URL="https://your-project.supabase.co"
fly secrets set SUPABASE_KEY="your-supabase-anon-key"
fly secrets set EXPO_PUSH_TOKEN="your-expo-push-token"
fly secrets set ENVIRONMENT="production"
```

5. **Deploy**:
```bash
fly deploy
```

6. **Scale (optional)**:
```bash
fly scale count 2  # Run 2 instances
```

### Docker Deployment

1. **Build Docker Image**:
```bash
docker build -t country-arrival-notifier .
```

2. **Run Container**:
```bash
docker run -p 8080:8080 \
  -e DATABASE_URL="your-db-url" \
  -e SUPABASE_JWT_SECRET="your-jwt-secret" \
  -e SUPABASE_URL="your-supabase-url" \
  -e SUPABASE_KEY="your-supabase-key" \
  -e EXPO_PUSH_TOKEN="your-expo-token" \
  -e ENVIRONMENT="production" \
  country-arrival-notifier
```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ENVIRONMENT` | Environment (development/production) | `development` |
| `DATABASE_URL` | PostgreSQL connection string | Required |
| `SUPABASE_JWT_SECRET` | Supabase JWT secret for token verification | Required |
| `SUPABASE_URL` | Supabase project URL | Required |
| `SUPABASE_KEY` | Supabase anon key | Required |
| `EXPO_PUSH_TOKEN` | Expo push notification token | Optional |

### Database Schema

The application uses the following database schema:

- **users**: User profiles and push tokens
- **groups**: Group information
- **group_members**: User-group relationships
- **user_locations**: Location history
- **notifications**: Notification records

## 🔒 Security

- JWT-based authentication via Supabase Auth
- Environment variable configuration
- Non-root Docker container execution
- Graceful shutdown handling
- Structured logging with zerolog

## 📱 Mobile App Integration

This backend is designed to work with a mobile app using Expo. Key integration points:

1. **Authentication**: Use Supabase Auth SDK in your mobile app
2. **Push Notifications**: Integrate with Expo Push API (stub provided)
3. **API Calls**: Use the documented endpoints with Bearer token authentication

## 🧪 Testing

```bash
# Run health check
curl http://localhost:8080/healthz

# Test authenticated endpoints (requires valid JWT token)
curl -H "Authorization: Bearer YOUR_JWT_TOKEN" http://localhost:8080/api/v1/groups
```

## 📝 Development Notes

- The authentication middleware currently uses a stub for development
- Push notifications are stubbed and log instead of actually sending
- CORS is enabled for development environment
- Structured logging with zerolog provides detailed request/response logging

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License.

## 🆘 Support

For issues and questions:
1. Check the documentation
2. Review the logs for error messages
3. Open an issue on the repository

## 🔄 Future Enhancements

- Background job processing for notifications
- WebSocket support for real-time updates
- Rate limiting and API throttling
- Enhanced monitoring and metrics
- Background country detection
- Enhanced notification feeds