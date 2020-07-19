package http

import (
	"fmt"
	"go-online/app/group/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/pic"
	"strconv"
	"time"
)

func groupList(c *bm.Context) {
	var (
		err   error
		count int
		list  []*model.Group
	)
	arg := new(model.GroupParam)
	if err = c.Bind(arg); err != nil {
		return
	}
	db := actSrv.DB
	if len(arg.GroupIDs) != 0 {
		db = db.Where("group_id in (?)", arg.GroupIDs)
	}
	if len(arg.UIDs) != 0 {
		db = db.Where("uid in (?)", arg.UIDs)
	}
	if len(arg.GroupNames) != 0 {
		db = db.Where("group_name in (?)", arg.GroupNames)
	}
	if len(arg.Owners) != 0 {
		db = db.Where("owner in (?)", arg.Owners)
	}
	if len(arg.Locations) != 0 {
		db = db.Where("location in (?)", arg.Locations)
	}
	if len(arg.Categorys) != 0 {
		db = db.Where("category in (?)", arg.Categorys)
	}
	if len(arg.Labels) != 0 {
		db = db.Where("label in (?)", arg.Labels)
	}
	if arg.IsFree != "" {
		if arg.IsFree == "是" {
			db = db.Where("free = ?", 1)
		} else if arg.IsFree == "否" {
			db = db.Where("free = ?", 0)
		}
	}

	if arg.Size == 0 { // 不分页
		if err = db.Find(&list).Error; err != nil {
			log.Error("groupList error(%v)", err)
			c.JSON(nil, err)
			return
		}
		c.JSON(list, nil)
		return
	}
	if arg.Page == 0 {
		arg.Page = 1
	}
	// 分页查看
	if err = db.
		Offset((arg.Page - 1) * arg.Size).Limit(arg.Size).Count(&count).
		Find(&list).Error; err != nil {
		log.Error("groupList(%d,%d) error(%v)", arg.Page, arg.Size, err)
		c.JSON(nil, err)
		return
	}

	data := map[string]interface{}{
		"data":  list,
		"pn":    arg.Page,
		"ps":    arg.Size,
		"total": count,
	}
	c.JSONMap(data, nil)
}

func groupInfo(c *bm.Context) {
	arg := new(struct {
		GroupID int64 `json:"group_id" form:"group_id" validate:"required,number,min=1"`
	})
	if err := c.Bind(arg); err != nil {
		return
	}
	if arg.GroupID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	group := new(model.Group)
	if err := actSrv.DB.Where("group_id = ?", arg.GroupID).First(group).Error; err != nil {
		log.Error("groupInfo(%d) error(%v)", arg.GroupID, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(group, nil)
}

func addGroup(c *bm.Context) {
	fmt.Println("addGroup")
	var (
		err error
		cc  string
	)
	arg := new(model.Group)
	if err := c.Bind(arg); err != nil {
		log.Error("addGroup Bind error(%v)", err)
		return
	}
	mid := c.Request.Header.Get("mid")
	userId, _ := strconv.ParseInt(mid, 10, 64)
	arg.UID = userId
	arg.CreatedTime = time.Now()
	// 图片base64编码
	if arg.GroupPhoto == "" {
		picPath := "./../../../JiqGstEfoWAOHiTxclqi.png"
		if cc, err = pic.Base64Encoding(picPath); err != nil {
			log.Error("pic.base64Encoding(%s,%d) error(%v)", picPath, arg.GroupID, err)
			c.JSON(nil, err)
			return
		}
		arg.GroupPhoto = cc
	}
	if err = actSrv.DB.Create(arg).Error; err != nil {
		log.Error("addGroup(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func saveGroup(c *bm.Context) {
	fmt.Println("saveGroup")
	var (
		err error
		cc  string
	)
	arg := new(model.Group)
	if err := c.Bind(arg); err != nil {
		log.Error("saveGroup Bind error(%v)", err)
		return
	}
	if arg.GroupID == -1 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	mid := c.Request.Header.Get("mid")
	userId, _ := strconv.ParseInt(mid, 10, 64)
	arg.UID = userId
	if arg.GroupPhoto == "" {
		picPath := "./../../../JiqGstEfoWAOHiTxclqi.png"
		if cc, err = pic.Base64Encoding(picPath); err != nil {
			log.Error("pic.base64Encoding(%s,%d) error(%v)", picPath, arg.GroupID, err)
			c.JSON(nil, err)
			return
		}
		arg.GroupPhoto = cc
	}
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

func deleteGroup(c *bm.Context) {
	var (
		err error
		uid int64
	)
	arg := new(struct {
		GroupId int `json:"group_id" form:"group_id" validate:"required,number,min=1"`
	})
	if err := c.Bind(arg); err != nil {
		log.Error("deleteGroup Bind error(%v)", err)
		return
	}
	mid := c.Request.Header.Get("mid")
	uid, _ = strconv.ParseInt(mid, 10, 64)
	if err = actSrv.DB.Where("uid = ? and group_id = ?", uid, arg.GroupId).
		Delete(model.Group{}).Error; err != nil {
		log.Error("deleteGroup delete group error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}
