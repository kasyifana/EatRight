# EatRight Backend

Production-ready backend for EatRight - a food waste reduction platform connecting restaurants with users through mystery boxes and discounted food listings.

## Tech Stack

- **Language**: Go 1.22+
- **Framework**: Fiber v2 (Fast HTTP framework)
- **Database**: Supabase PostgreSQL with GORM
- **Authentication**: Supabase OAuth + JWT
- **Deployment**: VPS (Ubuntu) with systemd + Nginx

## Features

- 🔐 Google OAuth authentication via Supabase
- 🏪 Restaurant management with geolocation
- 📦 Food listing system (Mystery Box & Reveal items)
- 🛒 Order management with automatic stock control
- 🔒 Role-based access control
- 📍 Location-based restaurant queries
- 🚀 Production-ready deployment configuration

## Project Structure

```
.
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   └── app/
│       ├── config/                 # Configuration management
│       ├── models/                 # Database models
│       ├── repositories/           # Data access layer
│       ├── services/               # Business logic
│       ├── handlers/               # HTTP handlers
│       ├── middlewares/            # HTTP middlewares
│       └── utils/                  # Utility functions
├── migrations/                     # SQL migration scripts
├── scripts/                        # Build and deployment scripts
├── deploy/                         # Deployment configurations
└── go.mod                         # Go dependencies
```

## Quick Start

### Prerequisites

- Go 1.22 or higher
- PostgreSQL (Supabase account)
- Supabase project with Google OAuth enabled

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd EatRight
```

2. Install dependencies:
```bash
go mod download
```

3. Configure environment variables:
```bash
cp .env.example .env
# Edit .env with your actual credentials
```

4. Run database migrations:
```bash
# Execute migrations/001_create_tables.sql in your Supabase SQL editor
```

5. Run the application:
```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication
- `POST /api/auth/verify` - Verify Supabase token and get JWT

### Users
- `GET /api/users/me` - Get current user profile (protected)

### Restaurants
- `POST /api/restaurants` - Create restaurant (restaurant role only)
- `GET /api/restaurants` - List nearby restaurants (with lat/lng params)
- `GET /api/restaurants/:id` - Get restaurant details

### Listings
- `POST /api/restaurants/:id/listings` - Create food listing
- `GET /api/listings` - List all active listings
- `GET /api/listings/:id` - Get listing details
- `PATCH /api/listings/:id/stock` - Update stock
- `PATCH /api/listings/:id/status` - Toggle active status

### Orders
- `POST /api/orders` - Create order
- `GET /api/orders/me` - Get user's order history
- `PATCH /api/orders/:id/status` - Update order status (restaurant owner)

## Database Schema

### Users
- `id` (UUID, PK)
- `name` (string)
- `email` (string, unique)
- `role` (enum: 'user', 'restaurant')
- `created_at` (timestamp)

### Restaurants
- `id` (UUID, PK)
- `owner_id` (UUID, FK → users)
- `name` (string)
- `address` (string)
- `lat`, `lng` (float - for geolocation)
- `closing_time` (time)
- `created_at` (timestamp)

### Listings
- `id` (UUID, PK)
- `restaurant_id` (UUID, FK → restaurants)
- `type` (enum: 'mystery_box', 'reveal')
- `name` (string, nullable for mystery box)
- `description` (string)
- `price` (integer, in smallest currency unit)
- `stock` (integer)
- `photo_url` (string)
- `pickup_time` (time)
- `is_active` (boolean)
- `created_at` (timestamp)

### Orders
- `id` (UUID, PK)
- `user_id` (UUID, FK → users)
- `listing_id` (UUID, FK → listings)
- `qty` (integer)
- `total_price` (integer)
- `status` (enum: 'pending', 'ready', 'completed', 'cancelled')
- `created_at` (timestamp)

## Environment Variables

See `.env.example` for all required configuration. Key variables:

- `DATABASE_URL` - PostgreSQL connection string from Supabase
- `SUPABASE_URL` - Your Supabase project URL
- `SUPABASE_KEY` - Supabase anon/public key
- `JWT_SECRET` - Secret key for JWT signing
- `PORT` - Server port (default: 8080)

## Building for Production

```bash
# Build binary
./scripts/build.sh

# Or manually
go build -o bin/server cmd/server/main.go
```

## Deployment

See [DEPLOYMENT.md](DEPLOYMENT.md) for detailed VPS deployment instructions including:
- Systemd service configuration
- Nginx reverse proxy setup
- SSL/TLS configuration
- Environment management

## Development

### Running in development mode:
```bash
go run cmd/server/main.go
```

### Running with auto-reload (using air):
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with air
air
```

## Security Features

- JWT-based authentication with configurable expiry
- Role-based access control (user vs restaurant)
- Supabase OAuth token verification
- Input validation and sanitization
- CORS protection
- Transaction-based stock management to prevent race conditions

## Business Logic

### Stock Management
- Stock automatically decrements when orders are created
- Prevents negative stock through validation
- Transaction-based updates prevent overselling
- Orders fail if requested quantity exceeds available stock

### Order Flow
1. User creates order with listing_id and quantity
2. System validates stock availability
3. Stock is decremented atomically
4. Order created with 'pending' status
5. Restaurant updates status: pending → ready → completed

## License

Proprietary - All rights reserved

## Support

For issues and questions, please contact the development team.
