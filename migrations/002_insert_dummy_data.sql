-- EatRight Dummy Data
-- Run this in Supabase SQL Editor after running migrations

-- ============================================
-- 1. INSERT DUMMY USERS (5 users)
-- ============================================

INSERT INTO users (id, name, email, role, created_at) VALUES
  ('550e8400-e29b-41d4-a716-446655440001', 'John Doe', 'john.doe@example.com', 'user', NOW()),
  ('550e8400-e29b-41d4-a716-446655440002', 'Jane Smith', 'jane.smith@example.com', 'user', NOW()),
  ('550e8400-e29b-41d4-a716-446655440003', 'Alice Johnson', 'alice.johnson@example.com', 'user', NOW()),
  ('550e8400-e29b-41d4-a716-446655440010', 'Warung Makan Bahari', 'warung.bahari@example.com', 'restaurant', NOW()),
  ('550e8400-e29b-41d4-a716-446655440011', 'Cafe Lestari', 'cafe.lestari@example.com', 'restaurant', NOW())
ON CONFLICT (email) DO NOTHING;

-- ============================================
-- 2. INSERT DUMMY RESTAURANTS (5 restaurants)
-- ============================================

INSERT INTO restaurants (id, owner_id, name, address, lat, lng, closing_time, created_at) VALUES
  ('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440010', 'Warung Makan Bahari', 'Jl. Sudirman No. 123, Jakarta', -6.208763, 106.845599, CURRENT_DATE + TIME '22:00:00', NOW()),
  ('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440011', 'Cafe Lestari', 'Jl. Gatot Subroto No. 45, Jakarta', -6.225014, 106.808331, CURRENT_DATE + TIME '23:00:00', NOW()),
  ('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440010', 'Rumah Makan Nusantara', 'Jl. Thamrin No. 67, Jakarta', -6.195003, 106.822502, CURRENT_DATE + TIME '21:30:00', NOW()),
  ('650e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440011', 'Bakery Corner', 'Jl. MH Thamrin No. 89, Jakarta', -6.186486, 106.824622, CURRENT_DATE + TIME '20:00:00', NOW()),
  ('650e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440010', 'Sushi Express', 'Jl. Rasuna Said No. 12, Jakarta', -6.223574, 106.830414, CURRENT_DATE + TIME '22:30:00', NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 3. INSERT DUMMY LISTINGS (5 listings)
-- ============================================

INSERT INTO listings (id, restaurant_id, type, name, description, price, stock, photo_url, pickup_time, is_active, created_at) VALUES
  -- Mystery Boxes
  ('750e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', 'mystery_box', NULL, 'Paket misteri berisi 3-4 menu pilihan chef. Cocok untuk makan malam!', 50000, 10, 'https://images.unsplash.com/photo-1546069901-ba9599a7e63c', CURRENT_DATE + TIME '20:00:00', TRUE, NOW()),
  ('750e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440002', 'mystery_box', NULL, 'Mystery box cafe - minuman + snack + dessert', 35000, 8, 'https://images.unsplash.com/photo-1495474472287-4d71bcdd2085', CURRENT_DATE + TIME '21:00:00', TRUE, NOW()),
  
  -- Reveal Items
  ('750e8400-e29b-41d4-a716-446655440003', '650e8400-e29b-41d4-a716-446655440003', 'reveal', 'Nasi Goreng Spesial + Es Teh', 'Nasi goreng dengan telur mata sapi, kerupuk, dan es teh manis', 25000, 15, 'https://images.unsplash.com/photo-1512058564366-18510be2db19', CURRENT_DATE + TIME '19:30:00', TRUE, NOW()),
  ('750e8400-e29b-41d4-a716-446655440004', '650e8400-e29b-41d4-a716-446655440004', 'reveal', 'Roti Sisa Hari Ini', 'Paket 5 roti campuran (croissant, donut, roti tawar)', 30000, 12, 'https://images.unsplash.com/photo-1509440159596-0249088772ff', CURRENT_DATE + TIME '19:00:00', TRUE, NOW()),
  ('750e8400-e29b-41d4-a716-446655440005', '650e8400-e29b-41d4-a716-446655440005', 'reveal', 'Sushi Combo Pack', '12 potong sushi campuran + miso soup', 75000, 5, 'https://images.unsplash.com/photo-1579584425555-c3ce17fd4351', CURRENT_DATE + TIME '21:30:00', TRUE, NOW())
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- 4. INSERT DUMMY ORDERS (5 orders)
-- ============================================

INSERT INTO orders (id, user_id, listing_id, qty, total_price, status, created_at) VALUES
  ('850e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 2, 100000, 'ready', NOW() - INTERVAL '2 hours'),
  ('850e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440003', 1, 25000, 'completed', NOW() - INTERVAL '1 day'),
  ('850e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440004', 3, 90000, 'pending', NOW() - INTERVAL '30 minutes'),
  ('850e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440002', 1, 35000, 'ready', NOW() - INTERVAL '3 hours'),
  ('850e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440005', 1, 75000, 'pending', NOW() - INTERVAL '15 minutes')
ON CONFLICT (id) DO NOTHING;

-- ============================================
-- VERIFICATION QUERIES
-- ============================================

-- Check inserted data
SELECT 'Users' as table_name, COUNT(*) as count FROM users
UNION ALL
SELECT 'Restaurants', COUNT(*) FROM restaurants
UNION ALL
SELECT 'Listings', COUNT(*) FROM listings
UNION ALL
SELECT 'Orders', COUNT(*) FROM orders;

-- Show summary
DO $$
BEGIN
    RAISE NOTICE '✅ Dummy data inserted successfully!';
    RAISE NOTICE '📊 Data summary:';
    RAISE NOTICE '   - 5 users (3 regular users + 2 restaurant owners)';
    RAISE NOTICE '   - 5 restaurants';
    RAISE NOTICE '   - 5 listings (2 mystery boxes + 3 reveal items)';
    RAISE NOTICE '   - 5 orders (various statuses)';
END $$;
