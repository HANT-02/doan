package material

type ReviewMaterialRequest struct {
	ComplianceOfficerID string `json:"compliance_officer_id" binding:"required"`
	Approved            bool   `json:"approved"`
	RejectReason        string `json:"reject_reason"`
	Notes               string `json:"notes"`
}
