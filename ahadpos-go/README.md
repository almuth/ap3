# AhadPOS 3 Go API

RESTful API for AhadPOS 3 Point of Sale system, built with Go and Gin framework.

## Features

- 🔐 JWT Authentication
- 📦 Product Management
- 📊 Inventory Tracking
- 💰 Sales & Purchase Management
- 🏦 Financial Tracking
- 🌐 RESTful API Design
- 📝 SQLite Database
- ⚡ High Performance

## Project Structure

```
ahadpos-go/
├── cmd/
│   └── api/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go       # Configuration management
│   ├── handlers/
│   │   ├── auth.go         # Authentication handlers
│   │   └── barang.go       # Product handlers
│   ├── middleware/
│   │   ├── auth.go         # JWT authentication middleware
│   │   └── cors.go        # CORS middleware
│   └── models/
│       └── models.go       # Data models
├── pkg/
│   ├── database/
│   │   └── database.go     # Database connection
│   └── utils/
│       └── utils.go        # Utility functions
├── configs/
│   └── config.yaml         # Configuration file
├── logs/                  # Application logs
├── Makefile              # Build automation
├── go.mod                # Go module definition
└── README.md             # This file
```

## Prerequisites

- Go 1.21 or higher
- SQLite 3

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd ahadpos-go
```

2. Install dependencies:
```bash
make deps
```

3. Configure the application:
Edit `configs/config.yaml` to match your environment.

4. Build the application:
```bash
make build
```

## Configuration

Edit `configs/config.yaml` to configure the application:

```yaml
server:
  port: 8080
  mode: debug  # debug, release, test

database:
  path: ../ahadpos3.db
  max_open_conns: 100
  max_idle_conns: 10

jwt:
  secret: your-secret-key-change-in-production
  expiration: 24  # hours
```

## Running

### Development Mode
```bash
make run
```

### Production Build
```bash
make build
./build/ahadpos-api
```

## API Documentation

### Base URL
```
http://localhost:8080/api/v1
```

### Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your-jwt-token>
```

### Endpoints

#### Health Check
```
GET /health
```

#### Authentication
```
POST   /api/v1/auth/login
GET    /api/v1/auth/me
```

#### Products
```
GET    /api/v1/barang          # List products (with pagination)
GET    /api/v1/barang/:id     # Get product by ID
POST   /api/v1/barang          # Create product
PUT    /api/v1/barang/:id     # Update product
DELETE /api/v1/barang/:id     # Delete product
```

#### Query Parameters

**List Products:**
- `page`: Page number (default: 1)
- `page_size`: Items per page (default: 20, max: 100)
- `kategori_id`: Filter by category ID
- `status`: Filter by status (0 or 1)
- `search`: Search by name or barcode

### Example Usage

#### Login
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "password": "admin123"
  }'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user_id": 1,
  "username": "admin",
  "role": "admin",
  "name": "Administrator"
}
```

#### Get Products
```bash
curl -X GET http://localhost:8080/api/v1/barang \
  -H "Authorization: Bearer <your-token>"
```

#### Create Product
```bash
curl -X POST http://localhost:8080/api/v1/barang \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "barcode": "8991234567890",
    "nama": "Sample Product",
    "kategori_id": 1,
    "satuan_id": 1,
    "rak_id": 1,
    "status": 1
  }'
```

## Default Credentials

- **Username**: admin
- **Password**: admin123

⚠️ **Important**: Change the default password and JWT secret in production!

## Makefile Commands

```bash
make deps     # Download dependencies
make build    # Build the application
make run      # Build and run
make test     # Run tests
make clean    # Clean build artifacts
make help     # Show help
```

## Development

### Adding New Endpoints

1. Create handler in `internal/handlers/`
2. Define models in `internal/models/`
3. Register routes in `cmd/api/main.go`
4. Add tests if needed

### Database Migrations

To modify the database schema:
1. Create migration SQL in `../migrations/`
2. Run migration manually or implement auto-migration
3. Update models in `internal/models/models.go`

## Testing

Run tests:
```bash
make test
```

## Deployment

### Build for Production
```bash
GOOS=linux GOARCH=amd64 make build
```

### Run as Service
Create a systemd service file:

```ini
[Unit]
Description=AhadPOS API
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/path/to/ahadpos-go
ExecStart=/path/to/ahadpos-go/build/ahadpos-api
Restart=always

[Install]
WantedBy=multi-user.target
```

Enable and start:
```bash
sudo systemctl enable ahadpos-api
sudo systemctl start ahadpos-api
```

## Security

- Change the JWT secret in production
- Use strong passwords
- Enable HTTPS in production
- Implement rate limiting
- Use environment variables for sensitive data
- Keep dependencies updated

## Performance

- SQLite database is embedded and fast
- Connection pooling configured
- Indexed queries for common operations
- Pagination for large datasets

## Troubleshooting

### Database Connection Error
- Check database path in config.yaml
- Ensure database file exists and is readable
- Check file permissions

### Port Already in Use
- Change port in config.yaml
- Kill process using the port

### JWT Token Invalid
- Check JWT secret matches
- Ensure token is not expired
- Verify Authorization header format

## License

This project is part of AhadPOS 3. See main project for license details.

## Support

For issues and questions, please refer to the main AhadPOS project documentation.

## Roadmap

- [ ] Complete all CRUD operations
- [ ] Add sales management endpoints
- [ ] Add purchase management endpoints
- [ ] Add inventory management
- [ ] Add reporting endpoints
- [ ] Add file upload for images
- [ ] Add WebSocket support for real-time updates
- [ ] Add API rate limiting
- [ ] Add caching layer
- [ ] Add comprehensive tests
- [ ] Add API documentation (Swagger)
- [ ] Docker support

## Contributing

Contributions are welcome! Please follow the project structure and coding standards.
