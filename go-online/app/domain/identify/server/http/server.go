package http

import (
	pb "go-online/app/domain/identify/api"
	"go-online/app/domain/identify/service"
	"go-online/lib/conf/paladin"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"net/http"
)

var (
	actSrv *service.Service
	svc    pb.IdentifyServer
)

// New new a bm server.
func New(s pb.IdentifyServer) (engine *bm.Engine, err error) {
	var (
		cfg bm.ServerConfig
		ct  paladin.TOML
	)
	if err = paladin.Get("http.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Server").UnmarshalTOML(&cfg); err != nil {
		return
	}
	actSrv = s.(*service.Service)
	svc = s
	engine = bm.DefaultServer(&cfg)
	pb.RegisterIdentifyBMServer(engine, s)
	initRouter(engine)
	if err = engine.Start(); err != nil {
		return
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	group := e.Group("/x/internal/identify")
	{
		group.GET("cookie", accessCookie)
		group.GET("token", accessToken)
		group.GET("cache/del", delCache)
	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		log.Error("ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}
