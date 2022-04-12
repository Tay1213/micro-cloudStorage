package main

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/registry"
	"github.com/micro/go-micro/v2/registry/etcd"
	"micro-cloudStorage/file/conf"
	"micro-cloudStorage/file/core"
	"micro-cloudStorage/file/model"
	"micro-cloudStorage/file/service"

	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	conf.Setup()
	model.Setup()

	etcdReg := etcd.NewRegistry(
		registry.Addrs("127.0.0.1:2379"),
	)

	microService := micro.NewService(
		micro.Name("rpcFileService"),
		micro.Address("127.0.0.1:8082"),
		micro.Registry(etcdReg),
	)

	microService.Init()

	_ = service.RegisterFileServiceHandler(microService.Server(), &core.FileService{})

	_ = microService.Run()
}
