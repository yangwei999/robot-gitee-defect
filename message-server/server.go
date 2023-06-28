package messageserver

import (
	kafka "github.com/opensourceways/kafka-lib/agent"

	"github.com/opensourceways/robot-gitee-defect/issue"
)

func Init(cfg *Config, handler issue.EventHandler) error {
	s := messageServer{
		handler: giteeEventHandler{
			handler:   handler,
			userAgent: cfg.UserAgent,
		},
	}

	return s.subscribe(cfg)
}

type messageServer struct {
	handler giteeEventHandler
}

func (m *messageServer) subscribe(cfg *Config) error {
	subscribers := map[string]kafka.Handler{
		cfg.Topics.DefectEvent: m.handler.handle,
	}

	return kafka.Subscribe(cfg.GroupName, subscribers)
}
