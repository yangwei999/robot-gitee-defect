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

	comment := func(content string) error {
		return impl.cli.CreateIssueComment(e.Repository.Namespace,
			e.Repository.Name, e.GetIssueNumber(), content,
		)
	}

	issueInfo, err := impl.parse(e.Issue.Body)
	if err != nil {
		return comment(err.Error())
	}

	affectedVersionSlice, err := impl.parseAffectedVersion(issueInfo[itemAffectedVersion])
	if err != nil {
		return comment(err.Error())
	}

	cmd := app.CmdToHandleDefect{
		IssueNumber:     e.GetIssueNumber(),
		IssueOrg:        e.Repository.Namespace,
		IssueRepo:       e.Repository.Name,
		IssueStatus:     *e.State,
		Kernel:          issueInfo[itemKernel],
		Component:       issueInfo[itemComponents],
		SystemVersion:   issueInfo[itemSystemVersion],
		Description:     issueInfo[itemDescription],
		ReferenceURL:    issueInfo[itemReferenceUrl],
		GuidanceURL:     issueInfo[itemGuidanceUrl],
		Influence:       issueInfo[itemInfluence],
		SeverityLevel:   issueInfo[itemSeverityLevel],
		AffectedVersion: affectedVersionSlice,
		ABI:             issueInfo[itemAbi],
	}

	return impl.service.HandleDefect(cmd)
}
