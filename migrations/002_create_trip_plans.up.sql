-- Trip Plans
CREATE TABLE IF NOT EXISTS trip_plans (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    origin VARCHAR(255) NOT NULL,
    destination VARCHAR(255) NOT NULL,
    depart_at TIMESTAMP NULL,
    status VARCHAR(50) DEFAULT 'planned' NOT NULL,
    selected_itinerary_id INTEGER NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Itineraries
CREATE TABLE IF NOT EXISTS itineraries (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES trip_plans(id) ON DELETE CASCADE,
    mode_mix VARCHAR(100) NOT NULL,
    total_minutes INTEGER NOT NULL,
    rough_cost_cents INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Legs
CREATE TABLE IF NOT EXISTS legs (
    id SERIAL PRIMARY KEY,
    itinerary_id INTEGER NOT NULL REFERENCES itineraries(id) ON DELETE CASCADE,
    index INTEGER NOT NULL,
    mode VARCHAR(20) NOT NULL,
    from_name VARCHAR(255) NOT NULL,
    to_name VARCHAR(255) NOT NULL,
    minutes INTEGER NOT NULL,
    distance_m BIGINT NOT NULL,
    provider VARCHAR(100) NULL
);

-- Ride Bookings
CREATE TABLE IF NOT EXISTS ride_bookings (
    id SERIAL PRIMARY KEY,
    plan_id INTEGER NOT NULL REFERENCES trip_plans(id) ON DELETE CASCADE,
    itinerary_id INTEGER NOT NULL REFERENCES itineraries(id) ON DELETE CASCADE,
    provider VARCHAR(100) NOT NULL,
    status VARCHAR(50) DEFAULT 'pending' NOT NULL,
    eta_minutes INTEGER NOT NULL,
    fare_cents INTEGER NOT NULL,
    payment_id INTEGER NULL,
    quote_id VARCHAR(255) NULL,
    external_ref VARCHAR(255) NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Payments
CREATE TABLE IF NOT EXISTS payments (
    id SERIAL PRIMARY KEY,
    amount_cents INTEGER NOT NULL,
    currency VARCHAR(10) DEFAULT 'THB' NOT NULL,
    status VARCHAR(50) DEFAULT 'authorized' NOT NULL,
    external_ref VARCHAR(255) DEFAULT 'stub' NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_trip_plans_user_id ON trip_plans(user_id);
CREATE INDEX IF NOT EXISTS idx_trip_plans_status ON trip_plans(status);
CREATE INDEX IF NOT EXISTS idx_itineraries_plan_id ON itineraries(plan_id);
CREATE INDEX IF NOT EXISTS idx_legs_itinerary_id ON legs(itinerary_id);
CREATE INDEX IF NOT EXISTS idx_ride_bookings_plan_id ON ride_bookings(plan_id);
CREATE INDEX IF NOT EXISTS idx_ride_bookings_itinerary_id ON ride_bookings(itinerary_id);
CREATE INDEX IF NOT EXISTS idx_ride_bookings_quote_id ON ride_bookings(quote_id);
CREATE INDEX IF NOT EXISTS idx_payments_status ON payments(status);