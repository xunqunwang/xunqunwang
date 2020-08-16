package service

import (
	"context"
	pb "go-online/app/domain/identify/api"
	"go-online/app/domain/identify/dao"
	"go-online/lib/conf/paladin"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.IdentifyServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao *dao.Dao
	DB  *gorm.DB
}

// New new a service and return.
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

// Ping ping the resource.
func (s *Service) Ping(ctx context.Context) error {
	return s.dao.Ping(ctx)
}

// Close close the resource.
func (s *Service) Close() {
}
