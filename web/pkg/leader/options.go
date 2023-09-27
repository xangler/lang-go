package leader

import (
	apicore "github.com/learn-go/web/common/api/core"
	log "github.com/sirupsen/logrus"
)

func (s *PodLeader) GetHealthServiceClient() apicore.HealthServiceClient {
	if s.IsLeader() || s.GetForwarder() == nil {
		return nil
	}
	conn, addr := s.GetForwarder().GetConnAddr()
	log.Info("HealthServiceClient forward >>> ", addr)
	return apicore.NewHealthServiceClient(conn)
}
