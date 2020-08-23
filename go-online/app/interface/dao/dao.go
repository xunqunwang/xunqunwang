package dao

import (
	"context"

	"go-online/lib/conf/paladin"
	"go-online/lib/sync/pipeline/fanout"
	xtime "go-online/lib/time"
	"time"

	"github.com/google/wire"

	"github.com/jinzhu/gorm"
)

var Provider = wire.NewSet(New, NewDB)

// Dao struct user of Dao.
type Dao struct {
	DB     *gorm.DB
	cache  *fanout.Fanout
	expire int32
}

// New create a instance of Dao and return.
func New(db *gorm.DB) (d *Dao, cf func(), err error) {
	var cfg struct {
		Expire xtime.Duration
	}
	if err = paladin.Get("application.toml").UnmarshalTOML(&cfg); err != nil {
		return
	}
	d = &Dao{
		DB:     db,
		cache:  fanout.New("cache"),
		expire: int32(time.Duration(cfg.Expire) / time.Second),
	}
	cf = d.Close
	return
}

// Ping check connection of db , mc.
func (d *Dao) Ping(c context.Context) (err error) {
	return nil
}

// Close close connection of db , mc.
func (d *Dao) Close() {
	d.cache.Close()
}
