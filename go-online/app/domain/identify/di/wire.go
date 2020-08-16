// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"go-online/app/domain/identify/dao"
	"go-online/app/domain/identify/server/grpc"
	"go-online/app/domain/identify/server/http"
	"go-online/app/domain/identify/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
