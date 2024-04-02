BEGIN;

-- Remove the foreign key constraint from the trips table
-- Note: This step may be unnecessary in some RDBMS like PostgreSQL, as dropping the column
-- will automatically remove the constraint. However, explicitly dropping the constraint
-- can be done for clarity and to ensure compatibility across different systems.
ALTER TABLE trips
DROP CONSTRAINT IF EXISTS fk_trips_created_by;

-- Remove the created_by column from the trips table
ALTER TABLE trips
DROP COLUMN IF EXISTS created_by;

COMMIT;
