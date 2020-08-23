package service

import (
	"context"
	"go-online/app/interface/dao"
	"go-online/lib/conf/paladin"

	"github.com/jinzhu/gorm"
)

// Service biz service def.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
	DB  *gorm.DB
}

// New new a Service and return.
func New(d *dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
		DB:  d.DB,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// Ping check dao health.
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close close all dao.
func (s *Service) Close() {
}
