package order

import (
	"context"
	"fmt"
	"log"
	"sync/atomic"
	"time"

	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go/config"
	"github.com/wtxue/k8s-demo/pkg/model"
	"istio.io/pkg/env"
)

func init() {
	config.SetProviderService(NewProvider())
	// ------for hessian2------
	hessian.RegisterPOJO(&model.Order{})
}

var (
	PodName      = env.RegisterStringVar("POD_NAME", "dubbo-local-order-xxx", "pod name")
	PodNamespace = env.RegisterStringVar("POD_NAMESPACE", "admin", "pod namespace")
	PodIp        = env.RegisterStringVar("POD_IP", "0.0.0.0", "pod ip")
)

type OrderProvider struct {
	Cache   map[string]*model.Order
	Default *model.Order
}

func NewProvider() *OrderProvider {
	return &OrderProvider{
		Cache: make(map[string]*model.Order),
		Default: &model.Order{
			Id:           1,
			Type:         "order",
			Name:         "hello",
			PodName:      PodName.Get(),
			PodNamespace: PodNamespace.Get(),
			PodIp:        PodIp.Get(),
			Time:         time.Now(),
		},
	}
}

func (o *OrderProvider) GetOrder(ctx context.Context, req []interface{}) (*model.Order, error) {
	log.Printf("req:%#v", req)
	var order *model.Order
	for _, item := range req {
		if _name, ok := item.(string); ok {
			log.Printf("parse name: %s", _name)
			if _tmp, ok := o.Cache[_name]; ok {
				order = _tmp
				break
			}
		}
	}

	if order == nil {
		log.Printf("use delault user")
		order = o.Default
	}

	order.Time = time.Now()
	log.Printf("order:%#v", order)
	atomic.AddInt64(&order.Id, 1)
	return order, nil
}

func (o *OrderProvider) SetOrder(ctx context.Context, req []interface{}) (*model.Order, error) {
	log.Printf("req:%#v", req)
	var order *model.Order
	for _, item := range req {
		if _tmp, ok := item.(*model.Order); ok {
			order = _tmp
			break
		}
	}

	if order == nil {
		log.Printf("user is nil")
		return nil, fmt.Errorf("can't parse user")
	}

	order.Type = "order"
	order.PodName = PodName.Get()
	order.PodNamespace = PodNamespace.Get()
	order.PodIp = PodIp.Get()
	order.Time = time.Now()
	o.Cache[order.Name] = order
	log.Printf("set order:%#v", order)
	return order, nil
}

func (o *OrderProvider) Reference() string {
	return "OrderProvider"
}
