package http

// import (
// 	"go-online/app/domain/identify/conf"
// 	"go-online/app/domain/identify/service"
// 	"go-online/lib/log"
// 	bm "go-online/lib/net/http/blademaster"
// 	"go-online/lib/net/http/blademaster/middleware/verify"
// )

// var (
// 	srv *service.Service
// 	vfy *verify.Verify
// )

// // Init init
// func Init(c *conf.Config, s *service.Service) {
// 	srv = s
// 	vfy = verify.New(c.VerifyConfig)

// 	// engine
// 	engIn := bm.DefaultServer(c.BM)
// 	innerRouter(engIn)
// 	// init inner server
// 	if err := engIn.Start(); err != nil {
// 		log.Error("engIn.Start error(%v)", err)
// 		panic(err)
// 	}
// }

// func innerRouter(e *bm.Engine) {
// 	e.Ping(ping)
// 	e.Register(register)
// 	group := e.Group("/x/internal/identify", vfy.Verify)
// 	{
// 		group.GET("cookie", accessCookie)
// 		group.GET("token", accessToken)
// 		group.GET("cache/del", delCache)
// 	}
// }

// // ping check server ok.
// func ping(c *bm.Context) {
// }

// // register support discovery.
// func register(c *bm.Context) {
// 	c.JSON(map[string]struct{}{}, nil)
// }
