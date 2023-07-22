package issue

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/opensourceways/robot-gitee-lib/client"
	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-defect/defect/app"
)

const checkCmd = "/check-issue"

type EventHandler interface {
	HandleIssueEvent(e *sdk.IssueEvent) error
	HandleNoteEvent(e *sdk.NoteEvent) error
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
	return impl.handIssue(e.Issue, e.Project)
}

func (impl eventHandler) HandleNoteEvent(e *sdk.NoteEvent) error {
	if !e.IsIssue() || !strings.Contains(e.Comment.Body, checkCmd) {
		return nil
	}

	return impl.handIssue(e.Issue, e.Project)
}

func (impl eventHandler) handIssue(issue *sdk.IssueHook, project *sdk.ProjectHook) error {
	if issue.TypeName != impl.cfg.IssueType {
		return nil
	}

	commentIssue := func(content string) error {
		return impl.cli.CreateIssueComment(project.Namespace,
			project.Name, issue.Number, content,
		)
	}

	issueInfo, err := impl.parse(issue.Body)
	if err != nil {
		return commentIssue(err.Error())
	}

	affectedVersionSlice, ok := issueInfo[itemAffectedVersion].([]string)
	if !ok {
		return errors.New("parse affected version error")
	}

	if issue.State == sdk.StatusClosed {
		if err := impl.checkPR(project.Namespace, issue.Number, affectedVersionSlice); err != nil {
			return commentIssue(err.Error())
		}
	}

	cmd := app.CmdToHandleDefect{
		IssueNumber:     issue.Number,
		IssueOrg:        project.Namespace,
		IssueRepo:       project.Name,
		IssueStatus:     issue.State,
		Kernel:          fmt.Sprint(issueInfo[itemKernel]),
		Component:       fmt.Sprint(issueInfo[itemComponents]),
		SystemVersion:   fmt.Sprint(issueInfo[itemSystemVersion]),
		Description:     fmt.Sprint(issueInfo[itemDescription]),
		ReferenceURL:    fmt.Sprint(issueInfo[itemReferenceUrl]),
		GuidanceURL:     fmt.Sprint(issueInfo[itemGuidanceUrl]),
		Influence:       fmt.Sprint(issueInfo[itemInfluence]),
		SeverityLevel:   fmt.Sprint(issueInfo[itemSeverityLevel]),
		AffectedVersion: affectedVersionSlice,
		ABI:             fmt.Sprint(issueInfo[itemAbi]),
	}

	err = impl.service.HandleDefect(cmd)
	if err == nil {
		return commentIssue("Your issue is accepted, thank you")
	}

	return err
}

func (impl eventHandler) checkPR(owner, number string, versions []string) error {
	endpoint := fmt.Sprintf("https://gitee.com/api/v5/repos/%v/issues/%v/pull_requests", owner, number)
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}

	var prs []sdk.PullRequest
	cli := utils.NewHttpClient(3)
	bytes, _, err := cli.Download(req)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(bytes, &prs); err != nil {
		return err
	}

	return nil
}