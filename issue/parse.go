package issue

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/opensourceways/robot-gitee-defect/utils"
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
	regexpOfItems = map[string]*regexp.Regexp{
		itemKernel:          regexp.MustCompile(`(内核信息)[:：]([\s\S]*?)缺陷归属组件`),
		itemComponents:      regexp.MustCompile(`(缺陷归属组件)[:：]([\s\S]*?)组件版本`),
		itemSystemVersion:   regexp.MustCompile(`(缺陷归属的版本)[:：]([\s\S]*?)缺陷简述`),
		itemDescription:     regexp.MustCompile(`(缺陷简述)[:：]([\s\S]*?)缺陷创建时间`),
		itemReferenceUrl:    regexp.MustCompile(`(缺陷详情参考链接)[:：]([\s\S]*?)缺陷分析指导链接`),
		itemGuidanceUrl:     regexp.MustCompile(`(缺陷分析指导链接)[:：]([\s\S]*?)二、缺陷分析结构反馈`),
		itemInfluence:       regexp.MustCompile(`(影响性分析说明)[:：]([\s\S]*?)缺陷严重等级`),
		itemSeverityLevel:   regexp.MustCompile(`(缺陷严重等级)[:：]([\s\S]*?)受影响版本排查`),
		itemAffectedVersion: regexp.MustCompile(`(受影响版本排查)\(受影响/不受影响\)[:：]([\s\S]*?)abi变化`),
		itemAbi:             regexp.MustCompile(`(abi变化)\(受影响/不受影响\)[:：]([\s\S]*?)$`),
	}

	sortOfItems = []string{
		itemKernel,
		itemComponents,
		itemSystemVersion,
		itemDescription,
		itemReferenceUrl,
		itemGuidanceUrl,
		itemInfluence,
		itemSeverityLevel,
		itemAffectedVersion,
		itemAbi,
	}

	noTrimItem = map[string]bool{
		itemDescription: true,
		itemInfluence:   true,
	}
)

func (impl eventHandler) parse(body string) (map[string]string, error) {
	issueInfo := make(map[string]string)
	for _, item := range sortOfItems {
		match := regexpOfItems[item].FindAllStringSubmatch(body, -1)
		if len(match) < 1 || len(match[0]) < 3 {
			return nil, fmt.Errorf("%s 解析失败", item)
		}

		trimItemInfo := utils.TrimString(match[0][2])
		if trimItemInfo == "" {
			return nil, fmt.Errorf("%s 不允许为空", match[0][1])
		}

		if _, ok := noTrimItem[item]; ok {
			issueInfo[item] = match[0][2]
		} else {
			issueInfo[item] = trimItemInfo
		}
	}

	return issueInfo, nil
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
		return nil, fmt.Errorf("受影响版本排查与当前维护版本不一致，当前维护版本:\n%s",
			strings.Join(impl.cfg.MaintainVersion, "\n"),
		)
	}

	return affectedVersion, nil
}
