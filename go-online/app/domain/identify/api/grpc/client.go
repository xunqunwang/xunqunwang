// Package v1 .
// NOTE: need registery discovery resovler into grpc before use this client
/*
import (
	"go-online/lib/naming/discovery"
	"go-online/lib/net/rpc/warden/resolver"
)

func main() {
	resolver.Register(discovery.New(nil))
}
*/
package v1

import (
	"context"
	"fmt"
	"go-online/lib/consul"
	"go-online/lib/net/rpc/warden"
	"go-online/lib/net/rpc/warden/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

// AppID unique app id for service discovery
// const AppID = "passport.service.identify"
const (
	// address = ":50051"
	appID = "domain.identify"
	host  = "127.0.0.1:8500"
	token = ""
)

func init() {
	consulBuilder, err := consul.NewConsulDiscovery(consul.Config{Host: host, Token: token})
	if err != nil {
		panic(err)
	}
	resolver.Register(consulBuilder)
}

// NewClient new identify grpc client
func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IdentifyClient, error) {
	opts = append(opts, grpc.WithBalancerName(roundrobin.Name))
	client := warden.NewClient(cfg, opts...)
	conn, err := client.Dial(context.Background(), fmt.Sprintf("consul://default/%s", appID))
	if err != nil {
		return nil, err
	}
	return NewIdentifyClient(conn), nil
}

// func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IdentifyClient, error) {
// 	conn, err := grpc.Dial(address, grpc.WithInsecure())
// 	if err != nil {
// 		panic(err)
// 	}
// 	return NewIdentifyClient(conn), nil
// }
