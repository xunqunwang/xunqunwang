package gorpc

import (
	"context"
	// "fmt"
	v1 "go-online/app/domain/identify/api/grpc"
	"go-online/app/domain/identify/service"
	"go-online/lib/consul"
	"go-online/lib/ecode"
	"go-online/lib/log"
	"go-online/lib/naming"
	"net"
	"os"

	"github.com/micro/go-micro/v2/util/addr"

	// "go-online/lib/net/metadata"
	// "go-online/lib/net/rpc/warden"

	"google.golang.org/grpc"
)

type server struct {
	svr *service.Service
}

const (
	port  = ":50051"
	host  = "127.0.0.1:8500"
	token = ""
)

// New Identify  rpc server
func New(s *service.Service) (ws *grpc.Server, cancelFunc context.CancelFunc, err error) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	ws = grpc.NewServer()
	v1.RegisterIdentifyServer(ws, &server{svr: s})
	go func() {
		if err := ws.Serve(lis); err != nil {
			panic(err)
		}
	}()

	addr, err := addr.Extract("0.0.0.0")
	if err != nil {
		panic(err)
	}

	hn, _ := os.Hostname()
	instance := &naming.Instance{
		Zone:     "zone1",
		Env:      "dev",
		AppID:    "domain.identify",
		Hostname: hn,
		Addrs: []string{
			"grpc://" + addr + port,
		},
	}

	consulBuilder, err := consul.NewConsulDiscovery(consul.Config{Host: host, Token: token})
	if err != nil {
		panic(err)
	}
	ctx, _ := context.WithCancel(context.Background())
	cancelFunc, err = consulBuilder.Register(ctx, instance)
	return ws, cancelFunc, err
}

var _ v1.IdentifyServer = &server{}
var (
	emptyCookieReply = &v1.GetCookieInfoReply{
		IsLogin: false,
	}

	emptyTokenReply = &v1.GetTokenInfoReply{
		IsLogin: false,
	}
)

// CookieInfo verify user info by cookie.
func (s *server) GetCookieInfo(ctx context.Context, req *v1.GetCookieInfoReq) (*v1.GetCookieInfoReply, error) {
	res, err := s.svr.GetCookieInfo(ctx, req.GetCookie())
	if err != nil {
		if err == ecode.NoLogin {
			return emptyCookieReply, nil
		}
		return nil, err
	}

	return &v1.GetCookieInfoReply{
		IsLogin: true,
		Mid:     res.Mid,
		Expires: res.Expires,
		Csrf:    res.Csrf,
	}, nil
}

// TokenInfo verify user info by token.
func (s *server) GetTokenInfo(ctx context.Context, req *v1.GetTokenInfoReq) (*v1.GetTokenInfoReply, error) {
	token := &v1.GetTokenInfoReq{
		Buvid: req.Buvid,
		Token: req.Token,
	}
	res, err := s.svr.GetTokenInfo(ctx, token)
	if err != nil {
		log.Error("GetTokenInfo(%v) error(%v)", token, err)
		if err == ecode.NoLogin {
			return emptyTokenReply, nil
		}
		return nil, err
	}
	return &v1.GetTokenInfoReply{
		IsLogin: true,
		Mid:     res.Mid,
		Expires: res.Expires,
		Csrf:    res.Csrf,
	}, nil
}
