package http

import (
	"go-online/app/interface/service"

	"go-online/lib/conf/paladin"
	bm "go-online/lib/net/http/blademaster"

	"go-online/lib/net/http/blademaster/middleware/auth"
	"go-online/lib/net/http/blademaster/middleware/proxy"
	"go-online/lib/net/http/blademaster/middleware/verify"
)

var (
	verifySvc *verify.Verify
	authSvc   *auth.Auth
	actSrv    *service.Service
	// kfcSrv  *kfc.Service
)

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
	initMiddleware()
	engine = bm.DefaultServer(&cfg)
	initRouter(engine)
	if err = engine.Start(); err != nil {
		return
	}
	return
}

func initMiddleware() {
	verifySvc = verify.New(nil)
	authSvc = auth.New(nil)
}

func initRouter(e *bm.Engine) {
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
			gapp.GET("my_list", authSvc.User, verifySvc.Verify, proxyHandler)
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
