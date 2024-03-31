CREATE TABLE IF NOT EXISTS trips (
    id bigserial PRIMARY KEY,  
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    city text NOT NULL,
    state_code text NOT NULL,
    google_place_id text NOT NULL,
    lat decimal(9, 6) NOT NULL,
    lng decimal(9, 6) NOT NULL,
    start_date timestamp(0) with time zone NOT NULL, 
    end_date timestamp(0) with time zone NOT NULL,
    version integer NOT NULL DEFAULT 1
);

