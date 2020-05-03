package http

import (
	"strconv"

	"go-common/app/service/main/archive/model/archive"
	"go-online/lib/ecode"
	"go-online/lib/log"
	bm "go-online/lib/net/http/blademaster"
	"go-online/lib/xstr"
)

// arcInfo write the archive data.
func arcInfo(c *bm.Context) {
	var (
		err error
		aid int64
	)
	params := c.Request.Form
	aidStr := params.Get("aid")
	// check params
	aid, err = strconv.ParseInt(aidStr, 10, 64)
	if err != nil {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.JSON(arcSvc.Archive3(c, aid))
}

// archives write the archives data.
func archives(c *bm.Context) {
	params := c.Request.Form
	aidsStr := params.Get("aids")
	// check params
	aids, err := xstr.SplitInts(aidsStr)
	if err != nil {
		log.Error("query aids(%s) split error(%v)", aidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if params.Get("appkey") == "fb06a25c6338edbc" && len(aids) > 50 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aids) > 50 {
		log.Error("Too many Args aids(%d) caller(%s)", len(aids), params.Get("appkey"))
	}
	c.JSON(arcSvc.Archives3(c, aids))
}

// archivesWithPlayer write the archives data.
func archivesWithPlayer(c *bm.Context) {
	params := c.Request.Form
	aidsStr := params.Get("aids")
	qnStr := params.Get("qn")
	pt := params.Get("platform")
	ip := params.Get("ip")
	fnver, _ := strconv.Atoi(params.Get("fnver"))
	fnval, _ := strconv.Atoi(params.Get("fnval"))
	forceHost, _ := strconv.Atoi(params.Get("force_host"))
	session := params.Get("session")
	containsPGC, _ := strconv.Atoi(params.Get("contains_pgc"))
	build, _ := strconv.Atoi(params.Get("build"))
	// check params
	aids, err := xstr.SplitInts(aidsStr)
	if err != nil {
		log.Error("query aids(%s) split error(%v)", aidsStr, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if len(aids) > 50 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	qn, _ := strconv.Atoi(qnStr)
	c.JSON(arcSvc.ArchivesWithPlayer(c, &archive.ArgPlayer{
		Aids:      aids,
		Qn:        qn,
		Platform:  pt,
		Build:     build,
		RealIP:    ip,
		Fnval:     fnval,
		Fnver:     fnver,
		Session:   session,
		ForceHost: forceHost,
	}, containsPGC == 1))
}

func typelist(c *bm.Context) {
	c.JSON(arcSvc.AllTypes(c), nil)
}

func maxAID(c *bm.Context) {
	c.JSON(arcSvc.MaxAID(c))
}
