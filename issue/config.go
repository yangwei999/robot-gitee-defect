package issue

type Config struct {
	RobotToken      string   `json:"robot_token"      required:"true"`
	MaintainVersion []string `json:"maintain_version" required:"true"`
}
