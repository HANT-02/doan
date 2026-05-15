BEGIN;

INSERT INTO campuses (code, name, address, created_at, updated_at)
VALUES
    ('XLB-CHUNG-CS1', 'Cơ sở trung tâm', 'Số 1 đường Trung tâm', NOW(), NOW()),
    ('XLB-CHUNG-CS2', 'Cơ sở vệ tinh', 'Số 2 đường Vệ tinh', NOW(), NOW()),
    ('XLB-CHUNG-CS3', 'Cơ sở mở rộng', 'Số 3 đường Mở rộng', NOW(), NOW())
ON CONFLICT (code) DO UPDATE
SET
    name = EXCLUDED.name,
    address = EXCLUDED.address,
    updated_at = NOW();

INSERT INTO campus_travel_times (from_campus_id, to_campus_id, travel_minutes, is_active, created_at, updated_at)
SELECT
    from_campus.id,
    to_campus.id,
    pair.travel_minutes,
    TRUE,
    NOW(),
    NOW()
FROM (
    VALUES
        ('XLB-CHUNG-CS1', 'XLB-CHUNG-CS2', 30),
        ('XLB-CHUNG-CS2', 'XLB-CHUNG-CS1', 30),
        ('XLB-CHUNG-CS1', 'XLB-CHUNG-CS3', 45),
        ('XLB-CHUNG-CS3', 'XLB-CHUNG-CS1', 45),
        ('XLB-CHUNG-CS2', 'XLB-CHUNG-CS3', 20),
        ('XLB-CHUNG-CS3', 'XLB-CHUNG-CS2', 20)
) AS pair(from_code, to_code, travel_minutes)
JOIN campuses AS from_campus ON from_campus.code = pair.from_code
JOIN campuses AS to_campus ON to_campus.code = pair.to_code
ON CONFLICT (from_campus_id, to_campus_id) DO UPDATE
SET
    travel_minutes = EXCLUDED.travel_minutes,
    is_active = TRUE,
    updated_at = NOW();

INSERT INTO shifts (
    code,
    name,
    start_time,
    end_time,
    duration_minutes,
    session_type,
    is_active,
    notes,
    created_at,
    updated_at
)
VALUES
    ('XLB-CHUNG-S1', 'Ca sáng 1', '08:00', '10:00', 120, 'MORNING', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW()),
    ('XLB-CHUNG-S2', 'Ca sáng 2', '10:15', '12:15', 120, 'MORNING', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW()),
    ('XLB-CHUNG-S3', 'Ca chiều 1', '13:30', '15:30', 120, 'AFTERNOON', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW()),
    ('XLB-CHUNG-S4', 'Ca chiều 2', '15:45', '17:45', 120, 'AFTERNOON', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW()),
    ('XLB-CHUNG-S5', 'Ca tối 1', '18:00', '20:00', 120, 'EVENING', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW()),
    ('XLB-CHUNG-S6', 'Ca tối 2', '20:15', '22:15', 120, 'EVENING', TRUE, 'Bộ ca chuẩn benchmark', NOW(), NOW())
ON CONFLICT (code) DO UPDATE
SET
    name = EXCLUDED.name,
    start_time = EXCLUDED.start_time,
    end_time = EXCLUDED.end_time,
    duration_minutes = EXCLUDED.duration_minutes,
    session_type = EXCLUDED.session_type,
    is_active = EXCLUDED.is_active,
    notes = EXCLUDED.notes,
    updated_at = NOW();

COMMIT;
