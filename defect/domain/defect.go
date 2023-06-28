package domain

type Defect struct {
	IssueNumber     string   `json:"issue_number"`
	IssueOrg        string   `json:"issue_org"`
	IssueRepo       string   `json:"issue_repo"`
	IssueStatus     string   `json:"issue_status"`
	Kernel          string   `json:"kernel"`
	Component       string   `json:"component"`
	SystemVersion   string   `json:"system_version"`
	Description     string   `json:"description"`
	ReferenceURL    string   `json:"reference_url"`
	GuidanceURL     string   `json:"guidance_url"`
	Influence       string   `json:"influence"`
	SeverityLevel   string   `json:"severity_level"`
	AffectedVersion []string `json:"affected_version"`
	ABI             string   `json:"abi"`
}
