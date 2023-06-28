package app

import (
	"github.com/opensourceways/robot-gitee-defect/defect/domain"
	"github.com/opensourceways/robot-gitee-defect/defect/domain/defectmanager"
)

type CmdToHandleDefect = domain.Defect

type DefectService interface {
	HandleDefect(defect CmdToHandleDefect) error
}

func NewDefectService(m defectmanager.Manager) *defectService {
	return &defectService{
		manager: m,
	}
}

type defectService struct {
	manager defectmanager.Manager
}

func (d defectService) HandleDefect(cmd CmdToHandleDefect) error {
	return d.manager.Save(cmd)
}
