package http

import (
	"go-online/app/group/service"
	"go-online/lib/conf/paladin"
	bm "go-online/lib/net/http/blademaster"
)

var (
	actSrv *service.Service
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
	engine = bm.DefaultServer(&cfg)
	initRouter(engine)
	if err = engine.Start(); err != nil {
		return
	}
	return
}

func initRouter(e *bm.Engine) {
	e.Ping(ping)
	g := e.Group("/v1/admin")
	{
		gapp := g.Group("/group")
		{
			gapp.GET("my_list", myGroupList)
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
