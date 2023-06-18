package notification

import "github.com/go-resty/resty/v2"

type Client interface {
	SendSMS(url string, parameter SendSMS) (result *resty.Response, err error)
}
