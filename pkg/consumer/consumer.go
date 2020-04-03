package consumer

import (
	"context"
	"log"
	"sync/atomic"

	"net/http"

	"time"

	"fmt"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/common/logger"
	"github.com/apache/dubbo-go/config"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	pb "github.com/wtxue/k8s-demo/pkg/api/echo"
	"github.com/wtxue/k8s-demo/pkg/model"
	"github.com/wtxue/k8s-demo/pkg/router"
	"google.golang.org/grpc"
	"istio.io/pkg/env"
)

var (
	PodName      = env.RegisterStringVar("POD_NAME", "consumer-local-order-xxx", "pod name")
	PodNamespace = env.RegisterStringVar("POD_NAMESPACE", "admin", "pod namespace")
	PodIp        = env.RegisterStringVar("POD_IP", "0.0.0.0", "pod ip")
)

var UserClient = new(UserConsumer)
var OrderClient = new(OrderConsumer)

func init() {
	config.SetConsumerService(UserClient)
	config.SetConsumerService(OrderClient)
	hessian.RegisterPOJO(&model.User{})
	hessian.RegisterPOJO(&model.Order{})
}

type UserConsumer struct {
	GetUser func(ctx context.Context, req []interface{}, rsp *model.User) error
	SetUser func(ctx context.Context, req []interface{}, rsp *model.User) error
}

type OrderConsumer struct {
	GetOrder func(ctx context.Context, req []interface{}, rsp *model.Order) error
	SetOrder func(ctx context.Context, req []interface{}, rsp *model.Order) error
}

type DownConsumer struct {
	IdCounter int64
	GrpcCli   pb.EchoClient
	*UserConsumer
	*OrderConsumer
}

func (u *UserConsumer) Reference() string {
	return "UserProvider"
}

func (u *OrderConsumer) Reference() string {
	return "OrderProvider"
}

func (u *DownConsumer) GetUserDubbo(c *gin.Context) {
	name := c.DefaultQuery("name", "hello")

	logger.Debugf("start to test dubbo get user name: %s", name)
	user := &model.User{}
	err := u.GetUser(context.TODO(), []interface{}{name}, user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   err.Error(),
			"resultMap": nil,
		})
		return
	}
	logger.Debugf("dubbo response result: %#v\n", user)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"resultMap": user,
	})
}

func (u *DownConsumer) SetUserDubbo(c *gin.Context) {
	name := c.DefaultQuery("name", "hello")
	logger.Debugf("start to test dubbo set user name: %s", name)
	var req = &model.User{
		Name:    name,
		Service: "client",
	}
	res := &model.User{}
	err := u.SetUser(context.TODO(), []interface{}{req}, res)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   err.Error(),
			"resultMap": nil,
		})
		return
	}
	logger.Debugf("dubbo response result: %#v\n", res)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"resultMap": res,
	})
}

func (u *DownConsumer) GetOrderDubbo(c *gin.Context) {
	name := c.DefaultQuery("name", "hello")

	logger.Debugf("start to test dubbo get Order name: %s", name)
	user := &model.Order{}
	err := u.GetOrder(context.TODO(), []interface{}{name}, user)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   err.Error(),
			"resultMap": nil,
		})
		return
	}
	logger.Debugf("dubbo response result: %#v\n", user)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"resultMap": user,
	})
}

func (u *DownConsumer) SetOrderDubbo(c *gin.Context) {
	name := c.DefaultQuery("name", "hello")
	logger.Debugf("start to test dubbo set Order name: %s", name)
	var req = &model.Order{
		Name:    name,
		Service: "client",
	}
	res := &model.Order{}
	err := u.SetOrder(context.TODO(), []interface{}{req}, res)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success":   false,
			"message":   err.Error(),
			"resultMap": nil,
		})
		return
	}
	logger.Debugf("dubbo response result: %#v\n", res)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"resultMap": res,
	})
}

func (u *DownConsumer) GetUserHttp(c *gin.Context) {
	logger.Debugf("start to test http get user")
	atomic.AddInt64(&u.IdCounter, 1)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success": true,
		"Id":      u.IdCounter,
		"podName": PodName.Get(),
		"podNs":   PodNamespace.Get(),
		"podIp":   PodIp.Get(),
		"Time":    time.Now(),
	})
}

func (u *DownConsumer) GetUserGrpc(c *gin.Context) {
	logger.Debugf("start to test grpc get id")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	res, err := u.GrpcCli.UnaryEcho(ctx, &pb.EchoRequest{
		Id:      fmt.Sprintf("%d", u.IdCounter),
		Message: c.Request.RemoteAddr + "keepalive demo",
	})
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": errors.Wrapf(err, "unexpected error from UnaryEcho").Error(),
		})
		return
	}

	atomic.AddInt64(&u.IdCounter, 1)
	c.IndentedJSON(http.StatusOK, gin.H{
		"success":   true,
		"resultMap": res,
	})
}

func (u *DownConsumer) Routes() []*router.Route {
	var routes []*router.Route

	ctlRoutes := []*router.Route{
		{
			Method:  "GET",
			Path:    "/userhttp",
			Handler: u.GetUserHttp,
		},
		{
			Method:  "GET",
			Path:    "/GetUserDubbo",
			Handler: u.GetUserDubbo,
		},
		{
			Method:  "GET",
			Path:    "/SetUserDubbo",
			Handler: u.SetUserDubbo,
		},
		{
			Method:  "GET",
			Path:    "/GetOrderDubbo",
			Handler: u.GetOrderDubbo,
		},
		{
			Method:  "GET",
			Path:    "/SetOrderDubbo",
			Handler: u.SetOrderDubbo,
		},
		{
			Method:  "GET",
			Path:    "/usergrpc",
			Handler: u.GetUserGrpc,
		},
	}

	routes = append(routes, ctlRoutes...)
	return routes
}

func GrpcInit(addr string) pb.EchoClient {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	c := pb.NewEchoClient(conn)
	return c
}

// DefaultHealhRoutes ...
func DefaultHealthRoutes() []*router.Route {
	var routes []*router.Route

	appRoutes := []*router.Route{
		{
			Method:  "GET",
			Path:    "/live",
			Handler: router.LiveHandler,
		},
		{
			Method:  "GET",
			Path:    "/ready",
			Handler: router.LiveHandler,
		},
	}

	routes = append(routes, appRoutes...)
	return routes
}

//
func GinInit() *router.Router {
	rt := router.NewRouter(router.DefaultOption())
	rt.AddRoutes("index", rt.DefaultRoutes())
	rt.AddRoutes("health", DefaultHealthRoutes())
	return rt
}

//
func HttpInit(addr string) *router.Router {
	cli := &DownConsumer{
		UserConsumer:  UserClient,
		OrderConsumer: OrderClient,
	}
	cli.GrpcCli = GrpcInit(addr)
	rt := router.NewRouter(&router.Options{
		Addr:           ":8090",
		GinLogEnabled:  true,
		PprofEnabled:   false,
		MetricsEnabled: true,
	})

	rt.AddRoutes("index", rt.DefaultRoutes())
	rt.AddRoutes("user", cli.Routes())
	return rt
}
