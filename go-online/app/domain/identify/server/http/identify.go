package http

import (
	v1 "go-online/app/domain/identify/api"
	"go-online/lib/ecode"
	bm "go-online/lib/net/http/blademaster"
)

// const (
// 	_actionChangePWD = "changePwd"
// 	_actionLoginOut  = "loginOut"
// )

func accessCookie(c *bm.Context) {
	cookie := c.Request.Header.Get("Cookie")
	if cookie == "" {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	// res, err := actSrv.GetCookieInfo(c, cookie)
	// if err == nil {
	// 	c.Set("mid", res.Mid)
	// }
	// c.JSON(res, err)
}

func accessToken(c *bm.Context) {
	token := new(v1.GetTokenInfoReq)
	if err := c.Bind(token); err != nil {
		c.JSON(nil, ecode.NoLogin)
		return
	}
	res, err := actSrv.GetTokenInfo(c, token)
	if err == nil {
		c.Set("mid", res.Mid)
	}
	c.JSON(res, err)
}

func delCache(c *bm.Context) {
	// query := c.Request.Form
	// action := query.Get("modifiedAttr")
	// if action != _actionChangePWD && action != _actionLoginOut {
	// 	return
	// }
	// key := query.Get("access_token")
	// if key == "" {
	// 	key = query.Get("session")
	// }
	// if key == "" {
	// 	return
	// }
	c.JSON(nil, nil)
}
