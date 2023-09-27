package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	apicore "github.com/learn-go/web/common/api/core"
	"github.com/learn-go/web/config"
	"github.com/learn-go/web/pkg/forward"
	"github.com/learn-go/web/pkg/frame"
	"github.com/learn-go/web/server"
	"github.com/learn-go/web/version"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var (
	configPath   = flag.String("c", "./config/go-proxy-server.yaml", "config file")
	printVersion = flag.Bool("version", false, "print version of this build")
)

func main() {
	flag.Parse()
	if *printVersion {
		version.PrintFullVersionInfo()
		return
	}
	c, err := config.LoadProxyConfig(*configPath)
	if err != nil {
		log.Fatal(err)
	}
	if c.MetricPort > 0 {
		setupMetrics(c.MetricPort)
	}

	opts := &frame.Options{}
	svcObj := server.NewPorxyServer()
	server := frame.NewProtocServer(opts).SetHttpPort(c.HTTPPort).SetRpcPort(c.RPCPort).
		SetRpcOptions([]grpc.UnaryServerInterceptor{forward.ForwardServerHandler(c.ForwardURL)}).
		RPC(apicore.RegisterHealthServiceServer, svcObj).
		HTTP(apicore.RegisterHealthServiceHandlerFromEndpoint)
	server.Run()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
Loop:
	for s := range sigs {
		log.Warnf("received signal %v", s)
		server.Close()
		break Loop
	}
}

func setupMetrics(ep int) {
	log.Println("metrics listening at: ", ep)
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", ep), nil); err != nil {
			log.Fatal(err)
		}
	}()
}
