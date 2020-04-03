package main

import (
	"context"
	"flag"
	"fmt"

	"log"

	"github.com/wtxue/k8s-demo/pkg/grpc"
	"github.com/wtxue/k8s-demo/pkg/logger"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// GRPCPort is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string

	// Log parameters section
	// LogLevel is global log level: Debug(-1), Info(0), Warn(1), Error(2), DPanic(3), Panic(4), Fatal(5)
	LogLevel int
	// LogTimeFormat is print time format for logger e.g. 2006-01-02T15:04:05Z07:00
	LogTimeFormat string
}

func main() {
	ctx := context.Background()

	// get configuration
	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "8000", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8090", "HTTP port to bind")
	flag.IntVar(&cfg.LogLevel, "log-level", 0, "Global log level")
	flag.StringVar(&cfg.LogTimeFormat, "log-time-format", "",
		"Print time format for logger e.g. 2006-01-02T15:04:05Z07:00")
	flag.Parse()

	// initialize logger
	if err := logger.Init(cfg.LogLevel, cfg.LogTimeFormat); err != nil {
		fmt.Printf("failed to initialize logger: %v", err)
		return
	}

	log.Printf("grpc RunServer port %s ... ", cfg.GRPCPort)
	if err := grpc.RunServer(ctx, grpc.NewServer(), cfg.GRPCPort); err != nil {
		panic(err)
	}
}
