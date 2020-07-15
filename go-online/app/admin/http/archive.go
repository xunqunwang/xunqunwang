package http

import (
	// "fmt"
	"go-online/app/admin/model"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
)

func testParam(c *bm.Context) {
	arg := new(struct {
		Ids []int64 `json:"ids" form:"ids,split" validate:"required,dive,min=1"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("testParam error(%v)", err)
		return
	}
	c.JSON(nil, nil)
}

func archives(c *bm.Context) {
	p := &model.ArchiveParam{}
	if err := c.Bind(p); err != nil {
		return
	}
	c.JSON(nil, nil)
	// c.JSON(actSrv.Archives(c, p))
}
