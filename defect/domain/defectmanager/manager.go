package defectmanager

import "github.com/opensourceways/robot-gitee-defect/defect/domain"

type Manager interface {
	Save(defect domain.Defect) error
}
