package metrics

import (
	"context"

	"github.com/learn-go/web/pkg/utils"
	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

func MetricsServerHandler() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		service, method := utils.ParseRpcInfo(info)
		if len(service)*len(method) > 0 {
			ServiceProcessTotalCount.WithLabelValues(service, method).Inc()
			timer := prometheus.NewTimer(prometheus.ObserverFunc(
				ServiceProcessDuration.WithLabelValues(service, method).Observe))
			defer timer.ObserveDuration()
		}
		resp, err := handler(ctx, req)
		if len(service)*len(method) > 0 {
			if err != nil {
				ServiceProcessFailedCount.WithLabelValues(service, method).Inc()
			} else {
				ServiceProcessSuccedCount.WithLabelValues(service, method).Inc()
			}
		}
		return resp, err
	}
}
