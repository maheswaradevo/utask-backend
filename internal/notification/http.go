package notification

import (
	"github.com/go-resty/resty/v2"
	"go.uber.org/zap"
)

type wavecellHandler struct {
	client *resty.Request
}

func NewClient(client *resty.Request) Client {
	return wavecellHandler{client: client}
}

func (w wavecellHandler) SendSMS(url string, parameter SendSMS) (result *resty.Response, err error) {
	result, err = w.client.SetBody(parameter).Post(url)
	if err != nil {
		zap.S().Error(err)
		return
	}
	return
}
