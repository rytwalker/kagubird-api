BEGIN;

-- Create the join table for the many-to-many relationship
CREATE TABLE trip_goers (
    trip_id bigint NOT NULL REFERENCES trips ON DELETE CASCADE,
    user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
    PRIMARY KEY (trip_id, user_id)
);

COMMIT;
