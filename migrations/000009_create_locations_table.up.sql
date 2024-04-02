
CREATE TABLE locations (
    id bigserial PRIMARY KEY,
    activity_id BIGINT NOT NULL REFERENCES activities ON DELETE CASCADE,
    name TEXT NOT NULL,
    address TEXT NOT NULL,
    google_place_id TEXT NOT NULL,
    lat decimal(9, 6) NOT NULL,
    lng decimal(9, 6) NOT NULL,
    website TEXT NOT NULL,
    phone TEXT NOT NULL,
    version integer NOT NULL DEFAULT 1,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);
