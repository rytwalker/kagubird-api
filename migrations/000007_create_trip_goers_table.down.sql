BEGIN;

-- Drop the trip_goers table to remove the many-to-many relationship
DROP TABLE IF EXISTS trip_goers;

COMMIT;
