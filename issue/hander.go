package issue

import (
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"

	"github.com/opensourceways/robot-gitee-defect/defect/app"
)

type EventHandler interface {
	HandleIssueEvent(e *sdk.IssueEvent) error
}

type iClient interface {
	CreateIssueComment(org, repo string, number string, comment string) error
}

func NewEventHandler(c *Config, s app.DefectService) *eventHandler {
	cli := client.NewClient(func() []byte {
		return []byte(c.RobotToken)
	})
	return &eventHandler{
		cfg:     c,
		cli:     cli,
		service: s,
	}
}

type eventHandler struct {
	cfg     *Config
	cli     iClient
	service app.DefectService
}

func (impl eventHandler) HandleIssueEvent(e *sdk.IssueEvent) error {
	if e.Issue.TypeName != impl.cfg.IssueType {
		return nil
	}

	cmd, err := impl.toCmd(e)
	if err != nil {
		return err
	}

	return impl.service.HandleDefect(cmd)
}
