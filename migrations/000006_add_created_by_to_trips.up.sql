BEGIN;

-- Add the created_by column with a default value of 0
ALTER TABLE trips
ADD COLUMN created_by BIGINT DEFAULT 0;

-- Update existing rows to set created_by to 0 where it is currently NULL
-- This step ensures that the NOT NULL constraint can be applied without issues
UPDATE trips
SET created_by = 0
WHERE created_by IS NULL;

-- Add a foreign key constraint to ensure created_by references a valid user_id in the users table
ALTER TABLE trips
ADD CONSTRAINT fk_trips_created_by
FOREIGN KEY (created_by)
REFERENCES users (id)
ON DELETE SET NULL
ON UPDATE CASCADE;

-- Alter the created_by column to set it as NOT NULL
ALTER TABLE trips
ALTER COLUMN created_by SET NOT NULL;

COMMIT;
