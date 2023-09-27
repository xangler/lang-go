package forward

import (
	"context"
	"fmt"
	"time"

	apicore "github.com/learn-go/web/common/api/core"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func ForwardServerHandler(forwardURL string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		conn, err := grpc.Dial(forwardURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		var resp interface{}
		switch info.FullMethod {
		case "/learngo.web.common.core.HealthService/HealthCheck":
			resp = new(apicore.HealthCheckResponse)
		default:
			log.Fatalf("not supprot %s", info.FullMethod)
		}

		fmt.Printf("demo:%v\n", req)
		cctx := metadata.AppendToOutgoingContext(ctx, "X-Forwarded-For", forwardURL)
		cctx, cancel := context.WithTimeout(cctx, 60*time.Second)
		defer cancel()
		log.Info("forward >>> ", forwardURL)
		err = grpc.Invoke(cctx, info.FullMethod, req, resp, conn)
		fmt.Printf("demo:%v, %v\n", resp, err)
		return resp, err
	}
}
