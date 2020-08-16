package dao

import (
	"context"
	"go-online/lib/conf/env"
	"go-online/lib/conf/paladin"
	"go-online/lib/consul"
	"go-online/lib/naming"
	"strings"

	"github.com/micro/go-micro/v2/util/addr"
)

func RegisterToConsul(schema, address string) (cf func(), err error) {
	var (
		cfg consul.Config
		ct  paladin.TOML
	)
	if err = paladin.Get("consul.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	addr, err := addr.Extract("0.0.0.0")
	if err != nil {
		return nil, err
	}
	s := strings.Split(address, ":")
	port := s[1]
	instance := &naming.Instance{
		Region:   env.Region,
		Zone:     env.Zone,
		Env:      env.DeployEnv,
		AppID:    "domain.identify",
		Hostname: env.Hostname,
		Addrs: []string{
			schema + "://" + addr + ":" + port,
		},
	}

	consulBuilder, err := consul.NewConsulDiscovery(cfg)
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithCancel(context.Background())
	cf, err = consulBuilder.Register(ctx, instance)
	return cf, err
}
