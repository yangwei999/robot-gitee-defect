package issue

import (
	"fmt"
	"regexp"
	"strings"

	sdk "github.com/opensourceways/go-gitee/gitee"
	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/robot-gitee-defect/defect/app"
)

const (
	itemKernel          = "kernel"
	itemComponents      = "components"
	itemSystemVersion   = "systemVersion"
	itemDescription     = "description"
	itemReferenceUrl    = "referenceUrl"
	itemGuidanceUrl     = "guidanceUrl"
	itemInfluence       = "influence"
	itemSeverityLevel   = "severityLevel"
	itemAffectedVersion = "affectedVersion"
	itemAbi             = "abi"
)

var (
	regexpKernel          = regexp.MustCompile(`内核信息[:：]([\s\S]*?)缺陷归属组件`)
	regexpComponents      = regexp.MustCompile(`缺陷归属组件[:：]([\s\S]*?)组件版本`)
	regexpSystemVersion   = regexp.MustCompile(`缺陷归属的版本[:：]([\s\S]*?)缺陷简述`)
	regexpDescription     = regexp.MustCompile(`缺陷简述[:：]([\s\S]*?)缺陷创建时间`)
	regexpReferenceURL    = regexp.MustCompile(`缺陷详情参考链接[:：]([\s\S]*?)缺陷分析指导链接`)
	regexpGuidanceURL     = regexp.MustCompile(`缺陷分析指导链接[:：]([\s\S]*?)二、缺陷分析结构反馈`)
	regexpInfluence       = regexp.MustCompile(`影响性分析说明[:：]([\s\S]*?)缺陷严重等级`)
	regexpSeverityLevel   = regexp.MustCompile(`缺陷严重等级[:：]([\s\S]*?)受影响版本排查`)
	regexpAffectedVersion = regexp.MustCompile(`受影响版本排查\(受影响/不受影响\)[:：]([\s\S]*?)abi变化`)
	regexpAbi             = regexp.MustCompile(`abi变化\(受影响/不受影响\)[:：]([\s\S]*?)$`)

	regexpOfItems = map[string]*regexp.Regexp{
		itemKernel:          regexpKernel,
		itemComponents:      regexpComponents,
		itemSystemVersion:   regexpSystemVersion,
		itemDescription:     regexpDescription,
		itemReferenceUrl:    regexpReferenceURL,
		itemGuidanceUrl:     regexpGuidanceURL,
		itemInfluence:       regexpInfluence,
		itemSeverityLevel:   regexpSeverityLevel,
		itemAffectedVersion: regexpAffectedVersion,
		itemAbi:             regexpAbi,
	}
)

func (impl eventHandler) parse(body string) (map[string]string, error) {
	trim := func(s string) string {
		t := strings.Replace(s, " ", "", -1)
		return strings.Replace(t, "\n", "", -1)
	}

	info := make(map[string]string)
	for item, reg := range regexpOfItems {
		match := reg.FindAllStringSubmatch(body, -1)
		if len(match) < 1 || len(match[0]) < 2 {
			return nil, fmt.Errorf("parse %s failed", item)
		}

		info[item] = trim(match[0][1])
	}

	return info, nil
}

func (impl eventHandler) parseAffectedVersion(s string) ([]string, error) {
	reg := regexp.MustCompile(`(openEuler.*?)[:：](不?受影响)`)
	matches := reg.FindAllStringSubmatch(s, -1)

	var affectedVersion []string
	var allVersion []string
	for _, v := range matches {
		allVersion = append(allVersion, v[1])

		if v[2] == "受影响" {
			affectedVersion = append(affectedVersion, v[1])
		}
	}

	av := sets.NewString(allVersion...)
	if !av.HasAll(impl.cfg.MaintainVersion...) {
		return nil, fmt.Errorf("affected version does not match MaintainVersion, "+
			"MaintainVersion are %s", impl.cfg.MaintainVersion)
	}

	return affectedVersion, nil
}

func (impl eventHandler) toCmd(e *sdk.IssueEvent) (cmd app.CmdToHandleDefect, err error) {
	issueInfo, err := impl.parse(e.Issue.Body)
	if err != nil {
		return
	}

	affectedVersionSlice, err := impl.parseAffectedVersion(issueInfo[itemAffectedVersion])
	if err != nil {
		impl.cli.CreateIssueComment(e.Repository.Namespace,
			e.Repository.Name, e.GetIssueNumber(), err.Error(),
		)

		return
	}

	cmd = app.CmdToHandleDefect{
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

	return
}