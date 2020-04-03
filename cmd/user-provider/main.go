package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	"github.com/wtxue/k8s-demo/pkg/model"

	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/protocol/dubbo"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	_ "github.com/wtxue/k8s-demo/pkg/provider/user"
)

// they are necessary:
// 		export CONF_PROVIDER_FILE_PATH="xxx"
// 		export APP_LOG_CONF_FILE="xxx"
func main() {

	hessian.RegisterPOJO(&model.User{})
	config.Load()

	initSignal()
}

func initSignal() {
	signals := make(chan os.Signal, 1)
	// It is not possible to block SIGKILL or syscall.SIGSTOP
	signal.Notify(signals, os.Interrupt, os.Kill, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		sig := <-signals
		logger.Infof("get signal %s", sig.String())
		switch sig {
		case syscall.SIGHUP:
			// reload()
		default:
			logger.Infof("provider start exit ...")
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			<-ctx.Done()
			cancel()

			// The program exits normally or timeout forcibly exits.
			logger.Infof("provider app exit now ...")
			return
		}
	}
}
