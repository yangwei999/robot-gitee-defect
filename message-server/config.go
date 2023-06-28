package messageserver

type Config struct {
	UserAgent string `json:"user_agent"    required:"true"`
	GroupName string `json:"group_name"    required:"true"`
	Topics    Topics `json:"topics"        required:"true" `
}

type Topics struct {
	DefectEvent string `json:"defect_event" required:"true"`
}
