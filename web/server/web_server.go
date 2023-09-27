package server

import (
	"context"

	apicore "github.com/learn-go/web/common/api/core"
	"github.com/learn-go/web/config"
	"github.com/learn-go/web/module/job"
	"github.com/learn-go/web/module/transaction"
	"github.com/learn-go/web/pkg/dbutils"
	"github.com/learn-go/web/service"
)

type WebServer struct {
	conf *config.WebConfig
	web  *service.WebService
	job  *job.WorkerJob
	apicore.UnimplementedHealthServiceServer
}

func NewWebServer(conf *config.WebConfig) (*WebServer, error) {
	dClient, err := transaction.NewTransaction(
		dbutils.ConnectDBbyCfg(conf.SQL),
		nil,
	)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	job := job.NewWorkerJob(conf, dClient)
	job.Start(ctx)
	web := service.NewWebService()
	return &WebServer{
		conf: conf,
		web:  web,
		job:  job,
	}, nil
}

func (s *WebServer) HealthCheck(ctx context.Context, in *apicore.HealthCheckRequest) (*apicore.HealthCheckResponse, error) {
	return s.web.HealthCheck(ctx, in)
}
