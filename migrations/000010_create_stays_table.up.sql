CREATE TABLE stays (
    id bigserial PRIMARY KEY,
    trip_id BIGINT NOT NULL REFERENCES trips ON DELETE CASCADE,
    name TEXT NOT NULL,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    address TEXT NOT NULL,
    lat decimal(9, 6) NOT NULL,
    lng decimal(9, 6) NOT NULL,
    phone TEXT NOT NULL,
    link TEXT NOT NULL,
    type TEXT NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
