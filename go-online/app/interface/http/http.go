package http

import (
	"go-online/app/interface/conf"
	"go-online/app/interface/service"

	// "go-online/app/interface/service/kfc"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"

	// "go-online/lib/net/http/blademaster/middleware/permit"
	"go-online/lib/net/http/blademaster/middleware/auth"
	"go-online/lib/net/http/blademaster/middleware/proxy"
	"go-online/lib/net/http/blademaster/middleware/verify"
)

var (
	verifySvc *verify.Verify
	authSvc   *auth.Auth
	actSrv    *service.Service
	// authSrv *permit.Permit
	// kfcSrv  *kfc.Service
)

// Init init http sever instance.
func Init(c *conf.Config, s *service.Service) {
	actSrv = s
	// kfcSrv = kfc.New(c)
	// authSrv = permit.New(c.Auth)
	initMiddleware(c)
	engine := bm.DefaultServer(c.HTTPServer)
	route(engine)
	if err := engine.Start(); err != nil {
		log.Error("httpx.Serve error(%v)", err)
		panic(err)
	}
}

func initMiddleware(c *conf.Config) {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
}

func route(e *bm.Engine) {
	e.Ping(ping)
	proxyHandler := proxy.NewZoneProxy("sh001", "http://127.0.0.1:8000")
	g := e.Group("/v1/admin")
	{
		testGroup := g.Group("/test")
		{
			testGroup.GET("/param", authSvc.Guest, verifySvc.Verify, proxyHandler)
		}
		gapp := g.Group("/group")
		{
			gapp.GET("/list", authSvc.Guest, verifySvc.Verify, proxyHandler)
			gapp.GET("/info", authSvc.Guest, verifySvc.Verify, proxyHandler)
			gapp.POST("/add", authSvc.User, verifySvc.Verify, proxyHandler)
			gapp.PUT("/save", authSvc.User, verifySvc.Verify, proxyHandler)
			gapp.DELETE("/delete", authSvc.User, verifySvc.Verify, proxyHandler)
		}
	}
	userProxyHandler := proxy.NewZoneProxy("sh001", "http://127.0.0.1:7742")
	userGroup := e.Group("/v1/user")
	{
		userGroup.POST("/register", authSvc.Guest, verifySvc.Verify, userProxyHandler)
		userGroup.PUT("/login", authSvc.Guest, verifySvc.Verify, userProxyHandler)
		userGroup.PUT("/logout", authSvc.User, verifySvc.Verify, userProxyHandler)
	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
