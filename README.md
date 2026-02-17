# Fiber Go Backend

## Setup

```bash
go mod download
go run main.go
```

## Project Structure

```
Backend(Go)/
├── main.go
├── config/
│   ├── config.go
│   └── firebase.go
├── models/
│   └── user.go
├── controllers/
│   └── user.go
├── routes/
│   └── routes.go
├── middleware/
│   └── auth.go
├── utils/
│   └── response.go
├── go.mod
└── .env.example
```

## API Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/users` - Get all users
- `GET /api/v1/users/:id` - Get user by ID
- `POST /api/v1/users` - Create user
- `PUT /api/v1/users/:id` - Update user (requires auth)
- `DELETE /api/v1/users/:id` - Delete user (requires auth)

## Environment Variables

Create `.env` file from `.env.example`
