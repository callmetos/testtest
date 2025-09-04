DROP INDEX IF EXISTS idx_payments_status;
DROP INDEX IF EXISTS idx_ride_bookings_quote_id;
DROP INDEX IF EXISTS idx_ride_bookings_itinerary_id;
DROP INDEX IF EXISTS idx_ride_bookings_plan_id;
DROP INDEX IF EXISTS idx_legs_itinerary_id;
DROP INDEX IF EXISTS idx_itineraries_plan_id;
DROP INDEX IF EXISTS idx_trip_plans_status;
DROP INDEX IF EXISTS idx_trip_plans_user_id;

DROP TABLE IF EXISTS payments;
DROP TABLE IF EXISTS ride_bookings;
DROP TABLE IF EXISTS legs;
DROP TABLE IF EXISTS itineraries;
DROP TABLE IF EXISTS trip_plans;