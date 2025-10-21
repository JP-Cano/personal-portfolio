package dto

// UploadCertificatesRequest represents the query parameters for file upload
type UploadCertificatesRequest struct {
	Workers int `form:"workers" validate:"omitempty,min=0,max=20"`
}

// CertificationMetadata represents optional metadata for each certification file
type CertificationMetadata struct {
	Title         string `form:"title" validate:"omitempty,max=255"`
	Issuer        string `form:"issuer" validate:"omitempty,max=255"`
	IssueDate     string `form:"issue_date" validate:"omitempty,date_format"`
	ExpiryDate    string `form:"expiry_date" validate:"omitempty,date_format"`
	CredentialID  string `form:"credential_id" validate:"omitempty,max=255"`
	CredentialURL string `form:"credential_url" validate:"omitempty,url,max=500"`
	Description   string `form:"description" validate:"omitempty"`
}
