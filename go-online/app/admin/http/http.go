package http

import (
	"go-online/app/admin/conf"
	"go-online/app/admin/service"

	// "go-online/app/admin/service/kfc"
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
	g := e.Group("/v1/admin")
	{
		gapp := g.Group("/group")
		{
			gapp.GET("/list", groupList)
			gapp.GET("/info", groupInfo)
			gapp.POST("/add", addGroup)
			gapp.PUT("/save", saveGroup)
			gapp.DELETE("/delete", deleteGroup)
		}
		testGroup := g.Group("/test")
		{
			testGroup.GET("/param", testParam)
		}
	}
}

func ping(c *bm.Context) {
	if err := actSrv.Ping(c); err != nil {
		c.Error = err
		c.AbortWithStatus(503)
	}
}
