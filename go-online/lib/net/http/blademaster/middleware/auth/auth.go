package auth

import (
	// idtv1 "go-common/app/service/main/identify/api/grpc"
	"fmt"
	idtv1 "go-online/app/domain/identify/api/grpc"
	"go-online/lib/ecode"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/net/metadata"
	"go-online/lib/net/rpc/warden"

	"github.com/pkg/errors"
)

// Config is the identify config model.
type Config struct {
	Identify *warden.ClientConfig
	// csrf switch.
	DisableCSRF bool
}

// Auth is the authorization middleware
type Auth struct {
	idtv1.IdentifyClient

	conf *Config
}

// authFunc will return mid and error by given context
type authFunc func(*bm.Context) (int64, error)

var _defaultConf = &Config{
	Identify:    nil,
	DisableCSRF: false,
}

// New is used to create an authorization middleware
func New(conf *Config) *Auth {
	if conf == nil {
		conf = _defaultConf
	}
	identify, err := idtv1.NewClient(conf.Identify)
	if err != nil {
		panic(errors.WithMessage(err, "Failed to dial identify service"))
	}
	fmt.Println(" idtv1.NewClient success")
	auth := &Auth{
		IdentifyClient: identify,
		conf:           conf,
	}
	return auth
}

// User is used to mark path as access required.
// If `access_key` is exist in request form, it will using mobile access policy.
// Otherwise to web access policy.
// func (a *Auth) User(ctx *bm.Context) {
// 	fmt.Println("Auth:User")
// 	req := ctx.Request
// 	if req.Form.Get("access_key") == "" {
// 		a.UserWeb(ctx)
// 		return
// 	}
// 	a.UserMobile(ctx)
// }

// modified by wangkai
func (a *Auth) User(ctx *bm.Context) {
	fmt.Println("Auth:User")
	req := ctx.Request
	if req.Header.Get("access_key") != "" {
		a.UserWeb(ctx)
		return
	}
	a.UserMobile(ctx)
}

// UserWeb is used to mark path as web access required.
// func (a *Auth) UserWeb(ctx *bm.Context) {
// 	fmt.Println("Auth:UserWeb")
// 	a.midAuth(ctx, a.AuthCookie)
// }

// modified by wangkai
func (a *Auth) UserWeb(ctx *bm.Context) {
	fmt.Println("Auth:UserWeb")
	// a.midAuth(ctx, a.AuthCookie)
	a.midAuth(ctx, a.AuthToken)
}

// UserMobile is used to mark path as mobile access required.
func (a *Auth) UserMobile(ctx *bm.Context) {
	fmt.Println("Auth:UserMobile")
	a.midAuth(ctx, a.AuthToken)
}

// Guest is used to mark path as guest policy.
// If `access_key` is exist in request form, it will using mobile access policy.
// Otherwise to web access policy.
func (a *Auth) Guest(ctx *bm.Context) {
	fmt.Println("Auth:Guest")
	req := ctx.Request
	if req.Form.Get("access_key") == "" {
		a.GuestWeb(ctx)
		return
	}
	a.GuestMobile(ctx)
}

// GuestWeb is used to mark path as web guest policy.
func (a *Auth) GuestWeb(ctx *bm.Context) {
	fmt.Println("Auth:GuestWeb")
	a.guestAuth(ctx, a.AuthCookie)
}

// GuestMobile is used to mark path as mobile guest policy.
func (a *Auth) GuestMobile(ctx *bm.Context) {
	a.guestAuth(ctx, a.AuthToken)
}

// AuthToken is used to authorize request by token
// func (a *Auth) AuthToken(ctx *bm.Context) (int64, error) {
// 	fmt.Println("AuthToken")
// 	req := ctx.Request
// 	key := req.Form.Get("access_key")
// 	if key == "" {
// 		return 0, ecode.NoLogin
// 	}
// 	buvid := req.Header.Get("buvid")

// 	reply, err := a.GetTokenInfo(ctx, &idtv1.GetTokenInfoReq{Token: key, Buvid: buvid})
// 	if err != nil {
// 		return 0, err
// 	}
// 	if !reply.IsLogin {
// 		return 0, ecode.NoLogin
// 	}

// 	return reply.Mid, nil
// }

// modified by wangkai
func (a *Auth) AuthToken(ctx *bm.Context) (int64, error) {
	fmt.Println("AuthToken")
	req := ctx.Request
	key := req.Header.Get("access_key")
	if key == "" {
		return 0, ecode.NoLogin
	}
	buvid := req.Header.Get("buvid") // buvid 移动端上报，再请求header里有，标识设备

	reply, err := a.GetTokenInfo(ctx, &idtv1.GetTokenInfoReq{Token: key, Buvid: buvid})
	if err != nil {
		return 0, err
	}

	if !reply.IsLogin {
		return 0, ecode.NoLogin
	}

	return reply.Mid, nil
}

// AuthCookie is used to authorize request by cookie
func (a *Auth) AuthCookie(ctx *bm.Context) (int64, error) {
	fmt.Println("Auth:AuthCookie")
	req := ctx.Request
	ssDaCk, _ := req.Cookie("SESSDATA")
	if ssDaCk == nil {
		return 0, ecode.NoLogin
	}

	cookie := req.Header.Get("Cookie")
	reply, err := a.GetCookieInfo(ctx, &idtv1.GetCookieInfoReq{Cookie: cookie})
	if err != nil {
		return 0, err
	}
	if !reply.IsLogin {
		return 0, ecode.NoLogin
	}

	// check csrf
	clientCsrf := req.FormValue("csrf")
	if a.conf != nil && !a.conf.DisableCSRF && req.Method == "POST" {
		if clientCsrf != reply.Csrf {
			return 0, ecode.CsrfNotMatchErr
		}
	}

	return reply.Mid, nil
}

func (a *Auth) midAuth(ctx *bm.Context, auth authFunc) {
	mid, err := auth(ctx)
	if err != nil {
		ctx.JSON(nil, err)
		ctx.Abort()
		return
	}
	ctx.Request.Header.Set("mid", fmt.Sprintf("%d", mid))
	setMid(ctx, mid)
}

func (a *Auth) guestAuth(ctx *bm.Context, auth authFunc) {
	fmt.Println("Auth:guestAuth")
	mid, err := auth(ctx)
	// no error happened and mid is valid
	if err == nil && mid > 0 {
		setMid(ctx, mid)
		return
	}

	ec := ecode.Cause(err)
	if ec.Equal(ecode.CsrfNotMatchErr) {
		ctx.JSON(nil, ec)
		ctx.Abort()
		return
	}
}

// set mid into context
// NOTE: This method is not thread safe.
func setMid(ctx *bm.Context, mid int64) {
	ctx.Set("mid", mid)
	if md, ok := metadata.FromContext(ctx); ok {
		md[metadata.Mid] = mid
		return
	}
}
