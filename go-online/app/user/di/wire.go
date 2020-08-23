// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"go-online/app/user/dao"
	"go-online/app/user/server/http"
	"go-online/app/user/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.New, http.New, NewApp))
}
