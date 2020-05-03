package http

import (
	// "fmt"
	"go-online/app/admin/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
)

func groupList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.Group
	)
	v := new(struct {
		Page int `form:"pn" default:"1"`
		Size int `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	db := actSrv.DB
	if v.Page == 0 {
		v.Page = 1
	}
	if v.Size == 0 {
		v.Size = 20
	}
	// if v.Status != -1 {
	// 	db = db.Where("status = ?", v.Status)
	// }
	// if v.SID != -1 {
	// 	db = db.Where("sid = ?", v.SID)
	// }
	// db = db.Where("status = ?", 0)
	if err = db.
		Offset((v.Page - 1) * v.Size).Limit(v.Size).
		Find(&list).Error; err != nil {
		log.Error("groupList(%d,%d) error(%v)", v.Page, v.Size, err)
		c.JSON(nil, err)
		return
	}
	if err = db.Model(&model.Group{}).Count(&count).Error; err != nil {
		log.Error("groupList count error(%v)", err)
		c.JSON(nil, err)
		return
	}

	data := map[string]interface{}{
		"data":  list,
		"pn":    v.Page,
		"ps":    v.Size,
		"total": count,
	}
	// c.JSON(data, nil)
	c.JSONMap(data, nil)
}

func groupInfo(c *bm.Context) {
	arg := new(model.GroupDetail)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.GroupID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err := actSrv.DB.Where("group_id = ?", arg.GroupID).First(arg).Error; err != nil {
		log.Error("groupInfo(%d) error(%v)", arg.GroupID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(arg, nil)
}
