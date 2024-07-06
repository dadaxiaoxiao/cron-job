package grpc

import (
	"context"
	cronjobv1 "github.com/dadaxiaoxiao/api-repository/api/proto/gen/cronjob/v1"
	"github.com/dadaxiaoxiao/cron-job/internal/domain"
	"github.com/dadaxiaoxiao/cron-job/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CronJobServiceServer struct {
	cronjobv1.UnimplementedCronJobServiceServer
	svc service.CronJobService
}

func NewCronJobServiceServer(svc service.CronJobService) *CronJobServiceServer {
	return &CronJobServiceServer{
		svc: svc,
	}
}

func (c *CronJobServiceServer) Register(server *grpc.Server) {
	cronjobv1.RegisterCronJobServiceServer(server, c)
}

func (c *CronJobServiceServer) Preempt(ctx context.Context, request *cronjobv1.PreemptRequest) (*cronjobv1.PreemptResponse, error) {
	job, err := c.svc.Preempt(ctx)
	return &cronjobv1.PreemptResponse{
		Cronjob: convertToV(job),
	}, err
}

func (c *CronJobServiceServer) ResetNextTime(ctx context.Context, request *cronjobv1.ResetNextTimeRequest) (*cronjobv1.ResetNextTimeResponse, error) {
	err := c.svc.ResetNextTime(ctx, convertToDomain(request.GetCronjob()))
	return &cronjobv1.ResetNextTimeResponse{}, err
}

func (c *CronJobServiceServer) AddJob(ctx context.Context, request *cronjobv1.AddJobRequest) (*cronjobv1.AddJobResponse, error) {
	err := c.svc.AddJob(ctx, convertToDomain(request.GetCronjob()))
	return &cronjobv1.AddJobResponse{}, err
}

func (c *CronJobServiceServer) StopJob(ctx context.Context, request *cronjobv1.StopJobRequest) (*cronjobv1.StopJobRequest, error) {
	err := c.svc.Stop(ctx, convertToDomain(request.GetCronjob()))
	return &cronjobv1.StopJobRequest{}, err
}

func convertToDomain(vCronJob *cronjobv1.CronJob) domain.CronJob {
	return domain.CronJob{
		Id:         vCronJob.GetId(),
		Name:       vCronJob.GetName(),
		Executor:   vCronJob.GetExecutor(),
		Cfg:        vCronJob.GetCfg(),
		Expression: vCronJob.GetExpression(),
		NextTime:   vCronJob.GetNextTime().AsTime(),
	}
}

func convertToV(domainCronJob domain.CronJob) *cronjobv1.CronJob {
	return &cronjobv1.CronJob{
		Id:         domainCronJob.Id,
		Name:       domainCronJob.Name,
		Executor:   domainCronJob.Executor,
		Cfg:        domainCronJob.Cfg,
		Expression: domainCronJob.Expression,
		NextTime:   timestamppb.New(domainCronJob.NextTime),
	}
}
