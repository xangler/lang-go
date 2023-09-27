package server

import (
	"context"
	"time"

	apicore "github.com/learn-go/web/common/api/core"
	"github.com/learn-go/web/pkg/leader"
)

type LeaderServer struct {
	poodLeader *leader.PodLeader
	apicore.UnimplementedHealthServiceServer
}

func NewLeaderServer(poodLeader *leader.PodLeader) *LeaderServer {
	return &LeaderServer{
		poodLeader: poodLeader,
	}
}

func (s *LeaderServer) HealthCheck(ctx context.Context, in *apicore.HealthCheckRequest) (*apicore.HealthCheckResponse, error) {
	if client := s.poodLeader.GetHealthServiceClient(); client != nil {
		cctx, cancel := context.WithTimeout(ctx, 60*time.Second)
		defer cancel()
		return client.HealthCheck(cctx, in)
	}
	return &apicore.HealthCheckResponse{}, nil
}
