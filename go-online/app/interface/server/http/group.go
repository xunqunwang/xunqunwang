package http

import (
	// "fmt"
	"go-online/app/interface/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/pic"
	"time"
)

func groupList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.Group
	)
	v := new(struct {
		Page int `form:"pn" default:"0"`
		Size int `form:"ps" default:"20"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	db := actSrv.DB
	if v.Page == 0 {
		if err = db.Find(&list).Error; err != nil {
			log.Error("groupList error(%v)", err)
			c.JSON(nil, err)
			return
		}
		c.JSON(list, nil)
		return
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
	c.JSONMap(data, nil)
}

func groupInfo(c *bm.Context) {
	arg := new(model.Group)
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

func addGroup(c *bm.Context) {
	var err error
	arg := new(model.Group)
	if err := c.Bind(arg); err != nil {
		return
	}
	arg.CreatedTime = time.Now()
	// 图片base64编码
	// picPath := "./../../../JiqGstEfoWAOHiTxclqi.png"
	// if cc, err = pic.Base64Encoding(picPath); err != nil {
	// 	log.Error("pic.base64Encoding(%s,%d) error(%v)", picPath, arg.GroupID, err)
	// 	c.JSON(nil, err)
	// 	return
	// }
	// arg.GroupPhoto = cc
	if err = actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addGroup(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func saveGroup(c *bm.Context) {
	var (
		err error
		cc  string
	)
	arg := new(model.Group)
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.GroupID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	picPath := "./../../../JiqGstEfoWAOHiTxclqi.png"
	if cc, err = pic.Base64Encoding(picPath); err != nil {
		log.Error("pic.base64Encoding(%s,%d) error(%v)", picPath, arg.GroupID, err)
		c.JSON(nil, err)
		return
	}
	arg.GroupPhoto = cc
	// if err = actSrv.DB.Model(&model.Group{GroupID:arg.GroupID}).
	if err = actSrv.DB.Model(arg).
		Omit("group_id", "uid", "created_at").
		Updates(map[string]interface{}{
			"group_name":        arg.GroupName,
			"owner":             arg.Owner,
			"location":          arg.Location,
			"category":          arg.Category,
			"label":             arg.Label,
			"created_time":      arg.CTime,
			"free":              arg.IsFree,
			"group_scale":       arg.GroupScale,
			"group_description": arg.GroupDescription,
			"group_activity":    arg.GroupActivity,
			"group_photo":       arg.GroupPhoto,
			"owner_description": arg.OwnerDescription,
			"entry_method":      arg.EntryMethod,
			"group_rule":        arg.GroupRule,
			"group_detail":      arg.GroupDetail,
			"updated_at":        time.Now(),
		}).Error; err != nil {
		log.Error("saveGroup(%d) error(%v)", arg.GroupID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
