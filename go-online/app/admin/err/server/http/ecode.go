package http

import (
	// "go-online/app/admin/err/model"
	bm "go-online/lib/net/http/blademaster"
)

func getEcodes(c *bm.Context) {
	arg := new(struct {
		Ver int64 `json:"ver" form:"ver" default:"0"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	c.JSON(actSrv.GetEcodeList(arg.Ver))
}
