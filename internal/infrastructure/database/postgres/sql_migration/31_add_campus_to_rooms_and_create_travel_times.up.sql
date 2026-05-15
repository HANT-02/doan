ALTER TABLE rooms
    ADD COLUMN IF NOT EXISTS campus_id UUID REFERENCES campuses(id);

UPDATE rooms
SET campus_id = (
    SELECT id
    FROM campuses
    WHERE code = 'MAIN'
    LIMIT 1
)
WHERE campus_id IS NULL;

CREATE INDEX IF NOT EXISTS idx_rooms_campus_id ON rooms(campus_id);

CREATE TABLE IF NOT EXISTS campus_travel_times (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    from_campus_id UUID NOT NULL REFERENCES campuses(id) ON DELETE CASCADE,
    to_campus_id UUID NOT NULL REFERENCES campuses(id) ON DELETE CASCADE,
    travel_minutes INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_at TIMESTAMP WITH TIME ZONE,
    CONSTRAINT uq_campus_travel_pair UNIQUE (from_campus_id, to_campus_id)
);

CREATE INDEX IF NOT EXISTS idx_campus_travel_times_deleted_at ON campus_travel_times(deleted_at);
CREATE INDEX IF NOT EXISTS idx_campus_travel_times_from_to ON campus_travel_times(from_campus_id, to_campus_id);
