package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"github.com/micro/go-micro/v2/web"
	"micro-cloudStorage/gateway/conf"
	"micro-cloudStorage/gateway/pkg/gredis"
	"micro-cloudStorage/gateway/service"
	"micro-cloudStorage/gateway/weblib"
	"micro-cloudStorage/gateway/weblib/handlers/v1"
	"micro-cloudStorage/gateway/wrappers"
	"time"
)

func main() {
	conf.Setup()
	gredis.Setup()
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	userMicroService := micro.NewService(
		micro.Name("userService.client"),
		micro.WrapClient(wrappers.NewUserWrapper),
	)
	userService := service.NewUserService("rpcUserService", userMicroService.Client())

	fileMicroService := micro.NewService(
		micro.Name("fileService.client"),
		micro.WrapClient(wrappers.NewFileWrapper),
	)
	fileService := service.NewFileService("rpcFileService", fileMicroService.Client())

	v1.UserServices = userService
	v1.FileServices = fileService

	server := web.NewService(
		web.Name("httpService"),
		web.Address("127.0.0.1:9000"),
		web.Registry(etcdReg),
		web.Handler(weblib.NewRouter()),
		web.RegisterTTL(time.Second*30),
		web.RegisterInterval(time.Second*15),
		web.Metadata(map[string]string{"protocol": "http"}),
	)

	_ = server.Init()
	_ = server.Run()

}
