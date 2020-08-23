package dao

import (
	"go-online/lib/conf/paladin"
	"go-online/lib/database/orm"

	"github.com/jinzhu/gorm"
)

func NewDB() (db *gorm.DB, cf func(), err error) {
	var (
		cfg orm.Config
		ct  paladin.TOML
	)
	if err = paladin.Get("db.toml").Unmarshal(&ct); err != nil {
		return
	}
	if err = ct.Get("Client").UnmarshalTOML(&cfg); err != nil {
		return
	}
	db = orm.NewPostgreSQL(&cfg)
	cf = func() { db.Close() }
	return
}
