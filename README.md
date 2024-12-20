# Casion Payment System

A payment system API with features like top-up, payment, transfer, and transaction history.

## Features

- User Registration and Authentication
- Profile Management
- Balance Top-up
- Payment Processing
- Money Transfer between Users
- Transaction History
- Dashboard Statistics

## Tech Stack

- Go 1.22
- Gin Web Framework
- GORM with MySQL
- JWT Authentication
- OpenAPI/Swagger Documentation

## Prerequisites

- Go 1.22 or higher
- MySQL 8.0 or higher
- Make sure MySQL service is running

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/casion.git
cd casion
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables in `.env`:
```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=casion_user
DB_PASS=casion_password
DB_NAME=casion_db

JWT_SECRET=your-256-bit-secret-key-here-make-it-long-and-secure

SERVER_PORT=8080
```

4. Create MySQL database and user:
```sql
CREATE DATABASE casion_db;
CREATE USER 'casion_user'@'localhost' IDENTIFIED BY 'casion_password';
GRANT ALL PRIVILEGES ON casion_db.* TO 'casion_user'@'localhost';
FLUSH PRIVILEGES;
```

5. Run the application:
```bash
go run cmd/main.go
```

## API Documentation

The API documentation is available in two formats:

1. OpenAPI Specification: `http://localhost:8080/openapi.yaml`
2. Swagger UI: `http://localhost:8080/swagger-ui/index.html`

### Main Endpoints

#### Authentication
- `POST /api/register` - Register new user
- `POST /api/login` - Login and get access token

#### Profile
- `PUT /api/profile` - Update user profile

#### Transactions
- `POST /api/topup` - Top up balance
- `POST /api/payment` - Make payment
- `POST /api/transfer` - Transfer money to another user
- `GET /api/transactions` - Get transaction history

#### Dashboard
- `GET /api/dashboard/stats` - Get dashboard statistics
- `GET /api/dashboard/transfers/recent` - Get recent transfers
- `GET /api/dashboard/transfers/failed` - Get failed transfers

## Authentication

The API uses JWT for authentication. Include the token in the Authorization header:
```
Authorization: Bearer <your_access_token>
```

## Example Requests

### Register User
```bash
curl -X 'POST' \
  'http://localhost:8080/api/register' \
  -H 'Content-Type: application/json' \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "phone_number": "0811111111",
    "address": "123 Main St",
    "pin": "123456"
  }'
```

### Login
```bash
curl -X 'POST' \
  'http://localhost:8080/api/login' \
  -H 'Content-Type: application/json' \
  -d '{
    "phone_number": "0811111111",
    "pin": "123456"
  }'
```

### Top Up
```bash
curl -X 'POST' \
  'http://localhost:8080/api/topup' \
  -H 'Authorization: Bearer <your_access_token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "amount": 1000000
  }'
```

### Transfer
```bash
curl -X 'POST' \
  'http://localhost:8080/api/transfer' \
  -H 'Authorization: Bearer <your_access_token>' \
  -H 'Content-Type: application/json' \
  -d '{
    "target_user": "0822222222",
    "amount": 500000,
    "remarks": "Monthly rent payment"
  }'
```

## Error Handling

The API returns consistent error responses in the following format:
```json
{
  "status": "error",
  "message": "Error message here"
}
```

Common HTTP status codes:
- 200: Success
- 400: Bad Request
- 401: Unauthorized
- 404: Not Found
- 500: Internal Server Error

## Development

### Project Structure
```
.
├── api/
│   ├── openapi.yaml        # OpenAPI specification
│   └── swagger-ui/         # Swagger UI files
├── cmd/
│   └── main.go            # Application entry point
├── internal/
│   ├── config/            # Configuration
│   ├── handlers/          # HTTP handlers
│   ├── middleware/        # Middleware functions
│   ├── models/            # Data models
│   ├── utils/             # Utility functions
│   └── worker/            # Background workers
├── .env                   # Environment variables
├── go.mod                 # Go modules file
└── README.md             # This file
```

