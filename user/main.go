package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"micro-cloudStorage/user/conf"
	"micro-cloudStorage/user/core"
	"micro-cloudStorage/user/model"
	"micro-cloudStorage/user/pkg/gredis"
	"micro-cloudStorage/user/service"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	conf.Setup()
	model.Setup()
	gredis.Setup()
	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	fileService := micro.NewService(
		micro.Registry(etcdReg))

	fileService.Init()

	newFileService := service.NewFileService("rpcFileService", fileService.Client())

	core.FileService = newFileService

	userService := micro.NewService(
		micro.Name("rpcUserService"),
		micro.Address("127.0.0.1:8081"),
		micro.Registry(etcdReg),
	)

	userService.Init()

	_ = service.RegisterUserServiceHandler(userService.Server(), &core.UserService{})

	_ = userService.Run()
}
