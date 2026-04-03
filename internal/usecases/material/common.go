package material

import (
	"encoding/json"

	"doan/internal/entities"
)

type LabelView struct {
	ID          string `json:"id"`
	Code        string `json:"code"`
	Name        string `json:"name"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
}

type AuditLogView struct {
	ID              string     `json:"id"`
	Status          string     `json:"status"`
	Provider        string     `json:"provider"`
	RawOCRText      string     `json:"raw_ocr_text"`
	ConfidenceScore float64    `json:"confidence_score"`
	Reasoning       string     `json:"reasoning"`
	DetectedIssues  []string   `json:"detected_issues"`
	TriggeredAt     string     `json:"triggered_at"`
	CompletedAt     *string    `json:"completed_at,omitempty"`
	Label           *LabelView `json:"label,omitempty"`
}

type ApprovalDecisionView struct {
	ID                  string `json:"id"`
	ComplianceOfficerID string `json:"compliance_officer_id"`
	Approved            bool   `json:"approved"`
	RejectReason        string `json:"reject_reason"`
	Notes               string `json:"notes"`
	DecidedAt           string `json:"decided_at"`
}

type MaterialView struct {
	ID           string                `json:"id"`
	TeacherID    string                `json:"teacher_id"`
	Title        string                `json:"title"`
	Description  string                `json:"description"`
	FileName     string                `json:"file_name"`
	FileType     string                `json:"file_type"`
	FileSize     int64                 `json:"file_size"`
	Status       string                `json:"status"`
	UploadedAt   string                `json:"uploaded_at"`
	LatestLabel  *LabelView            `json:"latest_label,omitempty"`
	LatestAudit  *AuditLogView         `json:"latest_audit,omitempty"`
	LastDecision *ApprovalDecisionView `json:"last_decision,omitempty"`
	AuditLogs    []AuditLogView        `json:"audit_logs"`
}

func mapLabel(label *entities.Label) *LabelView {
	if label == nil || label.ID == "" {
		return nil
	}
	return &LabelView{
		ID:          label.ID,
		Code:        label.Code,
		Name:        label.Name,
		Severity:    label.Severity,
		Description: label.Description,
	}
}

func mapMaterial(material *entities.Material) MaterialView {
	result := MaterialView{
		ID:          material.ID,
		TeacherID:   material.TeacherID,
		Title:       material.Title,
		Description: material.Description,
		FileName:    material.FileName,
		FileType:    material.FileType,
		FileSize:    material.FileSize,
		Status:      material.Status,
		UploadedAt:  material.UploadedAt.Format("2006-01-02T15:04:05Z07:00"),
		LatestLabel: mapLabel(&material.LatestLabel),
		AuditLogs:   make([]AuditLogView, 0, len(material.AuditLogs)),
	}

	for _, auditLog := range material.AuditLogs {
		var issues []string
		_ = json.Unmarshal([]byte(auditLog.DetectedIssues), &issues)

		var completedAt *string
		if auditLog.CompletedAt != nil {
			formatted := auditLog.CompletedAt.Format("2006-01-02T15:04:05Z07:00")
			completedAt = &formatted
		}

		view := AuditLogView{
			ID:              auditLog.ID,
			Status:          auditLog.Status,
			Provider:        auditLog.Provider,
			RawOCRText:      auditLog.RawOCRText,
			ConfidenceScore: auditLog.ConfidenceScore,
			Reasoning:       auditLog.Reasoning,
			DetectedIssues:  issues,
			TriggeredAt:     auditLog.TriggeredAt.Format("2006-01-02T15:04:05Z07:00"),
			CompletedAt:     completedAt,
			Label:           mapLabel(&auditLog.Label),
		}
		result.AuditLogs = append(result.AuditLogs, view)
		if result.LatestAudit == nil {
			result.LatestAudit = &view
		}
	}

	if len(material.ApprovalDecisions) > 0 {
		decision := material.ApprovalDecisions[0]
		result.LastDecision = &ApprovalDecisionView{
			ID:                  decision.ID,
			ComplianceOfficerID: decision.ComplianceOfficerID,
			Approved:            decision.Approved,
			RejectReason:        decision.RejectReason,
			Notes:               decision.Notes,
			DecidedAt:           decision.DecidedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result
}
