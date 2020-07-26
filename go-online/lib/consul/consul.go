package consul

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go-online/lib/log"
	"go-online/lib/naming"
	llog "log"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
)

var (
	ERR_INS_ADDRS_EMPTY = errors.New("len of ins.Addrs should not be 0")
)

type logWrapper struct {
}

func (wrapper logWrapper) Write(p []byte) (n int, err error) {
	log.Info(string(p))
	return len(p), nil
}

// Config discovery configures.
type Config struct {
	Nodes  []string
	Region string
	Zone   string
	Env    string
	Host   string
	Token  string
}

// Resolver resolve naming service
type Resolver struct {
	appID   string
	c       chan struct{}
	client  *api.Client
	agent   *api.Agent
	plan    *watch.Plan
	builder *Builder
	ins     atomic.Value
}

func (resolver *Resolver) watch() error {
	var params map[string]interface{}
	watchKey := fmt.Sprintf(`{"type":"service", "service":"%s"}`, resolver.appID)
	if err := json.Unmarshal([]byte(watchKey), &params); err != nil {
		return err
	}
	plan, err := watch.Parse(params)
	if err != nil {
		return err
	}
	plan.Handler = func(idx uint64, raw interface{}) {
		if raw == nil {
			return // ignore
		}
		v, ok := raw.([]*api.ServiceEntry)
		if !ok {
			return // ignore
		}
		log.Info("consul watch service %s notify, len %d", resolver.appID, len(v))
		ins := resolver.coverServiceEntry2Ins(v)
		resolver.ins.Store(ins)
		resolver.c <- struct{}{}
	}

	logger := llog.New(&logWrapper{}, "", llog.LstdFlags) // replace logger
	go func() {
		err := plan.RunWithClientAndLogger(resolver.client, logger)
		if err != nil {
			log.Error("watch service %s error %s", resolver.appID, err.Error())
		}
	}()
	resolver.plan = plan
	return nil
}

func (resolver *Resolver) Watch() <-chan struct{} {
	return resolver.c
}

func (resolver Resolver) coverServiceEntry2Ins(serviceArr []*api.ServiceEntry) map[string][]*naming.Instance {
	instances := make(map[string][]*naming.Instance)
	for _, service := range serviceArr {
		if service.Checks.AggregatedStatus() == api.HealthPassing {
			log.Info("appid %s ip %s port %d pass", resolver.appID, service.Service.Address, service.Service.Port)
			ins := resolver.coverService2Instance(service.Service)
			if _, ok := instances[ins.Zone]; !ok {
				instances[ins.Zone] = make([]*naming.Instance, 0, 10)
			}
			instances[ins.Zone] = append(instances[ins.Zone], ins)
		}
	}
	return instances
}

func (resolver *Resolver) Fetch(c context.Context) (ins map[string][]*naming.Instance, ok bool) {
	v := resolver.ins.Load()
	ins, ok = v.(map[string][]*naming.Instance)
	return
}

func (resolver Resolver) Close() error {
	if resolver.plan != nil && !resolver.plan.IsStopped() {
		resolver.plan.Stop()
	}
	return nil
}

func (resolver Resolver) coverService2Instance(service *api.AgentService) *naming.Instance {
	meta := service.Meta
	addr := []string{
		service.Address + ":" + strconv.Itoa(service.Port),
	}
	ins := &naming.Instance{
		Region:   meta["region"],
		Zone:     meta["zone"],
		Env:      meta["env"],
		Hostname: meta["hostname"],
		Version:  meta["version"],
		AppID:    service.Service,
		Addrs:    addr,
	}
	ins.Metadata = make(map[string]string)
	for key, value := range meta {
		if key == "region" || key == "env" || key == "zone" || key == "version" || key == "hostname" {
			continue
		}
		ins.Metadata[key] = value
	}
	ins.LastTs = time.Now().Unix()
	return ins
}

