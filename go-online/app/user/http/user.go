package http

import (
	"crypto/md5"
	"fmt"
	"go-online/app/user/model"
	"go-online/lib/ecode"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"

	// "go-online/lib/pic"
	"io"
	"strconv"
	"time"
)

func userRegister(c *bm.Context) {
	var (
		err   error
		count int
		user  *model.User
	)
	arg := new(model.User)
	if err = c.Bind(arg); err != nil {
		log.Error("userRegister error(%v)", err)
		return
	}
	if arg.Email == "" || arg.Password == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// find nothing about user who will register in the table of usertbl.
	db := actSrv.DB.Table("usertbl").Where("email = ?", arg.Email)
	if arg.NickName != "" {
		db = db.Or("nick_name = ?", arg.NickName)
	}
	if arg.LoginName != "" {
		db = db.Or("login_name = ?", arg.LoginName)
	}
	if err = db.Count(&count).Error; err != nil {
		log.Error("userRegister error(%v)", err)
		c.JSON(nil, err)
		return
	}
	if count != 0 {
		c.JSON(nil, ecode.AccountInexistence)
		return
	}
	arg.RegTime = time.Now()
	arg.CreatedAt = arg.RegTime
	arg.UpdatedAt = arg.CreatedAt
	// insert user record in the table of usertbl.
	if err = actSrv.DB.Create(arg).Error; err != nil {
		log.Error("userRegister(%v) error(%v)", arg, err)
		c.JSON(nil, err)
		return
	}
	// default permission
	user = arg
	if err = actSrv.DB.Model(user).Association("Roles").
		Append(model.Role{Id: 2, RoleName: "groupManager"}).Error; err != nil {
		log.Error("userRegister add default permission error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func userLogin(c *bm.Context) {
	var (
		err   error
		user  model.User
		token model.TokenObj
	)
	arg := new(struct {
		NickName  string `json:"nick_name" form:"nick_name"`
		Email     string `json:"email" form:"email" validate:"email"`
		LoginName string `json:"login_name" form:"login_name"`
		Password  string `json:"password" form:"password" validate:"required"`
	})
	if err = c.Bind(arg); err != nil {
		log.Error("userLogin error(%v)", err)
		return
	}
	if (arg.NickName == "" && arg.Email == "" && arg.LoginName == "") || arg.Password == "" {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	// get user id
	db := actSrv.DB.Table("usertbl").Where("password = ?", arg.Password)
	if arg.Email != "" {
		db = db.Where("email = ?", arg.Email)
	}
	if arg.NickName != "" {
		db = db.Where("nick_name = ?", arg.NickName)
	}
	if arg.LoginName != "" {
		db = db.Where("login_name = ?", arg.LoginName)
	}
	if err = db.First(&user).Error; err != nil {
		log.Error("userLogin find user error(%v)", err)
		c.JSON(nil, err)
		return
	}
	// generate token info
	token.UserId = user.Id
	token.Token = generateToken(user.Id)
	if err = actSrv.DB.Create(&token).Error; err != nil {
		log.Error("userLogin insert token(%v) error(%v)", token, err)
		c.JSON(nil, err)
		return
	}
	// return token and right list
	var roles []*model.Role
	if err = actSrv.DB.Model(&user).Association("Roles").Find(&roles).Error; err != nil {
		log.Error("userLogin query user roles error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var user_roles []*model.Role
	var role_id_list []int64
	for _, v := range roles {
		role_id_list = append(role_id_list, v.Id)
	}
	if err = actSrv.DB.Where("id in (?)", role_id_list).Preload("Rights").Find(&user_roles).Error; err != nil {
		log.Error("userLogin query user rights error(%v)", err)
		c.JSON(nil, err)
		return
	}
	var list []string
	for _, m := range user_roles {
		for _, v := range m.Rights {
			list = append(list, v.RightName)
		}
	}

	data := map[string]interface{}{
		"token":  token.Token,
		"rights": list,
	}
	c.JSON(data, nil)
}

func userLogout(c *bm.Context) {
	var (
		err    error
		userId int64
	)
	// delete token info
	mid := c.Request.Header.Get("mid")
	userId, _ = strconv.ParseInt(mid, 10, 64)
	if err = actSrv.DB.Where("user_id = ?", userId).
		Delete(model.TokenObj{}).Error; err != nil {
		log.Error("userLogout delete user token error(%v)", err)
		c.JSON(nil, err)
		return
	}
	c.JSON(nil, nil)
}

func generateToken(uid int64) string {
	currentTime := time.Now().Unix()
	h := md5.New()
	io.WriteString(h, strconv.FormatInt(currentTime+uid, 10))
	return fmt.Sprintf("%x", h.Sum(nil))
}
