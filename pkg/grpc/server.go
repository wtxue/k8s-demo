package grpc

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"sync/atomic"

	"google.golang.org/grpc"

	"time"

	"fmt"

	"github.com/go-resty/resty/v2"
	pb "github.com/wtxue/k8s-demo/pkg/api/echo"
	"github.com/wtxue/k8s-demo/pkg/grpc/mid"
	"github.com/wtxue/k8s-demo/pkg/logger"
	"istio.io/pkg/env"
)

var (
	PodName      = env.RegisterStringVar("POD_NAME", "grpc-test-xxx", "pod name")
	PodNamespace = env.RegisterStringVar("POD_NAMESPACE", "default", "pod namespace")
	PodIp        = env.RegisterStringVar("POD_IP", "0.0.0.0", "pod ip")
)

// server implements EchoServer.
type Server struct {
	Connter int64
	Cli     *resty.Client
	pb.UnimplementedEchoServer
}

func NewServer() *Server {
	// Create a Resty Client
	client := resty.New()

	// Retries are configured per client
	client.
		// Set retry count to non zero to enable retries
		SetRetryCount(3).
		// You can override initial retry wait time.
		// Default is 100 milliseconds.
		SetRetryWaitTime(5 * time.Second).
		// MaxWaitTime can be overridden as well.
		// Default is 2 seconds.
		SetRetryMaxWaitTime(20 * time.Second).
		// SetRetryAfter sets callback to calculate wait time between retries.
		// Default (nil) implies exponential backoff with jitter
		SetRetryAfter(func(client *resty.Client, resp *resty.Response) (time.Duration, error) {
			return 0, errors.New("quota exceeded")
		})

	return &Server{
		Cli: client,
	}
}
func (s *Server) UnaryEcho(ctx context.Context, req *pb.EchoRequest) (*pb.EchoResponse, error) {
	logger.Logger.Infof("req: %#v", req)

	body := "org body"
	atomic.AddInt64(&s.Connter, 1)
	if s.Connter%50 == 0 {
		res, err := s.Cli.R().Get("https://www.baidu.com/")
		if err != nil {
			logger.Logger.Warnf("get https://www.baidu.com err:#+v", err)
		}

		body = fmt.Sprintf("https://www.baidu.com body len:%d", res.Size())
	}
	return &pb.EchoResponse{
		Id:      req.Id,
		Message: req.Message,
		Meta: map[string]string{
			"PodName":      PodName.Get(),
			"PodNamespace": PodNamespace.Get(),
			"PodIp":        PodIp.Get(),
			"body":         body,
			"time":         time.Now().String(),
		},
	}, nil
}

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, s pb.EchoServer, port string) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// gRPC server statup options
	opts := []grpc.ServerOption{}

	// add middleware
	opts = mid.AddLogging(logger.Log, opts)

	// register service
	server := grpc.NewServer(opts...)
	pb.RegisterEchoServer(server, s)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			logger.Log.Warn("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	logger.Logger.Info("starting gRPC server...")
	return server.Serve(listen)
}
