package dao

import (
	"context"

	"go-common/app/admin/main/manager/conf"
	"go-online/lib/cache/memcache"
	"go-online/lib/database/orm"
	bm "go-online/lib/net/http/blademaster"

	"github.com/jinzhu/gorm"
)

// Dao .
type Dao struct {
	db         *gorm.DB
	mc         *memcache.Pool
	httpClient *bm.Client
	dsbClient  *bm.Client
}

// New new a instance
func New(c *conf.Config) (d *Dao) {
	d = &Dao{
		db:         orm.NewMySQL(c.ORM),
		mc:         memcache.NewPool(c.Memcache),
		httpClient: bm.NewClient(c.HTTPClient),
		dsbClient:  bm.NewClient(c.DsbClient),
	}
	d.initORM()
	return
}

func (d *Dao) initORM() {
	d.db.LogMode(true)
}

// DB .
func (d *Dao) DB() *gorm.DB {
	return d.db
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	if d.db != nil {
		err = d.db.DB().PingContext(c)
	}
	return
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	if d.db != nil {
		d.db.Close()
	}
}
