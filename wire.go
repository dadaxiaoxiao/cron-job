//go:build wireinject

package main

import (
	grpc2 "github.com/dadaxiaoxiao/cron-job/internal/grpc"
	"github.com/dadaxiaoxiao/cron-job/internal/repository"
	"github.com/dadaxiaoxiao/cron-job/internal/repository/dao"
	"github.com/dadaxiaoxiao/cron-job/internal/service"
	"github.com/dadaxiaoxiao/cron-job/ioc"
	"github.com/dadaxiaoxiao/go-pkg/customserver"
	"github.com/google/wire"
)

var thirdPartyProvider = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis,
	ioc.InitEtcdClient,
	ioc.InitLogger,
	ioc.InitOTEL,
)

func InitApp() *customserver.App {
	wire.Build(
		thirdPartyProvider,
		dao.NewGORMJobDAO,
		repository.NewPreemptCronJobRepository,
		service.NewCronJobService,
		grpc2.NewCronJobServiceServer,
		ioc.InitGRPCServer,
		wire.Struct(new(customserver.App), "GRPCServer"),
	)
	return new(customserver.App)
}
