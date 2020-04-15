package http

import (
	"net/http"

	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
)

// ping check server ok.
func ping(c *bm.Context) {
	if err := arcSvc.Ping(c); err != nil {
		log.Error("archive service ping error(%v)", err)
		c.AbortWithStatus(http.StatusServiceUnavailable)
	}
}

func register(c *bm.Context) {
	c.JSON(nil, nil)
}
