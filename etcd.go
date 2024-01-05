package etcd

import (
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/glog"
	etcd3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Option struct {
	Logger       glog.ILogger
	KeepaliveTTL time.Duration
}

type Registry struct {
	client       *etcd3.Client
	kv           etcd3.KV
	lease        etcd3.Lease
	keepaliveTTL time.Duration
	logger       glog.ILogger
}

const (
	// DefaultKeepAliveTTL is the default keepalive TTL.
	DefaultKeepAliveTTL = 10 * time.Second
)

func New(etcdCfg etcd3.Config, opts ...Option) gsvc.Registry {

	if len(etcdCfg.Endpoints) == 0 {
		panic(gerror.NewCodef(gcode.CodeInvalidParameter, `invalid etcd address "%s"`, etcdCfg.Endpoints))
	}
	client, err := etcd3.New(etcdCfg)
	if err != nil {
		panic(gerror.Wrap(err, `create etcd client failed`))
	}

	return NewWithClient(client, opts...)
}

func NewWithClient(client *etcd3.Client, option ...Option) *Registry {

	r := &Registry{
		client: client,
		kv:     etcd3.NewKV(client),
	}
	if len(option) > 0 {
		r.logger = option[0].Logger
		r.keepaliveTTL = option[0].KeepaliveTTL
	}
	if r.logger == nil {
		r.logger = g.Log()
	}
	if r.keepaliveTTL == 0 {
		r.keepaliveTTL = DefaultKeepAliveTTL
	}
	return r
}

// extractResponseToServices extracts etcd watch response context to service list.
func extractResponseToServices(res *etcd3.GetResponse) ([]gsvc.Service, error) {
	if res == nil || res.Kvs == nil {
		return nil, nil
	}
	var (
		services         []gsvc.Service
		servicePrefixMap = make(map[string]*Service)
	)
	for _, kv := range res.Kvs {
		service, err := gsvc.NewServiceWithKV(
			string(kv.Key), string(kv.Value),
		)
		if err != nil {
			return services, err
		}
		s := NewService(service)
		if v, ok := servicePrefixMap[service.GetPrefix()]; ok {
			v.Endpoints = append(v.Endpoints, service.GetEndpoints()...)
		} else {
			servicePrefixMap[s.GetPrefix()] = s
			services = append(services, s)
		}
	}
	return services, nil
}