func (builder Builder) coverIns2AgentService(ins *naming.Instance) ([]*api.AgentServiceRegistration, error) {
	if len(ins.Addrs) == 0 {
		return nil, ERR_INS_ADDRS_EMPTY
	}
	registrationArr := make([]*api.AgentServiceRegistration, len(ins.Addrs))
	meta := make(map[string]string)
	meta["region"] = ins.Region
	meta["zone"] = ins.Zone
	meta["env"] = ins.Env
	meta["hostname"] = ins.Hostname
	meta["version"] = ins.Version
	meta["last_ts"] = strconv.FormatInt(ins.LastTs, 10)

	for key, value := range ins.Metadata {
		meta[key] = value
	}
	for i, addr := range ins.Addrs {
		urlVal, err := url.Parse(addr)
		if err != nil {
			return nil, err
		}
		port, _ := strconv.Atoi(urlVal.Port())
		service := &api.AgentServiceRegistration{
			ID:      ins.AppID + "-" + urlVal.Hostname() + "-" + urlVal.Port(),
			Name:    ins.AppID,
			Kind:    api.ServiceKindTypical,
			Port:    port,
			Address: urlVal.Scheme + "://" + urlVal.Hostname(),
			Meta:    meta,
		}
		registrationArr[i] = service
	}
	return registrationArr, nil
}

func (builder Builder) Register(ctx context.Context, ins *naming.Instance) (cancelFunc context.CancelFunc, err error) {
	serviceArr, err := builder.coverIns2AgentService(ins)
	if err != nil {
		return
	}
	ch := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(ctx)
	// defer func() {
	// 	if err != nil { // avoid register partition
	// 		cancel()
	// 	}
	// }()
	cancelFunc = context.CancelFunc(func() {
		cancel()
		<-ch
	})
	for _, service := range serviceArr { //@todo 批量注册
		service.Check = &api.AgentServiceCheck{
			TTL:    "15s",
			Status: api.HealthPassing,
		}
		var status string
		var info *api.AgentServiceChecksInfo
		status, info, err = builder.agent.AgentHealthServiceByID(service.ID)
		if err != nil {
			return
		}
		if info == nil && status == api.HealthCritical {
			err = builder.agent.ServiceRegister(service) // @todo check had registered
			if err != nil {
				return
			}
		} else {
			err = builder.agent.PassTTL(fmt.Sprintf("service:%s", service.ID), "I am good :)")
			if err != nil {
				return
			}
		}

		go func(service *api.AgentServiceRegistration) {
			for {
				select {
				case <-ctx.Done():
					log.Info("ServiceDeregister %s", service.ID)
					err := builder.agent.ServiceDeregister(service.ID)
					if err != nil {
						log.Error("consul: ServiceDeregister %s err: %s", service.ID, err.Error())
					}
					ch <- struct{}{}
					return
				case <-time.After(time.Second * 5):
					err := builder.agent.PassTTL(fmt.Sprintf("service:%s", service.ID), "I am good :)")
					if err == nil {
						continue
					}
					log.Error("consul: PassTTL %s err: %s", service.ID, err.Error())
					if strings.Index(err.Error(), "does not have associated TTL") > 0 { // 注册已经失效
						err = builder.agent.ServiceRegister(service) // consul 下线会导致 有这个 error
						if err != nil {
							log.Error("consul: PassTTL %s reRegister err: %s", service.ID, err.Error())
						}
					}
				}
			}
		}(service)
	}
	return
}

func (builder Builder) Close() error {
	return nil
}

type Builder struct {
	client *api.Client
	agent  *api.Agent
	r      map[string]*Resolver
	locker sync.RWMutex
	c      *Config
}

func (builder Builder) Build(id string) naming.Resolver {
	builder.locker.RLock()
	if r, ok := builder.r[id]; ok {
		builder.locker.RUnlock()
		return r
	}
	builder.locker.RUnlock()
	builder.locker.Lock()
	r := &Resolver{
		appID:   id,
		client:  builder.client,
		agent:   builder.agent,
		builder: &builder,
	}
	r.c = make(chan struct{}, 10)
	builder.r[id] = r
	builder.locker.Unlock()
	err := r.watch()
	if err != nil {
		log.Error("watch error %s", err.Error())
	}
	return r
}

func (builder Builder) Scheme() string {
	return "consul"
}

func NewConsulDiscovery(c Config) (builder Builder, err error) {
	client, err := api.NewClient(&api.Config{Address: c.Host, Token: c.Token})
	if err != nil {
		return
	}
	builder.client = client
	builder.agent = client.Agent()
	builder.r = make(map[string]*Resolver)
	builder.c = &c
	return
}
