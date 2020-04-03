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
	hessian.RegisterPOJO(&model.User{})
}

var (
	PodName      = env.RegisterStringVar("POD_NAME", "dubbo-local-user-xxx", "pod name")
	PodNamespace = env.RegisterStringVar("POD_NAMESPACE", "admin", "pod namespace")
	PodIp        = env.RegisterStringVar("POD_IP", "0.0.0.0", "pod ip")
)

type UserProvider struct {
	Cache   map[string]*model.User
	Default *model.User
	Counter int64
}

func NewProvider() *UserProvider {
	return &UserProvider{
		Cache: make(map[string]*model.User),
		Default: &model.User{
			Id:           1,
			Name:         "hello",
			Type:         "user",
			PodName:      PodName.Get(),
			PodNamespace: PodNamespace.Get(),
			PodIp:        PodIp.Get(),
		},
	}
}

func (u *UserProvider) GetUser(ctx context.Context, req []interface{}) (*model.User, error) {
	log.Printf("req:%#v", req)
	var user *model.User
	for _, item := range req {
		if _name, ok := item.(string); ok {
			log.Printf("parse name: %s", _name)
			if _user, ok := u.Cache[_name]; ok {
				user = _user
				break
			}
		}
	}

	if user == nil {
		log.Printf("use delault user")
		user = u.Default
	}

	user.Time = time.Now()
	log.Printf("user:%#v", user)
	atomic.AddInt64(&user.Id, 1)
	return user, nil
}

func (u *UserProvider) SetUser(ctx context.Context, req []interface{}) (*model.User, error) {
	log.Printf("req:%#v", req)
	var user *model.User
	for _, item := range req {
		if _user, ok := item.(*model.User); ok {
			user = _user
			break
		}
	}

	if user == nil {
		log.Printf("user is nil")
		return nil, fmt.Errorf("can't parse user")
	}

	user.Type = "user"
	user.PodName = PodName.Get()
	user.PodNamespace = PodNamespace.Get()
	user.PodIp = PodIp.Get()
	u.Cache[user.Name] = user
	log.Printf("set user:%#v", user)
	return user, nil
}

func (u *UserProvider) Reference() string {
	return "UserProvider"
}
