package server

import (
	"context"
	"fmt"

	apicore "github.com/learn-go/web/common/api/core"
)

type PorxyServer struct {
	apicore.UnimplementedHealthServiceServer
}

func NewPorxyServer() *PorxyServer {
	return &PorxyServer{}
}

func (s *PorxyServer) HealthCheck(ctx context.Context, in *apicore.HealthCheckRequest) (*apicore.HealthCheckResponse, error) {
	return nil, fmt.Errorf("not implemented")
}
