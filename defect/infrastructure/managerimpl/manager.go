package managerimpl

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/opensourceways/server-common-lib/utils"

	"github.com/opensourceways/robot-gitee-defect/defect/domain"
)

func NewManagerImpl(c Config) *managerImpl {
	return &managerImpl{
		cli: utils.NewHttpClient(3),
		cfg: c,
	}
}

type managerImpl struct {
	cli utils.HttpClient
	cfg Config
}

func (impl managerImpl) Save(defect domain.Defect) error {
	body, err := json.Marshal(defect)
	if err != nil {
		return err
	}

	url := impl.cfg.Endpoint + "/api/v1/defect"
	request, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}

	_, err = impl.cli.ForwardTo(request, nil)

	return err
}
