package gorpc

import (
	"context"
	// "fmt"
	v1 "go-online/app/domain/identify/api/grpc"
	"go-online/app/domain/identify/service"
	"net"

	"go-online/lib/ecode"
	// "go-online/lib/log"

	// "go-online/lib/net/metadata"
	// "go-online/lib/net/rpc/warden"

	"google.golang.org/grpc"
)

type server struct {
	svr *service.Service
}

const (
	port = ":50051"
)

// New Identify  rpc server
func New(s *service.Service) *grpc.Server {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		panic(err)
	}
	ws := grpc.NewServer()
	v1.RegisterIdentifyServer(ws, &server{svr: s})
	ws.Serve(lis)
	return ws
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
