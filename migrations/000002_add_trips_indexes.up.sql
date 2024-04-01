CREATE INDEX IF NOT EXISTS trips_name_idx ON trips USING GIN (to_tsvector('simple', name));
