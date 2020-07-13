package http

import (
	"go-online/app/user/conf"
	"go-online/app/user/service"

	// "go-online/app/user/service/kfc"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/net/http/blademaster/middleware/permit"
)

var (
	actSrv  *service.Service
	authSrv *permit.Permit
	// kfcSrv  *kfc.Service
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	actSrv = s
	// kfcSrv = kfc.New(c)
	authSrv = permit.New(c.Auth)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve error(%v)", err)
		panic(err)
	}
}

func route(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/v1/user")
	{
		g.POST("/register", userRegister)
		g.PUT("/login", userLogin)
		g.PUT("/logout", userLogout)
	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
