package http

import (
	"go-online/app/admin/err/service"
	"go-online/lib/conf/paladin"
	bm "go-online/lib/net/http/blademaster"
)

var actSrv *service.Service

// New new a bm server.
func New(s *service.Service) (engine *bm.Engine, err error) {
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
	actSrv = s
	engine = bm.DefaultServer(&cfg)
	initRouter(engine)
	if err = engine.Start(); err != nil {
		return
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/v1/admin/err")
	{
		g.GET("/list", getEcodes)
	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
