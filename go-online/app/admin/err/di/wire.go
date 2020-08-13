// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"go-online/app/admin/err/dao"
	"go-online/app/admin/err/server/http"
	"go-online/app/admin/err/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.New, http.New, NewApp))
}
