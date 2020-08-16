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

type Dao struct {
	DB     *gorm.DB
	cache  *fanout.Fanout
	expire int32
}

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

// Close close the resource.
func (d *Dao) Close() {
	d.cache.Close()
}

// Ping ping the resource.
func (d *Dao) Ping(ctx context.Context) (err error) {
	return nil
}
