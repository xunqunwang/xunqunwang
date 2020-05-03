package http

import (
	// "fmt"
	"go-online/app/admin/model"
	bm "go-online/lib/net/http/blademaster"
)

func archives(c *bm.Context) {
	p := &model.ArchiveParam{}
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, nil)
	// c.JSON(actSrv.Archives(c, p))
}
