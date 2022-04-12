package wrappers

import (
	"context"
	"github.com/afex/hystrix-go/hystrix"
	"github.com/micro/go-micro/v2/client"
	"micro-cloudStorage/gateway/pkg/logging"
)

type FileWrapper struct {
	client.Client
}

func (wrapper *FileWrapper) Call(ctx context.Context, req client.Request, resp interface{}, opt ...client.CallOption) error {
	cdmName := req.Service() + req.Endpoint()
	config := hystrix.CommandConfig{
		Timeout:                30000,
		RequestVolumeThreshold: 20,
		ErrorPercentThreshold:  50,
		SleepWindow:            5000,
	}
	hystrix.ConfigureCommand(cdmName, config)
	return hystrix.Do(cdmName, func() error {
		return wrapper.Client.Call(ctx, req, resp)
	}, func(err error) error {
		resp = map[string]interface{}{
			"code": 502,
			"msg":  "服务器忙，稍后再试",
		}
		logging.Info("fileService-- :" + req.Endpoint() + ":服务降级")
		return nil
	})
}

func NewFileWrapper(c client.Client) client.Client {
	return &FileWrapper{c}
}
