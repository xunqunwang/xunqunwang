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
	// "context"

	"go-online/lib/net/rpc/warden"
	"google.golang.org/grpc"
)

// AppID unique app id for service discovery
// const AppID = "passport.service.identify"

// NewClient new identify grpc client
// func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IdentifyClient, error) {
// 	client := warden.NewClient(cfg, opts...)
// 	conn, err := client.Dial(context.Background(), "discovery://default/"+AppID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return NewIdentifyClient(conn), nil
// }

const (
	address = ":50051"
)

func NewClient(cfg *warden.ClientConfig, opts ...grpc.DialOption) (IdentifyClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	return NewIdentifyClient(conn), nil
}
