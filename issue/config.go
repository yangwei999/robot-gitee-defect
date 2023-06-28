package issue

type Config struct {
	RobotToken      string   `json:"robot_token"      required:"true"`
	IssueType       string   `json:"issue_type"       required:"true"`
	MaintainVersion []string `json:"maintain_version" required:"true"`
}
