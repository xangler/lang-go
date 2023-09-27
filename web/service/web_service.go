package service

import (
	"context"

	apicore "github.com/learn-go/web/common/api/core"
)

type WebService struct {
}

func NewWebService() *WebService {
	return &WebService{}
}

func (s *WebService) HealthCheck(ctx context.Context, in *apicore.HealthCheckRequest) (*apicore.HealthCheckResponse, error) {
	return &apicore.HealthCheckResponse{}, nil
}
