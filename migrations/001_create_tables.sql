-- EatRight Database Schema
-- Run this script in your Supabase SQL Editor

-- Create custom types for enums
DO $$ BEGIN
    CREATE TYPE user_role AS ENUM ('user', 'restaurant');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE listing_type AS ENUM ('mystery_box', 'reveal');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

DO $$ BEGIN
    CREATE TYPE order_status AS ENUM ('pending', 'ready', 'completed', 'cancelled');
EXCEPTION
    WHEN duplicate_object THEN null;
END $$;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    role user_role NOT NULL DEFAULT 'user',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Restaurants table
CREATE TABLE IF NOT EXISTS restaurants (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    lat DECIMAL(10, 8) NOT NULL,
    lng DECIMAL(11, 8) NOT NULL,
    closing_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for restaurants
CREATE INDEX IF NOT EXISTS idx_restaurants_owner_id ON restaurants(owner_id);
CREATE INDEX IF NOT EXISTS idx_restaurants_location ON restaurants(lat, lng);

-- Listings table
CREATE TABLE IF NOT EXISTS listings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    restaurant_id UUID NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    type listing_type NOT NULL,
    name VARCHAR(255),
    description TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    stock INTEGER NOT NULL DEFAULT 0 CHECK (stock >= 0),
    photo_url TEXT,
    pickup_time TIME NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for listings
CREATE INDEX IF NOT EXISTS idx_listings_restaurant_id ON listings(restaurant_id);
CREATE INDEX IF NOT EXISTS idx_listings_is_active ON listings(is_active);
CREATE INDEX IF NOT EXISTS idx_listings_type ON listings(type);

-- Orders table
CREATE TABLE IF NOT EXISTS orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    listing_id UUID NOT NULL REFERENCES listings(id) ON DELETE CASCADE,
    qty INTEGER NOT NULL CHECK (qty > 0),
    total_price INTEGER NOT NULL CHECK (total_price >= 0),
    status order_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for orders
CREATE INDEX IF NOT EXISTS idx_orders_user_id ON orders(user_id);
CREATE INDEX IF NOT EXISTS idx_orders_listing_id ON orders(listing_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);

-- Comments for documentation
COMMENT ON TABLE users IS 'Stores user information including role (user or restaurant owner)';
COMMENT ON TABLE restaurants IS 'Stores restaurant information with geolocation data';
COMMENT ON TABLE listings IS 'Stores food listings (mystery boxes and reveal items)';
COMMENT ON TABLE orders IS 'Stores customer orders with status tracking';

COMMENT ON COLUMN listings.type IS 'Type of listing: mystery_box or reveal';
COMMENT ON COLUMN listings.name IS 'Name of the item (nullable for mystery boxes)';
COMMENT ON COLUMN listings.price IS 'Price in smallest currency unit (e.g., cents)';
COMMENT ON COLUMN orders.total_price IS 'Total price in smallest currency unit (e.g., cents)';

-- Grant permissions (adjust based on your Supabase setup)
-- These are handled by Supabase RLS policies, but included for reference
-- ALTER TABLE users ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE restaurants ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE listings ENABLE ROW LEVEL SECURITY;
-- ALTER TABLE orders ENABLE ROW LEVEL SECURITY;

-- Success message
DO $$
BEGIN
    RAISE NOTICE '✅ EatRight database schema created successfully!';
    RAISE NOTICE '📊 Tables created: users, restaurants, listings, orders';
    RAISE NOTICE '🔍 Indexes created for optimal query performance';
END $$;
