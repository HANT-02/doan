CREATE TABLE IF NOT EXISTS labels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    severity VARCHAR(50) NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

INSERT INTO labels (code, name, severity, description)
VALUES
    ('SAFE', 'An toàn', 'SAFE', 'Nội dung phù hợp để phát hành'),
    ('WARNING', 'Cảnh báo', 'WARNING', 'Nội dung cần compliance officer rà soát'),
    ('DANGER', 'Nguy hiểm', 'DANGER', 'Nội dung rủi ro cao cần chặn/phê duyệt')
ON CONFLICT (code) DO NOTHING;

CREATE TABLE IF NOT EXISTS materials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    teacher_id UUID REFERENCES teachers(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    file_name VARCHAR(255) NOT NULL,
    file_path TEXT NOT NULL,
    file_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'UPLOADED',
    latest_label_id UUID REFERENCES labels(id) ON DELETE SET NULL,
    uploaded_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_materials_teacher_id ON materials(teacher_id);
CREATE INDEX IF NOT EXISTS idx_materials_status ON materials(status);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    material_id UUID REFERENCES materials(id) ON DELETE CASCADE,
    label_id UUID REFERENCES labels(id) ON DELETE SET NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'COMPLETED',
    provider VARCHAR(100) NOT NULL,
    raw_ocr_text TEXT,
    confidence_score NUMERIC(5,4) DEFAULT 0,
    reasoning TEXT,
    detected_issues JSONB NOT NULL DEFAULT '[]'::jsonb,
    triggered_at TIMESTAMPTZ DEFAULT NOW(),
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_material_id ON audit_logs(material_id);

CREATE TABLE IF NOT EXISTS approval_decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    material_id UUID REFERENCES materials(id) ON DELETE CASCADE,
    audit_log_id UUID REFERENCES audit_logs(id) ON DELETE SET NULL,
    compliance_officer_id UUID REFERENCES users(id) ON DELETE SET NULL,
    approved BOOLEAN NOT NULL,
    reject_reason TEXT,
    notes TEXT,
    decided_at TIMESTAMPTZ DEFAULT NOW(),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_approval_decisions_material_id ON approval_decisions(material_id);
